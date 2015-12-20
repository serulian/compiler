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

	referredType := typeRef.ReferredTypeDecl()
	if referredType.TypeKind() == typegraph.GenericType {
		return referredType.Name()
	}

	// Add the module path.
	modulePath := "$g." + p.GetModulePath(referredType.ParentModule())

	// Add the type name.
	typePath := modulePath + "." + referredType.Name()

	// If the type has generics, invoke them under the given reference type.
	generics := referredType.Generics()
	if len(generics) > 0 {

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

	return typePath
}

// GetStaticTypePath returns the global path for the given defined type.
func (p *Pather) GetStaticTypePath(typedecl typegraph.TGTypeDecl, referenceType typegraph.TypeReference) string {
	instanceTypeRef := typedecl.GetTypeReference().TransformUnder(referenceType)
	return p.TypeReferenceCall(instanceTypeRef)
}

// GetStaticMemberPath returns the global path for the given statically defined type member.
func (p *Pather) GetStaticMemberPath(member typegraph.TGMember, referenceType typegraph.TypeReference) string {
	name := strings.Replace(unidecode.Unidecode(member.Name()), "*", "$", 1)
	parent := member.Parent()
	if parent.IsType() {
		return p.GetStaticTypePath(parent.(typegraph.TGTypeDecl), referenceType) + "." + name
	} else {
		return "$g." + p.GetModulePath(parent.(typegraph.TGModule)) + "." + name
	}
}

// GetModulePath returns the global path for the given module.
func (p *Pather) GetModulePath(module typegraph.TGModule) string {
	// We create the exported path based on the location of this module's source file relative
	// to the entrypoint file.
	srgModule, _ := module.SRGModule()

	basePath := filepath.Dir(p.graph.RootSourceFilePath)
	rel, err := filepath.Rel(basePath, string(srgModule.InputSource()))
	if err != nil {
		panic(err)
	}

	rel = strings.Replace(rel, "../", "_", -1)
	rel = strings.Replace(rel, "/", ".", -1)
	rel = rel[0 : len(rel)-5]
	return rel
}
