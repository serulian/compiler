// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
)

type typerefGetter func(containingType TypeReference) TypeReference

// operatorParameter represents a single expected parameter on an operator.
type operatorParameter struct {
	Name             string        // The name of the parameter.
	getParameterType typerefGetter // The expected type.
}

// ExpectedType returns the type expected for this parameter.
func (op *operatorParameter) ExpectedType(containingType TypeReference) TypeReference {
	return op.getParameterType(containingType)
}

// operatorDefinition represents the definition of a supported operator on a Serulian type.
type operatorDefinition struct {
	Name          string              // The name of the operator.
	getReturnType typerefGetter       // The expected return type.
	Parameters    []operatorParameter // The expected parameters.
}

// ExpectedReturnType returns the return type expected for this operator.
func (od *operatorDefinition) ExpectedReturnType(containingType TypeReference) TypeReference {
	return od.getReturnType(containingType)
}

// GetMemberType returns the member type for this operator definition.
func (od *operatorDefinition) GetMemberType(containingType TypeReference, declaredReturnType TypeReference) TypeReference {
	// The member type for an operator is a function that takes in the expected parameters
	// and returns the declared return type.
	typegraph := containingType.tdg

	// Add the operator's parameters.
	var funcType = typegraph.NewTypeReference(typegraph.FunctionType()).WithGeneric(declaredReturnType)
	for _, param := range od.Parameters {
		funcType = funcType.WithParameter(param.getParameterType(containingType))
	}

	return funcType
}

// buildOperatorDefinitions sets the defined operators supported in the type system.
func (t *TypeGraph) buildOperatorDefinitions() {
	containingTypeGetter := func(containingType TypeReference) TypeReference {
		return containingType
	}

	streamContainingTypeGetter := func(containingType TypeReference) TypeReference {
		return t.NewTypeReference(t.StreamType(), containingType)
	}

	staticTypeGetter := func(staticType compilergraph.GraphNode) typerefGetter {
		return func(containingType TypeReference) TypeReference {
			return t.NewTypeReference(staticType)
		}
	}

	anyTypeGetter := func(containingType TypeReference) TypeReference {
		return t.AnyTypeReference()
	}

	streamAnyTypeGetter := func(containingType TypeReference) TypeReference {
		return t.NewTypeReference(t.StreamType(), t.AnyTypeReference())
	}

	binaryParameters := []operatorParameter{
		operatorParameter{"left", containingTypeGetter},
		operatorParameter{"right", containingTypeGetter},
	}

	operators := []operatorDefinition{
		// Binary operators: +, -, *, /, %
		operatorDefinition{"plus", containingTypeGetter, binaryParameters},
		operatorDefinition{"minus", containingTypeGetter, binaryParameters},
		operatorDefinition{"times", containingTypeGetter, binaryParameters},
		operatorDefinition{"div", containingTypeGetter, binaryParameters},
		operatorDefinition{"mod", containingTypeGetter, binaryParameters},

		// Bitwise operators: ^, |, &, <<, >>, ~
		operatorDefinition{"xor", containingTypeGetter, binaryParameters},
		operatorDefinition{"or", containingTypeGetter, binaryParameters},
		operatorDefinition{"and", containingTypeGetter, binaryParameters},
		operatorDefinition{"leftshift", containingTypeGetter, binaryParameters},
		operatorDefinition{"rightshift", containingTypeGetter, binaryParameters},
		operatorDefinition{"not", containingTypeGetter, binaryParameters},

		// Equality.
		operatorDefinition{"equals", staticTypeGetter(t.BoolType()), binaryParameters},

		// Comparison.
		operatorDefinition{"compare", staticTypeGetter(t.IntType()), binaryParameters},

		// Range.
		operatorDefinition{"range", streamContainingTypeGetter, binaryParameters},

		// Slice.
		operatorDefinition{"slice", streamAnyTypeGetter, []operatorParameter{
			operatorParameter{"startindex", staticTypeGetter(t.IntType())},
			operatorParameter{"endindex", staticTypeGetter(t.IntType())},
		}},

		// Index.
		operatorDefinition{"index", anyTypeGetter, []operatorParameter{
			operatorParameter{"index", anyTypeGetter},
		}},
	}

	for _, operator := range operators {
		t.operators[operator.Name] = operator
	}
}
