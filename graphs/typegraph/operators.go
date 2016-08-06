// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

const ASSIGNABLE_OP_VALUE = "value"

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
	IsStatic      bool                // Whether the operator is static.
	IsAssignable  bool                // Whether the operator is assignable.
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

	staticTypeGetter := func(staticType TGTypeDecl) typerefGetter {
		return func(containingType TypeReference) TypeReference {
			return t.NewTypeReference(staticType)
		}
	}

	staticNullableTypeGetter := func(staticType TGTypeDecl) typerefGetter {
		return func(containingType TypeReference) TypeReference {
			return t.NewTypeReference(staticType).AsNullable()
		}
	}

	anyTypeGetter := func(containingType TypeReference) TypeReference {
		return t.AnyTypeReference()
	}

	voidTypeGetter := func(containingType TypeReference) TypeReference {
		return t.VoidTypeReference()
	}

	unaryParameters := []operatorParameter{
		operatorParameter{"value", containingTypeGetter},
	}

	binaryParameters := []operatorParameter{
		operatorParameter{"left", containingTypeGetter},
		operatorParameter{"right", containingTypeGetter},
	}

	operators := []operatorDefinition{
		// Binary operators: +, -, *, /, %
		operatorDefinition{"plus", true, false, containingTypeGetter, binaryParameters},
		operatorDefinition{"minus", true, false, containingTypeGetter, binaryParameters},
		operatorDefinition{"times", true, false, containingTypeGetter, binaryParameters},
		operatorDefinition{"div", true, false, containingTypeGetter, binaryParameters},
		operatorDefinition{"mod", true, false, containingTypeGetter, binaryParameters},

		// Bitwise operators: ^, |, &, <<, >>, ~
		operatorDefinition{"xor", true, false, containingTypeGetter, binaryParameters},
		operatorDefinition{"or", true, false, containingTypeGetter, binaryParameters},
		operatorDefinition{"and", true, false, containingTypeGetter, binaryParameters},
		operatorDefinition{"leftshift", true, false, containingTypeGetter, binaryParameters},
		operatorDefinition{"rightshift", true, false, containingTypeGetter, binaryParameters},

		operatorDefinition{"not", true, false, containingTypeGetter, unaryParameters},

		operatorDefinition{"bool", true, false, staticTypeGetter(t.BoolType()), unaryParameters},

		// Equality.
		operatorDefinition{"equals", true, false, staticTypeGetter(t.BoolType()), binaryParameters},

		// Comparison.s
		operatorDefinition{"compare", true, false, staticTypeGetter(t.IntType()), binaryParameters},

		// Range.
		operatorDefinition{"range", true, false, streamContainingTypeGetter, binaryParameters},

		// Contains.
		operatorDefinition{"contains", false, false, staticTypeGetter(t.BoolType()), []operatorParameter{
			operatorParameter{"item", anyTypeGetter},
		}},

		// Slice.
		operatorDefinition{"slice", false, false, anyTypeGetter, []operatorParameter{
			operatorParameter{"startindex", staticNullableTypeGetter(t.IntType())},
			operatorParameter{"endindex", staticNullableTypeGetter(t.IntType())},
		}},

		// Index.
		operatorDefinition{"index", false, false, anyTypeGetter, []operatorParameter{
			operatorParameter{"index", anyTypeGetter},
		}},

		// SetIndex.
		operatorDefinition{"setindex", false, true, voidTypeGetter, []operatorParameter{
			operatorParameter{"index", anyTypeGetter},
			operatorParameter{ASSIGNABLE_OP_VALUE, anyTypeGetter},
		}},
	}

	for _, operator := range operators {
		t.operators[operator.Name] = operator
	}
}
