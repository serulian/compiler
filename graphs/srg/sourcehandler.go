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

func (sh *srgSourceHandler) Apply(packageMap packageloader.LoadedPackageMap) {
	// Save the package map for later resolution.
	sh.srg.packageMap = packageMap

	// Apply the changes to the graph.
	sh.modifier.Apply()
}

func (sh *srgSourceHandler) Verify(errorReporter packageloader.ErrorReporter, warningReporter packageloader.WarningReporter) {
	g := sh.srg

	// Collect any parse errors found and add them to the result.
	eit := g.findAllNodes(parser.NodeTypeError).BuildNodeIterator(
		parser.NodePredicateErrorMessage,
		parser.NodePredicateSource,
		parser.NodePredicateStartRune)

	for eit.Next() {
		sal := salForPredicates(eit.Values())
		errorReporter(compilercommon.NewSourceError(sal, eit.Values()[parser.NodePredicateErrorMessage]))
	}

	// Verify all 'from ... import ...' are valid.
	fit := g.findAllNodes(parser.NodeTypeImportPackage).
		Has(parser.NodeImportPredicateSubsource).
		BuildNodeIterator(parser.NodeImportPredicateSubsource,
		parser.NodePredicateSource,
		parser.NodePredicateStartRune)

	for fit.Next() {
		// Load the package information.
		packageInfo := g.getPackageForImport(fit.Node())
		if !packageInfo.IsSRGPackage() {
			continue
		}

		// Search for the subsource.
		subsource := fit.Values()[parser.NodeImportPredicateSubsource]
		_, found := packageInfo.FindTypeOrMemberByName(subsource, ModuleResolveExportedOnly)
		if !found {
			source := fit.Node().GetIncomingNode(parser.NodeImportPredicatePackageRef).Get(parser.NodeImportPredicateSource)
			sal := salForPredicates(fit.Values())
			errorReporter(compilercommon.SourceErrorf(sal, "Import '%s' not found under package '%s'", subsource, source))
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
			sal := salForNode(decorator)
			errorReporter(compilercommon.SourceErrorf(sal, "Alias decorator requires a single string literal parameter"))
			continue
		}

		var aliasName = parameter.Get(parser.NodeStringLiteralExpressionValue)
		aliasName = aliasName[1 : len(aliasName)-1] // Remove the quotes.

		aliasedType := SRGType{decorator.GetIncomingNode(parser.NodeTypeDefinitionDecorator), g}
		g.aliasMap[aliasName] = aliasedType
	}
}
