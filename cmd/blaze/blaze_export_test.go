package main

import (
	"strings"
	"testing"

	"github.com/cirss/go-cli/pkg/util"
)

func TestBlazegraphCmd_export_empty_store(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	t.Run("jsonld", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format jsonld", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`[ ]
			`)
	})

	t.Run("nt", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format nt", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			``)
	})

	t.Run("ttl", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format ttl", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
			@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
			@prefix sesame: <http://www.openrdf.org/schema/sesame#> .
			@prefix owl: <http://www.w3.org/2002/07/owl#> .
			@prefix xsd: <http://www.w3.org/2001/XMLSchema#> .
			@prefix fn: <http://www.w3.org/2005/xpath-functions#> .
			@prefix foaf: <http://xmlns.com/foaf/0.1/> .
			@prefix dc: <http://purl.org/dc/elements/1.1/> .
			@prefix hint: <http://www.bigdata.com/queryHints#> .
			@prefix bd: <http://www.bigdata.com/rdf#> .
			@prefix bds: <http://www.bigdata.com/rdf/search#> .
			`)
	})

	t.Run("xml", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format xml", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`<?xml version="1.0" encoding="UTF-8"?>
	        <rdf:RDF
	        	xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	        	xmlns:rdfs="http://www.w3.org/2000/01/rdf-schema#"
	        	xmlns:sesame="http://www.openrdf.org/schema/sesame#"
	        	xmlns:owl="http://www.w3.org/2002/07/owl#"
	        	xmlns:xsd="http://www.w3.org/2001/XMLSchema#"
	        	xmlns:fn="http://www.w3.org/2005/xpath-functions#"
	        	xmlns:foaf="http://xmlns.com/foaf/0.1/"
	        	xmlns:dc="http://purl.org/dc/elements/1.1/"
	        	xmlns:hint="http://www.bigdata.com/queryHints#"
	        	xmlns:bd="http://www.bigdata.com/rdf#"
	        	xmlns:bds="http://www.bigdata.com/rdf/search#">

			</rdf:RDF>
			`)
	})
}

func TestBlazegraphCmd_export_two_triples(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.InReader = strings.NewReader(`
		<http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
		<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
	`)
	Program.Invoke("blaze import --format ttl")

	t.Run("jsonld", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format jsonld", 0)
		util.JSONEquals(t, outputBuffer.String(),
			`[
			{
			  "@id": "http://tmcphill.net/data#x",
			  "http://tmcphill.net/tags#tag": [
				{
				  "@value": "seven"
				}
			  ]
			},
			{
			  "@id": "http://tmcphill.net/data#y",
			  "http://tmcphill.net/tags#tag": [
				{
				  "@value": "eight"
				}
			  ]
			}
		  ]
		`)
	})

	t.Run("nt", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format nt --sort=true", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
			<http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
			`)
	})

	t.Run("ttl", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format ttl", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
			@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
			@prefix sesame: <http://www.openrdf.org/schema/sesame#> .
			@prefix owl: <http://www.w3.org/2002/07/owl#> .
			@prefix xsd: <http://www.w3.org/2001/XMLSchema#> .
			@prefix fn: <http://www.w3.org/2005/xpath-functions#> .
			@prefix foaf: <http://xmlns.com/foaf/0.1/> .
			@prefix dc: <http://purl.org/dc/elements/1.1/> .
			@prefix hint: <http://www.bigdata.com/queryHints#> .
			@prefix bd: <http://www.bigdata.com/rdf#> .
			@prefix bds: <http://www.bigdata.com/rdf/search#> .

			<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .

			<http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
			`)
	})

	t.Run("xml", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format xml", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`<?xml version="1.0" encoding="UTF-8"?>
            <rdf:RDF
                xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
                xmlns:rdfs="http://www.w3.org/2000/01/rdf-schema#"
                xmlns:sesame="http://www.openrdf.org/schema/sesame#"
                xmlns:owl="http://www.w3.org/2002/07/owl#"
                xmlns:xsd="http://www.w3.org/2001/XMLSchema#"
                xmlns:fn="http://www.w3.org/2005/xpath-functions#"
                xmlns:foaf="http://xmlns.com/foaf/0.1/"
                xmlns:dc="http://purl.org/dc/elements/1.1/"
                xmlns:hint="http://www.bigdata.com/queryHints#"
                xmlns:bd="http://www.bigdata.com/rdf#"
                xmlns:bds="http://www.bigdata.com/rdf/search#">

            <rdf:Description rdf:about="http://tmcphill.net/data#x">
                <tag xmlns="http://tmcphill.net/tags#">seven</tag>
            </rdf:Description>

            <rdf:Description rdf:about="http://tmcphill.net/data#y">
                <tag xmlns="http://tmcphill.net/tags#">eight</tag>
            </rdf:Description>

            </rdf:RDF>
			`)
	})

}

func TestBlazegraphCmd_export_address_book(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb")

	Program.Invoke("blaze import --file testdata/address-book.jsonld --format jsonld")

	t.Run("jsonld", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format jsonld", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`[ {
			"@id" : "http://learningsparql.com/ns/data#i0432",
			"http://learningsparql.com/ns/addressbook#email" : [ {
				"@value" : "richard49@hotmail.com"
			} ],
			"http://learningsparql.com/ns/addressbook#firstname" : [ {
				"@value" : "Richard"
			} ],
			"http://learningsparql.com/ns/addressbook#homeTel" : [ {
				"@value" : "(229) 276-5135"
			} ],
			"http://learningsparql.com/ns/addressbook#lastname" : [ {
				"@value" : "Mutt"
			} ],
			"http://learningsparql.com/ns/addressbook#nickname" : [ {
				"@value" : "Dick"
			} ]
			}, {
			"@id" : "http://learningsparql.com/ns/data#i8301",
			"http://learningsparql.com/ns/addressbook#email" : [ {
				"@value" : "c.ellis@usairwaysgroup.com"
			}, {
				"@value" : "craigellis@yahoo.com"
			} ],
			"http://learningsparql.com/ns/addressbook#firstname" : [ {
				"@value" : "Craig"
			} ],
			"http://learningsparql.com/ns/addressbook#homeTel" : [ {
				"@value" : "(194) 966-1505"
			} ],
			"http://learningsparql.com/ns/addressbook#lastname" : [ {
				"@value" : "Ellis"
			} ]
			}, {
			"@id" : "http://learningsparql.com/ns/data#i9771",
			"http://learningsparql.com/ns/addressbook#email" : [ {
				"@value" : "cindym@gmail.com"
			} ],
			"http://learningsparql.com/ns/addressbook#firstname" : [ {
				"@value" : "Cindy"
			} ],
			"http://learningsparql.com/ns/addressbook#homeTel" : [ {
				"@value" : "(245) 646-5488"
			} ],
			"http://learningsparql.com/ns/addressbook#lastname" : [ {
				"@value" : "Marshall"
			} ],
			"http://learningsparql.com/ns/addressbook#mobileTel" : [ {
				"@value" : "(245) 732-8991"
			} ]
			} ]
			`)
	})

	t.Run("nt", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format nt", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`<http://learningsparql.com/ns/data#i0432> <http://learningsparql.com/ns/addressbook#email> "richard49@hotmail.com"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i0432> <http://learningsparql.com/ns/addressbook#firstname> "Richard"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i0432> <http://learningsparql.com/ns/addressbook#homeTel> "(229) 276-5135"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i0432> <http://learningsparql.com/ns/addressbook#lastname> "Mutt"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i0432> <http://learningsparql.com/ns/addressbook#nickname> "Dick"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i8301> <http://learningsparql.com/ns/addressbook#email> "c.ellis@usairwaysgroup.com"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i8301> <http://learningsparql.com/ns/addressbook#email> "craigellis@yahoo.com"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i8301> <http://learningsparql.com/ns/addressbook#firstname> "Craig"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i8301> <http://learningsparql.com/ns/addressbook#homeTel> "(194) 966-1505"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i8301> <http://learningsparql.com/ns/addressbook#lastname> "Ellis"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i9771> <http://learningsparql.com/ns/addressbook#email> "cindym@gmail.com"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i9771> <http://learningsparql.com/ns/addressbook#firstname> "Cindy"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i9771> <http://learningsparql.com/ns/addressbook#homeTel> "(245) 646-5488"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i9771> <http://learningsparql.com/ns/addressbook#lastname> "Marshall"^^<http://www.w3.org/2001/XMLSchema#string> .
			<http://learningsparql.com/ns/data#i9771> <http://learningsparql.com/ns/addressbook#mobileTel> "(245) 732-8991"^^<http://www.w3.org/2001/XMLSchema#string> .
		`)
	})

	t.Run("ttl", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format ttl", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
			@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
			@prefix sesame: <http://www.openrdf.org/schema/sesame#> .
			@prefix owl: <http://www.w3.org/2002/07/owl#> .
			@prefix xsd: <http://www.w3.org/2001/XMLSchema#> .
			@prefix fn: <http://www.w3.org/2005/xpath-functions#> .
			@prefix foaf: <http://xmlns.com/foaf/0.1/> .
			@prefix dc: <http://purl.org/dc/elements/1.1/> .
			@prefix hint: <http://www.bigdata.com/queryHints#> .
			@prefix bd: <http://www.bigdata.com/rdf#> .
			@prefix bds: <http://www.bigdata.com/rdf/search#> .

			<http://learningsparql.com/ns/data#i0432> <http://learningsparql.com/ns/addressbook#email> "richard49@hotmail.com"^^xsd:string ;
				<http://learningsparql.com/ns/addressbook#firstname> "Richard"^^xsd:string ;
				<http://learningsparql.com/ns/addressbook#homeTel> "(229) 276-5135"^^xsd:string ;
				<http://learningsparql.com/ns/addressbook#lastname> "Mutt"^^xsd:string ;
				<http://learningsparql.com/ns/addressbook#nickname> "Dick"^^xsd:string .

			<http://learningsparql.com/ns/data#i8301> <http://learningsparql.com/ns/addressbook#email> "c.ellis@usairwaysgroup.com"^^xsd:string , "craigellis@yahoo.com"^^xsd:string ;
				<http://learningsparql.com/ns/addressbook#firstname> "Craig"^^xsd:string ;
				<http://learningsparql.com/ns/addressbook#homeTel> "(194) 966-1505"^^xsd:string ;
				<http://learningsparql.com/ns/addressbook#lastname> "Ellis"^^xsd:string .

			<http://learningsparql.com/ns/data#i9771> <http://learningsparql.com/ns/addressbook#email> "cindym@gmail.com"^^xsd:string ;
				<http://learningsparql.com/ns/addressbook#firstname> "Cindy"^^xsd:string ;
				<http://learningsparql.com/ns/addressbook#homeTel> "(245) 646-5488"^^xsd:string ;
				<http://learningsparql.com/ns/addressbook#lastname> "Marshall"^^xsd:string ;
				<http://learningsparql.com/ns/addressbook#mobileTel> "(245) 732-8991"^^xsd:string .
			`)
	})

	t.Run("xml", func(t *testing.T) {
		outputBuffer.Reset()
		Program.AssertExitCode(t, "blaze export --format xml", 0)
		util.LineContentsEqual(t, outputBuffer.String(),
			`<?xml version="1.0" encoding="UTF-8"?>
			<rdf:RDF
				xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
				xmlns:rdfs="http://www.w3.org/2000/01/rdf-schema#"
				xmlns:sesame="http://www.openrdf.org/schema/sesame#"
				xmlns:owl="http://www.w3.org/2002/07/owl#"
				xmlns:xsd="http://www.w3.org/2001/XMLSchema#"
				xmlns:fn="http://www.w3.org/2005/xpath-functions#"
				xmlns:foaf="http://xmlns.com/foaf/0.1/"
				xmlns:dc="http://purl.org/dc/elements/1.1/"
				xmlns:hint="http://www.bigdata.com/queryHints#"
				xmlns:bd="http://www.bigdata.com/rdf#"
				xmlns:bds="http://www.bigdata.com/rdf/search#">

			<rdf:Description rdf:about="http://learningsparql.com/ns/data#i0432">
				<email xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">richard49@hotmail.com</email>
				<firstname xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">Richard</firstname>
				<homeTel xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">(229) 276-5135</homeTel>
				<lastname xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">Mutt</lastname>
				<nickname xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">Dick</nickname>
			</rdf:Description>

			<rdf:Description rdf:about="http://learningsparql.com/ns/data#i8301">
				<email xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">c.ellis@usairwaysgroup.com</email>
				<email xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">craigellis@yahoo.com</email>
				<firstname xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">Craig</firstname>
				<homeTel xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">(194) 966-1505</homeTel>
				<lastname xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">Ellis</lastname>
			</rdf:Description>

			<rdf:Description rdf:about="http://learningsparql.com/ns/data#i9771">
				<email xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">cindym@gmail.com</email>
				<firstname xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">Cindy</firstname>
				<homeTel xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">(245) 646-5488</homeTel>
				<lastname xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">Marshall</lastname>
				<mobileTel xmlns="http://learningsparql.com/ns/addressbook#" rdf:datatype="http://www.w3.org/2001/XMLSchema#string">(245) 732-8991</mobileTel>
			</rdf:Description>

			</rdf:RDF>
			`)
	})
}

