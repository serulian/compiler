// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esbuilder

import (
	"bytes"
	"text/template"

	"github.com/serulian/compiler/sourcemap"
)

// TemplateBasedBuilder defines an interface for a templated code block.
type TemplateBasedBuilder interface {
	// Mapping returns the source mapping for the block being built.
	Mapping() (sourcemap.SourceMapping, bool)

	// Code returns the code of the block being built.
	Code() string

	emitSource(sb *sourceBuilder)
}

// templateBuilder defines a wrapper for templates.
type templateBuilder struct {
	name           string
	templateSource string
	data           interface{}

	// The source mapping for the template as a whole, if any.
	mapping *sourcemap.SourceMapping
}

// offsetedSourceMap represents a source map offseted by some index
// into the string generated by the template.
type offsetedSourceMap struct {
	// The 0-based offset into the templated string.
	offset int

	// The source map at the offset.
	sourcemap *sourcemap.SourceMap
}

func (builder templateBuilder) Mapping() (sourcemap.SourceMapping, bool) {
	if builder.mapping == nil {
		return sourcemap.SourceMapping{}, false
	}

	return *builder.mapping, true
}

func (builder templateBuilder) emitSource(sb *sourceBuilder) {
	var maps []offsetedSourceMap = make([]offsetedSourceMap, 0)
	var source bytes.Buffer

	// Register an `emit` function which does two things:
	//
	// 1) Builds and emits the source for the node at the place the function is called
	// 2) Saves the node's source map at the current template location for later appending
	emitNode := func(node ExpressionOrStatementBuilder) string {
		currentIndex := source.Len()
		nsb := buildSource(node)
		maps = append(maps, offsetedSourceMap{currentIndex, nsb.sourcemap})
		return nsb.buf.String()
	}

	funcMap := template.FuncMap{
		"emit": emitNode,
	}

	t := template.New(builder.name).Funcs(funcMap)
	parsed, err := t.Parse(builder.templateSource)
	if err != nil {
		panic(err)
	}

	// Execute the template.
	eerr := parsed.Execute(&source, builder.data)
	if eerr != nil {
		panic(eerr)
	}

	// Append the generated source to the builder.
	generatedSource := source.String()
	sb.append(generatedSource)

	// Append any offsetted source mappings.
	for _, osm := range maps {
		sb.sourcemap.AppendMap(osm.sourcemap.OffsetBy(generatedSource[0:osm.offset]))
	}
}

// Code returns the generated code for the expression being built.
func (builder templateBuilder) Code() string {
	return buildSource(builder).buf.String()
}

// WithMapping adds a source mapping to a builder.
func (builder templateBuilder) WithMapping(mapping sourcemap.SourceMapping) TemplateBasedBuilder {
	builder.mapping = &mapping
	return builder
}

// Template returns an inline templated code block.
func Template(name string, template string, data interface{}) TemplateBasedBuilder {
	return templateBuilder{name, template, data, nil}
}
