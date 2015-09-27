// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"strconv"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/parser"
)

// salForPredicates returns a SourceAndLocation for the loaded predicate map. Note that
// the map *must* contain the NodePredicateSource and NodePredicateStartRune predicates.
func salForPredicates(predicateValues map[string]string) compilercommon.SourceAndLocation {
	source := compilercommon.InputSource(predicateValues[parser.NodePredicateSource])
	bytePosition, err := strconv.Atoi(predicateValues[parser.NodePredicateStartRune])
	if err != nil {
		panic(fmt.Sprintf("Expected int value for byte position, found: %v", bytePosition))
	}

	return compilercommon.NewSourceAndLocation(source, bytePosition)
}
