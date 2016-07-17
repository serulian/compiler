// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"fmt"

	"github.com/serulian/compiler/compilerutil"
)

// clientQueryFilter represents a single filter applied on the client side to a set of nodes.
type clientQueryFilter struct {
	operation clientQueryOperation // The operation to perform.
	predicate Predicate            // The predicate to lookup.
	value     string               // The value to compare.
}

// clientQueryOperation defines the set of supported client-side operations.
type clientQueryOperation int

const (
	WhereLTE clientQueryOperation = iota // Floating point: Less than or equals
	WhereGTE                             // Floating point: Greater than or equals
	WhereLT                              // Floating point: Less than
	WhereGT                              // Floating point: Greater than
)

// apply applies this filter to the given node, returning true if the node meets the criteria and false
// otherwise. Note: All predicates needed by this filter must be in the values map.
func (cqf clientQueryFilter) apply(it NodeIterator) bool {
	strValue := valueToOriginalString(it.getRequestedPredicate(cqf.predicate))

	switch cqf.operation {
	case WhereGTE:
		return compilerutil.ParseFloat(strValue) >= compilerutil.ParseFloat(cqf.value)

	case WhereLTE:
		return compilerutil.ParseFloat(strValue) <= compilerutil.ParseFloat(cqf.value)

	case WhereGT:
		return compilerutil.ParseFloat(strValue) > compilerutil.ParseFloat(cqf.value)

	case WhereLT:
		return compilerutil.ParseFloat(strValue) < compilerutil.ParseFloat(cqf.value)

	default:
		panic(fmt.Sprintf("Unknown client filter operation kind: %v", cqf.operation))
	}
}
