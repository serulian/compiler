// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
)

// newEmptyScopeInfo returns a ScopeInfo block representing empty scope, with some validation.
func newEmptyScopeInfo(valid bool) proto.ScopeInfo {
	return proto.ScopeInfo{
		IsValid: &valid,
	}
}

// newReturningScopeInfo returns a ScopeInfo block representing scope that returns an instance of
// some type.
func newReturningScopeInfo(valid bool, returning typegraph.TypeReference) proto.ScopeInfo {
	returnedValue := returning.Value()

	return proto.ScopeInfo{
		IsValid:      &valid,
		ReturnedType: &returnedValue,
	}
}
