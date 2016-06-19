// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/generator/escommon/esbuilder"
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
func (gt generatingType) GenerateComposition() *ordered_map.OrderedMap {
	typeMap := ordered_map.NewOrderedMap()
	parentTypes := gt.Type.ParentTypes()
	for _, parentTypeRef := range parentTypes {
		data := struct {
			ComposedTypeLocation string
			InnerInstanceName    string
			RequiredFields       []typegraph.TGMember
		}{
			gt.Generator.pather.TypeReferenceCall(parentTypeRef),
			gt.Generator.pather.InnerInstanceName(parentTypeRef),
			parentTypeRef.ReferredType().RequiredFields(),
		}

		source := esbuilder.Template("composition", compositionTemplateStr, data)
		typeMap.Set(parentTypeRef, source)
	}

	return typeMap
}

// TypeSignatureMethod generates the $typesig method on a type definition.
func (gt generatingType) TypeSignatureMethod() string {
	return gt.Generator.templater.Execute("typesig", typeSignatureTemplateStr, gt)
}

// typeSignatureTemplateStr defines a template for generating the signature for a type.
const typeSignatureTemplateStr = `
	this.$typesig = function() {
		return $t.createtypesig(
		{{ $parent := . }}
		{{ range $midx, $member := .Type.Members }}
			{{ if $midx }},{{ end }}
			['{{ $member.Name }}', {{ $member.Signature.MemberKind }}, ({{ $parent.TypeReferenceCall $member.MemberType }}).$typeref()]
		{{ end }}
		);
	};
`

// compositionTemplateStr defines a template for instantiating a composed type.
const compositionTemplateStr = `
	({{ .ComposedTypeLocation }}).new({{ range $ridx, $field := .RequiredFields }}{{ if $ridx }}, {{ end }}{{ $field.Name }}{{ end }}).then(function(value) {
	  instance.{{ .InnerInstanceName }} = value;
	})
`

// genericsTemplateStr defines a template for generating generics.
const genericsTemplateStr = `{{ range $index, $generic := .Generics }}{{ if $index }}, {{ end }}{{ $generic.Name }}{{ end }}`

