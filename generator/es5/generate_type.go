// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/cevaris/ordered_map"
)

// generateTypes generates all the types under the  given modules into ES5.
func (gen *es5generator) generateTypes(module typegraph.TGModule) *ordered_map.OrderedMap {
	typeMap := ordered_map.NewOrderedMap()
	types := module.Types()
	for _, typedecl := range types {
		typeMap.Set(typedecl, gen.generateType(typedecl))
	}

	return typeMap
}

// generateType generates the given type into ES5.
func (gen *es5generator) generateType(typedef typegraph.TGTypeDecl) string {
	generating := generatingType{typedef, gen}

	switch typedef.TypeKind() {
	case typegraph.ClassType:
		return gen.templater.Execute("class", classTemplateStr, generating)

	case typegraph.InterfaceType:
		return gen.templater.Execute("interface", interfaceTemplateStr, generating)

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
	return gt.Generator.templater.Execute("generics", genericsTemplateStr, gt.Type)
}

// GenerateImplementedMembers generates the source for all the members defined under the type that
// have implementations.
func (gt generatingType) GenerateImplementedMembers() *ordered_map.OrderedMap {
	return gt.Generator.generateImplementedMembers(gt.Type)
}

// GenerateVariables generates the source for all the variables defined under the type.
func (gt generatingType) GenerateVariables() *ordered_map.OrderedMap {
	return gt.Generator.generateVariables(gt.Type)
}

// GenerateComposition generates the source for all the composed types structurually inherited by the type.
func (gt generatingType) GenerateComposition() *ordered_map.OrderedMap {
	typeMap := ordered_map.NewOrderedMap()
	parentTypes := gt.Type.ParentTypes()
	for _, parentTypeRef := range parentTypes {
		data := struct {
			ComposedTypeLocation string
			InnerInstanceName    string
		}{
			gt.Generator.pather.TypeReferenceCall(parentTypeRef),
			gt.Generator.pather.InnerInstanceName(parentTypeRef),
		}

		source := gt.Generator.templater.Execute("composition", compositionTemplateStr, data)
		typeMap.Set(parentTypeRef, source)
	}

	return typeMap
}

// compositionTemplateStr defines a template for instantiating a composed type.
const compositionTemplateStr = `
	function() {
		var $this = this;
		return ({{ .ComposedTypeLocation }}).new().then(function(value) {
			$this.{{ .InnerInstanceName }} = value;
		});
	}
`

// genericsTemplateStr defines a template for generating generics.
const genericsTemplateStr = `{{ range $index, $generic := .Generics }}{{ if $index }}, {{ end }}{{ $generic.Name }}{{ end }}`

// classTemplateStr defines the template for generating a class type.
const classTemplateStr = `
this.$class('{{ .Type.Name }}', function({{ .Generics }}) {
	var $static = this;
    var $instance = this.prototype;

    {{ $vars := .GenerateVariables }}
    {{ $composed := .GenerateComposition }}
    {{ if or $vars.Iter $composed.Iter }}
	$static.new = function() {
		var instance = new $static();
		var init = [];
		{{range $idx, $kv := $composed.Iter }}
			init.push(({{ $kv.Value }})());
  		{{ end }}
		{{ range $idx, $kv := $vars.Iter }}
			init.push(({{ $kv.Value }})());
		{{ end }}
		return $promise.all(init).then(function() {
			return instance;
		});
	};
	{{ end }}

	{{range $idx, $kv := .GenerateImplementedMembers.Iter }}
  	  {{ $kv.Value }}
  	{{end}}
});
`

// interfaceTemplateStr defines the template for generating an interface type.
const interfaceTemplateStr = `
this.$interface('{{ .Type.Name }}', function({{ .Generics }}) {
	{{range $idx, $kv := .GenerateImplementedMembers.Iter }}
  	  {{ $kv.Value }}
  	{{end}}
});
`
