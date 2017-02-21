// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"fmt"

	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/cevaris/ordered_map"
)

var _ = fmt.Printf

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
func (gen *es5generator) generateType(typedef typegraph.TGTypeDecl) esbuilder.SourceBuilder {
	generating := generatingType{typedef, gen}

	switch typedef.TypeKind() {
	case typegraph.ClassType:
		return esbuilder.Template("class", classTemplateStr, generating)

	case typegraph.ImplicitInterfaceType:
		return esbuilder.Template("interface", interfaceTemplateStr, generating)

	case typegraph.NominalType:
		return esbuilder.Template("nominal", nominalTemplateStr, generating)

	case typegraph.StructType:
		return esbuilder.Template("struct", structTemplateStr, generating)

	case typegraph.ExternalInternalType:
		return esbuilder.Snippet("")

	default:
		panic("Unknown typedef kind")
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
func (gt generatingType) GenerateVariables() *generatedInitMap {
	varMap := newGeneratedInitMap()
	gt.Generator.generateVariables(gt.Type, varMap)
	return varMap
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

// NominalDataType returns the root data type behind this nominal type.
func (gt generatingType) NominalDataType() typegraph.TypeReference {
	if gt.Type.TypeKind() != typegraph.NominalType {
		panic("Cannot call WrappedType on non-nominal type")
	}

	return gt.Type.GetTypeReference().NominalDataType()
}

func (gt generatingType) BoolType() typegraph.TypeReference {
	return gt.Generator.scopegraph.TypeGraph().BoolTypeReference()
}

func (gt generatingType) MappingAnyType() typegraph.TypeReference {
	return gt.Generator.scopegraph.TypeGraph().MappingTypeReference(gt.Generator.scopegraph.TypeGraph().AnyTypeReference())
}

// GenerateComposition generates the source for all the composed types structurually inherited by the type.
func (gt generatingType) GenerateComposition() *generatedInitMap {
	initMap := newGeneratedInitMap()
	parentTypes := gt.Type.ParentTypes()
	constructor, _ := gt.Type.GetMember("new")

	for _, parentTypeRef := range parentTypes {
		data := struct {
			ComposedTypeLocation             string
			InnerInstanceName                string
			RequiredFields                   []typegraph.TGMember
			ComposedTypeConstructorPromising bool
		}{
			gt.Generator.pather.TypeReferenceCall(parentTypeRef),
			gt.Generator.pather.InnerInstanceName(parentTypeRef),
			parentTypeRef.ReferredType().RequiredFields(),
			gt.Generator.scopegraph.IsPromisingMember(constructor, scopegraph.PromisingAccessFunctionCall),
		}

		source := esbuilder.Template("composition", compositionTemplateStr, data)
		initMap.Set(parentTypeRef, generatedSourceResult{source, data.ComposedTypeConstructorPromising})
	}

	return initMap
}

// TypeSignatureMethod generates the $typesig method on a type definition.
func (gt generatingType) TypeSignatureMethod() string {
	return gt.Generator.templater.Execute("typesig", typeSignatureTemplateStr, gt)
}

// TypeSignature generates and returns the type signature for the given type.
func (gt generatingType) TypeSignature() typeSignature {
	return computeTypeSignature(gt.Type)
}

// typeSignatureTemplateStr defines a template for generating the signature for a type.
const typeSignatureTemplateStr = `
	this.$typesig = function() {
		{{ $sig := .TypeSignature }}
		{{ if $sig.IsEmpty }}
			return {};
		{{ else }}
			if (this.$cachedtypesig) { return this.$cachedtypesig; }

			var computed = {
				{{ range $sidx, $static := $sig.StaticSignatures }}
				{{ if $sidx }},{{ end }}
				{{ $static.ESCode }}: true
				{{ end }}
			};

			{{ range $sidx, $dynamic := $sig.DynamicSignatures }}
				computed[{{ $dynamic.ESCode }}] = true;
			{{ end }}

			return this.$cachedtypesig = computed;
		{{ end }}
	};
`

// compositionTemplateStr defines a template for instantiating a composed type.
const compositionTemplateStr = `
	{{ if .ComposedTypeConstructorPromising }}
	$promise.maybe({{ .ComposedTypeLocation }}.new({{ range $ridx, $field := .RequiredFields }}{{ if $ridx }}, {{ end }}{{ $field.Name }}{{ end }})).then(function(value) {
	  instance.{{ .InnerInstanceName }} = value;
	})
	{{ else }}
	instance.{{ .InnerInstanceName }} = {{ .ComposedTypeLocation }}.new({{ range $ridx, $field := .RequiredFields }}{{ if $ridx }}, {{ end }}{{ $field.Name }}{{ end }})
	{{ end }}
`

// genericsTemplateStr defines a template for generating generics.
const genericsTemplateStr = `{{ range $index, $generic := .Generics }}{{ if $index }}, {{ end }}{{ $generic.Name }}{{ end }}`

// interfaceTemplateStr defines the template for generating an interface type.
const interfaceTemplateStr = `
this.$interface('{{ .Type.GlobalUniqueId }}', '{{ .Type.Name }}', {{ .HasGenerics }}, '{{ .Alias }}', function({{ .Generics }}) {
	var $static = this;
	{{range $idx, $kv := .GenerateImplementedMembers.UnsafeIter }}
  	  {{ emit $kv.Value }}
  	{{end}}

  	{{ .TypeSignatureMethod }}
});
`

// classTemplateStr defines the template for generating a class type.
const classTemplateStr = `
this.$class('{{ .Type.GlobalUniqueId }}', '{{ .Type.Name }}', {{ .HasGenerics }}, '{{ .Alias }}', function({{ .Generics }}) {
	var $static = this;
    var $instance = this.prototype;

    {{ $vars := .GenerateVariables }}
    {{ $composed := .GenerateComposition }}
    {{ $combined := $vars.CombineWith $composed }}

	$static.new = function({{ range $ridx, $field := .RequiredFields }}{{ if $ridx }}, {{ end }}{{ $field.Name }}{{ end }}) {
		var instance = new $static();
		{{ range $idx, $field := .RequiredFields }}
		{{ if not $field.HasBaseMember }}
			instance.{{ $field.Name }} = {{ $field.Name }};
		{{ end }}
		{{ end }}

		{{ if $combined.Promising }}
		var init = [];
		{{ end }}

		{{ range $idx, $kv := $combined.Iter }}
			{{ emit $kv.Value }};
  		{{ end }}

		{{ if $combined.Promising }}
		return $promise.all(init).then(function() {
			return instance;
		});
		{{ else }}
		return instance;
		{{ end }}
	};

	{{range $idx, $kv := .GenerateImplementedMembers.UnsafeIter }}
  	  {{ emit $kv.Value }}
  	{{end}}

  	{{ .TypeSignatureMethod }}
});
`

// structTemplateStr defines the template for generating a struct type.
const structTemplateStr = `
this.$struct('{{ .Type.GlobalUniqueId }}', '{{ .Type.Name }}', {{ .HasGenerics }}, '{{ .Alias }}', function({{ .Generics }}) {
	var $static = this;
	var $instance = this.prototype;

	// new is the constructor called from Serulian code to construct the struct instance.
    {{ $vars := .GenerateVariables }}
	$static.new = function({{ range $ridx, $field := .RequiredFields }}{{ if $ridx }}, {{ end }}{{ $field.Name }}{{ end }}) {
		var instance = new $static();		
		instance[BOXED_DATA_PROPERTY] = {
			{{ range $idx, $field := .RequiredFields }}
			'{{ $field.SerializableName }}': {{ $field.Name }},
			{{ end }}
		};
		instance.$markruntimecreated();

		{{ if $vars.HasEntries }}
		return $static.$initDefaults(instance, true);
		{{ else }}
		return instance;
		{{ end }}
	};

	{{ if $vars.HasEntries }}
	$static.$initDefaults = function(instance, isRuntimeCreated) {
		var boxed = instance[BOXED_DATA_PROPERTY];

		{{ if $vars.Promising }}
		var init = [];
		{{ end }}
		{{ range $idx, $kv := $vars.Iter }}
			if (isRuntimeCreated || boxed['{{ $kv.Key.Name }}'] === undefined) {
				{{ emit $kv.Value }};
			}
		{{ end }}

		{{ if $vars.Promising }}
		return $promise.all(init).then(function() {
			return instance;
		});
		{{ else }}
		return instance;
		{{ end }}
	};
	{{ end }}

	$static.$fields = [];

	{{ $parent := . }}

	{{ range $idx, $field := .Fields }}
		$t.defineStructField($static,
							 '{{ $field.Name }}',
							 '{{ $field.SerializableName }}',
							 function() { return {{ $parent.TypeReferenceCall $field.MemberType }} },
							 function() { return {{ $parent.TypeReferenceCall $field.MemberType.NominalRootType }} },
							 {{ $field.MemberType.NullValueAllowed }});
	{{ end }}

  	{{ .TypeSignatureMethod }}
});
`

// nominalTemplateStr defines the template for generating a nominal type.
const nominalTemplateStr = `
this.$type('{{ .Type.GlobalUniqueId }}', '{{ .Type.Name }}', {{ .HasGenerics }}, '{{ .Alias }}', function({{ .Generics }}) {
	var $instance = this.prototype;
	var $static = this;

	this.$box = function($wrapped) {
		var instance = new this();
		instance[BOXED_DATA_PROPERTY] = $wrapped;
		return instance;
	};

	this.$roottype = function() {
		return {{ .TypeReferenceCall .NominalDataType }};
	};

	{{ range $idx, $kv := .GenerateImplementedMembers.UnsafeIter }}
  	  {{ emit $kv.Value }}
  	{{ end }}

  	{{ .TypeSignatureMethod }}
});
`
