// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser/shared"
	"github.com/serulian/compiler/sourceshape"
)

type noopNode struct{}

func (nn noopNode) Connect(predicate string, other shared.AstNode) shared.AstNode { return nn }
func (nn noopNode) Decorate(property string, value string) shared.AstNode         { return nn }
func (nn noopNode) DecorateWithInt(property string, value int) shared.AstNode     { return nn }

func noopBuilder(source compilercommon.InputSource, kind sourceshape.NodeType) shared.AstNode {
	return noopNode{}
}

func noopImportHandler(sourceKind string, importPath string, importType packageloader.PackageImportType, importSource compilercommon.InputSource, runePosition int) string {
	return ""
}
