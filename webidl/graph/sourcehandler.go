// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

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

func (sh *irgSourceHandler) Apply(packageMap packageloader.LoadedPackageMap, sourceTracker packageloader.SourceTracker) {
	// Apply the changes to the graph.
	sh.modifier.Apply()

	// Make sure we didn't encounter any errors.
	if sh.irg.findAllNodes(parser.NodeTypeError).BuildNodeIterator().Next() {
		return
	}

	// Perform type collapsing.
	modifier := sh.irg.layer.NewModifier()
	defer modifier.Apply()

	sh.irg.sourceTracker = sourceTracker
	sh.irg.typeCollapser = createTypeCollapser(sh.irg, modifier)
}

func (sh *irgSourceHandler) Verify(errorReporter packageloader.ErrorReporter, warningReporter packageloader.WarningReporter) {
	g := sh.irg

	// Collect any parse errors found and add them to the result.
	eit := g.findAllNodes(parser.NodeTypeError).BuildNodeIterator(
		parser.NodePredicateErrorMessage)

	for eit.Next() {
		sourceRange, hasSourceRange := sh.irg.SourceRangeOf(eit.Node())
		if !hasSourceRange {
			continue
		}

		errorReporter(compilercommon.NewSourceError(sourceRange, eit.GetPredicate(parser.NodePredicateErrorMessage).String()))
	}
}

// buildASTNode constructs a new node in the IRG.
func (sh *irgSourceHandler) buildASTNode(source compilercommon.InputSource, kind parser.NodeType) parser.AstNode {
	graphNode := sh.modifier.CreateNode(kind)
	return &irgASTNode{
		graphNode: graphNode,
	}
}
