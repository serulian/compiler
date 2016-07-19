// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shared

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/parser"
	"github.com/serulian/compiler/sourcemap"
)

// SourceMapWrapExpr wraps a builder expression with its source mapping location.
func SourceMapWrapExpr(builder esbuilder.ExpressionBuilder, expression codedom.Expression, positionMapper *compilercommon.PositionMapper) esbuilder.ExpressionBuilder {
	mapping, hasMapping := getMapping(expression, positionMapper)
	if !hasMapping {
		return builder
	}

	return builder.WithMappingAsExpr(mapping)
}

// SourceMapWrap wraps a builder statement or expression with its source mapping location.
func SourceMapWrap(builder esbuilder.SourceBuilder, dom codedom.StatementOrExpression, positionMapper *compilercommon.PositionMapper) esbuilder.SourceBuilder {
	mapping, hasMapping := getMapping(dom, positionMapper)
	if !hasMapping {
		return builder
	}

	return builder.WithMapping(mapping)
}

func getMapping(dom codedom.StatementOrExpression, positionMapper *compilercommon.PositionMapper) (sourcemap.SourceMapping, bool) {
	basisNode := dom.BasisNode()
	inputSource, hasInputSource := basisNode.TryGet(parser.NodePredicateSource)
	if !hasInputSource {
		return sourcemap.SourceMapping{}, false
	}

	startRune := basisNode.GetValue(parser.NodePredicateStartRune).Int()
	originalLine, originalCol, err := positionMapper.Map(compilercommon.InputSource(inputSource), startRune)
	if err != nil {
		panic(err)
	}

	var name = ""
	if named, ok := dom.(codedom.Named); ok {
		name = named.ExprName()
	}

	return sourcemap.SourceMapping{inputSource, originalLine, originalCol, name}, true
}
