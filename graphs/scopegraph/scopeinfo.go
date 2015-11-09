// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
)

var _ = fmt.Print

type scopeInfoBuilder struct {
	info *proto.ScopeInfo
}

func newScope() *scopeInfoBuilder {
	return &scopeInfoBuilder{
		info: &proto.ScopeInfo{},
	}
}

// Valid marks the scope as valid.
func (sib *scopeInfoBuilder) Valid() *scopeInfoBuilder {
	return sib.IsValid(true)
}

// Invalid marks the scope as invalid.
func (sib *scopeInfoBuilder) Invalid() *scopeInfoBuilder {
	return sib.IsValid(false)
}

// IsValid marks the scope as valid or invalid.
func (sib *scopeInfoBuilder) IsValid(isValid bool) *scopeInfoBuilder {
	sib.info.IsValid = &isValid
	return sib
}

// ResolvingTypeOf marks the scope as resolving the type of the given scope.
func (sib *scopeInfoBuilder) ResolvingTypeOf(scope *proto.ScopeInfo) *scopeInfoBuilder {
	resolvedValue := scope.GetResolvedType()
	sib.info.ResolvedType = &resolvedValue
	return sib
}

// Assignable marks the scope as being assignable with a value of the given type.
func (sib *scopeInfoBuilder) Assignable(assignable typegraph.TypeReference) *scopeInfoBuilder {
	assignableValue := assignable.Value()
	sib.info.AssignableType = &assignableValue
	return sib
}

// Resolving marks the scope as resolving a value of the given type.
func (sib *scopeInfoBuilder) Resolving(resolved typegraph.TypeReference) *scopeInfoBuilder {
	resolvedValue := resolved.Value()
	sib.info.ResolvedType = &resolvedValue
	return sib
}

// WithStaticType marks the scope as having the given static type.
func (sib *scopeInfoBuilder) WithStaticType(static typegraph.TypeReference) *scopeInfoBuilder {
	staticValue := static.Value()
	sib.info.StaticType = &staticValue
	return sib
}

// AssignableResolvedTypeOf marks the scope as being assignable of the *resolved* type of the given scope.
func (sib *scopeInfoBuilder) AssignableResolvedTypeOf(scope *proto.ScopeInfo) *scopeInfoBuilder {
	resolvedValue := scope.GetResolvedType()
	sib.info.AssignableType = &resolvedValue
	return sib
}

// ReturningTypeOf marks the scope as returning the return type of the given scope.
func (sib *scopeInfoBuilder) ReturningTypeOf(scope *proto.ScopeInfo) *scopeInfoBuilder {
	returnedValue := scope.GetReturnedType()
	sib.info.ReturnedType = &returnedValue
	return sib
}

// ReturningResolvedTypeOf marks the scope as returning the *resolved* type of the given scope.
func (sib *scopeInfoBuilder) ReturningResolvedTypeOf(scope *proto.ScopeInfo) *scopeInfoBuilder {
	resolvedValue := scope.GetResolvedType()
	sib.info.ReturnedType = &resolvedValue
	return sib
}

// Returning marks the scope as returning a value of the given type.
func (sib *scopeInfoBuilder) Returning(returning typegraph.TypeReference) *scopeInfoBuilder {
	returnedValue := returning.Value()
	sib.info.ReturnedType = &returnedValue
	return sib
}

// IsTerminatingStatement marks the scope as containing a terminating statement.
func (sib *scopeInfoBuilder) IsTerminatingStatement() *scopeInfoBuilder {
	trueValue := true
	sib.info.IsTerminatingStatement = &trueValue
	return sib
}

type typeModifier func(typeRef typegraph.TypeReference) typegraph.TypeReference

// ForNamedScopeUnderModifiedType points the scope to the referred named scope, with its value
// type being transformed under the given parent type and then transformed by the modifier.
func (sib *scopeInfoBuilder) ForNamedScopeUnderModifiedType(info namedScopeInfo, parentType typegraph.TypeReference, modifier typeModifier) *scopeInfoBuilder {
	transformedValueType := info.ValueType().TransformUnder(parentType)
	return sib.ForNamedScope(info).Resolving(modifier(transformedValueType))
}

// ForNamedScopeUnderType points the scope to the referred named scope, with its value
// type being transformed under the given parent type.
func (sib *scopeInfoBuilder) ForNamedScopeUnderType(info namedScopeInfo, parentType typegraph.TypeReference) *scopeInfoBuilder {
	transformedValueType := info.ValueType().TransformUnder(parentType)
	return sib.ForNamedScope(info).Resolving(transformedValueType)
}

// ForNamedScope points the scope to the referred named scope.
func (sib *scopeInfoBuilder) ForNamedScope(info namedScopeInfo) *scopeInfoBuilder {
	if info.IsGeneric() {
		genericKind := proto.ScopeKind_GENERIC
		sib.info.Kind = &genericKind
	} else if info.IsStatic() {
		staticKind := proto.ScopeKind_STATIC
		sib.info.Kind = &staticKind
		sib.WithStaticType(info.StaticType())
	}

	sib.info.NamedReference = &proto.ScopeReference{}

	if info.typeInfo != nil {
		falseValue := false
		namedId := string(info.typeInfo.Node().NodeId)
		sib.info.NamedReference.ReferencedNode = &namedId
		sib.info.NamedReference.IsSRGNode = &falseValue
	} else {
		trueValue := true
		namedId := string(info.srgInfo.GraphNode.NodeId)
		sib.info.NamedReference.ReferencedNode = &namedId
		sib.info.NamedReference.IsSRGNode = &trueValue
	}

	return sib.Resolving(info.ValueType()).Valid()
}

// WithKind sets the kind of this scope to the given kind.
func (sib *scopeInfoBuilder) WithKind(kind proto.ScopeKind) *scopeInfoBuilder {
	sib.info.Kind = &kind
	return sib
}

// GetScope returns the scope constructed.
func (sib *scopeInfoBuilder) GetScope() proto.ScopeInfo {
	return *sib.info
}
