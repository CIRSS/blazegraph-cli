package main

import (
	"strings"
	"testing"

	"github.com/cirss/go-cli/pkg/util"
)

func TestBlazegraphCmd_destroy_default_dataset(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")
	Program.Invoke("blaze create --quiet")

	Program.AssertExitCode(t, "blaze destroy", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`Successfully destroyed dataset kb
		`)
}

func TestBlazegraphCmd_destroy_default_dataset_quiet(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")
	Program.Invoke("blaze create --quiet")

	Program.AssertExitCode(t, "blaze destroy --quiet", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		``)
}

func TestBlazegraphCmd_destroy_nonexistent_default_dataset(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")

	Program.AssertExitCode(t, "blaze destroy", 1)
	util.LineContentsEqual(t, outputBuffer.String(),
		`blaze destroy: destroy dataset failed: dataset kb does not exist
		`)
}

func TestBlazegraphCmd_destroy_nonexistent_default_dataset_quiet(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")

	Program.AssertExitCode(t, "blaze destroy --quiet", 1)
	util.LineContentsEqual(t, outputBuffer.String(),
		`blaze destroy: destroy dataset failed: dataset kb does not exist
		`)
}

func TestBlazegraphCmd_destroy_custom_dataset(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")
	Program.Invoke("blaze create --dataset foo --quiet")

	Program.AssertExitCode(t, "blaze destroy --dataset foo", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`Successfully destroyed dataset foo
		`)
}

func TestBlazegraphCmd_destroy_nonexistent_custom_dataset(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")

	Program.AssertExitCode(t, "blaze destroy --dataset foo", 1)
	util.LineContentsEqual(t, outputBuffer.String(),
		`blaze destroy: destroy dataset failed: dataset foo does not exist
		`)
}

func TestBlazegraphCmd_destroy_all_of_several_datasets(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")
	Program.Invoke("blaze create --dataset foo --quiet")
	Program.Invoke("blaze create --dataset bar --quiet")
	Program.Invoke("blaze create --dataset baz --quiet")

	Program.AssertExitCode(t, "blaze destroy --all", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`Successfully destroyed dataset bar
		 Successfully destroyed dataset baz
		 Successfully destroyed dataset foo
		`)
}

func TestBlazegraphCmd_destroy_one_of_several_datasets(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --all --quiet")
	Program.Invoke("blaze create --dataset foo --quiet")
	Program.Invoke("blaze create --dataset bar --quiet")
	Program.Invoke("blaze create --dataset baz --quiet")

	Program.AssertExitCode(t, "blaze destroy --dataset bar", 0)
	util.LineContentsEqual(t, outputBuffer.String(),
		`Successfully destroyed dataset bar
		`)

	outputBuffer.Reset()
	Program.Invoke("blaze list")
	util.LineContentsEqual(t, outputBuffer.String(),
		`baz        0
		 foo        0
		`)
}

func TestBlazegraphCmd_destroy_missing_dataset_name(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.AssertExitCode(t, "blaze destroy --dataset", 1)
	util.LineContentsEqual(t, outputBuffer.String(),
		`blaze destroy: flag needs an argument: -dataset

        usage: blaze destroy [<flags>]

        flags:
        	-all
            	destroy ALL datasets in the Blazegraph instance
          	-dataset name
            	name of RDF dataset to destroy (default "kb")
          	-instance URL
            	URL of Blazegraph instance (default "http://127.0.0.1:9999/blazegraph")
          	-quiet
            	Discard normal command output
		 	-silent
				Discard normal and error command output
		`)
}

var expectedDestroyHelpOutput = string(
	`blaze destroy: Deletes an RDF dataset and corresponding Blazegraph namespace, all RDF graphs
        in the dataset, and all triples in each of those graphs.

        usage: blaze destroy [<flags>]

        flags:
			-all
				destroy ALL datasets in the Blazegraph instance
			-dataset name
				name of RDF dataset to destroy (default "kb")
			-instance URL
				URL of Blazegraph instance (default "http://127.0.0.1:9999/blazegraph")
			-quiet
				Discard normal command output
			-silent
				Discard normal and error command output
		`)

func TestBlazegraphCmd_destroy_help(t *testing.T) {
	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer
	Program.AssertExitCode(t, "blaze destroy help", 0)
	util.LineContentsEqual(t, outputBuffer.String(), expectedDestroyHelpOutput)
}

func TestBlazegraphCmd_help_destroy(t *testing.T) {
	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer
	Program.AssertExitCode(t, "blaze help destroy", 0)
	util.LineContentsEqual(t, outputBuffer.String(), expectedDestroyHelpOutput)
}

func TestBlazegraphCmd_destroy_no_dataset_argument(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.AssertExitCode(t, "blaze destroy -dataset", 1)
	util.LineContentsEqual(t, outputBuffer.String(),
		`blaze destroy: flag needs an argument: -dataset

		usage: blaze destroy [<flags>]

		flags:
			-all
					destroy ALL datasets in the Blazegraph instance
			-dataset name
					name of RDF dataset to destroy (default "kb")
			-instance URL
					URL of Blazegraph instance (default "http://127.0.0.1:9999/blazegraph")
			-quiet
					Discard normal command output
			-silent
					Discard normal and error command output
	`)
}

func TestBlazegraphCmd_destroy_bad_flag(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.AssertExitCode(t, "blaze destroy --not-a-flag", 1)
	util.LineContentsEqual(t, outputBuffer.String(),
		`blaze destroy: flag provided but not defined: -not-a-flag

		usage: blaze destroy [<flags>]

		flags:
			-all
					destroy ALL datasets in the Blazegraph instance
			-dataset name
					name of RDF dataset to destroy (default "kb")
			-instance URL
					URL of Blazegraph instance (default "http://127.0.0.1:9999/blazegraph")
			-quiet
					Discard normal command output
			-silent
					Discard normal and error command output
	`)
}
