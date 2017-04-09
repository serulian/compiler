// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
)

var _ = fmt.Printf

// typeCompositionProcessor defines a helper for copying over members from agents composed
// into a class/agent, as well as validating that this composition doesn't conflict.
type typeCompositionProcessor struct {
	composingType TGTypeDecl
	modifier      compilergraph.GraphLayerModifier
	tdg           *TypeGraph
}

type memberStatus struct {
	Name       string
	Title      string
	IsRequired bool
}

// processComposition performs all composition declaration, building and validation.
func (p typeCompositionProcessor) processComposition() bool {
	// Build a map of all existing members. If the boolean is `true`, then it was
	// defined on the original type, vs on a composed type.
	members := map[string]memberStatus{}
	for _, member := range p.composingType.Members() {
		members[member.Name()] = memberStatus{member.Name(), member.Title(), member.IsRequiredField()}
	}

	// For each composed agent type, add a field representing the agent itself
	// and then copy over its members.
	composedAgents := p.composingType.ComposedAgents()
	for _, agent := range composedAgents {
		// Add the field itself.
		if !p.addCompositionField(agent, members) {
			return false
		}
	}

	for _, agent := range composedAgents {
		// Copy over members and operators.
		p.buildComposedMembership(agent, members, NodePredicateMember, MemberResolutionInstance)
		p.buildComposedMembership(agent, members, NodePredicateTypeOperator, MemberResolutionOperator)
	}

	return true
}

// addCompositionField adds a composition field to the composingType for the given agent, checking
// to make sure the field won't be shadowed.
func (p typeCompositionProcessor) addCompositionField(agent TGAgentReference, members map[string]memberStatus) bool {
	agentTypeRef := agent.AgentType()
	agentType := agentTypeRef.ReferredType()

	// Make sure the composition field will not shadow another field in the type.
	status, ok := members[agent.CompositionName()]
	if ok {
		p.tdg.decorateWithError(p.modifier.Modify(p.composingType.GraphNode),
			"%v %v cannot compose %v %v as its compositon name '%v' shadows %v '%v'",
			p.composingType.Title(),
			p.composingType.Name(),
			agentType.Title(),
			agentType.Name(),
			agent.CompositionName(),
			status.Title,
			status.Name)
		return false
	}

	// Add the composition field to the type.
	builder := &MemberBuilder{
		tdg:        p.tdg,
		modifier:   p.modifier,
		parent:     p.composingType,
		isOperator: false,
	}

	builtMember := builder.Name(agent.CompositionName()).Define()
	decorator := &MemberDecorator{
		tdg:                p.tdg,
		modifier:           p.modifier,
		member:             builtMember,
		memberName:         agent.CompositionName(),
		genericConstraints: map[compilergraph.GraphNode]TypeReference{},
		tags:               map[string]string{},
	}

	decorator.
		MemberKind(FieldMemberSignature).
		MemberType(agentTypeRef).
		ReadOnly(true).
		Decorate()

	members[agent.CompositionName()] = memberStatus{agent.CompositionName(), "type member", true}
	return true
}

// buildComposedMembership copies any applicable type members from the composed type references to the given type.
func (p typeCompositionProcessor) buildComposedMembership(agent TGAgentReference, members map[string]memberStatus, childPredicate compilergraph.Predicate, memberResolutionKind MemberResolutionKind) {
	// Copy over any composed members, skipping those that are already defined (shadowing)
	// or inaccessible. `globallyValidate` has already ensured the agent type refers to a valid agent.
	agentTypeRef := agent.AgentType()
	agentType := agentTypeRef.ReferredType()

	composingTypeModulePath := compilercommon.InputSource(p.composingType.ParentModule().Path())

	for _, member := range agentType.Members() {
		// If the member is already defined, skip.
		if _, exists := members[member.Name()]; exists {
			continue
		}

		// If the member is inaccessible from the composing type, skip.
		_, err := agentTypeRef.ResolveAccessibleMember(member.Name(), composingTypeModulePath, memberResolutionKind)
		if err != nil {
			continue
		}

		// Mark the member as added.
		members[member.Name()] = memberStatus{member.Name(), member.Title(), false}

		// Create a new node of the same kind and copy over any predicates except the type.
		clonedMemberNode := member.GraphNode.CloneExcept(p.modifier, NodePredicateMemberType)
		clonedMemberNode.Connect(NodePredicateMemberBaseMember, member.GraphNode)
		clonedMemberNode.DecorateWithTagged(NodePredicateMemberBaseSource, agentTypeRef)

		p.modifier.Modify(p.composingType.GraphNode).Connect(childPredicate, clonedMemberNode)

		// If the node is an operator, nothing more to do.
		if clonedMemberNode.Kind == NodeTypeOperator {
			continue
		}

		// Decorate the cloned member with the member type of the base member, transformed under
		// any generics found.
		memberType := member.GraphNode.
			GetTagged(NodePredicateMemberType, p.tdg.AnyTypeReference()).(TypeReference).
			TransformUnder(agentTypeRef)

		clonedMemberNode.DecorateWithTagged(NodePredicateMemberType, memberType)
	}
}
