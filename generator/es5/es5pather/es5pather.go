// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The es5pather package implements rules for generating access paths and names.
package es5pather

import (
	"path/filepath"
	"strings"

	"github.com/rainycape/unidecode"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
)

// Pather defines a helper type for generating paths.
type Pather struct {
	graph *compilergraph.SerulianGraph // The root graph.
}

// New creates a new path generator for the given graph.
func New(graph *compilergraph.SerulianGraph) *Pather {
	return &Pather{graph}
}

// TypeReferenceCall returns source for retrieving an object reference to the type defined by the given
// type reference.
func (p *Pather) TypeReferenceCall(typeRef typegraph.TypeReference) string {
	if typeRef.IsAny() {
		return "$t.any"
	}

	if typeRef.IsNull() {
		return "$t.null"
	}

	if typeRef.IsVoid() {
		return "$t.void"
	}

	referredType := typeRef.ReferredType()
	if referredType.TypeKind() == typegraph.GenericType {
		return referredType.Name()
	}

	// Add the type name.
	typePath := p.GetTypePath(referredType)

	// If there are no generics, then simply return the type path.
	if !typeRef.HasGenerics() {
		return typePath
	}

	// Invoke the type with generics (if any).
	var genericsString = "("
	for index, generic := range typeRef.Generics() {
		if index > 0 {
			genericsString = genericsString + ", "
		}

		genericsString = genericsString + p.TypeReferenceCall(generic)
	}

	genericsString = genericsString + ")"
	return typePath + genericsString
}

// InnerInstanceName returns the name of an inner instance of the given type, when accessed under a
// type instance which structurally composes it.
func (p *Pather) InnerInstanceName(innerType typegraph.TypeReference) string {
	var name = unidecode.Unidecode(innerType.ReferredType().Name())
	if !innerType.HasGenerics() {
		return name
	}

	for _, generic := range innerType.Generics() {
		name = name + "$"
		name = name + p.InnerInstanceName(generic)
	}

	return name
}

// GetMemberName returns the name of the given member.
func (p *Pather) GetMemberName(member typegraph.TGMember) string {
	return strings.Replace(unidecode.Unidecode(member.ChildName()), "*", "$", 1)
}

// GetStaticTypePath returns the global path for the given defined type.
func (p *Pather) GetStaticTypePath(typedecl typegraph.TGTypeDecl, referenceType typegraph.TypeReference) string {
	instanceTypeRef := typedecl.GetTypeReference().TransformUnder(referenceType)
	return p.TypeReferenceCall(instanceTypeRef)
}

// GetStaticMemberPath returns the global path for the given statically defined type member.
func (p *Pather) GetStaticMemberPath(member typegraph.TGMember, referenceType typegraph.TypeReference) string {
	// TODO(jschorr): We should generalize non-SRG support.
	if member.Name() == "new" && !member.IsPromising() {
		parentType, _ := member.ParentType()
		if strings.HasSuffix(parentType.ParentModule().Path(), ".webidl") {
			return "$t.nativenew(" + p.GetTypePath(parentType) + ")"
		}
	}

	name := p.GetMemberName(member)
	parent := member.Parent()
	if parent.IsType() {
		return p.GetStaticTypePath(parent.(typegraph.TGTypeDecl), referenceType) + "." + name
	} else {
		return p.GetModulePath(parent.(typegraph.TGModule)) + "." + name
	}
}

// GetTypePath returns the global path for the given type.
func (p *Pather) GetTypePath(typedecl typegraph.TGTypeDecl) string {
	if typedecl.TypeKind() == typegraph.GenericType {
		return typedecl.Name()
	}

	return p.GetModulePath(typedecl.ParentModule()) + "." + typedecl.Name()
}

// GetModulePath returns the global path for the given module.
func (p *Pather) GetModulePath(module typegraph.TGModule) string {
	// TODO(jschorr): We should generalize non-SRG support.
	if strings.HasSuffix(module.Path(), ".webidl") {
		return "$global"
	}

	return "$g." + p.GetRelativeModulePath(module)
}

// GetRelativeModulePath returns the relative path for the given module.
func (p *Pather) GetRelativeModulePath(module typegraph.TGModule) string {
	// We create the exported path based on the location of this module's source file relative
	// to the entrypoint file.
	basePath := filepath.Dir(p.graph.RootSourceFilePath)
	rel, err := filepath.Rel(basePath, module.Path())
	if err != nil {
		panic(err)
	}

	rel = strings.Replace(rel, "../", "_", -1)
	rel = strings.Replace(rel, "/", ".", -1)
	rel = rel[0 : len(rel)-5]

	if rel[0] == '.' {
		rel = rel[1 : len(rel)-1]
	}

	return rel
}
