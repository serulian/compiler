// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/serulian/compiler/packageloader"
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

	source    string
	subsource string
	name      string
	kind      string

	comparisonKey string
	sortKey       string

	isVCS      bool
	isSerulian bool
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
		subsource := importNode.getProperty(parser.NodeImportPredicateSubsource)
		kind := importNode.getProperty(parser.NodeImportPredicateKind)
		name := importNode.getProperty(parser.NodeImportPredicateName)

		// If the import has no subsource, then check for a distinct name.
		if subsource == "" {
			packageName := importNode.getProperty(parser.NodeImportPredicatePackageName)
			if packageName != source || kind != "" {
				name = packageName
			}
		}

		isVCS := strings.Contains(importNode.getProperty(parser.NodeImportPredicateSource), "/")
		isSerulian := kind == ""

		// Determine the various runes for the sorting key.
		var subsourceRune = 'z'
		if subsource != "" {
			subsourceRune = 'a'
		}

		var vcsRune = 'a'
		if isVCS {
			vcsRune = 'z'
		}

		var serulianRune = 'z'
		if isSerulian {
			serulianRune = 'a'
		}

		info := importInfo{
			node: importNode,

			source:    source,
			subsource: subsource,
			kind:      kind,
			name:      name,

			sortKey:       fmt.Sprintf("%s/%s/%s/%s/%s/%s", vcsRune, serulianRune, kind, subsourceRune, source, name),
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
		log.Printf("Could not parse VCS path '%v': %v", info.source, err)
		return false
	}

	// Check if the import's URL was specified to be modified.
	if !sf.importHandling.hasImport(parsed.URL()) {
		return false
	}

	switch sf.importHandling.option {
	case importHandlingUnfreeze:
		// For unfreezing, append the HEAD form of the VCS path.
		sf.append(parsed.AsHEAD().String())
		return true

	case importHandlingFreeze:
		// For freezing, perform VCS checkout and append the commit of
		// the checked out info.
		if commitSha, exists := sf.vcsCommitCache[parsed.URL()]; exists {
			sf.append(parsed.WithCommit(commitSha).String())
			return true
		}

		log.Printf("Performing checkout and inspection of '%v'", parsed.URL())
		commitSha, err, _ := vcs.PerformVCSCheckoutAndInspect(
			info.source, packageloader.SerulianPackageDirectory,
			sf.importHandling.vcsDevelopmentDirectories...)

		if err != nil {
			log.Printf("Could not checkout and inspect %v: %v", parsed.URL(), err)
			return false
		}

		sf.vcsCommitCache[parsed.URL()] = commitSha
		sf.append(parsed.WithCommit(commitSha).String())
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

	if info.subsource != "" {
		sf.append("from ")
	} else {
		sf.append("import ")
	}

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

	if info.subsource != "" {
		sf.append(" import ")
		sf.append(info.subsource)
	}

	if info.name != info.subsource {
		sf.append(" as ")
		sf.append(info.name)
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
