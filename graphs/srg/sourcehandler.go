// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
	"github.com/serulian/compiler/parser/shared"
	"github.com/serulian/compiler/sourceshape"
)

// srgSourceHandler implements the SourceHandler interface from the packageloader for
// populating the SRG from source files and packages.
type srgSourceHandler struct {
	srg *SRG // The SRG being populated.
}

func (sh srgSourceHandler) Kind() string {
	return srgSourceKind
}

func (sh srgSourceHandler) PackageFileExtension() string {
	return ".seru"
}

func (sh srgSourceHandler) NewParser() packageloader.SourceHandlerParser {
	return srgSourceHandlerParser{sh.srg, sh.srg.layer.NewModifier()}
}

type srgSourceHandlerParser struct {
	srg      *SRG // The SRG being populated.
	modifier compilergraph.GraphLayerModifier
}

// buildASTNode constructs a new node in the SRG.
func (sh srgSourceHandlerParser) buildASTNode(source compilercommon.InputSource, kind sourceshape.NodeType) shared.AstNode {
	graphNode := sh.modifier.CreateNode(kind)
	return &srgASTNode{
		graphNode: graphNode,
	}
}

func (sh srgSourceHandlerParser) Parse(source compilercommon.InputSource, input string, importHandler packageloader.ImportHandler) {
	parser.Parse(sh.buildASTNode, importHandler, source, input)
}

func (sh srgSourceHandlerParser) Apply(packageMap packageloader.LoadedPackageMap, sourceTracker packageloader.SourceTracker) {
	// Save the package map and source tracker for later resolution.
	sh.srg.packageMap = packageMap
	sh.srg.sourceTracker = sourceTracker

	// Apply the changes to the graph.
	sh.modifier.Apply()
}

func (sh srgSourceHandlerParser) Verify(errorReporter packageloader.ErrorReporter, warningReporter packageloader.WarningReporter) {
	g := sh.srg

	// Collect any parse errors found and add them to the result.
	eit := g.findAllNodes(sourceshape.NodeTypeError).BuildNodeIterator(
		sourceshape.NodePredicateErrorMessage)

	for eit.Next() {
		sourceRange, hasRange := sh.srg.SourceRangeOf(eit.Node())
		if !hasRange {
			panic("Missing source range")
		}

		errorReporter(compilercommon.NewSourceError(sourceRange, eit.GetPredicate(sourceshape.NodePredicateErrorMessage).String()))
	}

	// Verify all 'from ... import ...' are valid.
	fit := g.findAllNodes(sourceshape.NodeTypeImportPackage).
		Has(sourceshape.NodeImportPredicateSubsource).
		BuildNodeIterator(sourceshape.NodeImportPredicateSubsource)

	for fit.Next() {
		// Load the package information.
		packageInfo, err := g.getPackageForImport(fit.Node())
		if err != nil || !packageInfo.IsSRGPackage() {
			continue
		}

		// Search for the subsource.
		subsource := fit.GetPredicate(sourceshape.NodeImportPredicateSubsource).String()
		_, found := packageInfo.FindTypeOrMemberByName(subsource)
		if !found {
			source := fit.Node().GetIncomingNode(sourceshape.NodeImportPredicatePackageRef).Get(sourceshape.NodeImportPredicateSource)
			sourceRange, hasRange := sh.srg.SourceRangeOf(fit.Node())
			if !hasRange {
				panic("Missing source range")
			}

			errorReporter(compilercommon.SourceErrorf(sourceRange, "Import '%s' not found under package '%s'", subsource, source))
		}
	}

	// Build the map for globally aliased types.
	ait := g.findAllNodes(sourceshape.NodeTypeDecorator).
		Has(sourceshape.NodeDecoratorPredicateInternal, aliasInternalDecoratorName).
		BuildNodeIterator()

	for ait.Next() {
		// Find the name of the alias.
		decorator := ait.Node()
		parameter, ok := decorator.TryGetNode(sourceshape.NodeDecoratorPredicateParameter)
		if !ok || parameter.Kind() != sourceshape.NodeStringLiteralExpression {
			sourceRange, hasRange := sh.srg.SourceRangeOf(decorator)
			if !hasRange {
				panic("Missing source range")
			}

			errorReporter(compilercommon.SourceErrorf(sourceRange, "Alias decorator requires a single string literal parameter"))
			continue
		}

		var aliasName = parameter.Get(sourceshape.NodeStringLiteralExpressionValue)
		aliasName = aliasName[1 : len(aliasName)-1] // Remove the quotes.

		aliasedType := SRGType{decorator.GetIncomingNode(sourceshape.NodeTypeDefinitionDecorator), g}
		g.aliasMap[aliasName] = aliasedType
	}
}
