// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
)

// srgSourceHandler implements the SourceHandler interface from the packageloader for
// populating the SRG from source files and packages.
type srgSourceHandler struct {
	srg      *SRG                             // The SRG being populated.
	modifier compilergraph.GraphLayerModifier // Modifier used to write the parsed AST.
}

func (sh *srgSourceHandler) Kind() string {
	return srgSourceKind
}

func (sh *srgSourceHandler) PackageFileExtension() string {
	return ".seru"
}

// buildASTNode constructs a new node in the SRG.
func (h *srgSourceHandler) buildASTNode(source compilercommon.InputSource, kind parser.NodeType) parser.AstNode {
	graphNode := h.modifier.CreateNode(kind)
	return &srgASTNode{
		graphNode: graphNode,
	}
}

func (sh *srgSourceHandler) Parse(source compilercommon.InputSource, input string, importHandler packageloader.ImportHandler) {
	parser.Parse(sh.buildASTNode, importHandler, source, input)
}

func (sh *srgSourceHandler) Apply(packageMap packageloader.LoadedPackageMap, sourceTracker packageloader.SourceTracker) {
	// Save the package map and source tracker for later resolution.
	sh.srg.packageMap = packageMap
	sh.srg.sourceTracker = sourceTracker

	// Apply the changes to the graph.
	sh.modifier.Apply()
}

func (sh *srgSourceHandler) Verify(errorReporter packageloader.ErrorReporter, warningReporter packageloader.WarningReporter) {
	g := sh.srg

	// Collect any parse errors found and add them to the result.
	eit := g.findAllNodes(parser.NodeTypeError).BuildNodeIterator(
		parser.NodePredicateErrorMessage)

	for eit.Next() {
		sourceRange, hasRange := sh.srg.SourceRangeOf(eit.Node())
		if !hasRange {
			panic("Missing source range")
		}

		errorReporter(compilercommon.NewSourceError(sourceRange, eit.GetPredicate(parser.NodePredicateErrorMessage).String()))
	}

	// Verify all 'from ... import ...' are valid.
	fit := g.findAllNodes(parser.NodeTypeImportPackage).
		Has(parser.NodeImportPredicateSubsource).
		BuildNodeIterator(parser.NodeImportPredicateSubsource)

	for fit.Next() {
		// Load the package information.
		packageInfo, err := g.getPackageForImport(fit.Node())
		if err != nil || !packageInfo.IsSRGPackage() {
			continue
		}

		// Search for the subsource.
		subsource := fit.GetPredicate(parser.NodeImportPredicateSubsource).String()
		_, found := packageInfo.FindTypeOrMemberByName(subsource)
		if !found {
			source := fit.Node().GetIncomingNode(parser.NodeImportPredicatePackageRef).Get(parser.NodeImportPredicateSource)
			sourceRange, hasRange := sh.srg.SourceRangeOf(fit.Node())
			if !hasRange {
				panic("Missing source range")
			}

			errorReporter(compilercommon.SourceErrorf(sourceRange, "Import '%s' not found under package '%s'", subsource, source))
		}
	}

	// Build the map for globally aliased types.
	ait := g.findAllNodes(parser.NodeTypeDecorator).
		Has(parser.NodeDecoratorPredicateInternal, aliasInternalDecoratorName).
		BuildNodeIterator()

	for ait.Next() {
		// Find the name of the alias.
		decorator := ait.Node()
		parameter, ok := decorator.TryGetNode(parser.NodeDecoratorPredicateParameter)
		if !ok || parameter.Kind() != parser.NodeStringLiteralExpression {
			sourceRange, hasRange := sh.srg.SourceRangeOf(decorator)
			if !hasRange {
				panic("Missing source range")
			}

			errorReporter(compilercommon.SourceErrorf(sourceRange, "Alias decorator requires a single string literal parameter"))
			continue
		}

		var aliasName = parameter.Get(parser.NodeStringLiteralExpressionValue)
		aliasName = aliasName[1 : len(aliasName)-1] // Remove the quotes.

		aliasedType := SRGType{decorator.GetIncomingNode(parser.NodeTypeDefinitionDecorator), g}
		g.aliasMap[aliasName] = aliasedType
	}
}
