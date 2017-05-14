// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/serulian/compiler/packageloader"
)

var _ = fmt.Printf

type searchTest struct {
	name     string
	subtests []searchSubTest
}

type searchSubTest struct {
	query           string
	expectedSymbols []expectedSymbol
}

type expectedSymbol struct {
	name string
	kind SymbolKind
}

var searchTests = []searchTest{
	searchTest{
		"basic",
		[]searchSubTest{
			searchSubTest{"basi", []expectedSymbol{
				expectedSymbol{"basic.seru", ModuleSymbol},
				expectedSymbol{"basic.webidl", ModuleSymbol},
				expectedSymbol{"basictypes.seru", ModuleSymbol},
			}},
			searchSubTest{"Some", []expectedSymbol{
				expectedSymbol{"SomeClass", TypeSymbol},
				expectedSymbol{"DoSomething", MemberSymbol},
			}},
			searchSubTest{"some", []expectedSymbol{
				expectedSymbol{"SomeClass", TypeSymbol},
				expectedSymbol{"DoSomething", MemberSymbol},
			}},
			searchSubTest{"me", []expectedSymbol{
				expectedSymbol{"Message", MemberSymbol},
				expectedSymbol{"Message", MemberSymbol},
				expectedSymbol{"message", MemberSymbol},
				expectedSymbol{"DoSomething", MemberSymbol},
				expectedSymbol{"fileName", MemberSymbol},
				expectedSymbol{"name", MemberSymbol},
				expectedSymbol{"SomeClass", TypeSymbol},
			}},
		},
	},
}

func TestSymbolSearch(t *testing.T) {
	for _, test := range searchTests {
		testSourcePath := "tests/" + test.name + "/" + test.name + ".seru"
		groker := NewGroker(testSourcePath, []string{}, []packageloader.Library{packageloader.Library{TESTLIB_PATH, false, ""}})
		handle, err := groker.GetHandle()

		// Ensure we have a valid groker.
		if !assert.Nil(t, err, "Expected no error for test %s", test.name) {
			continue
		}

		// Conduct the searches.
		for _, searchTest := range test.subtests {
			symbols, err := handle.FindSymbols(searchTest.query)
			if !assert.Nil(t, err, "Expected no error for query %s under test %s", searchTest.query, test.name) {
				continue
			}

			if !assert.Equal(t, len(searchTest.expectedSymbols), len(symbols), "Expected symbol count mismatch for query %s under test %s: %v", searchTest.query, test.name, symbols) {
				continue
			}

			for index, expectedSymbol := range searchTest.expectedSymbols {
				if !assert.Equal(t, expectedSymbol.name, symbols[index].Name, "Name mismatch for query %s under test %s", searchTest.query, test.name) {
					continue
				}

				if !assert.Equal(t, expectedSymbol.kind, symbols[index].Kind, "Kind mismatch for query %s under test %s", searchTest.query, test.name) {
					continue
				}
			}
		}
	}
}
