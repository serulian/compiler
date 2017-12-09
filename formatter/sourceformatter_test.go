// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
	"github.com/serulian/compiler/sourceshape"
	"github.com/stretchr/testify/assert"
)

// ignoredNodeTypes defines the list of node types that can be skipped in the formatting
// tests.
var ignoredNodeTypes = []sourceshape.NodeType{
	// SKIPPED: Parser currently cannot parse >>, as it conflicts with generic type refs.
	sourceshape.NodeBitwiseShiftRightExpression,
}

type goldenTest struct {
	name     string
	filename string
}

func (ft *goldenTest) input() []byte {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s.in.seru", ft.filename))
	if err != nil {
		panic(err)
	}

	return b
}

func (ft *goldenTest) output() []byte {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s.out.seru", ft.filename))
	if err != nil {
		panic(err)
	}

	return b
}

func (ft *goldenTest) writeFormatted(value []byte) {
	err := ioutil.WriteFile(fmt.Sprintf("tests/%s.out.seru", ft.filename), value, 0644)
	if err != nil {
		panic(err)
	}
}

var goldenTests = []goldenTest{
	{"basic test", "basic"},
	{"comments test", "comments"},
	{"unary precedence test", "unary"},
	{"binary precedence test", "binary"},
	{"imports test", "imports"},
	{"class test", "class"},
	{"agent test", "agent"},
	{"interface test", "interface"},
	{"struct test", "struct"},
	{"nominal test", "nominal"},
	{"relative imports test", "relative"},
	{"template strings test", "templatestrings"},
	{"match test", "match"},
	{"switch test", "switch"},
	{"typerefs test", "typerefs"},
	{"statements test", "statements"},
	{"expressions test", "expressions"},
	{"nullable precedence test", "nullable"},
	{"sml test", "sml"},
	{"nested sml test", "nestedsml"},
	{"sml text with tag test", "smltextwithtag"},
	{"sml long text test", "smllongtext"},
	{"sml long text 2 test", "smllongtext2"},
	{"simple return test", "simplereturn"},
	{"if-else return test", "ifelsereturn"},
	{"null assert test", "nullassert"},
	{"mappings test", "mappings"},
	{"nested sml loop test", "nestedsmlloop"},
	{"known issue test", "knownissue"},
	{"precedence test", "precedence"},
}

func conductParsing(t *testing.T, test goldenTest, source []byte) (*parseTree, formatterNode) {
	parseTree := newParseTree(source)
	inputSource := compilercommon.InputSource(test.filename)
	rootNode := parser.Parse(parseTree.createAstNode, nil, inputSource, string(source))
	if !assert.Empty(t, parseTree.errors, "Expected no parse errors for test %s", test.name) {
		return nil, formatterNode{}
	}

	return parseTree, rootNode.(formatterNode)
}

func addEncounteredNodeTypes(node formatterNode, encounteredNodeTypes map[sourceshape.NodeType]bool) {
	encounteredNodeTypes[node.GetType()] = true
	for _, child := range node.getAllChildren() {
		addEncounteredNodeTypes(child, encounteredNodeTypes)
	}
}

func TestGolden(t *testing.T) {
	encounteredNodeTypes := map[sourceshape.NodeType]bool{}

	for _, test := range goldenTests {
		if os.Getenv("FILTER") != "" {
			if !strings.Contains(test.name, os.Getenv("FILTER")) {
				continue
			} else {
				fmt.Printf("Matched Test: %v\n", test.name)
			}
		}

		parseTree, rootNode := conductParsing(t, test, test.input())
		if parseTree == nil {
			continue
		}

		if os.Getenv("FILTER") == "" {
			addEncounteredNodeTypes(rootNode, encounteredNodeTypes)
		}

		formatted := buildFormattedSource(parseTree, rootNode, importHandlingInfo{})
		if os.Getenv("REGEN") == "true" {
			test.writeFormatted(formatted)
		} else {
			expected := string(test.output())

			// Make sure the output is equal to that expected.
			if !assert.Equal(t, expected, string(formatted), test.name) {
				t.Log(string(formatted))
			}

			// Process the formatted source again and ensure it doesn't change.
			reparsedTree, reparsedRootNode := conductParsing(t, test, formatted)
			if reparsedTree == nil {
				continue
			}

			formattedAgain := buildFormattedSource(reparsedTree, reparsedRootNode, importHandlingInfo{})
			if !assert.Equal(t, string(formatted), string(formattedAgain), test.name) {
				t.Log(string(formattedAgain))
			}
		}
	}

	// Ensure that all parser node types were encountered. This makes sure that we have handled
	// all formatting cases in our tests. Note that we only run this check if we are not filtering
	// tests, as the filter will almost certainly skip node types we need.
	if os.Getenv("FILTER") == "" {
	outer:
		for i := int(sourceshape.NodeTypeError) + 1; i < int(sourceshape.NodeTypeTagged); i++ {
			current := sourceshape.NodeType(i)
			for _, skipped := range ignoredNodeTypes {
				if current == skipped {
					continue outer
				}
			}

			if _, ok := encounteredNodeTypes[current]; !ok {
				t.Errorf("Missing formatting of %v", current)
			}
		}
	}
}

func TestAll(t *testing.T) {
	// For every .seru file found in the entire compiler package, parse, and if parseable, format
	// and parse again. This ensures that formatting never breaks on any test files.
	compilerutil.WalkSourcePath("../...", func(filePath string, info os.FileInfo) (bool, error) {
		if !strings.HasSuffix(info.Name(), sourceshape.SerulianFileExtension) {
			return false, nil
		}

		// Load the source file.
		source, err := ioutil.ReadFile(filePath)
		if err != nil {
			return false, err
		}

		// Parse the source file.
		parseTree := newParseTree(source)
		inputSource := compilercommon.InputSource(filePath)
		rootNode := parser.Parse(parseTree.createAstNode, nil, inputSource, string(source))
		if len(parseTree.errors) > 0 {
			t.Logf("Skipping file %s\n", filePath)
			return true, nil
		}

		t.Logf("Checking file %s\n", filePath)

		// Format the source file.
		formatted := buildFormattedSource(parseTree, rootNode.(formatterNode), importHandlingInfo{})

		// Re-parse the file, and ensure there are no errors.
		formattedParseTree := newParseTree(source)
		parser.Parse(formattedParseTree.createAstNode, nil, inputSource, string(formatted))
		assert.Equal(t, 0, len(parseTree.errors), "Got parse error on formatted source file %s: %v", inputSource, parseTree.errors)
		return true, nil
	}, packageloader.SerulianPackageDirectory)
}
