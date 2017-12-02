// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"bytes"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/sourceshape"
)

// IdentifierPathString returns the string form of the identifier path referenced
// by the given node. Will return false if the node is not an identifier path.
func IdentifierPathString(node compilergraph.GraphNode) (string, bool) {
	switch node.Kind() {
	case sourceshape.NodeTypeIdentifierExpression:
		return node.Get(sourceshape.NodeIdentifierExpressionName), true

	case sourceshape.NodeThisLiteralExpression:
		return "this", true

	case sourceshape.NodePrincipalLiteralExpression:
		return "principal", true

	case sourceshape.NodeMemberAccessExpression:
		parentPath, ok := IdentifierPathString(node.GetNode(sourceshape.NodeMemberAccessChildExpr))
		if !ok {
			return "", false
		}

		return parentPath + "." + node.Get(sourceshape.NodeMemberAccessIdentifier), true

	default:
		return "", false
	}
}

type named interface {
	Name() (string, bool)
	Node() compilergraph.GraphNode
}

// getParameterDocumentation returns the documentation associated with the given parameterized entity
// (parameter or generic).
func getParameterDocumentation(srg *SRG, named named, parentPredicate compilergraph.Predicate) (SRGDocumentation, bool) {
	name, hasName := named.Name()
	if !hasName {
		return SRGDocumentation{}, false
	}

	parentNode, hasParentNode := named.Node().TryGetIncomingNode(parentPredicate)
	if !hasParentNode {
		return SRGDocumentation{}, false
	}

	documentation, hasDocumentation := srg.getDocumentationForNode(parentNode)
	if !hasDocumentation {
		return SRGDocumentation{}, false
	}

	return documentation.ForParameter(name)
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

		cs, hasSummary := parameter.Code()
		if hasSummary {
			buffer.WriteString(cs.Code)
		} else {
			buffer.WriteString("?")
		}
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

		cs, hasSummary := generic.Code()
		if hasSummary {
			buffer.WriteString(cs.Code)
		} else {
			buffer.WriteString("?")
		}
	}

	buffer.WriteString(">")
}
