// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type importSnippetTest struct {
	snippetString   string
	expectedValid   bool
	expectedSnippet importSnippet
}

var importSnippetTests = []importSnippetTest{
	// Success tests.
	importSnippetTest{"import ", true, importSnippet{"import", "", "", importKindPackage}},
	importSnippetTest{"import foobar.", true, importSnippet{"import", "", "foobar.", importKindPackage}},
	importSnippetTest{"from ", true, importSnippet{"from", "", "", importKindPackage}},
	importSnippetTest{"from foobar.", true, importSnippet{"from", "", "foobar.", importKindPackage}},
	importSnippetTest{"from foobar import ", true, importSnippet{"from", "", "foobar", importKindFromPackage}},
	importSnippetTest{"from foobar.baz import ", true, importSnippet{"from", "", "foobar.baz", importKindFromPackage}},
	importSnippetTest{"from \"foobar\" import ", true, importSnippet{"from", "", "\"foobar\"", importKindFromPackage}},
	importSnippetTest{"from \"a/b/c/d\" import ", true, importSnippet{"from", "", "\"a/b/c/d\"", importKindFromPackage}},
	importSnippetTest{"from webidl`foobarbaz` import ", true, importSnippet{"from", "webidl", "`foobarbaz`", importKindFromPackage}},
	importSnippetTest{"from webidl`foo/bar/baz` import ", true, importSnippet{"from", "webidl", "`foo/bar/baz`", importKindFromPackage}},

	importSnippetTest{"import foo", true, importSnippet{"import", "", "foo", importKindPackage}},
	importSnippetTest{"from foo", true, importSnippet{"from", "", "foo", importKindPackage}},

	// Failure tests.
	importSnippetTest{"", false, importSnippet{}},
	importSnippetTest{"import", false, importSnippet{}},
	importSnippetTest{"impart ", false, importSnippet{}},
	importSnippetTest{"from", false, importSnippet{}},
	importSnippetTest{"fram", false, importSnippet{}},
}

func TestImportSnippets(t *testing.T) {
	for _, test := range importSnippetTests {
		snippet, ok := buildImportSnippet(test.snippetString)
		if !assert.Equal(t, ok, test.expectedValid, "Mismatched validity on test %s", test.snippetString) {
			continue
		}

		if !test.expectedValid {
			continue
		}

		if !assert.Equal(t, snippet, test.expectedSnippet, "Mismatched snippet on test %s", test.snippetString) {
			continue
		}
	}
}
