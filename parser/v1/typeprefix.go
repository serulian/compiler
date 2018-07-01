// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"container/list"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser/shared"
	"github.com/serulian/compiler/sourceshape"
)

// IsTypePrefix returns whether the given input string is a prefix that supports a type reference declared right
// after it. For example, the string `function DoSomething() ` will return `true`, as a type can be specified right
// after that code snippet.
func IsTypePrefix(input string) bool {
	return checkForTypePrefix(input,
		// Check for type members.
		func(p *sourceParser) *typePrefixAstNode {
			rootNode := typePrefixAstNodeBuilder(compilercommon.InputSource(""), sourceshape.NodeTypeFile)
			p.consumeToken()
			p.consumeImplementedTypeMembers(rootNode)
			return rootNode.(*typePrefixAstNode)
		},

		// Check for types.
		func(p *sourceParser) *typePrefixAstNode {
			p.consumeToken()
			return p.consumeTypeDefinition().(*typePrefixAstNode)
		},

		// Check within expressions.
		func(p *sourceParser) *typePrefixAstNode {
			p.consumeToken()
			return p.consumeExpression(consumeExpressionAllowBraces).(*typePrefixAstNode)
		})
}

type consumeCaller func(p *sourceParser) *typePrefixAstNode

// checkForTypePrefix invokes each of the given caller functions to construct a small AST with the sentinal. The callers each return the
// root node from the parse tree, which is then checked for the presence of the `sentinal` under a type reference. If the `sentinal` is
// found under a type reference, then that token (which is placed right after the input snippet) is supported as a type reference.
func checkForTypePrefix(input string, callers ...consumeCaller) bool {
	for _, caller := range callers {
		p := buildParser(typePrefixAstNodeBuilder, noopReportImport, compilercommon.InputSource(""), bytePosition(0), input+"sentinal")
		rootNode := caller(p)
		p.lex.lex.consumeRemainder()
		if rootNode.hasSentinaledTypeReference() {
			return true
		}
	}

	return false
}

// typePrefixAstNode defines an AST node constructed by the parser for finding the sentinal under a type reference.
type typePrefixAstNode struct {
	nodeType sourceshape.NodeType
	children *list.List

	hasSentinal bool
}

func typePrefixAstNodeBuilder(source compilercommon.InputSource, kind sourceshape.NodeType) shared.AstNode {
	return &typePrefixAstNode{
		nodeType: kind,
		children: list.New(),
	}
}

func (tn *typePrefixAstNode) GetType() sourceshape.NodeType {
	return tn.nodeType
}

func (tn *typePrefixAstNode) Connect(predicate string, other shared.AstNode) shared.AstNode {
	tn.children.PushBack(other)
	return tn
}

func (tn *typePrefixAstNode) Decorate(property string, value string) shared.AstNode {
	// We only care if we find the `sentinal` as part of an identifier access.
	if property == sourceshape.NodeIdentifierAccessName && value == "sentinal" {
		tn.hasSentinal = true
	}

	return tn
}

func (tn *typePrefixAstNode) DecorateWithInt(property string, value int) shared.AstNode {
	return tn
}

// hasSentinaledTypeReference returns true if the parse tree contains a sentinal-ed type reference
// anywhere in the tree.
func (tn *typePrefixAstNode) hasSentinaledTypeReference() bool {
	if tn.nodeType == sourceshape.NodeTypeTypeReference && tn.hasSentinalIdentifierOrChild() {
		return true
	}

	for e := tn.children.Front(); e != nil; e = e.Next() {
		if e.Value.(*typePrefixAstNode).hasSentinaledTypeReference() {
			return true
		}
	}

	return false
}

// hasSentinalOrChild returns true if the parse tree contains an identifier with the `sentinal` value.
func (tn *typePrefixAstNode) hasSentinalIdentifierOrChild() bool {
	if tn.hasSentinal {
		return true
	}

	for e := tn.children.Front(); e != nil; e = e.Next() {
		if e.Value.(*typePrefixAstNode).hasSentinalIdentifierOrChild() {
			return true
		}
	}

	return false
}

func noopReportImport(sourceKind string, importPath string, importType packageloader.PackageImportType, importSource compilercommon.InputSource, runePosition int) string {
	return ""
}
