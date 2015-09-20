// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

func TestBasicLoading(t *testing.T) {
	graph, err := compilergraph.NewGraph("tests/basic/basic.seru")
	if err != nil {
		t.Errorf("%v", err)
	}

	testSRG := NewSRG(graph)
	result := testSRG.LoadAndParse()
	if !result.Status {
		t.Errorf("Expected successful parse")
	}

	// Ensure that both classes were loaded.
	iterator := testSRG.FindAllNodes(parser.NodeTypeClass).BuildNodeIterator(parser.NodeClassPredicateName)

	var classesFound []string = make([]string, 0, 2)
	for iterator.Next() {
		classesFound = append(classesFound, iterator.Values[parser.NodeClassPredicateName])
	}

	if len(classesFound) != 2 {
		t.Errorf("Expected 2 classes found, found: %v", classesFound)
	}
}

func TestSyntaxError(t *testing.T) {
	graph, err := compilergraph.NewGraph("tests/syntaxerror/start.seru")
	if err != nil {
		t.Errorf("%v", err)
	}

	testSRG := NewSRG(graph)
	result := testSRG.LoadAndParse()
	if result.Status {
		t.Errorf("Expected failed parse")
	}

	// Ensure the syntax error was reported.
	if !strings.Contains(result.Errors[0].Error(), "Expected identifier") {
		t.Errorf("Expected parse error, found: %v", result.Errors)
	}
}
