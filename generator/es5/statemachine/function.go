// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/es5pather"
	"github.com/serulian/compiler/generator/es5/templater"
	"github.com/serulian/compiler/graphs/scopegraph"
)

// FunctionDef defines the interface for a function accepted by GenerateFunctionSource.
type FunctionDef interface {
	Generics() []string                // Returns the names of the generics on the function, if any.
	Parameters() []string              // Returns the names of the parameters on the function, if any.
	RequiresThis() bool                // Returns if this function is requires the "this" var to be added.
	IsExtension() bool                 // Returns true if this function is an extension function.
	BodyNode() compilergraph.GraphNode // The parser root node for the function body.
}

// GenerateFunctionSource generates the source code for a function, including its internal state machine.
func GenerateFunctionSource(functionDef FunctionDef, templater *templater.Templater, pather *es5pather.Pather, scopegraph *scopegraph.ScopeGraph) string {
	data := functionData{
		Generics:        functionDef.Generics(),
		Parameters:      functionDef.Parameters(),
		RequiresThis:    functionDef.RequiresThis() && !functionDef.IsExtension(),
		Body:            Build(functionDef.BodyNode(), templater, pather, scopegraph),
		ExtensionMethod: functionDef.IsExtension(),
	}

	return templater.Execute("functionSource", functionTemplateStr, data)
}

type functionData struct {
	Generics        []string
	Parameters      []string
	RequiresThis    bool
	Body            GeneratedMachine
	ExtensionMethod bool
}

// functionTemplateStr defines the template for generating inline functions.
const functionTemplateStr = `
{{ if .Generics }}
  function({{ range $index, $generic := .Generics }}{{ if $index }}, {{ end }}{{ $generic }}{{ end }}) {
	var $f =
{{ end }}
		function({{ if .ExtensionMethod }}$this, {{ end }}{{ range $index, $parameter := .Parameters }}{{ if $index }}, {{ end }}{{ $parameter }}{{ end }}) {
			{{ $body := .Body.Source }}
			{{ $hasBody := .Body.HasSource }}
			{{ if $hasBody }}
				{{ if .RequiresThis }}var $this = this;{{ end }}
				{{ $body }}
				return $promise.build($state);
			{{ else }}
				return $promise.empty();
			{{ end }}
		}
{{ if .Generics }}
	return $f;
  }
{{ end }}
`
