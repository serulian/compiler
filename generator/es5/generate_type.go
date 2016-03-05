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

	case typegraph.ImplicitInterfaceType:
		return gen.templater.Execute("interface", interfaceTemplateStr, generating)

	case typegraph.NominalType:
		return gen.templater.Execute("nominal", nominalTemplateStr, generating)

	case typegraph.StructType:
		return gen.templater.Execute("struct", structTemplateStr, generating)

	case typegraph.ExternalInternalType:
		return ""

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

// HasGenerics returns whether the type has generics.
func (gt generatingType) HasGenerics() bool {
	return gt.Type.HasGenerics()
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

// RequiredFields returns the fields required to be initialized by this type.
func (gt generatingType) RequiredFields() []typegraph.TGMember {
	return gt.Type.RequiredFields()
}

// Fields returns the all the fields on this type.
func (gt generatingType) Fields() []typegraph.TGMember {
	return gt.Type.Fields()
}

// Alias returns the alias for this type, if any.
func (gt generatingType) Alias() string {
	alias, hasAlias := gt.Type.Alias()
	if !hasAlias {
		return ""
	}

	return alias
}

// TypeReferenceCall returns the expression for calling a type ref.
func (gt generatingType) TypeReferenceCall(typeref typegraph.TypeReference) string {
	return gt.Generator.pather.TypeReferenceCall(typeref)
}

// WrappedType returns the type wrapped by this nominal type.
func (gt generatingType) WrappedType() typegraph.TypeReference {
	if gt.Type.TypeKind() != typegraph.NominalType {
		panic("Cannot call WrappedType on non-nominal type")
	}

	return gt.Type.ParentTypes()[0]
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
	({{ .ComposedTypeLocation }}).new().then(function(value) {
	  instance.{{ .InnerInstanceName }} = value;
	})
`

// genericsTemplateStr defines a template for generating generics.
const genericsTemplateStr = `{{ range $index, $generic := .Generics }}{{ if $index }}, {{ end }}{{ $generic.Name }}{{ end }}`

// classTemplateStr defines the template for generating a class type.
const classTemplateStr = `
this.$class('{{ .Type.Name }}', {{ .HasGenerics }}, '{{ .Alias }}', function({{ .Generics }}) {
	var $static = this;
    var $instance = this.prototype;

    {{ $vars := .GenerateVariables }}
    {{ $composed := .GenerateComposition }}
	$static.new = function({{ range $ridx, $field := .RequiredFields }}{{ if $ridx }}, {{ end }}{{ $field.Name }}{{ end }}) {
		var instance = new $static();
		var init = [];
		{{ range $idx, $field := .RequiredFields }}
			instance.{{ $field.Name }} = {{ $field.Name }};
		{{ end }}
		{{ range $idx, $kv := $composed.Iter }}
			init.push({{ $kv.Value }});
  		{{ end }}
		{{ range $idx, $kv := $vars.Iter }}
			init.push({{ $kv.Value }});
		{{ end }}
		return $promise.all(init).then(function() {
			return instance;
		});
	};

	{{range $idx, $kv := .GenerateImplementedMembers.Iter }}
  	  {{ $kv.Value }}
  	{{end}}
});
`

// structTemplateStr defines the template for generating a struct type.
const structTemplateStr = `
this.$struct('{{ .Type.Name }}', {{ .HasGenerics }}, '{{ .Alias }}', function({{ .Generics }}) {
	var $static = this;
	var $instance = this.prototype;

	$static.new = function({{ range $ridx, $field := .RequiredFields }}{{ if $ridx }}, {{ end }}{{ $field.Name }}{{ end }}) {
		var instance = new $static();
		instance.$data = {};
		instance.$lazycheck = false;

		{{ range $idx, $field := .RequiredFields }}
			instance.{{ $field.Name }} = {{ $field.Name }};
		{{ end }}
		return $promise.resolve(instance);
	};

	{{ $parent := . }}

	{{ range $idx, $field := .Fields }}
	  {{ $boxed := $field.MemberType.IsNominalOrStruct }}
	  Object.defineProperty($instance, '{{ $field.Name }}', {
	    get: function() {
	    	if (this.$lazycheck) {
	    		$t.ensurevalue(this.$data.{{ $field.Name }}, {{ $parent.TypeReferenceCall $field.MemberType }});
	    	}

	    	{{ if $boxed }}
	    	return $t.box(this.$data.{{ $field.Name }}, {{ $parent.TypeReferenceCall $field.MemberType }});
	    	{{ else }}
	    	return this.$data.{{ $field.Name }};
	    	{{ end }}
	    },

	    set: function(val) {
	    	{{ if $boxed }}
	    	this.$data.{{ $field.Name }} = $t.unbox(val, {{ $parent.TypeReferenceCall $field.MemberType }});
	    	{{ else }}
	    	this.$data.{{ $field.Name }} = val;
	    	{{ end }}
	    }
	  });
	{{ end }}
});
`

// interfaceTemplateStr defines the template for generating an interface type.
const interfaceTemplateStr = `
this.$interface('{{ .Type.Name }}', {{ .HasGenerics }}, '{{ .Alias }}', function({{ .Generics }}) {
	var $static = this;
	{{range $idx, $kv := .GenerateImplementedMembers.Iter }}
  	  {{ $kv.Value }}
  	{{end}}
});
`

// nominalTemplateStr defines the template for generating a nominal type.
const nominalTemplateStr = `
this.$type('{{ .Type.Name }}', {{ .HasGenerics }}, '{{ .Alias }}', function({{ .Generics }}) {
	var $instance = this.prototype;
	var $static = this;

	this.new = function($wrapped) {
		var instance = new this();
		instance.$wrapped = $wrapped;
		return instance;
	};

	this.$apply = function(data) {
		var instance = new this();
		{{ if .WrappedType.IsNominalOrStruct }}
		instance.$wrapped = {{ .TypeReferenceCall .WrappedType }}.$apply(data.$wrapped);
		{{ else }}
		instance.$wrapped = data.$wrapped;
		{{ end }}
		return instance;
	};
	
	{{range $idx, $kv := .GenerateImplementedMembers.Iter }}
  	  {{ $kv.Value }}
  	{{end}}
});
`
