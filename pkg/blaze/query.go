package blaze

import (
	"fmt"

	"github.com/cirss/geist/pkg/geist"
	"github.com/cirss/go-cli/pkg/cli"
)

func Query(cc *cli.CommandContext) (err error) {

	// declare command flags
	dataset := cc.Flags.String("dataset", "kb", "`name` of RDF dataset to query")
	dryrun := cc.Flags.Bool("dryrun", false, "Output query but do not execute it")
	file := cc.Flags.String("file", "-", "File containing the SPARQL query to execute")
	format := cc.Flags.String("format", "json", "Format of result set to produce [csv, json, table, or xml]")
	separators := cc.Flags.Bool("columnseparators", true, "Display column separators in table format")
	includeinferred := cc.Flags.Bool("includeinferred", true, "Include inferred triples in result set")

	// parse flags
	var helped bool
	if helped, err = cc.ParseFlags(); helped || err != nil {
		return
	}

	bc := cc.Resource("BlazegraphClient").(*BlazegraphClient)
	bc.SetDataset(*dataset)
	bc.IncludeInferred = *includeinferred

	queryText, err := cc.ReadFileOrStdin(*file)
	if err != nil {
		fmt.Fprintf(cc.ErrWriter, err.Error())
		return
	}

	queryTemplate := geist.NewTemplate("query", string(queryText), nil, bc)
	err = queryTemplate.Parse()
	if err != nil {
		fmt.Fprintf(cc.ErrWriter, "Error expanding query template:\n")
		fmt.Fprintf(cc.ErrWriter, "%s\n", err.Error())
		return
	}

	q, err := queryTemplate.Expand(nil)

	if err != nil {
		fmt.Fprintf(cc.ErrWriter, "Error expanding query template: ")
		fmt.Fprintf(cc.ErrWriter, "%s\n", err.Error())
		return
	}

	if *dryrun {
		fmt.Print(string(q))
		return
	}

	switch *format {

	case "csv":
		resultCSV, _ := bc.SelectCSV(string(q))
		if err != nil {
			break
		}
		fmt.Fprintf(cc.OutWriter, resultCSV)
		return

	case "json":
		rs, e := bc.Select(string(q))
		err = e
		if err != nil {
			break
		}
		resultJSON, _ := rs.JSONString()
		fmt.Fprintf(cc.OutWriter, resultJSON)
		return

	case "table":
		rs, e := bc.Select(string(q))
		err = e
		if err != nil {
			break
		}
		table := rs.FormattedTable(*separators)
		fmt.Fprintf(cc.OutWriter, table)
		return

	case "xml":
		resultXML, e := bc.SelectXML(string(q))
		err = e
		if err != nil {
			break
		}
		fmt.Fprintf(cc.OutWriter, resultXML)
		return
	}

	if err != nil {
		fmt.Fprintf(cc.ErrWriter, err.Error())
	}

	return

}
