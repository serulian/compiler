// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webidl

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/webidl/parser"
)

// irgSourceHandler implements the SourceHandler interface from the packageloader for
// populating the WebIDL IRG from webidl files.
type irgSourceHandler struct {
	irg      *WebIRG                          // The IRG being populated.
	modifier compilergraph.GraphLayerModifier // Modifier used to write the parsed AST.
}

func (sh *irgSourceHandler) Kind() string {
	return "webidl"
}

func (sh *irgSourceHandler) PackageFileExtension() string {
	return ".webidl"
}

func (sh *irgSourceHandler) Parse(source compilercommon.InputSource, input string, importHandler packageloader.ImportHandler) {
	rootNode := sh.modifier.Modify(sh.irg.rootModuleNode)
	parser.Parse(&irgASTNode{rootNode}, sh.buildASTNode, source, input)
}

func (sh *irgSourceHandler) Apply(packageMap map[string]packageloader.PackageInfo) {
	// Apply the changes to the graph.
	sh.modifier.Apply()
}

func (sh *irgSourceHandler) Verify(errorReporter packageloader.ErrorReporter, warningReporter packageloader.WarningReporter) {

}

// buildASTNode constructs a new node in the IRG.
func (sh *irgSourceHandler) buildASTNode(source compilercommon.InputSource, kind parser.NodeType) parser.AstNode {
	graphNode := sh.modifier.CreateNode(kind)
	return &irgASTNode{
		graphNode: graphNode,
	}
}
