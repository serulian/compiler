// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"path"
	"regexp"
	"strings"

	"github.com/serulian/compiler/compilercommon"
)

// fromImportRegex matches a `from ... import` snippet.
var fromImportRegex, _ = regexp.Compile("^from ((\"[^\"]+\")|(([a-zA-Z0-9]+)(`[^`]+`))|([a-zA-Z0-9]+(\\.[a-zA-Z0-9]+)*)) import $")

// importRegex matches either `import `, `import foo.`, `from ` or `from foo.`
var importRegex, _ = regexp.Compile("^(from|import) ([0-9a-zA-Z]+(\\.[0-9a-zA-Z]+)*\\.?)?$")

type importKind int

const (
	importKindPackage importKind = iota
	importKindFromPackage
)

// importSnippet represents a snippet of an import statement in one form or another. Allows for
// easier handling of completions of imports.
type importSnippet struct {
	// prefix is the prefix for the import (either 'from' or 'import')
	prefix string

	// sourceKind is the kind of source for the package or module. Will be empty string for SRG.
	sourceKind string

	// importSource is the source module/package referenced by the import.
	importSource string

	// importKind defines the kind of import.
	importKind importKind
}

// buildImportSnippet returns an import snippet for the given snippet string. If a snippet could
// not be parsed, returns false.
func buildImportSnippet(snippetString string) (importSnippet, bool) {
	fromPackageSubmatch := fromImportRegex.FindStringSubmatch(snippetString)
	if len(fromPackageSubmatch) > 0 {
		importSource := fromPackageSubmatch[1]
		if fromPackageSubmatch[4] != "" {
			importSource = fromPackageSubmatch[5]
		}

		return importSnippet{
			importKind:   importKindFromPackage,
			sourceKind:   fromPackageSubmatch[4],
			importSource: importSource,
			prefix:       "from",
		}, true
	}

	importSubmatch := importRegex.FindStringSubmatch(snippetString)
	if len(importSubmatch) > 0 {
		return importSnippet{
			importKind:   importKindPackage,
			sourceKind:   "",
			importSource: importSubmatch[2],
			prefix:       importSubmatch[1],
		}, true
	}

	return importSnippet{}, false
}

// populateCompletions populates the completions for this import snippet into the given
// builder, if applicable.
func (is importSnippet) populateCompletions(builder *completionBuilder, parentModule compilercommon.InputSource) bool {
	if is.importKind == importKindPackage {
		return is.populatePackages(builder, parentModule)
	}

	return is.populateTypesAndMembers(builder, parentModule)
}

// populatePackages populates all packages found under or matching the source package.
func (is importSnippet) populatePackages(builder *completionBuilder, parentModule compilercommon.InputSource) bool {
	// Check if this is a sub-package/module by looking for the presence of a dot at the end.
	packageDir := path.Dir(string(parentModule))
	isSubPackage := false
	prefix := is.importSource

	if strings.HasSuffix(is.importSource, ".") {
		// Sub-package/module.
		packageDir = path.Join(packageDir, is.importSource[0:len(is.importSource)-2])
		isSubPackage = true
		prefix = ""
	}

	modulesOrPackages, err := builder.handle.scopeResult.Graph.PackageLoader().ListSubModulesAndPackages(packageDir)
	if err != nil {
		return false
	}

	for _, moduleOrPackage := range modulesOrPackages {
		// If this is a subpackage that we obviously cannot import a non-Serulian file, so
		// we skip.
		if (isSubPackage || prefix != "") && moduleOrPackage.SourceKind != "" {
			continue
		}

		// Skip the parent module.
		if moduleOrPackage.Path == string(parentModule) {
			continue
		}

		// Make sure the import matches the current prefix.
		if strings.HasPrefix(moduleOrPackage.Name, prefix) {
			builder.addImport(moduleOrPackage.Name, moduleOrPackage.SourceKind)
		}
	}

	return true
}

// populateTypesAndMembers populates the types and members of the source package.
func (is importSnippet) populateTypesAndMembers(builder *completionBuilder, parentModule compilercommon.InputSource) bool {
	packagePath := is.importSource
	isVCSPath := false

	// Construct the path and VCS info for lookup in the package loader.
	if strings.HasPrefix(packagePath, "\"") || strings.HasPrefix(packagePath, "`") {
		packagePath = packagePath[1 : len(packagePath)-2]
	}

	if strings.Contains(packagePath, "/") {
		isVCSPath = true
	} else {
		pieces := strings.Split(packagePath, ".")
		packagePath = path.Join(pieces...)
		packagePath = path.Join(path.Dir(string(parentModule)), packagePath)
	}

	// Retrieve the package information for the path (if any).
	packageInfo, err := builder.handle.scopeResult.Graph.PackageLoader().LocalPackageInfoForPath(packagePath, is.sourceKind, isVCSPath)
	if err != nil {
		return false
	}

	// Lookup all types and members under the package from the typegraph.
	for _, completion := range builder.handle.scopeResult.Graph.TypeGraph().TypeOrMembersUnderPackage(packageInfo) {
		if completion.IsAccessibleTo(parentModule) {
			builder.addTypeOrMember(completion)
		}
	}
	return true
}
