// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esbuilder

import "github.com/serulian/compiler/sourcemap"

// StatementBuilder defines an interface for all statements.
type StatementBuilder interface {
	// Mapping returns the source mapping for the statement being built.
	Mapping() (sourcemap.SourceMapping, bool)

	// Code returns the code of the statement being built.
	Code() string

	emitSource(sb *sourceBuilder)
}

// statementNode is an interface for all nodes representing statements.
type statementNode interface {
	// emit emits the source for this node.
	emit(sb *sourceBuilder)
}

// statementBuilder defines a wrapper for all statements.
type statementBuilder struct {
	// The actual statement.
	statement statementNode

	// The source mapping for this statement, if any.
	mapping *sourcemap.SourceMapping
}

func (builder statementBuilder) Mapping() (sourcemap.SourceMapping, bool) {
	if builder.mapping == nil {
		return sourcemap.SourceMapping{}, false
	}

	return *builder.mapping, true
}

func (builder statementBuilder) emitSource(sb *sourceBuilder) {
	builder.statement.emit(sb)
}

// Code returns the generated code for the expression being built.
func (builder statementBuilder) Code() string {
	return buildSource(builder).buf.String()
}

// WithMapping adds a source mapping to a builder.
func (builder statementBuilder) WithMapping(mapping sourcemap.SourceMapping) StatementBuilder {
	builder.mapping = &mapping
	return builder
}
