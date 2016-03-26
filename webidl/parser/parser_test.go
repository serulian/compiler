// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
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
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s.webidl", pt.filename))
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

func (pt *parserTest) writeTree(value string) {
	err := ioutil.WriteFile(fmt.Sprintf("tests/%s.tree", pt.filename), []byte(value), 0644)
	if err != nil {
		panic(err)
	}
}

func createAstNode(source compilercommon.InputSource, kind NodeType) AstNode {
	return &testNode{
		nodeType:   kind,
		properties: make(map[string]string),
		children:   make(map[string]*list.List),
	}
}

func (tn *testNode) GetType() NodeType {
	return tn.nodeType
}

func (tn *testNode) Connect(predicate string, other AstNode) AstNode {
	if tn.children[predicate] == nil {
		tn.children[predicate] = list.New()
	}

	tn.children[predicate].PushBack(other)
	return tn
}

func (tn *testNode) Decorate(property string, value string) AstNode {
	if _, ok := tn.properties[property]; ok {
		panic(fmt.Sprintf("Existing key for property %s\n\tNode: %v", property, tn.properties))
	}

	tn.properties[property] = value
	return tn
}

var parserTests = []parserTest{
	parserTest{"empty interface test", "interface"},
	parserTest{"implements test", "implements"},
	parserTest{"inheritance test", "inheritance"},
	parserTest{"annotated interface test", "annotatedinterface"},
	parserTest{"members test", "members"},
	parserTest{"annotation value test", "annotationvalue"},
	parserTest{"indexer test", "indexer"},
	parserTest{"custom operation test", "customop"},

	parserTest{"known issue test", "knownissue"},
	parserTest{"full file test", "fullfile"},
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		if os.Getenv("FILTER") != "" {
			if !strings.Contains(test.name, os.Getenv("FILTER")) {
				continue
			} else {
				fmt.Printf("Matched Test: %v\n", test.name)
			}
		}

		moduleNode := createAstNode(compilercommon.InputSource(test.name), NodeTypeGlobalModule)

		Parse(moduleNode, createAstNode, compilercommon.InputSource(test.name), test.input())
		parseTree := getParseTree((moduleNode).(*testNode), 0)
		assert := assert.New(t)

		expected := strings.TrimSpace(test.tree())
		found := strings.TrimSpace(parseTree)

		if os.Getenv("REGEN") == "true" {
			test.writeTree(found)
		} else {
			if !assert.Equal(expected, found, test.name) {
				t.Log(parseTree)
			}
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

	keys = make([]string, 0)

	for key, _ := range currentNode.children {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		value := currentNode.children[key]
		parseTree = parseTree + fmt.Sprintf("%s%v =>", strings.Repeat(" ", indentation+2), key)
		parseTree = parseTree + "\n"

		for e := value.Front(); e != nil; e = e.Next() {
			parseTree = parseTree + getParseTree(e.Value.(*testNode), indentation+4)
		}
	}

	return parseTree
}
