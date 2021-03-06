package blaze

import (
	"fmt"

	"github.com/cirss/geist/pkg/geist"
	"github.com/cirss/go-cli/pkg/cli"
)

func Report(cc *cli.CommandContext) (err error) {

	// declare command flags
	dataset := cc.Flags.String("dataset", "kb", "`name` of RDF dataset to create report from")
	file := cc.Flags.String("file", "-", "File containing report template to expand")
	includeinferred := cc.Flags.Bool("includeinferred", true, "Include inferred triples in result set")

	// parse flags
	var helped bool
	if helped, err = cc.ParseFlags(); helped || err != nil {
		return
	}

	bc := cc.Resource("BlazegraphClient").(*BlazegraphClient)
	bc.SetDataset(*dataset)
	bc.IncludeInferred = *includeinferred

	reportTemplate, err := cc.ReadFileOrStdin(*file)
	if err != nil {
		fmt.Fprintf(cc.ErrWriter, err.Error())
		return
	}

	rt := geist.NewTemplate("main", string(reportTemplate), nil, bc)

	report, err := bc.ExpandReport(rt)
	if err != nil {
		fmt.Fprintf(cc.ErrWriter, err.Error())
		return
	}

	fmt.Fprint(cc.OutWriter, report)
	return
}
