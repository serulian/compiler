// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/parser"
	"github.com/serulian/compiler/vcs"
)

// emitFile emits the source code for a file.
func (sf *sourceFormatter) emitFile(node formatterNode) {
	// Emit the imports for the file.
	sf.emitImports(node)

	// Emit the module-level definitions.
	sf.emitOrderedNodes(node.getChildren(parser.NodePredicateChild))
}

// importInfo is a struct that represents the parsed information about an import.
type importInfo struct {
	node formatterNode

	source   string
	kind     string
	packages []importPackageInfo

	comparisonKey string
	sortKey       string

	isVCS      bool
	isSerulian bool
}

type importPackageInfo struct {
	subsource   string
	name        string
	packageName string
}

type byImportSortKey []importInfo

func (s byImportSortKey) Len() int {
	return len(s)
}

func (s byImportSortKey) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byImportSortKey) Less(i, j int) bool {
	return s[i].sortKey < s[j].sortKey
}

// emitImports emits the import statements for the source file.
func (sf *sourceFormatter) emitImports(node formatterNode) {
	var sortedImports = make([]importInfo, 0)
	for _, importNode := range node.getChildrenOfType(parser.NodePredicateChild, parser.NodeTypeImport) {
		// Remove any padding around the source name for VCS or non-Serulian imports.
		var source = importNode.getProperty(parser.NodeImportPredicateSource)
		if strings.HasPrefix(source, "`") || strings.HasPrefix(source, "\"") {
			source = source[1 : len(source)-1]
		}

		// Pull out the various pieces of the import.
		kind := importNode.getProperty(parser.NodeImportPredicateKind)

		// Pull out the package name(s).
		var packages = make([]importPackageInfo, 0)
		var packagesKey = ""
		for _, packageNode := range importNode.getChildren(parser.NodeImportPredicatePackageRef) {
			subsource := packageNode.getProperty(parser.NodeImportPredicateSubsource)
			name := packageNode.getProperty(parser.NodeImportPredicateName)
			packageName := packageNode.getProperty(parser.NodeImportPredicatePackageName)

			packagesKey += ":" + subsource + ":" + name + ":" + packageName

			packages = append(packages, importPackageInfo{subsource, name, packageName})
		}

		// Check for VCS and Serulian imports.
		isVCS := strings.Contains(source, "/")
		isSerulian := kind == ""

		// Determine the various runes for the sorting key.
		var vcsRune = 'a'
		if isVCS {
			vcsRune = 'z'
		}

		var serulianRune = 'z'
		if isSerulian {
			serulianRune = 'a'
		}

		var directRune = 'a'
		if len(packages) == 1 && packages[0].subsource == "" {
			directRune = 'z'
		}

		info := importInfo{
			node: importNode,

			source:   source,
			kind:     kind,
			packages: packages,

			sortKey:       fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s", vcsRune, serulianRune, kind, directRune, source, packagesKey),
			comparisonKey: kind,

			isVCS:      isVCS,
			isSerulian: isSerulian,
		}

		sortedImports = append(sortedImports, info)
	}

	// Sort the imports:
	// - VCS imports (non-Serulian)
	// - VCS imports (Serulian)
	// - Local imports (non-Serulian)
	// - Local imports (Serulian)
	sort.Sort(byImportSortKey(sortedImports))

	// Emit the imports.
	sf.emitImportInfos(sortedImports)
}

// emitImportInfos emits the importInfo structs as imports.
func (sf *sourceFormatter) emitImportInfos(infos []importInfo) {
	var lastKey = ""
	for index, info := range infos {
		if index > 0 && info.comparisonKey != lastKey {
			sf.appendLine()
		}

		sf.emitImport(info)
		lastKey = info.comparisonKey
	}
}