var expectedExportHelpOutput = string(
	`blaze export: Exports all triples in an RDF dataset in the requested format.

	usage: blaze export [<flags>]

	flags:
		-dataset name
				name of RDF dataset to export (default "kb")
		-format string
				Format for exported triples [jsonld, nt, ttl, or xml] (default "nt")
		-includeinferred
			Include inferred triples in result set (default true)
		-instance URL
				URL of Blazegraph instance (default "http://127.0.0.1:9999/blazegraph")
		-quiet
				Discard normal command output
		-silent
				Discard normal and error command output
		-sort
				Sort the exported triples if true
	`)

func TestBlazegraphCmd_export_help(t *testing.T) {
	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer
	Program.AssertExitCode(t, "blaze export help", 0)
	util.LineContentsEqual(t, outputBuffer.String(), expectedExportHelpOutput)
}

func TestBlazegraphCmd_help_export(t *testing.T) {
	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer
	Program.AssertExitCode(t, "blaze help export", 0)
	util.LineContentsEqual(t, outputBuffer.String(), expectedExportHelpOutput)
}

func TestBlazegraphCmd_export_bad_flag(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.AssertExitCode(t, "blaze export --not-a-flag", 1)
	util.LineContentsEqual(t, outputBuffer.String(),
		`blaze export: flag provided but not defined: -not-a-flag

		usage: blaze export [<flags>]

		flags:
			-dataset name
					name of RDF dataset to export (default "kb")
			-format string
					Format for exported triples [jsonld, nt, ttl, or xml] (default "nt")
			-includeinferred
				Include inferred triples in result set (default true)
			-instance URL
				URL of Blazegraph instance (default "http://127.0.0.1:9999/blazegraph")
			-quiet
					Discard normal command output
			-silent
					Discard normal and error command output
			-sort
					Sort the exported triples if true
		`)
}
