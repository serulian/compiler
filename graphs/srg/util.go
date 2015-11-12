// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"strconv"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// salForPredicates returns a SourceAndLocation for the loaded predicate map. Note that
// the map *must* contain the NodePredicateSource and NodePredicateStartRune predicates.
func salForPredicates(predicateValues map[string]string) compilercommon.SourceAndLocation {
	return salForValues(predicateValues[parser.NodePredicateSource], predicateValues[parser.NodePredicateStartRune])
}

// salForNode returns a SourceAndLocation for the given graph node.
func salForNode(node compilergraph.GraphNode) compilercommon.SourceAndLocation {
	return salForValues(node.Get(parser.NodePredicateSource), node.Get(parser.NodePredicateStartRune))
}

// salForValues returns a SourceAndLocation for the given string predicate values.
func salForValues(sourceStr string, bytePositionStr string) compilercommon.SourceAndLocation {
	source := compilercommon.InputSource(sourceStr)
	bytePosition, err := strconv.Atoi(bytePositionStr)
	if err != nil {
		panic(fmt.Sprintf("Expected int value for byte position, found: %v", bytePositionStr))
	}

	return compilercommon.NewSourceAndLocation(source, bytePosition)
}
