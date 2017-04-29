// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
)

var _ = fmt.Print

type scopeInfoBuilder struct {
	info *proto.ScopeInfo
}

func newScope() *scopeInfoBuilder {
	return &scopeInfoBuilder{
		info: &proto.ScopeInfo{
			Kind: proto.ScopeKind_VALUE,
		},
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
	sib.info.IsValid = isValid
	return sib
}

// ResolvingTypeOf marks the scope as resolving the type of the given scope.
func (sib *scopeInfoBuilder) ResolvingTypeOf(scope *proto.ScopeInfo) *scopeInfoBuilder {
	sib.info.ResolvedType = scope.GetResolvedType()
	return sib
}

// Assignable marks the scope as being assignable with a value of the given type.
func (sib *scopeInfoBuilder) Assignable(assignable typegraph.TypeReference) *scopeInfoBuilder {
	sib.info.AssignableType = assignable.Value()
	return sib
}

// Resolving marks the scope as resolving a value of the given type.
func (sib *scopeInfoBuilder) Resolving(resolved typegraph.TypeReference) *scopeInfoBuilder {
	sib.info.ResolvedType = resolved.Value()
	return sib
}

// WithStaticType marks the scope as having the given static type.
func (sib *scopeInfoBuilder) WithStaticType(static typegraph.TypeReference) *scopeInfoBuilder {
	sib.info.StaticType = static.Value()
	return sib
}

// WithGenericType marks the scope as having the given generic type.
func (sib *scopeInfoBuilder) WithGenericType(generic typegraph.TypeReference) *scopeInfoBuilder {
	sib.info.GenericType = generic.Value()
	return sib
}

// AssignableResolvedTypeOf marks the scope as being assignable of the *resolved* type of the given scope.
func (sib *scopeInfoBuilder) AssignableResolvedTypeOf(scope *proto.ScopeInfo) *scopeInfoBuilder {
	sib.info.AssignableType = scope.GetResolvedType()
	return sib
}

// ReturningTypeOf marks the scope as returning the return type of the given scope.
func (sib *scopeInfoBuilder) ReturningTypeOf(scope *proto.ScopeInfo) *scopeInfoBuilder {
	sib.info.ReturnedType = scope.GetReturnedType()
	sib.info.IsSettlingScope = scope.GetIsSettlingScope()
	return sib
}

// ReturningResolvedTypeOf marks the scope as returning the *resolved* type of the given scope.
func (sib *scopeInfoBuilder) ReturningResolvedTypeOf(scope *proto.ScopeInfo) *scopeInfoBuilder {
	sib.info.ReturnedType = scope.GetResolvedType()
	sib.info.IsSettlingScope = true
	return sib
}

// Returning marks the scope as returning a value of the given type.
func (sib *scopeInfoBuilder) Returning(returning typegraph.TypeReference, settlesScope bool) *scopeInfoBuilder {
	sib.info.ReturnedType = returning.Value()
	sib.info.IsSettlingScope = settlesScope
	return sib
}

// IsSettlingScope marks the scope as settling the function.
func (sib *scopeInfoBuilder) IsSettlingScope() *scopeInfoBuilder {
	sib.info.IsSettlingScope = true
	return sib
}

// IsTerminatingStatement marks the scope as containing a terminating statement.
func (sib *scopeInfoBuilder) IsTerminatingStatement() *scopeInfoBuilder {
	sib.info.IsTerminatingStatement = true
	return sib
}

type typeModifier func(typeRef typegraph.TypeReference) typegraph.TypeReference

// ForNamedScopeUnderModifiedType points the scope to the referred named scope, with its value
// type being transformed under the given parent type and then transformed by the modifier.
func (sib *scopeInfoBuilder) ForNamedScopeUnderModifiedType(info namedScopeInfo, parentType typegraph.TypeReference, modifier typeModifier, context scopeContext) *scopeInfoBuilder {
	// If the named scope is generic, then mark the scope with the *generic* type transformed under
	// the parent type. This ensures that chained, static generic calls (like SomeType<int>.SomeConstructor<bool>)
	// will resolve the first generic (e.g. `int` in this example) properly. Otherwise, we transform the value
	// type.
	if info.IsGeneric() {
		genericType, isValid := info.ValueOrGenericType(context)
		if !isValid {
			panic("Expected valid scope")
		}

		transformedGenericType := genericType.TransformUnder(parentType)
		sib.ForNamedScope(info, context).WithGenericType(transformedGenericType)
	} else {
		valueType, isValid := info.ValueType(context)
		if !isValid {
			panic("Expected valid scope")
		}

		transformedValueType := valueType.TransformUnder(parentType)
		sib.ForNamedScope(info, context).Resolving(modifier(transformedValueType))

		if info.IsAssignable() {
			assignableType, _ := info.AssignableType(context)
			transformedAssignableType := assignableType.TransformUnder(parentType)
			sib.Assignable(modifier(transformedAssignableType))
		}
	}

	return sib
}

// ForNamedScopeUnderType points the scope to the referred named scope, with its value
// type being transformed under the given parent type.
func (sib *scopeInfoBuilder) ForNamedScopeUnderType(info namedScopeInfo, parentType typegraph.TypeReference, context scopeContext) *scopeInfoBuilder {
	modifier := func(typeRef typegraph.TypeReference) typegraph.TypeReference {
		return typeRef
	}

	return sib.ForNamedScopeUnderModifiedType(info, parentType, modifier, context)
}

// ForAnonymousScope points the scope to an anonymously scope.
func (sib *scopeInfoBuilder) ForAnonymousScope(typegraph *typegraph.TypeGraph) *scopeInfoBuilder {
	sib.info.IsAnonymousReference = true
	return sib.Resolving(typegraph.VoidTypeReference()).Valid()
}

// CallsOperator marks the scope as being the result of a call to the specified operator.
func (sib *scopeInfoBuilder) CallsOperator(op typegraph.TGMember) *scopeInfoBuilder {
	sib.info.CalledOpReference = &proto.ScopeReference{}
	sib.info.CalledOpReference.ReferencedNode = string(op.GraphNode.NodeId)
	sib.info.CalledOpReference.IsSRGNode = false
	return sib
}

// ForNamedScope points the scope to the referred named scope.
func (sib *scopeInfoBuilder) ForNamedScope(info namedScopeInfo, context scopeContext) *scopeInfoBuilder {
	if info.IsGeneric() {
		genericType, _ := info.ValueOrGenericType(context)
		sib.info.Kind = proto.ScopeKind_GENERIC
		sib.WithGenericType(genericType)
	} else if info.IsStatic() {
		sib.info.Kind = proto.ScopeKind_STATIC
		sib.WithStaticType(info.StaticType(context))
	}

	sib.info.NamedReference = &proto.ScopeReference{}

	if info.typeInfo != nil {
		sib.info.NamedReference.ReferencedNode = string(info.typeInfo.Node().NodeId)
		sib.info.NamedReference.IsSRGNode = false
	} else {
		sib.info.NamedReference.ReferencedNode = string(info.srgInfo.GraphNode.NodeId)
		sib.info.NamedReference.IsSRGNode = true
	}

	if info.IsAssignable() {
		assignableType, _ := info.AssignableType(context)
		sib.Assignable(assignableType)
	}

	valueType, _ := info.ValueType(context)
	return sib.Resolving(valueType).Valid()
}

// WithKind sets the kind of this scope to the given kind.
func (sib *scopeInfoBuilder) WithKind(kind proto.ScopeKind) *scopeInfoBuilder {
	sib.info.Kind = kind
	return sib
}

// WithLabel adds the label to this scope.
func (sib *scopeInfoBuilder) WithLabel(label proto.ScopeLabel) *scopeInfoBuilder {
	sib.info.Labels = append(sib.info.Labels, label)
	return sib
}

// WithLabelSet adds all the labels found in the given set to this scope.
func (sib *scopeInfoBuilder) WithLabelSet(labelSet *statementLabelSet) *scopeInfoBuilder {
	sib.info.Labels = labelSet.GetLabels()
	return sib
}

// LabelSetOf sets the label set for this scope to the set found on the other scope.
func (sib *scopeInfoBuilder) LabelSetOf(scope *proto.ScopeInfo) *scopeInfoBuilder {
	sib.info.Labels = scope.Labels
	return sib
}

// LabelSetOf sets the label set for this scope to the set found on the other scope, except those
// labels specified.
func (sib *scopeInfoBuilder) LabelSetOfExcept(scope *proto.ScopeInfo, except ...proto.ScopeLabel) *scopeInfoBuilder {
	labelSet := newLabelSet()
	labelSet.AppendLabelsOf(scope)
	labelSet.RemoveLabels(except...)
	return sib.WithLabelSet(labelSet)
}

// Targets sets the targetted node for a break or continue statement.
func (sib *scopeInfoBuilder) Targets(node compilergraph.GraphNode) *scopeInfoBuilder {
	sib.info.TargetedReference = &proto.ScopeReference{}
	sib.info.TargetedReference.ReferencedNode = string(node.NodeId)
	sib.info.TargetedReference.IsSRGNode = true
	return sib
}

// GetScope returns the scope constructed.
func (sib *scopeInfoBuilder) GetScope() proto.ScopeInfo {
	return *sib.info
}