// emitModifiedImportSource emits the source for the import,
// modified per the freezing/unfreezing options.
func (sf *sourceFormatter) emitModifiedImportSource(info importInfo) bool {
	parsed, err := vcs.ParseVCSPath(info.source)
	if err != nil {
		sf.importHandling.logError(info.node, "Could not parse VCS path '%v': %v", info.source, err)
		return false
	}

	// Check if the import's URL was specified to be modified.
	if !sf.importHandling.matchesImport(parsed.URL()) {
		return false
	}

	switch sf.importHandling.option {
	case importHandlingUnfreeze:
		// For unfreezing, append the HEAD form of the VCS path.
		sf.importHandling.logSuccess(info.node, "Unfreezing '%v'", info.source)
		sf.append(parsed.AsGeneric().String())
		return true

	case importHandlingUpdate:
		// Make sure the import refers to a tag that has a semver.
		if parsed.Tag() == "" {
			sf.importHandling.logInfo(info.node, "Import '%v' doesn't refer to a version; skipped", info.source)
			return false
		}

		// For updating, perform VCS checkout and append the latest applicable minor version of the
		// import, as per semver. If none, then we don't change the import.
		inspectInfo, err := sf.getVCSInfo(info)
		if err != nil {
			return false
		}

		// Find the latest *minor* version, and update to it.
		currentTag, err := compilerutil.SemanticUpdateVersion(parsed.Tag(), inspectInfo.Tags, compilerutil.UpdateVersionMinor)
		if err != nil {
			sf.importHandling.logInfo(info.node, "Import '%v' doesn't refer to a semantic version; skipped", info.source)
			return false
		}

		if currentTag == "" {
			sf.importHandling.logInfo(info.node, "No updated version found for '%v'", info.source)
			return false
		}

		sf.importHandling.logSuccess(info.node, "Updating '%v' to version '%v'", info.source, currentTag)
		sf.append(parsed.WithTag(currentTag).String())
		return true

	case importHandlingUpgrade:
		// For upgrading, perform VCS checkout and append the latest stable version of the
		// import, as per semver. If none, then we don't change the import.
		inspectInfo, err := sf.getVCSInfo(info)
		if err != nil {
			return false
		}

		// Find the latest *major* version, and update to it.
		currentTag, err := compilerutil.SemanticUpdateVersion(parsed.Tag(), inspectInfo.Tags, compilerutil.UpdateVersionMajor)
		if err != nil {
			sf.importHandling.logInfo(info.node, "Import '%v' doesn't refer to a semantic version; skipped", info.source)
			return false
		}

		if currentTag == "" {
			sf.importHandling.logInfo(info.node, "No updated version found for '%v'", info.source)
			return false
		}

		sf.importHandling.logSuccess(info.node, "Upgrading '%v' to version '%v'", info.source, currentTag)
		sf.append(parsed.WithTag(currentTag).String())
		return true

	case importHandlingCustomFreeze:
		if sf.importHandling.customTagOrCommit.isTag {
			sf.importHandling.logSuccess(info.node, "Freezing '%v' at tag '%v'", info.source, sf.importHandling.customTagOrCommit.commitOrTag)
			sf.append(parsed.WithTag(sf.importHandling.customTagOrCommit.commitOrTag).String())
		} else {
			sf.importHandling.logSuccess(info.node, "Freezing '%v' at commit '%v'", info.source, sf.importHandling.customTagOrCommit.commitOrTag)
			sf.append(parsed.WithCommit(sf.importHandling.customTagOrCommit.commitOrTag).String())
		}

		return true

	case importHandlingFreeze:
		// For freezing, perform VCS checkout and append the commit of
		// the checked out info.
		inspectInfo, err := sf.getVCSInfo(info)
		if err != nil {
			return false
		}

		sf.importHandling.logSuccess(info.node, "Freezing '%v' at commit '%v'", info.source, inspectInfo.CommitSHA)
		sf.append(parsed.WithCommit(inspectInfo.CommitSHA).String())
		return true

	case importHandlingNone:
		// No changes needed.
		return false
	}

	panic("Unknown import handling option")
	return false
}

// emitImport emits the formatted form of the import.
func (sf *sourceFormatter) emitImport(info importInfo) {
	sf.emitCommentsForNode(info.node)

	if len(info.packages) == 1 && info.packages[0].subsource == "" {
		sf.emitDirectImport(info)
	} else {
		sf.emitSourcedImported(info)
	}
}

func (sf *sourceFormatter) emitDirectImport(info importInfo) {
	// import somepackage
	// import somepackage as somethingelse
	// import somekind`somepackage` as something

	sf.append("import ")
	if !info.isSerulian {
		sf.append(info.kind)
		sf.append("`")
	} else if info.isVCS {
		sf.append("\"")
	}

	// Handle freezing, unfreezing of VCS.
	if info.isVCS {
		if !sf.emitModifiedImportSource(info) {
			sf.append(info.source)
		}
	} else {
		sf.append(info.source)
	}

	if !info.isSerulian {
		sf.append("`")
	} else if info.isVCS {
		sf.append("\"")
	}

	packageInfo := info.packages[0]
	if packageInfo.packageName != info.source || !info.isSerulian {
		sf.append(" as ")
		sf.append(packageInfo.packageName)
	}

	sf.appendLine()
}

func (sf *sourceFormatter) emitSourcedImported(info importInfo) {
	// from "something" import something
	// from something import something
	// from somekind`something` import something
	// from somekind`something` import something as somethingelse

	sf.append("from ")

	if !info.isSerulian {
		sf.append(info.kind)
		sf.append("`")
	} else if info.isVCS {
		sf.append("\"")
	}

	// Handle freezing, unfreezing of VCS.
	if info.isVCS {
		if !sf.emitModifiedImportSource(info) {
			sf.append(info.source)
		}
	} else {
		sf.append(info.source)
	}

	if !info.isSerulian {
		sf.append("`")
	} else if info.isVCS {
		sf.append("\"")
	}

	sf.append(" import ")

	for index, packageInfo := range info.packages {
		if index > 0 {
			sf.append(", ")
		}

		if packageInfo.subsource != "" {
			sf.append(packageInfo.subsource)
		}

		if packageInfo.name != packageInfo.subsource {
			sf.append(" as ")
			sf.append(packageInfo.name)
		}
	}

	sf.appendLine()
}

type byNamePredicate []formatterNode

func (s byNamePredicate) Len() int {
	return len(s)
}

func (s byNamePredicate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byNamePredicate) Less(i, j int) bool {
	return s[i].getProperty("named") < s[j].getProperty("named")
}
