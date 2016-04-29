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
	"github.com/serulian/compiler/parser"
	"github.com/stretchr/testify/assert"
)

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
	{"template strings test", "templatestrings"},
}

func TestGolden(t *testing.T) {
	for _, test := range goldenTests {
		if os.Getenv("FILTER") != "" {
			if !strings.Contains(test.name, os.Getenv("FILTER")) {
				continue
			} else {
				fmt.Printf("Matched Test: %v\n", test.name)
			}
		}

		parseTree := newParseTree(test.input())
		inputSource := compilercommon.InputSource(test.filename)
		rootNode := parser.Parse(parseTree.createAstNode, nil, inputSource, string(test.input()))

		if !assert.Empty(t, parseTree.errors, "Expected no parse errors for test %s", test.name) {
			continue
		}

		formatted := buildFormattedSource(parseTree, rootNode.(formatterNode), importHandlingInfo{})

		if os.Getenv("REGEN") == "true" {
			test.writeFormatted(formatted)
		} else {
			expected := string(test.output())

			if !assert.Equal(t, expected, string(formatted), test.name) {
				t.Log(string(formatted))
			}
		}
	}
}
