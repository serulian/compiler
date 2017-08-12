// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shared

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rainycape/unidecode"

	"unicode"

	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/typegraph"
)

// Pather defines a helper type for generating paths.
type Pather struct {
	scopegraph *scopegraph.ScopeGraph
}

// NewPather creates a new path generator for the given graph.
func NewPather(scopegraph *scopegraph.ScopeGraph) Pather {
	return Pather{scopegraph}
}

// TypeReferenceCall returns source for retrieving an object reference to the type defined by the given
// type reference.
func (p Pather) TypeReferenceCall(typeRef typegraph.TypeReference) string {
	if typeRef.IsAny() {
		return "$t.any"
	}

	if typeRef.IsStruct() {
		return "$t.struct"
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
func (p Pather) InnerInstanceName(innerType typegraph.TypeReference) string {
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
func (p Pather) GetMemberName(member typegraph.TGMember) string {
	return strings.Replace(unidecode.Unidecode(member.ChildName()), "*", "$", 1)
}

// GetSetterName returns the name of the setter for the given member.
func (p Pather) GetSetterName(member typegraph.TGMember) string {
	return "set$" + p.GetMemberName(member)
}

// GetStaticTypePath returns the global path for the given defined type.
func (p Pather) GetStaticTypePath(typedecl typegraph.TGTypeDecl, referenceType typegraph.TypeReference) string {
	instanceTypeRef := typedecl.GetTypeReference().TransformUnder(referenceType)
	return p.TypeReferenceCall(instanceTypeRef)
}

// GetStaticMemberPath returns the global path for the given statically defined type member.
func (p Pather) GetStaticMemberPath(member typegraph.TGMember, referenceType typegraph.TypeReference) string {
	sourceGraphID := member.SourceGraphId()
	if sourceGraphID != "srg" {
		integration := p.scopegraph.MustGetLanguageIntegration(sourceGraphID)
		pathHandler := integration.PathHandler()
		if pathHandler != nil {
			staticPath := integration.PathHandler().GetStaticMemberPath(member, referenceType)
			if staticPath != "" {
				return staticPath
			}
		}
	}

	name := p.GetMemberName(member)
	parent := member.Parent()
	if parent.IsType() {
		return p.GetStaticTypePath(parent.(typegraph.TGTypeDecl), referenceType) + "." + name
	}

	return p.GetModulePath(parent.(typegraph.TGModule)) + "." + name
}

// GetTypePath returns the global path for the given type.
func (p Pather) GetTypePath(typedecl typegraph.TGTypeDecl) string {
	if typedecl.TypeKind() == typegraph.GenericType {
		return typedecl.Name()
	}

	return p.GetModulePath(typedecl.ParentModule()) + "." + typedecl.Name()
}

// GetModulePath returns the global path for the given module.
func (p Pather) GetModulePath(module typegraph.TGModule) string {
	sourceGraphID := module.SourceGraphId()
	if sourceGraphID != "srg" {
		integration := p.scopegraph.MustGetLanguageIntegration(sourceGraphID)
		pathHandler := integration.PathHandler()
		if pathHandler != nil {
			modulePath := pathHandler.GetModulePath(module)
			if modulePath != "" {
				return modulePath
			}
		}
	}

	return "$g." + p.GetRelativeModulePath(module)
}

// GetRelativeModulePath returns the relative path for the given module.
func (p Pather) GetRelativeModulePath(module typegraph.TGModule) string {
	// We create the exported path based on the location of this module's source file relative
	// to the entrypoint file.
	basePath := filepath.Dir(p.scopegraph.RootSourceFilePath())
	rel, err := filepath.Rel(basePath, module.Path())
	if err != nil {
		rel = module.Path()
	}

	return normalizeModulePath(rel)
}

var allowedModulePathCharacters, _ = regexp.Compile("[^a-zA-Z_0-9\\$\\.]")

func normalizeModulePath(rel string) string {
	if rel == "" {
		return ""
	}

	// Replace any relative pathing with underscores.
	rel = strings.Replace(rel, "../", "__", -1)

	// Strip off any starting .
	if rel[0] == '.' {
		rel = rel[1:len(rel)]
	}

	// Remove any Serulian source file suffix.
	if strings.HasSuffix(rel, ".seru") {
		rel = rel[0 : len(rel)-5]
	}

	// Replace any dot parts and change slashes into submodules.
	parts := strings.Split(rel, "/")
	updatedParts := make([]string, 0, len(parts))

	for _, part := range parts {
		if part == "." || part == "" {
			continue
		}

		// If the part does not start with a letter or underscore, prepend a prefix, as literals in ES5 cannot start
		// with numbers.
		newPart := part
		firstRune := rune(part[0])
		if firstRune != '_' && !unicode.IsLetter(firstRune) {
			newPart = "m$" + part
		}

		updatedParts = append(updatedParts, newPart)
	}

	path := strings.Join(updatedParts, ".")

	// Replace any other bad characters.
	path = allowedModulePathCharacters.ReplaceAllLiteralString(path, "_")
	return path
}
