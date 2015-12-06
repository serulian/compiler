// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/typegraph"
)

// generateTypes generates all the types under the  given modules into ES5.
func (gen *es5generator) generateTypes(module typegraph.TGModule) map[typegraph.TGTypeDecl]string {
	// Queue all the types to be generated.
	types := module.Types()
	generatedSource := make([]string, len(types))
	queue := compilerutil.Queue()
	for index, typedef := range types {
		fn := func(key interface{}, value interface{}) bool {
			generatedSource[key.(int)] = gen.generateType(value.(typegraph.TGTypeDecl))
			return true
		}

		queue.Enqueue(index, typedef, fn)
	}

	// Generate the full source tree for each type.
	queue.Run()

	// Build a map from type to source tree.
	typeMap := map[typegraph.TGTypeDecl]string{}
	for index, typedef := range types {
		typeMap[typedef] = generatedSource[index]
	}

	return typeMap
}

// generateType generates the given type into ES5.
func (gen *es5generator) generateType(typedef typegraph.TGTypeDecl) string {
	generating := generatingType{typedef, gen}

	switch typedef.TypeKind() {
	case typegraph.ClassType:
		return gen.runTemplate("class", classTemplateStr, generating)

	case typegraph.InterfaceType:
		return gen.runTemplate("interface", interfaceTemplateStr, generating)

	default:
		panic("Unknown typedef kind")
		return ""
	}
}

// generatingType represents a type being generated.
type generatingType struct {
	Type      typegraph.TGTypeDecl
	Generator *es5generator // The parent generator.
}

// Generics returns a string representing the named generics on this type.
func (gt generatingType) Generics() string {
	return gt.Generator.runTemplate("generics", genericsTemplateStr, gt.Type)
}

// GenerateImplementedMembers generates the source for all the members defined under the type that
// have implementations.
func (gt generatingType) GenerateImplementedMembers() map[typegraph.TGMember]string {
	return gt.Generator.generateImplementedMembers(gt.Type)
}

// GenerateVariables generates the source for all the variables defined under the type.
func (gt generatingType) GenerateVariables() map[typegraph.TGMember]string {
	return gt.Generator.generateVariables(gt.Type)
}

// genericsTemplateStr defines a template for generating generics.
const genericsTemplateStr = `{{ range $index, $generic := .Context.Generics }}{{ if $index }}, {{ end }}{{ $generic.Name }}{{ end }}`

// classTemplateStr defines the template for generating a class type.
const classTemplateStr = `
this.cls('{{ .Context.Type.Name }}', function({{ .Context.Generics }}) {
	var $static = this;
    var $instance = this.prototype;

	$static.$new = function() {
		var instance = new $static();
		(function() {
		{{range $member, $source := .Context.GenerateVariables }}
	  	  {{ $source }}
  		{{end}}
  		}).call(instance);
		return instance;
	};

	{{range $member, $source := .Context.GenerateImplementedMembers }}
  	  {{ $source }}
  	{{end}}
});
`

// interfaceTemplateStr defines the template for generating an interface type.
const interfaceTemplateStr = `
this.interface('{{ .Context.Type.Name }}', function({{ .Context.Generics }}) {
	{{range $member, $source := .Context.GenerateImplementedMembers }}
  	  {{ $source }}
  	{{end}}
});
`
