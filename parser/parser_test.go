// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testNode struct {
	nodeType   NodeType
	properties map[string]string
	children   map[string]*list.List
}

type parserTest struct {
	name     string
	filename string
}

func (pt *parserTest) input() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s.seru", pt.filename))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (pt *parserTest) tree() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s.tree", pt.filename))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func createAstNode(source InputSource, kind NodeType) AstNode {
	return &testNode{
		nodeType:   kind,
		properties: make(map[string]string),
		children:   make(map[string]*list.List),
	}
}

func (tn *testNode) Connect(predicate string, other AstNode) AstNode {
	if tn.children[predicate] == nil {
		tn.children[predicate] = list.New()
	}

	tn.children[predicate].PushBack(other)
	return tn
}

func (tn *testNode) Decorate(property string, value string) AstNode {
	tn.properties[property] = value
	return tn
}

var parserTests = []parserTest{
	// Import success tests.
	{"basic import test", "import/basic"},
	{"from import test", "import/from"},
	{"string import test", "import/string_import"},
	{"multiple imports test", "import/multiple"},
	{"complex imports test", "import/complex"},

	// Import failure tests.
	{"missing source import test", "import/missing_source"},
	{"invalid source", "import/invalid_source"},
	{"invalid subsource", "import/invalid_subsource"},
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		rootNode := parse(createAstNode, InputSource(test.name), test.input())
		parseTree := getParseTree((rootNode).(*testNode), 0)
		assert := assert.New(t)
		assert.Equal(test.tree(), parseTree, test.name)
		if test.tree() != parseTree {
			t.Log(parseTree)
		}
	}
}

func getParseTree(currentNode *testNode, indentation int) string {
	parseTree := ""
	parseTree = parseTree + strings.Repeat(" ", indentation)
	parseTree = parseTree + fmt.Sprintf("%v", currentNode.nodeType)
	parseTree = parseTree + "\n"

	keys := make([]string, 0)

	for key, _ := range currentNode.properties {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		parseTree = parseTree + strings.Repeat(" ", indentation+2)
		parseTree = parseTree + fmt.Sprintf("%s = %s", key, currentNode.properties[key])
		parseTree = parseTree + "\n"
	}

	for key, value := range currentNode.children {
		parseTree = parseTree + fmt.Sprintf("%s%v =>", strings.Repeat(" ", indentation+2), key)
		parseTree = parseTree + "\n"

		for e := value.Front(); e != nil; e = e.Next() {
			parseTree = parseTree + getParseTree(e.Value.(*testNode), indentation+4)
		}
	}

	return parseTree
}
