// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esbuilder

import "github.com/serulian/compiler/sourcemap"

// StatementBuilder defines an interface for all statements.
type StatementBuilder interface {
	// WithMapping adds a source mapping to the statement being built.
	WithMapping(mapping sourcemap.SourceMapping) SourceBuilder

	// mapping returns the source mapping for the statement being built.
	mapping() (sourcemap.SourceMapping, bool)

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
	sourceMapping *sourcemap.SourceMapping
}

func (builder statementBuilder) mapping() (sourcemap.SourceMapping, bool) {
	if builder.sourceMapping == nil {
		return sourcemap.SourceMapping{}, false
	}

	return *builder.sourceMapping, true
}

func (builder statementBuilder) emitSource(sb *sourceBuilder) {
	builder.statement.emit(sb)
}

// WithMapping adds a source mapping to a builder.
func (builder statementBuilder) WithMapping(mapping sourcemap.SourceMapping) SourceBuilder {
	builder.sourceMapping = &mapping
	return builder
}
