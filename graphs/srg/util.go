// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"bytes"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// salForIterator returns a SourceAndLocation for the given iterator. Note that
// the iterator *must* contain the NodePredicateSource and NodePredicateStartRune predicates.
func salForIterator(iterator compilergraph.NodeIterator) compilercommon.SourceAndLocation {
	return compilercommon.NewSourceAndLocation(
		compilercommon.InputSource(iterator.GetPredicate(parser.NodePredicateSource).String()),
		iterator.GetPredicate(parser.NodePredicateStartRune).Int())
}

// SourceLocationForNode returns a SourceAndLocation for the given graph node.
func SourceLocationForNode(node compilergraph.GraphNode) compilercommon.SourceAndLocation {
	return compilercommon.NewSourceAndLocation(
		compilercommon.InputSource(node.Get(parser.NodePredicateSource)),
		node.GetValue(parser.NodePredicateStartRune).Int())
}

// salForNode returns a SourceAndLocation for the given graph node.
func salForNode(node compilergraph.GraphNode) compilercommon.SourceAndLocation {
	return compilercommon.NewSourceAndLocation(
		compilercommon.InputSource(node.Get(parser.NodePredicateSource)),
		node.GetValue(parser.NodePredicateStartRune).Int())
}

// IdentifierPathString returns the string form of the identifier path referenced
// by the given node. Will return false if the node is not an identifier path.
func IdentifierPathString(node compilergraph.GraphNode) (string, bool) {
	switch node.Kind() {
	case parser.NodeTypeIdentifierExpression:
		return node.Get(parser.NodeIdentifierExpressionName), true

	case parser.NodeThisLiteralExpression:
		return "this", true

	case parser.NodePrincipalLiteralExpression:
		return "principal", true

	case parser.NodeMemberAccessExpression:
		parentPath, ok := IdentifierPathString(node.GetNode(parser.NodeMemberAccessChildExpr))
		if !ok {
			return "", false
		}

		return parentPath + "." + node.Get(parser.NodeMemberAccessIdentifier), true

	default:
		return "", false
	}
}

type documentable interface {
	Documentation() (SRGDocumentation, bool)
}

// getSummarizedDocumentation returns the summarized documentation for the given documentable
// instance, or empty string if none.
func getSummarizedDocumentation(documentable documentable) string {
	documentation, hasDocumentation := documentable.Documentation()
	if !hasDocumentation {
		return ""
	}

	value := documentation.String()
	if len(value) == 0 {
		return ""
	}

	sentences := strings.Split(value, ". ")
	if len(sentences) < 1 {
		return "// " + value
	}

	firstSentence := strings.TrimSpace(sentences[0])
	if documentation.IsDocComment() {
		if !strings.HasSuffix(firstSentence, ".") && !strings.HasSuffix(firstSentence, "!") {
			firstSentence = firstSentence + "."
		}
	}

	return "// " + firstSentence
}

type parameterable interface {
	Parameters() []SRGParameter
}

type genericable interface {
	Generics() []SRGGeneric
}

// writeCodeParameters writes the code representation of the parameters found to the
// given buffer.
func writeCodeParameters(m parameterable, buffer *bytes.Buffer) {
	parameters := m.Parameters()
	buffer.WriteString("(")
	for index, parameter := range parameters {
		if index > 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteString(parameter.Code())
	}

	buffer.WriteString(")")
}

// writeCodeGenerics writes the code representation of the generics found to the
// given buffer.
func writeCodeGenerics(m genericable, buffer *bytes.Buffer) {
	generics := m.Generics()
	if len(generics) == 0 {
		return
	}

	buffer.WriteString("<")
	for index, generic := range generics {
		if index > 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteString(generic.Name())

		constraint, hasConstraint := generic.GetConstraint()
		if hasConstraint {
			buffer.WriteString(" : ")
			buffer.WriteString(constraint.String())
		}
	}

	buffer.WriteString(">")
}
