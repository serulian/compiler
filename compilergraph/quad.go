// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"fmt"

	"github.com/cayleygraph/cayley/quad"
)

// Cayley type mappings:
//   GraphNodeId <-> quad.Raw
//   Predicate <-> quad.IRI
//   TaggedValue <-> quad.Raw
//   Other values <-> quad.Value

// nodeIdToValue returns a Cayley value for a Graph Node ID.
func nodeIdToValue(nodeId GraphNodeId) quad.Value {
	return quad.Raw(nodeId)
}

// valueToNodeId returns a Graph Node ID for a Cayley value.
func valueToNodeId(value quad.Value) GraphNodeId {
	return GraphNodeId(value.(quad.Raw))
}

// predicateToValue converts a Predicate to a Cayley value.
func predicateToValue(predicate Predicate) quad.Value {
	return quad.IRI(string(predicate))
}

// valueToPredicateString returns the string form of the predicate for a Cayley value.
func valueToPredicateString(predicateValue quad.Value) string {
	return string(valueToPredicate(predicateValue))
}

// valueToPredicate returns a Predicate for a Cayley value.
func valueToPredicate(predicateValue quad.Value) Predicate {
	return Predicate(iriToString(predicateValue))
}

// iriToString returns the string form of an IRI.
func iriToString(iri quad.Value) string {
	return string(iri.(quad.IRI))
}

// valueToTaggedValueData returns the internal string presentation of a
// tagged value for a Cayley value.
func valueToTaggedValueData(value quad.Value) string {
	return string(value.(quad.Raw))
}

// taggedValueDataToValue returns a Cayley value for the given internal
// tagged value representation.
func taggedValueDataToValue(taggedValue string) quad.Value {
	return quad.Raw(taggedValue)
}

// valueToOriginalString converts the Cayley value back into its original string value.
func valueToOriginalString(value quad.Value) string {
	return string(value.(quad.String))
}

// buildGraphValueForValue returns the GraphValue for a Cayley value.
func buildGraphValueForValue(value quad.Value) GraphValue {
	return GraphValue{value}
}

// taggedToQuadValues converts a slice of TaggedValue's under a layer into their
// Cayley values.
func taggedToQuadValues(values []TaggedValue, gl *GraphLayer) []quad.Value {
	quadValues := make([]quad.Value, len(values))
	for index, v := range values {
		quadValues[index] = gl.getTaggedKey(v)
	}
	return quadValues
}

// graphIdsToQuadValues converts a slice of Graph Node IDs into their Cayley values.
func graphIdsToQuadValues(values []GraphNodeId) []quad.Value {
	quadValues := make([]quad.Value, len(values))
	for index, v := range values {
		quadValues[index] = nodeIdToValue(v)
	}
	return quadValues
}

// toQuadValues converts a slice of arbitrary values into their Cayley values.
func toQuadValues(values []interface{}, gl *GraphLayer) []quad.Value {
	quadValues := make([]quad.Value, len(values))
	for index, v := range values {
		switch v := v.(type) {
		case GraphNodeId:
			quadValues[index] = nodeIdToValue(v)

		case TaggedValue:
			quadValues[index] = gl.getTaggedKey(v)

		default:
			vle, ok := quad.AsValue(v)
			if !ok {
				panic(fmt.Sprintf("Unknown value %v", v))
			}

			quadValues[index] = vle
		}

	}

	return quadValues
}
