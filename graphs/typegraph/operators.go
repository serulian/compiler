// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

// operatorParameter represents a single expected parameter on an operator.
type operatorParameter struct {
	Name             string        // The name of the parameter.
	isContainingType bool          // If true, the expected type of the parameter is the type containing type.
	predefinedType   TypeReference // The expected type (if not dependent on the defining type.)
}

// ExpectedType returns the type expected for this parameter.
func (op *operatorParameter) ExpectedType(containingType TypeReference) TypeReference {
	if op.isContainingType {
		return containingType
	}

	return op.predefinedType
}

// operatorDefinition represents the definition of a supported operator on a Serulian type.
type operatorDefinition struct {
	Name       string              // The name of the operator.
	Parameters []operatorParameter // The expected parameters.
}

// buildOperatorDefinitions sets the defined operators supported in the type system.
func (t *TypeGraph) buildOperatorDefinitions() {
	binaryParameters := []operatorParameter{
		operatorParameter{"left", true, TypeReference{}},
		operatorParameter{"right", true, TypeReference{}},
	}

	operators := []operatorDefinition{
		// Binary operators: +, -, *, /, %
		operatorDefinition{"plus", binaryParameters},
		operatorDefinition{"minus", binaryParameters},
		operatorDefinition{"times", binaryParameters},
		operatorDefinition{"div", binaryParameters},
		operatorDefinition{"mod", binaryParameters},

		// Bitwise operators: ^, |, &, <<, >>, ~
		operatorDefinition{"xor", binaryParameters},
		operatorDefinition{"or", binaryParameters},
		operatorDefinition{"and", binaryParameters},
		operatorDefinition{"leftshift", binaryParameters},
		operatorDefinition{"rightshift", binaryParameters},
		operatorDefinition{"not", binaryParameters},

		// Equality.
		operatorDefinition{"equals", binaryParameters},

		// Comparison.
		operatorDefinition{"compare", binaryParameters},

		// Range.
		operatorDefinition{"range", binaryParameters},

		// Slice.
		operatorDefinition{"slice", []operatorParameter{
			operatorParameter{"startindex", false, t.NewTypeReference(t.IntType())},
			operatorParameter{"endindex", false, t.NewTypeReference(t.IntType())},
		}},

		// Index.
		operatorDefinition{"index", []operatorParameter{
			operatorParameter{"index", false, t.AnyTypeReference()},
		}},
	}

	for _, operator := range operators {
		t.operators[operator.Name] = operator
	}
}
