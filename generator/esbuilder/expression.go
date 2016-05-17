// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esbuilder

import (
	"fmt"

	"github.com/serulian/compiler/sourcemap"
)

// ExpressionBuilder defines an interface for all expressions.
type ExpressionBuilder interface {
	// WithMapping adds a source mapping to the expression being built.
	WithMapping(mapping sourcemap.SourceMapping) ExpressionBuilder

	// Member returns an expression builder of an expression of a member under
	// this expression.
	Member(name interface{}) ExpressionBuilder

	// Call returns an expression builder of an expression calling this expression.
	Call(arguments ...interface{}) ExpressionBuilder

	// Mapping returns the source mapping for the expression being built.
	Mapping() (sourcemap.SourceMapping, bool)

	// Code returns the code of the expression being built.
	Code() string

	emitSource(sb *sourceBuilder)
}

// expressionNode is an interface for all nodes representing expressions.
type expressionNode interface {
	// emit emits the source for this node.
	emit(sb *sourceBuilder)
}

// expressionBuilder defines a wrapper for all expressions.
type expressionBuilder struct {
	// The actual expression.
	expression expressionNode

	// The source mapping for this expression, if any.
	mapping *sourcemap.SourceMapping
}

func (builder expressionBuilder) emitSource(sb *sourceBuilder) {
	builder.expression.emit(sb)
}

func (builder expressionBuilder) Mapping() (sourcemap.SourceMapping, bool) {
	if builder.mapping == nil {
		return sourcemap.SourceMapping{}, false
	}

	return *builder.mapping, true
}

// Code returns the generated code for the expression being built.
func (builder expressionBuilder) Code() string {
	return buildSource(builder).buf.String()
}

// WithMapping adds a source mapping to a builder.
func (builder expressionBuilder) WithMapping(mapping sourcemap.SourceMapping) ExpressionBuilder {
	builder.mapping = &mapping
	return builder
}

// Member returns the member with the given name under the expression.
func (builder expressionBuilder) Member(name interface{}) ExpressionBuilder {
	if strName, ok := name.(string); ok {
		return Member(builder, strName)
	}

	if exprName, ok := name.(ExpressionBuilder); ok {
		return ExprMember(builder, exprName)
	}

	panic(fmt.Sprintf("Unsupported name type: %T", name))
}

// Call returns a function call on the expression.
func (builder expressionBuilder) Call(arguments ...interface{}) ExpressionBuilder {
	args := make([]ExpressionBuilder, len(arguments))
	for index, arg := range arguments {
		if exprArg, ok := arg.(ExpressionBuilder); ok {
			args[index] = exprArg
		} else {
			args[index] = Value(arg)
		}
	}

	return Call(builder, args...)
}
