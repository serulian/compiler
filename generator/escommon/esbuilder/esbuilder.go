// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The esbuilder package implements an ECMAScript AST for easier generation of
// code with source mapping references.
package esbuilder

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"github.com/serulian/compiler/sourcemap"
)

// SourceBuilder defines an interface for all expressions, statements and templates.
type SourceBuilder interface {
	// WithMapping adds a source mapping to the item being built.
	WithMapping(mapping sourcemap.SourceMapping) SourceBuilder

	// mapping returns the source mapping for the item being built.
	mapping() (sourcemap.SourceMapping, bool)

	emitSource(sb *sourceBuilder)
}

// sourceBuilder defines a helper type for constructing the source and source map from ES Builder
// objects.
type sourceBuilder struct {
	buf              bytes.Buffer         // The buffer for the new source code.
	indentationLevel int                  // The current indentation level.
	charactersOnLine int                  // The number of characters on the current line.
	newlineCount     int                  // The number of newlines.
	hasNewline       bool                 // Whether we are on a newline.
	sourcemap        *sourcemap.SourceMap // The source map being constructed
}

// BuildSource builds the full source for the given builder.
func BuildSource(builder SourceBuilder) bytes.Buffer {
	return BuildSourceAndMap(builder, nil)
}

// BuildSourceAndMap builds the full source and sourcemap for the given builder.
func BuildSourceAndMap(builder SourceBuilder, sm *sourcemap.SourceMap) bytes.Buffer {
	sb := &sourceBuilder{
		indentationLevel: 0,
		hasNewline:       true,
		newlineCount:     0,
		charactersOnLine: 0,
		sourcemap:        sm,
	}

	sb.emit(builder)
	return sb.buf
}

// emitSeparated emits the source for each of the given builders.
func (sb *sourceBuilder) emitSeparated(builders []ExpressionBuilder, sep string) {
	for index, builder := range builders {
		if index > 0 {
			sb.append(",")
		}

		sb.emit(builder)
	}
}

// emitWrapped emits the given builder node's source at the current location, wrapped
// in parens.
func (sb *sourceBuilder) emitWrapped(builder ExpressionBuilder) {
	sb.append("(")
	sb.emit(builder)
	sb.append(")")
}

// emit emits the given builder node's source at the current location.
func (sb *sourceBuilder) emit(builder SourceBuilder) {
	if sb.sourcemap != nil {
		// Add the builder's mapping, if any.
		mapping, hasMapping := builder.mapping()
		if hasMapping {
			sb.sourcemap.AddMapping(sb.newlineCount, sb.charactersOnLine, mapping)
		}
	}

	// Generate the code for the builder.
	builder.emitSource(sb)
}

// indent increases the current indentation.
func (sb *sourceBuilder) indent() {
	sb.indentationLevel = sb.indentationLevel + 1
}

// dedent decreases the current indentation.
func (sb *sourceBuilder) dedent() {
	sb.indentationLevel = sb.indentationLevel - 1
}

// append adds the given value to the buffer, indenting as necessary.
func (sb *sourceBuilder) append(value string) {
	for _, currentRune := range value {
		if currentRune == '\n' {
			sb.buf.WriteRune('\n')
			sb.newlineCount++
			sb.charactersOnLine = 0
			sb.hasNewline = true
			continue
		}

		if sb.hasNewline {
			sb.buf.WriteString(strings.Repeat("  ", sb.indentationLevel))
			sb.charactersOnLine += sb.indentationLevel * 2
			sb.hasNewline = false
		}

		sb.buf.WriteRune(currentRune)
		sb.charactersOnLine += utf8.RuneLen(currentRune)
	}
}

// appendLine adds a newline.
func (sb *sourceBuilder) appendLine() {
	sb.append("\n")
}
