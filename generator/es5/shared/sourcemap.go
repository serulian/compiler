// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shared

import (
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/sourcemap"
)

// SourceMapWrapExpr wraps a builder expression with its source mapping location.
func SourceMapWrapExpr(builder esbuilder.ExpressionBuilder, expression codedom.Expression, srg *srg.SRG) esbuilder.ExpressionBuilder {
	mapping, hasMapping := getMapping(expression, srg)
	if !hasMapping {
		return builder
	}

	return builder.WithMappingAsExpr(mapping)
}

// SourceMapWrap wraps a builder statement or expression with its source mapping location.
func SourceMapWrap(builder esbuilder.SourceBuilder, dom codedom.StatementOrExpression, srg *srg.SRG) esbuilder.SourceBuilder {
	mapping, hasMapping := getMapping(dom, srg)
	if !hasMapping {
		return builder
	}

	return builder.WithMapping(mapping)
}

func getMapping(dom codedom.StatementOrExpression, srg *srg.SRG) (sourcemap.SourceMapping, bool) {
	sourceRange, hasSourceRange := srg.SourceRangeOf(dom.BasisNode())
	if !hasSourceRange {
		return sourcemap.SourceMapping{}, false
	}

	originalLine, originalCol, err := sourceRange.Start().LineAndColumn()
	if err != nil {
		panic(err)
	}

	var name = ""
	if named, ok := dom.(codedom.Named); ok {
		name = named.ExprName()
	}

	return sourcemap.SourceMapping{string(sourceRange.Source()), originalLine, originalCol, name}, true
}
