// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/packageloader"
)

var _ = fmt.Printf

type grokSignaturesTest struct {
	name     string
	subtests []grokSignaturesSubTest
}

type grokSignaturesSubTest struct {
	rangeName         string
	activationString  string
	expectedSignature expectedSignature
}

type expectedSignature struct {
	name           string
	doc            string
	paramIndex     int
	expectedParams []expectedParam
}

type expectedParam struct {
	name      string
	paramType string
	doc       string
}

var grokSignaturesTests = []grokSignaturesTest{
	grokSignaturesTest{"signatures",
		[]grokSignaturesSubTest{
			grokSignaturesSubTest{
				"ids",
				"SomeClass.CreateMe(",
				expectedSignature{
					"CreateMe",
					"",
					0,
					[]expectedParam{},
				},
			},
			grokSignaturesSubTest{
				"ids",
				"DoSomething(",
				expectedSignature{
					"DoSomething",
					"DoSomething does something. When `someParam` is specified, good things happen. `anotherParam` controls other things.",
					0,
					[]expectedParam{
						expectedParam{"someParam", "Integer", "When ***someParam*** is specified, good things happen"},
						expectedParam{"anotherParam", "String", "***anotherParam*** controls other things"},
					},
				},
			},
			grokSignaturesSubTest{
				"ids",
				"DoSomething(42,",
				expectedSignature{
					"DoSomething",
					"DoSomething does something. When `someParam` is specified, good things happen. `anotherParam` controls other things.",
					1,
					[]expectedParam{
						expectedParam{"someParam", "Integer", "When ***someParam*** is specified, good things happen"},
						expectedParam{"anotherParam", "String", "***anotherParam*** controls other things"},
					},
				},
			},
			grokSignaturesSubTest{
				"ids",
				"sc[",
				expectedSignature{
					"index",
					"A cool indexer, that takes `index`.",
					0,
					[]expectedParam{
						expectedParam{"index", "Integer", "A cool indexer, that takes ***index***"},
					},
				},
			},
			grokSignaturesSubTest{
				"ids",
				"SomeClass[",
				expectedSignature{
					"",
					"",
					0,
					[]expectedParam{},
				},
			},
			grokSignaturesSubTest{
				"ids",
				"UnknownFunction(",
				expectedSignature{
					"",
					"",
					0,
					[]expectedParam{},
				},
			},
			grokSignaturesSubTest{
				"ids",
				"InvalidExpression",
				expectedSignature{
					"",
					"",
					0,
					[]expectedParam{},
				},
			},
		}},
}

func TestGrokSignatures(t *testing.T) {
	for _, grokSignaturesTest := range grokSignaturesTests {
		testSourcePath := "tests/" + grokSignaturesTest.name + "/" + grokSignaturesTest.name + ".seru"
		groker := NewGroker(testSourcePath, []string{}, []packageloader.Library{packageloader.Library{TESTLIB_PATH, false, "", "testcore"}})
		handle, err := groker.GetHandle()

		// Ensure we have a valid groker.
		if !assert.Nil(t, err, "Expected no error for test %s", grokSignaturesTest.name) {
			continue
		}

		// Find the ranges.
		ranges, err := getAllNamedRanges(handle)
		if !assert.Nil(t, err, "Error when looking up named ranges") {
			continue
		}

		for _, grokSignaturesSubTest := range grokSignaturesTest.subtests {
			// Find the named range.
			testRange, ok := ranges[grokSignaturesSubTest.rangeName]
			if !assert.True(t, ok, "Could not find range %s under signatures test %s", grokSignaturesSubTest.rangeName, grokSignaturesTest.name) {
				continue
			}

			// Retrieve the signatures for the range's location.
			pm := compilercommon.LocalFilePositionMapper{}
			sourcePosition := compilercommon.InputSource(testSourcePath).PositionForRunePosition(testRange.startIndex, pm)
			signatureInfo, err := handle.GetSignature(grokSignaturesSubTest.activationString, sourcePosition)
			if !assert.Nil(t, err, "Error when looking up signatures: %s => %v", testRange.name, sourcePosition) {
				continue
			}

			if !assert.Equal(t, grokSignaturesSubTest.expectedSignature.name, signatureInfo.Name, "Mismatch on signature name on test `%s` under range %s", grokSignaturesSubTest.activationString, grokSignaturesSubTest.rangeName) {
				continue
			}

			if !assert.Equal(t, grokSignaturesSubTest.expectedSignature.doc, signatureInfo.Documentation, "Mismatch on signature documentation on test `%s` under range %s", grokSignaturesSubTest.activationString, grokSignaturesSubTest.rangeName) {
				continue
			}

			if !assert.Equal(t, grokSignaturesSubTest.expectedSignature.paramIndex, signatureInfo.ActiveParameterIndex, "Mismatch on parameter index on test `%s` under range %s", grokSignaturesSubTest.activationString, grokSignaturesSubTest.rangeName) {
				continue
			}

			if !assert.Equal(t, len(grokSignaturesSubTest.expectedSignature.expectedParams), len(signatureInfo.Parameters), "Mismatch on number of parameters on test `%s` under range %s", grokSignaturesSubTest.activationString, grokSignaturesSubTest.rangeName) {
				continue
			}

			for index, expectedParameter := range grokSignaturesSubTest.expectedSignature.expectedParams {
				if !assert.Equal(t, expectedParameter.name, signatureInfo.Parameters[index].Name, "Mismatch on parameter name for param %v on test `%s` under range %s", index, grokSignaturesSubTest.activationString, grokSignaturesSubTest.rangeName) {
					continue
				}

				if !assert.Equal(t, expectedParameter.doc, signatureInfo.Parameters[index].Documentation, "Mismatch on parameter docs for param %v on test `%s` under range %s", index, grokSignaturesSubTest.activationString, grokSignaturesSubTest.rangeName) {
					continue
				}

				if !assert.Equal(t, expectedParameter.paramType, signatureInfo.Parameters[index].TypeReference.String(), "Mismatch on parameter type for param %v on test `%s` under range %s", index, grokSignaturesSubTest.activationString, grokSignaturesSubTest.rangeName) {
					continue
				}
			}
		}
	}
}
