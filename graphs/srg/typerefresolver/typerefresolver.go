// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typerefresolver

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/streamrail/concurrent-map"
)

// TypeReferenceResolver defines a helper type for resolving SRG type references into type
// system type references.
type TypeReferenceResolver struct {
	srg            *srg.SRG               // The SRG of the references being resolved.
	readWriteCache cmap.ConcurrentMap     // Cache of resolved type references.
	readOnlyCache  map[string]interface{} // Read-only cache of resolved type references.
	frozen         bool                   // Whether the cache has been frozen.
}

// NewResolver returns a new type reference resolver for the given SRG.
func NewResolver(srg *srg.SRG) *TypeReferenceResolver {
	return &TypeReferenceResolver{
		srg:            srg,
		readWriteCache: cmap.New(),
		readOnlyCache:  map[string]interface{}{},
		frozen:         false,
	}
}

// FreezeCache freezes the resolver's cache, no longer allowing any new references to be added to the cache.
func (trr *TypeReferenceResolver) FreezeCache() {
	trr.readOnlyCache = trr.readWriteCache.Items()
	trr.frozen = true
}

// ResolveTypeRef builds a type graph type reference from the SRG type reference, fully resolving it.
func (trr *TypeReferenceResolver) ResolveTypeRef(typeref srg.SRGTypeRef, tdg *typegraph.TypeGraph) (typegraph.TypeReference, error) {
	if trr.frozen {
		if found, ok := trr.readOnlyCache[string(typeref.NodeId)]; ok {
			return found.(typegraph.TypeReference), nil
		}

		return trr.resolveTypeRef(typeref, tdg)
	} else {
		if found, ok := trr.readWriteCache.Get(string(typeref.NodeId)); ok {
			return found.(typegraph.TypeReference), nil
		}

		tr, err := trr.resolveTypeRef(typeref, tdg)
		if err == nil {
			trr.readWriteCache.Set(string(typeref.NodeId), tr)
		}

		return tr, err
	}
}

func (trr *TypeReferenceResolver) resolveTypeRef(typeref srg.SRGTypeRef, tdg *typegraph.TypeGraph) (typegraph.TypeReference, error) {
	switch typeref.RefKind() {
	case srg.TypeRefVoid:
		return tdg.VoidTypeReference(), nil

	case srg.TypeRefAny:
		return tdg.AnyTypeReference(), nil

	case srg.TypeRefStruct:
		return tdg.StructTypeReference(), nil

	case srg.TypeRefMapping:
		innerType, err := trr.ResolveTypeRef(typeref.InnerReference(), tdg)
		if err != nil {
			return tdg.AnyTypeReference(), err
		}

		return tdg.NewTypeReference(tdg.MappingType(), innerType), nil

	case srg.TypeRefSlice:
		innerType, err := trr.ResolveTypeRef(typeref.InnerReference(), tdg)
		if err != nil {
			return tdg.AnyTypeReference(), err
		}

		return tdg.NewTypeReference(tdg.SliceType(), innerType), nil

	case srg.TypeRefStream:
		innerType, err := trr.ResolveTypeRef(typeref.InnerReference(), tdg)
		if err != nil {
			return tdg.AnyTypeReference(), err
		}

		return tdg.NewTypeReference(tdg.StreamType(), innerType), nil

	case srg.TypeRefNullable:
		innerType, err := trr.ResolveTypeRef(typeref.InnerReference(), tdg)
		if err != nil {
			return tdg.AnyTypeReference(), err
		}

		return innerType.AsNullable(), nil

	case srg.TypeRefPath:
		// Resolve the package type for the type ref.
		resolvedTypeInfo, found := typeref.ResolveType()
		if !found {
			sourceError := compilercommon.SourceErrorf(typeref.Location(),
				"Type '%s' could not be found",
				typeref.ResolutionPath())

			return tdg.AnyTypeReference(), sourceError
		}

		// If the type information refers to an SRG type or generic, find the node directly
		// in the type graph.
		var constructedRef = tdg.AnyTypeReference()
		if !resolvedTypeInfo.IsExternalPackage {
			// Get the type in the type graph.
			resolvedType, hasResolvedType := tdg.GetTypeForSourceNode(resolvedTypeInfo.ResolvedType.Node())
			if !hasResolvedType {
				panic(fmt.Sprintf("Could not find typegraph type for SRG type %v", resolvedTypeInfo.ResolvedType.Name()))
			}

			constructedRef = tdg.NewTypeReference(resolvedType)
		} else {
			// Otherwise, we search for the type in the type graph based on the package from which it
			// was imported.
			resolvedType, hasResolvedType := tdg.ResolveTypeUnderPackage(resolvedTypeInfo.ExternalPackageTypePath, resolvedTypeInfo.ExternalPackage)
			if !hasResolvedType {
				sourceError := compilercommon.SourceErrorf(typeref.Location(),
					"Type '%s' could not be found",
					typeref.ResolutionPath())

				return tdg.AnyTypeReference(), sourceError
			}

			constructedRef = tdg.NewTypeReference(resolvedType)
		}

		// Add the generics.
		if typeref.HasGenerics() {
			for _, srgGeneric := range typeref.Generics() {
				genericTypeRef, err := trr.ResolveTypeRef(srgGeneric, tdg)
				if err != nil {
					return tdg.AnyTypeReference(), err
				}

				constructedRef = constructedRef.WithGeneric(genericTypeRef)
			}
		}

		// Add the parameters.
		if typeref.HasParameters() {
			for _, srgParameter := range typeref.Parameters() {
				parameterTypeRef, err := trr.ResolveTypeRef(srgParameter, tdg)
				if err != nil {
					return tdg.AnyTypeReference(), err
				}
				constructedRef = constructedRef.WithParameter(parameterTypeRef)
			}
		}

		return constructedRef, nil

	default:
		panic(fmt.Sprintf("Unknown kind of SRG type ref: %v", typeref.RefKind()))
		return tdg.AnyTypeReference(), nil
	}
}
