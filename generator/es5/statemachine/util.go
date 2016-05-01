// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/parser"
)

// getSourceMappingComment returns the source mapping comment for the CodeDOM node, if any.
func getSourceMappingComment(dom codedom.StatementOrExpression) string {
	basisNode := dom.BasisNode()
	inputSource, hasInputSource := basisNode.TryGet(parser.NodePredicateSource)
	if !hasInputSource {
		return ""
	}

	startRune := basisNode.Get(parser.NodePredicateStartRune)
	endRune := basisNode.Get(parser.NodePredicateEndRune)
	var name = ""
	if named, ok := dom.(codedom.Named); ok {
		name = named.ExprName()
	}

	return "/*@{" + inputSource + "," + startRune + "," + endRune + "," + name + "}*/"
}
