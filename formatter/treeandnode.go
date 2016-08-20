// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/parser"
)

// parseTree represents an in-memory parsing tree.
type parseTree struct {
	// The lines of the source.
	lines []string

	// The errors found during parsing.
	errors []formatterNode
}

// formatterNode defines a fully in-memory node for parser output. We don't use the SRG
// here as it will also bring in imports, which is not needed for formatting.
type formatterNode struct {
	nodeType   parser.NodeType
	properties map[string]string
	children   map[string]*list.List
}

func newParseTree(source []byte) *parseTree {
	return &parseTree{
		lines:  strings.Split(string(source), "\n"),
		errors: []formatterNode{},
	}
}

func (pt *parseTree) getLineNumber(characterIndex int) int {
	var current = 0
	for index, line := range pt.lines {
		if characterIndex <= current+len(line) {
			return index
		}

		current += len(line) + 1
	}

	return len(pt.lines)
}

func (pt *parseTree) createAstNode(source compilercommon.InputSource, kind parser.NodeType) parser.AstNode {
	node := formatterNode{
		nodeType:   kind,
		properties: make(map[string]string),
		children:   make(map[string]*list.List),
	}

	if kind == parser.NodeTypeError {
		pt.errors = append(pt.errors, node)
	}

	return node
}

func (fn formatterNode) hasType(types ...parser.NodeType) bool {
	for _, nodeType := range types {
		if nodeType == fn.GetType() {
			return true
		}
	}

	return false
}

func (fn formatterNode) getAllChildren() []formatterNode {
	var children = make([]formatterNode, 0)
	for _, childList := range fn.children {
		for e := childList.Front(); e != nil; e = e.Next() {
			children = append(children, e.Value.(formatterNode))
		}
	}
	return children
}

func (fn formatterNode) hasProperty(name string) bool {
	_, hasProp := fn.properties[name]
	return hasProp
}

func (fn formatterNode) hasChild(predicate string) bool {
	children, hasChildren := fn.children[predicate]
	if !hasChildren {
		return false
	}

	return children.Len() > 0
}

func (fn formatterNode) getProperty(name string) string {
	return fn.properties[name]
}

func (fn formatterNode) tryGetProperty(name string) (string, bool) {
	prop, ok := fn.properties[name]
	return prop, ok
}

func (fn formatterNode) getChildren(predicates ...string) []formatterNode {
	if len(predicates) == 1 {
		return fn.getSpecificChildren(predicates[0])
	}

	var children = make([]formatterNode, 0)
	for _, predicate := range predicates {
		children = append(children, fn.getSpecificChildren(predicate)...)
	}

	return children
}

func (fn formatterNode) getSpecificChildren(predicate string) []formatterNode {
	children, hasChildren := fn.children[predicate]
	if !hasChildren {
		return []formatterNode{}
	}

	var casted = make([]formatterNode, 0, children.Len())
	for e := children.Front(); e != nil; e = e.Next() {
		casted = append(casted, e.Value.(formatterNode))
	}

	return casted
}

func (fn formatterNode) getChild(predicate string) formatterNode {
	return fn.children[predicate].Front().Value.(formatterNode)
}

func (fn formatterNode) tryGetChild(predicate string) (formatterNode, bool) {
	children, hasChildren := fn.children[predicate]
	if !hasChildren {
		return formatterNode{}, false
	}

	return children.Front().Value.(formatterNode), true
}

func (fn formatterNode) getChildrenOfType(predicate string, childType parser.NodeType) []formatterNode {
	children, hasChildren := fn.children[predicate]
	if !hasChildren {
		return []formatterNode{}
	}

	var filtered = make([]formatterNode, 0, children.Len())
	for e := children.Front(); e != nil; e = e.Next() {
		if e.Value.(formatterNode).nodeType == childType {
			filtered = append(filtered, e.Value.(formatterNode))
		}
	}

	return filtered
}

func (fn formatterNode) GetType() parser.NodeType {
	return fn.nodeType
}

func (fn formatterNode) Connect(predicate string, other parser.AstNode) parser.AstNode {
	if fn.children[predicate] == nil {
		fn.children[predicate] = list.New()
	}

	fn.children[predicate].PushBack(other)
	return fn
}

func (fn formatterNode) Decorate(property string, value string) parser.AstNode {
	if _, ok := fn.properties[property]; ok {
		panic(fmt.Sprintf("Existing key for property %s\n\tNode: %v", property, fn.properties))
	}

	fn.properties[property] = value
	return fn
}

func (fn formatterNode) DecorateWithInt(property string, value int) parser.AstNode {
	return fn.Decorate(property, strconv.Itoa(value))
}
