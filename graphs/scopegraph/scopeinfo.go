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

// WithGenericType marks the scope as having the given generic type.
func (sib *scopeInfoBuilder) WithGenericType(generic typegraph.TypeReference) *scopeInfoBuilder {
	genericValue := generic.Value()
	sib.info.GenericType = &genericValue
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
	settlesValue := scope.GetIsSettlingScope()

	sib.info.ReturnedType = &returnedValue
	sib.info.IsSettlingScope = &settlesValue
	return sib
}

// ReturningResolvedTypeOf marks the scope as returning the *resolved* type of the given scope.
func (sib *scopeInfoBuilder) ReturningResolvedTypeOf(scope *proto.ScopeInfo) *scopeInfoBuilder {
	resolvedValue := scope.GetResolvedType()
	settlesScope := true

	sib.info.ReturnedType = &resolvedValue
	sib.info.IsSettlingScope = &settlesScope
	return sib
}

// Returning marks the scope as returning a value of the given type.
func (sib *scopeInfoBuilder) Returning(returning typegraph.TypeReference, settlesScope bool) *scopeInfoBuilder {
	returnedValue := returning.Value()
	sib.info.ReturnedType = &returnedValue
	sib.info.IsSettlingScope = &settlesScope
	return sib
}

// IsSettlingScope marks the scope as settling the function.
func (sib *scopeInfoBuilder) IsSettlingScope() *scopeInfoBuilder {
	trueValue := true
	sib.info.IsSettlingScope = &trueValue
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
func (sib *scopeInfoBuilder) ForNamedScopeUnderModifiedType(info namedScopeInfo, parentType typegraph.TypeReference, modifier typeModifier, context scopeContext) *scopeInfoBuilder {
	// If the named scope is generic, then mark the scope with the *generic* type transformed under
	// the parent type. This ensures that chained, static generic calls (like SomeType<int>.SomeConstructor<bool>)
	// will resolve the first generic (e.g. `int` in this example) properly. Otherwise, we transform the value
	// type.
	if info.IsGeneric() {
		transformedGenericType := info.ValueOrGenericType(context).TransformUnder(parentType)
		sib.ForNamedScope(info, context).WithGenericType(transformedGenericType)
	} else {
		transformedValueType := info.ValueType(context).TransformUnder(parentType)
		sib.ForNamedScope(info, context).Resolving(modifier(transformedValueType))

		if info.IsAssignable() {
			transformedAssignableType := info.AssignableType(context).TransformUnder(parentType)
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
	trueValue := true
	sib.info.IsAnonymousReference = &trueValue
	return sib.Resolving(typegraph.VoidTypeReference()).Valid()
}

// CallsOperator marks the scope as being the result of a call to the specified operator.
func (sib *scopeInfoBuilder) CallsOperator(op typegraph.TGMember) *scopeInfoBuilder {
	sib.info.CalledOpReference = &proto.ScopeReference{}

	falseValue := false
	namedId := string(op.GraphNode.NodeId)
	sib.info.CalledOpReference.ReferencedNode = &namedId
	sib.info.CalledOpReference.IsSRGNode = &falseValue
	return sib
}

// ForNamedScope points the scope to the referred named scope.
func (sib *scopeInfoBuilder) ForNamedScope(info namedScopeInfo, context scopeContext) *scopeInfoBuilder {
	if info.IsGeneric() {
		genericKind := proto.ScopeKind_GENERIC
		sib.info.Kind = &genericKind
		sib.WithGenericType(info.ValueOrGenericType(context))
	} else if info.IsStatic() {
		staticKind := proto.ScopeKind_STATIC
		sib.info.Kind = &staticKind
		sib.WithStaticType(info.StaticType(context))
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

	if info.IsAssignable() {
		sib.Assignable(info.AssignableType(context))
	}

	return sib.Resolving(info.ValueType(context)).Valid()
}

// WithKind sets the kind of this scope to the given kind.
func (sib *scopeInfoBuilder) WithKind(kind proto.ScopeKind) *scopeInfoBuilder {
	sib.info.Kind = &kind
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

// GetScope returns the scope constructed.
func (sib *scopeInfoBuilder) GetScope() proto.ScopeInfo {
	return *sib.info
}
