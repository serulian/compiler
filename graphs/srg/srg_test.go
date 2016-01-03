// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/serulian/compiler/parser"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestBasicLoading(t *testing.T) {
	testSRG := getSRG(t, "tests/basic/basic.seru")

	// Ensure that both classes were loaded.
	iterator := testSRG.findAllNodes(parser.NodeTypeClass).BuildNodeIterator(parser.NodeTypeDefinitionName)

	var classesFound []string = make([]string, 0, 2)
	for iterator.Next() {
		classesFound = append(classesFound, iterator.Values()[parser.NodeTypeDefinitionName])
	}

	if len(classesFound) != 2 {
		t.Errorf("Expected 2 classes found, found: %v", classesFound)
	}
}

func TestSyntaxError(t *testing.T) {
	_, result := loadSRG(t, "tests/syntaxerror/start.seru")
	if result.Status {
		t.Errorf("Expected failed parse")
	}

	// Ensure the syntax error was reported.
	if !strings.Contains(result.Errors[0].Error(), "Expected identifier") {
		t.Errorf("Expected parse error, found: %v", result.Errors)
	}
}

func TestImportError(t *testing.T) {
	_, result := loadSRG(t, "tests/syntaxerror/start.seru")
	if result.Status {
		t.Errorf("Expected failed parse")
	}

	// Ensure the parse error was reported.
	assert.Equal(t, "Expected identifier, found token tokenTypeLeftBrace", result.Errors[0].Error())
}
