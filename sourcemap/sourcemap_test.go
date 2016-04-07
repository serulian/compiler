// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sourcemap

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type OriginalLocation struct {
	lineNumber     int
	columnPosition int
}

type sourceMapTest struct {
	name     string
	encoded  string
	mappings map[OriginalLocation]SourceMapping
}

var sourceMapTests = []sourceMapTest{
	sourceMapTest{
		name: "example",
		encoded: `{
		  "version": 3,
		  "file": "out.js",
		  "sources": [
		    "example.seru",
		    "github.com/Serulian/corelib:master"
		  ],
		  "names": [],
		  "mappings": "ACAA;ADkcE;AA4BsB;AACb"
		}`,
		mappings: map[OriginalLocation]SourceMapping{
			OriginalLocation{0, 0}: SourceMapping{"github.com/Serulian/corelib:master", 0, 0, ""},
			OriginalLocation{1, 0}: SourceMapping{"example.seru", 450, 2, ""},
			OriginalLocation{2, 0}: SourceMapping{"example.seru", 478, 24, ""},
			OriginalLocation{3, 0}: SourceMapping{"example.seru", 479, 11, ""},
		},
	},
}

func TestSourceMap(t *testing.T) {
	for _, test := range sourceMapTests {
		// Test decoding.
		parsed, err := Parse([]byte(test.encoded))
		if !assert.Nil(t, err, "Got unexpected parse error on test: %s", test.name) {
			continue
		}

		for original, mapping := range test.mappings {
			found, ok := parsed.LookupMapping(original.lineNumber, original.columnPosition)
			if !assert.True(t, ok, "Could not find original location %v in test %s", original, test.name) {
				continue
			}

			if !assert.Equal(t, mapping.SourcePath, found.SourcePath, "SourcePath mismatch in test %s for %v", test.name, original) {
				continue
			}

			if !assert.Equal(t, mapping.Name, found.Name, "Name mismatch in test %s for %v", test.name, original) {
				continue
			}

			if !assert.Equal(t, mapping.LineNumber, found.LineNumber, "LineNumber mismatch in test %s for %v", test.name, original) {
				continue
			}

			if !assert.Equal(t, mapping.ColumnPosition, found.ColumnPosition, "ColumnPosition mismatch in test %s for %v", test.name, original) {
				continue
			}
		}

		// Test encoding.
		sourceMap := NewSourceMap("out.js", "")
		for original, mapping := range test.mappings {
			sourceMap.AddMapping(original.lineNumber, original.columnPosition, mapping)
		}

		jsonValue, err := sourceMap.Build().Marshal()
		if !assert.Nil(t, err, "Error when marshaling source map in test %s", test.name) {
			continue
		}

		// Compare the JSON output.
		unmarshalled := json.RawMessage{}
		uerr := json.Unmarshal([]byte(test.encoded), &unmarshalled)
		if !assert.Nil(t, uerr, "Error when unmarshaling source map in test %s", test.name) {
			continue
		}

		marshalled, merr := json.MarshalIndent(&unmarshalled, "", "  ")
		if !assert.Nil(t, merr, "Error when marshaling encoded source map in test %s", test.name) {
			continue
		}

		assert.Equal(t, string(marshalled), string(jsonValue), "Encoded source map mismatch in test %s", test.name)
	}
}
