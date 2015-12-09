// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The es5 package implements a generator for compiling Serulian into ECMAScript 5.
package es5

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/typegraph"
)

// es5generator defines a generator for producing ECMAScript 5 code.
type es5generator struct {
	graph      *compilergraph.SerulianGraph // The root graph.
	scopegraph *scopegraph.ScopeGraph       // The scope graph.
}

// generateModules generates all the modules found in the given scope graph into source.
func generateModules(sg *scopegraph.ScopeGraph, format bool) map[typegraph.TGModule]string {
	generator := es5generator{
		graph:      sg.SourceGraph().Graph,
		scopegraph: sg,
	}

	// Generate the code for each of the modules.
	generated := generator.generateModules(sg.TypeGraph().Modules())
	if !format {
		return generated
	}

	formatted := map[typegraph.TGModule]string{}
	for module, source := range generated {
		formatted[module] = formatSource(source)
	}

	return formatted
}

// GenerateES5 produces ES5 code from the given scope graph.
func GenerateES5(sg *scopegraph.ScopeGraph) (string, error) {
	generated := generateModules(sg, false)

	// Collect the generated modules into their final source.
	for _, source := range generated {
		fmt.Printf("%v", source)
	}

	return "", nil
}

// templateContext defines context given when invoking a template.
type templateContext struct {
	Context   interface{}   // The inner context.
	Generator *es5generator // The code generator.
}

// runTemplate runs the given go-template over the given context.
func (gen *es5generator) runTemplate(name string, templateStr string, context interface{}) string {
	// TODO: Cache this?
	t := template.New(name)
	parsed, err := t.Parse(templateStr)
	if err != nil {
		panic(err)
	}

	var source bytes.Buffer
	eerr := parsed.Execute(&source, templateContext{context, gen})
	if eerr != nil {
		panic(eerr)
	}

	return source.String()
}
