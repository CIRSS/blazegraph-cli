package main

import (
	"strings"
	"testing"

	"github.com/cirss/go-cli/pkg/util"
)

func TestBlazegraphCmd_static_macro_in_select(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.InReader = strings.NewReader(`
		<:x> <:foo> <:y> .
		<:y> <:bar> <:z> .
		<:z> <:baz> "baz" .
	`)
	Program.Invoke("blaze import --format ttl")

	outputBuffer.Reset()
	template := `

		{{{
			{{ macro "foo" '''
				<:foo>
			''' }}
		}}}

		SELECT DISTINCT ?s ?o
		WHERE
		{ ?s {{foo}} ?o }
		ORDER BY ?o
	`
	Program.InReader = strings.NewReader(template)
	Program.AssertExitCode(t, "blaze query --format table", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`s                                                | o
		 ===================================================================================================
		 http://127.0.0.1:9999/blazegraph/namespace/kb/:x | http://127.0.0.1:9999/blazegraph/namespace/kb/:y
		`)
}

func TestBlazegraphCmd_included_static_macro_in_select(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.InReader = strings.NewReader(`
		<:x> <:foo> <:y> .
		<:y> <:bar> <:z> .
		<:z> <:baz> "baz" .
	`)
	Program.Invoke("blaze import --format ttl")

	outputBuffer.Reset()
	template := `

		{{{
			{{ include "./testdata/rules.g" }}
		}}}

		SELECT DISTINCT ?s ?o
		WHERE
		{ ?s {{foo}} ?o }
		ORDER BY ?o
	`
	Program.InReader = strings.NewReader(template)
	Program.AssertExitCode(t, "blaze query --format table", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`s                                                | o
		===================================================================================================
		http://127.0.0.1:9999/blazegraph/namespace/kb/:x | http://127.0.0.1:9999/blazegraph/namespace/kb/:y
		`)
}

func TestBlazegraphCmd_dynamic_macro_in_select(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.InReader = strings.NewReader(`
		<:x> <:foo> <:y> .
		<:y> <:bar> <:z> .
		<:z> <:baz> "baz" .
	`)
	Program.Invoke("blaze import --format ttl")

	outputBuffer.Reset()
	template := `

		{{{
			{{ macro "bar" "Sub" "Obj" '''
				{{_subject $Sub}} <:bar> {{_object $Obj}}
			''' }}
		}}}

		SELECT DISTINCT ?s ?o
		WHERE
		{ {{ bar "?s" "?o" }} }
		ORDER BY ?o
	`
	Program.InReader = strings.NewReader(template)
	Program.AssertExitCode(t, "blaze query --format table", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`s                                                | o
		===================================================================================================
		http://127.0.0.1:9999/blazegraph/namespace/kb/:y | http://127.0.0.1:9999/blazegraph/namespace/kb/:z
	   `)
}

func TestBlazegraphCmd_included_dynamic_macro_in_select(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.InReader = strings.NewReader(`
		<:x> <:foo> <:y> .
		<:y> <:bar> <:z> .
		<:z> <:baz> "baz" .
	`)
	Program.Invoke("blaze import --format ttl")

	outputBuffer.Reset()
	template := `

		{{{
			{{ include "./testdata/rules.g" }}
		}}}

		SELECT DISTINCT ?s ?o
		WHERE
		{ {{ bar "?s" "?o" }} }
		ORDER BY ?o
	`
	Program.InReader = strings.NewReader(template)
	Program.AssertExitCode(t, "blaze query --format table", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`s                                                | o
		===================================================================================================
		http://127.0.0.1:9999/blazegraph/namespace/kb/:y | http://127.0.0.1:9999/blazegraph/namespace/kb/:z
		`)
}

func TestBlazegraphCmd_rule_in_select(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.InReader = strings.NewReader(`
		<:x> <:foo> <:y> .
		<:y> <:bar> <:z> .
		<:z> <:baz> "baz" .
	`)
	Program.Invoke("blaze import --format ttl")

	outputBuffer.Reset()
	template := `

		{{{
			{{ rule "foo_bar_baz" "s" "o" '''
				{{_subject $s}} <:foo> ?y .
				?y <:bar> ?z .
				?z <:baz> {{_object $o}} .
			''' }}
		}}}

		SELECT DISTINCT ?o
		WHERE
		{ {{ foo_bar_baz "?s" "?o" }} }
		ORDER BY ?o
	`
	Program.InReader = strings.NewReader(template)
	Program.AssertExitCode(t, "blaze query --format table", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`o
		==
		baz
	`)
}

