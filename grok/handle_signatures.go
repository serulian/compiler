// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
)

// GetSignatureForPosition returns the parameter signature information for the given activationString at the given position.
// If the activation string does not end in a signature activation character, returns nothing.
func (gh Handle) GetSignatureForPosition(activationString string, source compilercommon.InputSource, lineNumber int, colPosition int) (SignatureInformation, error) {
	// Note: We cannot use PositionFromLineAndColumn here, as it will lookup the position in the *tracked* source, which
	// may be different than the current live source.
	liveRune, err := gh.scopeResult.SourceTracker.LineAndColToRunePosition(lineNumber, colPosition, source, compilercommon.SourceMapCurrent)
	if err != nil {
		return SignatureInformation{}, err
	}

	sourcePosition := source.PositionForRunePosition(liveRune, gh.scopeResult.SourceTracker)
	return gh.GetSignature(activationString, sourcePosition)
}

// GetSignature returns the parameter signature information for the given activationString at the given position.
// If the activation string does not end in a signature activation character, returns nothing.
func (gh Handle) GetSignature(activationString string, sourcePosition compilercommon.SourcePosition) (SignatureInformation, error) {
	updatedPosition, err := gh.scopeResult.SourceTracker.GetPositionOffset(sourcePosition, packageloader.CurrentFilePosition)
	if err != nil {
		return SignatureInformation{}, err
	}

	// Find the node or a nearby node at the position.
	sourceGraph := gh.scopeResult.Graph.SourceGraph()
	node, found := sourceGraph.FindNearbyNodeForPosition(updatedPosition)
	if !found {
		return SignatureInformation{}, nil
	}

	// Make sure we are under an implementable.
	parentImplementable, hasParentImplementable := gh.structureFinder.TryGetContainingImplemented(node)
	if !hasParentImplementable {
		return SignatureInformation{}, nil
	}

	// Extract the expression over which we are activating and the parameter position.
	expressionString, signatureIndex, ok := extractCalled(activationString)
	if !ok {
		return SignatureInformation{}, nil
	}

	// Parse the expression string into an expression.
	source := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
	startRune := node.GetValue(parser.NodePredicateStartRune).Int()
	parsed, ok := srg.ParseExpression(expressionString, source, startRune)
	if !ok {
		return SignatureInformation{}, nil
	}

	// Scope the expression underneath the context of the node found.
	scopeInfo, isValid := gh.scopeResult.Graph.BuildTransientScope(parsed, parentImplementable)
	if !isValid {
		return SignatureInformation{}, nil
	}

	// Find the parameter information for the scoped item, which must either refer to a type member
	// with parameter types *or* must have type function.
	var functionType = gh.scopeResult.Graph.TypeGraph().VoidTypeReference()

	switch scopeInfo.Kind {
	case proto.ScopeKind_VALUE:
		functionType = scopeInfo.ResolvedTypeRef(gh.scopeResult.Graph.TypeGraph())

	case proto.ScopeKind_STATIC:
		return SignatureInformation{}, nil

	case proto.ScopeKind_GENERIC:
		return SignatureInformation{}, nil

	default:
		panic("Unknown scope kind")
	}

	// Check for a named reference. If we find a member, use its information directly. This allows
	// us to get parameter names and documentation, in addition to just the type information.
	// If a value referring to a type, then the access is treated as an indexer.
	referencedName, isNamedScope := gh.scopeResult.Graph.GetReferencedName(scopeInfo)
	if isNamedScope {
		member, isMember := referencedName.Member()

		if !isMember {
			if functionType.IsNormal() {
				member, isMember = functionType.ReferredType().GetOperator("index")
			}
		}

		if isMember {
			name := member.Name()
			documentation, _ := member.Documentation()

			parameters := member.Parameters()
			parameterInfo := make([]ParameterInformation, len(parameters))
			for index, parameter := range parameters {
				parameterName, _ := parameter.Name()
				parameterDocumentation, _ := parameter.Documentation()
				parameterInfo[index] = ParameterInformation{
					Name:          parameterName,
					TypeReference: parameter.DeclaredType(),
					Documentation: trimDocumentation(highlightParameter(parameterDocumentation, parameterName)),
				}
			}

			return SignatureInformation{
				TypeReference:        functionType,
				Name:                 name,
				Documentation:        trimDocumentation(documentation),
				ActiveParameterIndex: signatureIndex,
				Member:               &member,
				Parameters:           parameterInfo,
			}, nil
		}
	}

	// Otherwise, ensure we have a function type.
	if !functionType.IsDirectReferenceTo(gh.scopeResult.Graph.TypeGraph().FunctionType()) {
		// Not a function, so not callable.
		return SignatureInformation{}, nil
	}

	// Generate signature information based on the type information found.
	parameterTypes := functionType.Parameters()
	parameterInfo := make([]ParameterInformation, len(parameterTypes))
	for index, parameterType := range parameterTypes {
		parameterInfo[index] = ParameterInformation{
			Name:          "", // Anonymous
			TypeReference: parameterType,
			Documentation: "",
		}
	}

	return SignatureInformation{
		TypeReference:        functionType,
		Name:                 "", // Anonymous
		Documentation:        "",
		ActiveParameterIndex: signatureIndex,
		Member:               nil,
		Parameters:           parameterInfo,
	}, nil
}
