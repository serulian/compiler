// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
)

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

// ReturningTypeOf marks the scope as returning the type of the given scope.
func (sib *scopeInfoBuilder) ReturningTypeOf(scope *proto.ScopeInfo) *scopeInfoBuilder {
	returnedValue := scope.GetReturnedType()
	sib.info.ReturnedType = &returnedValue
	return sib
}

// Returning marks the scope as returning a value of the given type.
func (sib *scopeInfoBuilder) Returning(returning typegraph.TypeReference) *scopeInfoBuilder {
	returnedValue := returning.Value()
	sib.info.ReturnedType = &returnedValue
	return sib
}

// IsTerminatingStatement marks the scope as containing a terminating statement.
func (sib *scopeInfoBuilder) IsTerminatingStatement() *scopeInfoBuilder {
	trueValue := true
	sib.info.IsTerminatingStatement = &trueValue
	return sib
}

// GetScope returns the scope constructed.
func (sib *scopeInfoBuilder) GetScope() proto.ScopeInfo {
	return *sib.info
}