func TestBlazegraphCmd_included_rule_in_select(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.InReader = strings.NewReader(`
		<:x> <:foo> <:y> .
		<:y> <:bar> <:z> .
		<:z> <:baz> "baz" .
	`)
	Program.Invoke("blaze import --format ttl")

	outputBuffer.Reset()
	template := `

		{{{
			{{ include "./testdata/rules.g" }}
		}}}

		SELECT DISTINCT ?o
		WHERE
		{ {{ foo_bar_baz "?s" "?o" }} }
		ORDER BY ?o
	`
	Program.InReader = strings.NewReader(template)
	Program.AssertExitCode(t, "blaze query --format table", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`o
		==
		baz
	`)
}

func TestBlazegraphCmd_rule_in_select_in_report(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.InReader = strings.NewReader(`
		<:x> <:foo> <:y> .
		<:y> <:bar> <:z> .
		<:z> <:baz> "baz" .
	`)
	Program.Invoke("blaze import --format ttl")

	outputBuffer.Reset()
	template := `
		{{{
			{{ rule "foo_bar_baz" "s" "o" '''
				{{_subject $s}} <:foo> ?y .
				?y <:bar> ?z .
				?z <:baz> {{_object $o}} .
			''' }}
		}}}

		{{ select '''
			SELECT DISTINCT ?o
			WHERE
			{ {{ foo_bar_baz "?s" "?o" }} }
			ORDER BY ?o
		''' | value }}
	`
	Program.InReader = strings.NewReader(template)
	Program.AssertExitCode(t, "blaze report", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`
		baz
	`)
}

func TestBlazegraphCmd_rule_in_query(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.InReader = strings.NewReader(`
		<:x> <:foo> <:y> .
		<:y> <:bar> <:z> .
		<:z> <:baz> "baz" .
	`)
	Program.Invoke("blaze import --format ttl")

	outputBuffer.Reset()
	template := `
		{{{
			{{ rule "foo_bar_baz" "s" "o" '''
				{{_subject $s}} <:foo> ?y .
				?y <:bar> ?z .
				?z <:baz> {{_object $o}} .
			''' }}

			{{query "select_foo_bar_baz" '''
				SELECT DISTINCT ?o
				WHERE
				{ {{ foo_bar_baz "?s" "?o" }} }
				ORDER BY ?o
			''' }}
		}}}

		{{ select_foo_bar_baz ":x" | value }}
`
	Program.InReader = strings.NewReader(template)
	Program.AssertExitCode(t, "blaze report", 0)
	util.LineContentsEqual(t, outputBuffer.String(), `
		baz
	`)
}

func TestBlazegraphCmd_rule_in_rule(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.InReader = strings.NewReader(`
		<:x> <:foo> <:y> .
		<:y> <:bar> <:z> .
		<:z> <:baz> "baz" .
	`)
	Program.Invoke("blaze import --format ttl")

	outputBuffer.Reset()
	template := `
		{{{
			{{ rule "foo_bar_baz_rule_1" "s" "o" '''
				{{_subject $s}} <:foo> ?y .
				?y <:bar> ?z .
				?z <:baz> {{_object $o}} .
			''' }}

			{{ rule "foo_bar_baz_rule_2" "s" "o" '''
				{{ foo_bar_baz_rule_1 $s $o }}
			'''}}

			{{query "foo_bar_baz_query" '''
				SELECT DISTINCT ?s ?o
				WHERE
				{ {{ foo_bar_baz_rule_2 "?s" "?o" }} }
				ORDER BY ?o
			''' }}

		}}}
		{{ foo_bar_baz_query ":x" | tabulate }}
`
	Program.InReader = strings.NewReader(template)
	Program.AssertExitCode(t, "blaze report", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`s                                                | o
	 	 ======================================================
		 http://127.0.0.1:9999/blazegraph/namespace/kb/:x | baz

	`)
}

func TestBlazegraphCmd_rule_in_query_called_by_macro(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.InReader = strings.NewReader(`
		<:y> <:tag> "eight" .
		<:x> <:tag> "seven" .
	`)
	Program.Invoke("blaze import --format ttl")

	outputBuffer.Reset()
	template := `
		{{{
			{{ rule "hasTag" "s" "o" '''
				{{_subject $s}} <:tag> {{_object $o}}
			''' }}

			{{query "select_subjects" '''
				SELECT DISTINCT ?s
				WHERE
				{ {{ hasTag "?s" "?o" }} }
				ORDER BY ?s
			''' }}

			{{query "select_tags_for_subject" "Subject" '''
				SELECT ?tag
				WHERE { {{ hasTag $Subject "?tag" }} }
				ORDER BY ?tag
			''' }}

			{{macro "tabulate_tags_for_subject" "Subject" '''
				{{ select_tags_for_subject $Subject | tabulate }}
			''' }}
		}}}
																	\
		{{ range $Subject := select_subjects | vector }}
			{{ tabulate_tags_for_subject $Subject }}
		{{ end }}
`
	Program.InReader = strings.NewReader(template)
	Program.AssertExitCode(t, "blaze report", 0)
	util.LineContentsEqual(t, outputBuffer.String(), `
		tag
		====
		seven

		tag
		====
		eight

	`)
}
