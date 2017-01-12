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

func (sh *irgSourceHandler) Apply(packageMap packageloader.LoadedPackageMap) {
	// Apply the changes to the graph.
	sh.modifier.Apply()
}

func (sh *irgSourceHandler) Verify(errorReporter packageloader.ErrorReporter, warningReporter packageloader.WarningReporter) {
	g := sh.irg

	// Collect any parse errors found and add them to the result.
	eit := g.findAllNodes(parser.NodeTypeError).BuildNodeIterator(
		parser.NodePredicateErrorMessage,
		parser.NodePredicateSource,
		parser.NodePredicateStartRune)

	for eit.Next() {
		sal := salForIterator(eit)
		errorReporter(compilercommon.NewSourceError(sal, eit.GetPredicate(parser.NodePredicateErrorMessage).String()))
	}
}

// salForIterator returns a SourceAndLocation for the given iterator. Note that
// the iterator *must* contain the NodePredicateSource and NodePredicateStartRune predicates.
func salForIterator(iterator compilergraph.NodeIterator) compilercommon.SourceAndLocation {
	return compilercommon.NewSourceAndLocation(
		compilercommon.InputSource(iterator.GetPredicate(parser.NodePredicateSource).String()),
		iterator.GetPredicate(parser.NodePredicateStartRune).Int())
}

// buildASTNode constructs a new node in the IRG.
func (sh *irgSourceHandler) buildASTNode(source compilercommon.InputSource, kind parser.NodeType) parser.AstNode {
	graphNode := sh.modifier.CreateNode(kind)
	return &irgASTNode{
		graphNode: graphNode,
	}
}
