// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeconstructor

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/webidl"
)

// GetConstructor returns a TypeGraph constructor for the given IRG.
func GetConstructor(irg *webidl.WebIRG) *irgTypeConstructor {
	return &irgTypeConstructor{
		irg: irg,
	}
}

// irgTypeConstructor defines a type for populating a type graph from the IRG.
type irgTypeConstructor struct {
	irg *webidl.WebIRG // The IRG being transformed.
}

func (itc *irgTypeConstructor) DefineModules(builder typegraph.GetModuleBuilder) {
	for _, module := range itc.irg.GetModules() {
		builder().
			Name(module.Name()).
			Path(string(module.InputSource())).
			SourceNode(module.Node()).
			Define()
	}
}

func (itc *irgTypeConstructor) DefineTypes(builder typegraph.GetTypeBuilder) {
	for _, module := range itc.irg.GetModules() {
		for _, declaration := range module.Declarations() {
			typeBuilder := builder(module.Node())
			typeBuilder.Name(declaration.Name()).
				SourceNode(declaration.GraphNode).
				TypeKind(typegraph.ExternalInternalType).
				Define()
		}
	}
}

func (itc *irgTypeConstructor) DefineDependencies(annotator *typegraph.Annotator, graph *typegraph.TypeGraph) {

}

func (itc *irgTypeConstructor) DefineMembers(builder typegraph.GetMemberBuilder, graph *typegraph.TypeGraph) {
	for _, declaration := range itc.irg.Declarations() {
		// TODO: define constructors.

		for _, member := range declaration.Members() {
			ibuilder, reporter := builder(declaration.GraphNode, false).
				Name(member.Name()).
				SourceNode(member.GraphNode).
				InitialDefine()

			declaredType, err := itc.ResolveType(member.DeclaredType(), graph)
			if err != nil {
				reporter.ReportError(member.GraphNode, "%v", err)
				continue
			}

			var memberType = declaredType
			var isReadonly = member.IsReadonly()

			switch member.Kind() {
			case webidl.FunctionMember:
				isReadonly = true
				memberType = graph.FunctionTypeReference(memberType)

				// Add the parameter types.
				for _, parameter := range member.Parameters() {
					parameterType, err := itc.ResolveType(parameter.DeclaredType(), graph)
					if err != nil {
						reporter.ReportError(member.GraphNode, "%v", err)
						continue
					}

					memberType = memberType.WithParameter(parameterType)
				}

			case webidl.AttributeMember:
				if len(member.Parameters()) > 0 {
					reporter.ReportError(member.GraphNode, "Attributes cannot have parameters")
				}

			default:
				panic("Unknown WebIDL member kind")
			}

			ibuilder.Exported(true).
				Static(member.IsStatic()).
				Synchronous(true).
				ReadOnly(isReadonly).
				MemberKind(uint64(member.Kind())).
				MemberType(memberType).
				Define()
		}
	}
}

func (itc *irgTypeConstructor) Validate(reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	seen := map[string]bool{}

	for _, module := range itc.irg.GetModules() {
		for _, declaration := range module.Declarations() {
			if _, ok := seen[declaration.Name()]; ok {
				reporter.ReportError(declaration.GraphNode, "'%s' is already declared in WebIDL", declaration.Name())
			}
			seen[declaration.Name()] = true
		}
	}
}

func (itc *irgTypeConstructor) GetLocation(sourceNodeId compilergraph.GraphNodeId) (compilercommon.SourceAndLocation, bool) {
	layerNode, found := itc.irg.TryGetNode(sourceNodeId)
	if !found {
		return compilercommon.SourceAndLocation{}, false
	}

	return itc.irg.NodeLocation(layerNode), true
}

// ResolveType attempts to resolve the given type string.
func (itc *irgTypeConstructor) ResolveType(typeString string, graph *typegraph.TypeGraph) (typegraph.TypeReference, error) {
	if typeString == "any" {
		return graph.AnyTypeReference(), nil
	}

	if typeString == "void" {
		return graph.VoidTypeReference(), nil
	}

	declaration, hasDeclaration := itc.irg.FindDeclaration(typeString)
	if !hasDeclaration {
		return graph.AnyTypeReference(), fmt.Errorf("Could not find WebIDL type %v", typeString)
	}

	typeDecl, hasType := graph.GetTypeForSourceNode(declaration.GraphNode)
	if !hasType {
		panic("Type not found for WebIDL type declaration")
	}

	return typeDecl.GetTypeReference(), nil
}
