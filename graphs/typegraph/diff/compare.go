// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"github.com/serulian/compiler/graphs/typegraph"
)

func compareTypes(original typegraph.TypeReference, updated typegraph.TypeReference) bool {
	// Fast path check: If both types are exactly the same, nothing more to do.
	if original.PackageQualifiedString() == updated.PackageQualifiedString() {
		return true
	}

	// Otherwise, compare the types using their shape. Only types with the exact shape are
	// considered the same.
	if !original.IsNormal() || !updated.IsNormal() {
		return false
	}

	// TODO(jschorr): This.
	return false
}

func compareParameters(original []typegraph.TGParameter, updated []typegraph.TGParameter) MemberDiffReason {
	if len(updated) < len(original) {
		return MemberDiffReasonParametersNotCompatible
	}

	// For each of the original parameters, ensure that a updated parameter exists and has an assignable type.
	for index := range original {
		originalType := original[index].DeclaredType()
		updatedType := updated[index].DeclaredType()
		if !compareTypes(originalType, updatedType) {
			return MemberDiffReasonParametersNotCompatible
		}
	}

	// Ensure each of the additional updated parameters (if any) are nullable, making them optional.
	for _, updatedParam := range updated[len(original):] {
		if !updatedParam.DeclaredType().NullValueAllowed() {
			return MemberDiffReasonParametersNotCompatible
		}
	}

	if len(updated) > len(original) {
		return MemberDiffReasonParametersCompatible
	}

	return MemberDiffReasonNotApplicable
}

func compareGenerics(original []typegraph.TGGeneric, updated []typegraph.TGGeneric) bool {
	if len(original) != len(updated) {
		return false
	}

	for index, originalGeneric := range original {
		if originalGeneric.Name() != updated[index].Name() {
			return false
		}

		if originalGeneric.Constraint().String() != updated[index].Constraint().String() {
			return false
		}
	}

	return true
}
