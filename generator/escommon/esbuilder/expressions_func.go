// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esbuilder

// functionNode defines a function literal.
type functionNode struct {
	// name is the optional name for the function.
	name string

	// parameters are the parameters to the function.
	parameters []string

	// body is the body of the function.
	body SourceBuilder
}

func (node functionNode) emit(sb *sourceBuilder) {
	if node.name == "" {
		sb.append("(")
	}

	sb.append("function ")

	if node.name != "" {
		sb.append(node.name)
	}

	sb.append("(")
	for index, name := range node.parameters {
		if index > 0 {
			sb.append(",")
		}

		sb.append(name)
	}

	sb.append(")")

	sb.append("{")
	sb.appendLine()
	sb.indent()
	sb.emit(node.body)
	sb.dedent()
	sb.appendLine()
	sb.append("}")

	if node.name == "" {
		sb.append(")")
	}
}

func (node functionNode) isStateless() bool {
	return node.name == ""
}

// Closure returns an anonymous closure.
func Closure(body SourceBuilder, parameters ...string) ExpressionBuilder {
	return expressionBuilder{functionNode{"", parameters, body}, nil}
}

// Function returns a named function.
func Function(name string, body SourceBuilder, parameters ...string) ExpressionBuilder {
	return expressionBuilder{functionNode{name, parameters, body}, nil}
}
