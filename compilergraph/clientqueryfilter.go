// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"fmt"
)

// clientQueryFilter represents a single filter applied on the client side to a set of nodes.
type clientQueryFilter struct {
	operation clientQueryOperation // The operation to perform.
	predicate Predicate            // The predicate to lookup.
	value     interface{}          // The value to compare.
}

// clientQueryOperation defines the set of supported client-side operations.
type clientQueryOperation int

const (
	WhereLTE clientQueryOperation = iota // Integer: Less than or equals
	WhereGTE                             // Integer: Greater than or equals
	WhereLT                              // Integer: Less than
	WhereGT                              // Integer: Greater than
)

// apply applies this filter to the given node, returning true if the node meets the criteria and false
// otherwise. Note: All predicates needed by this filter must be in the values map.
func (cqf clientQueryFilter) apply(it NodeIterator) bool {
	graphValue := it.GetPredicate(cqf.predicate)

	switch cqf.operation {
	case WhereGTE:
		return graphValue.Int() >= cqf.value.(int)

	case WhereLTE:
		return graphValue.Int() <= cqf.value.(int)

	case WhereGT:
		return graphValue.Int() > cqf.value.(int)

	case WhereLT:
		return graphValue.Int() < cqf.value.(int)

	default:
		panic(fmt.Sprintf("Unknown client filter operation kind: %v", cqf.operation))
	}
}