// interfaceTemplateStr defines the template for generating an interface type.
const interfaceTemplateStr = `
this.$interface('{{ .Type.Name }}', {{ .HasGenerics }}, '{{ .Alias }}', function({{ .Generics }}) {
	var $static = this;
	{{range $idx, $kv := .GenerateImplementedMembers.UnsafeIter }}
  	  {{ emit $kv.Value }}
  	{{end}}

  	{{ .TypeSignatureMethod }}
});
`

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
		{{ range $idx, $kv := $composed.UnsafeIter }}
			init.push({{ emit $kv.Value }});
  		{{ end }}
		{{ range $idx, $field := .RequiredFields }}
		{{ if not $field.HasBaseMember }}
			init.push($promise.new(function(resolve) {
				instance.{{ $field.Name }} = {{ $field.Name }};
				resolve();
			}));
		{{ end }}
		{{ end }}
		{{ range $idx, $kv := $vars.UnsafeIter }}
			init.push({{ emit $kv.Value }});
		{{ end }}
		return $promise.all(init).then(function() {
			return instance;
		});
	};

	{{range $idx, $kv := .GenerateImplementedMembers.UnsafeIter }}
  	  {{ emit $kv.Value }}
  	{{end}}

  	{{ .TypeSignatureMethod }}
});
`

// structTemplateStr defines the template for generating a struct type.
const structTemplateStr = `
this.$struct('{{ .Type.Name }}', {{ .HasGenerics }}, '{{ .Alias }}', function({{ .Generics }}) {
	var $static = this;
	var $instance = this.prototype;

	// new is the constructor called from Serulian code to construct the struct instance.
	$static.new = function({{ range $ridx, $field := .RequiredFields }}{{ if $ridx }}, {{ end }}{{ $field.Name }}{{ end }}) {
		var instance = new $static();
		instance[BOXED_DATA_PROPERTY] = {
			{{ range $idx, $field := .RequiredFields }}
			'{{ $field.SerializableName }}': {{ $field.Name }},
			{{ end }}
		};

		return $promise.resolve(instance);
	};

	{{ $parent := . }}

	// $box is the constructor called when deserializing a struct from JSON or another form of Mapping.
	$static.$box = function(data) {		
		var instance = new $static();
		instance[BOXED_DATA_PROPERTY] = data;
		instance.$lazychecked = {};

		// Override the properties to ensure we box as necessary.
		{{ range $idx, $field := .Fields }}
	 	{{ $boxed := $field.MemberType.IsNominalOrStruct }}
		  Object.defineProperty(instance, '{{ $field.Name }}', {
		    get: function() {
		    	if (this.$lazychecked['{{ $field.SerializableName }}']) {
		    		$t.ensurevalue(this[BOXED_DATA_PROPERTY]['{{ $field.SerializableName }}'], {{ $parent.TypeReferenceCall $field.MemberType.NominalRootType }}, {{ $field.MemberType.NullValueAllowed }}, '{{ $field.Name }}');
		    		this.$lazychecked['{{ $field.SerializableName }}'] = true;
		    	}

		    	{{ if $boxed }}
		    	return $t.box(this[BOXED_DATA_PROPERTY]['{{ $field.SerializableName }}'], {{ $parent.TypeReferenceCall $field.MemberType }});
	    		{{ else }}
		    	return this[BOXED_DATA_PROPERTY]['{{ $field.SerializableName }}'];
		    	{{ end }}
		    }
		  });
		{{ end }}

		// Override Mapping to make sure it returns boxed values.
		instance.Mapping = function() {
			var mapped = {};

			{{ range $idx, $field := .Fields }}
				mapped['{{ $field.SerializableName }}'] = this.{{ $field.Name }};
			{{ end }}

			return $promise.resolve($t.box(mapped, {{ .TypeReferenceCall .MappingAnyType }}));
		};

		return instance;
	};

	$instance.Mapping = function() {
		return $promise.resolve($t.box(this[BOXED_DATA_PROPERTY], {{ .TypeReferenceCall .MappingAnyType }}));
	};

	$static.$equals = function(left, right) {
		if (left === right) {
			return $promise.resolve($t.box(true, {{ .TypeReferenceCall .BoolType }}));
		}

		// TODO: find a way to do this without checking *all* fields.
		var promises = [];
		{{ range $idx, $field := .Fields }}
		promises.push($t.equals(left[BOXED_DATA_PROPERTY]['{{ $field.SerializableName }}'], 
		 					    right[BOXED_DATA_PROPERTY]['{{ $field.SerializableName }}'],
		 					    {{ $parent.TypeReferenceCall $field.MemberType }}));
		{{ end }}

		return Promise.all(promises).then(function(values) {
		  for (var i = 0; i < values.length; i++) {
		  	if (!$t.unbox(values[i])) {
	   		  return values[i];
		  	}
		  }

   		  return $t.box(true, {{ .TypeReferenceCall .BoolType }});
		});
	};

	{{ range $idx, $field := .Fields }}
	  Object.defineProperty($instance, '{{ $field.Name }}', {
	    get: function() {
	    	return this[BOXED_DATA_PROPERTY]['{{ $field.SerializableName }}'];
	    },

	    set: function(value) {
	    	this[BOXED_DATA_PROPERTY]['{{ $field.SerializableName }}'] = value;
	    }
	  });
	{{ end }}

  	{{ .TypeSignatureMethod }}
});
`

// nominalTemplateStr defines the template for generating a nominal type.
const nominalTemplateStr = `
this.$type('{{ .Type.Name }}', {{ .HasGenerics }}, '{{ .Alias }}', function({{ .Generics }}) {
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

	{{range $idx, $kv := .GenerateImplementedMembers.UnsafeIter }}
  	  {{ emit $kv.Value }}
  	{{end}}

  	{{ .TypeSignatureMethod }}
});
`
