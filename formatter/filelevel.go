// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import (
	"sort"
	"strings"

	"github.com/serulian/compiler/parser"
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

	sortKey string

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
	// Load the imports for the node, categorized by VCS vs local and Serulian vs non-serulian.
	var vcsNonSerulian = make([]importInfo, 0)
	var vcsSerulian = make([]importInfo, 0)
	var localNonSerulian = make([]importInfo, 0)
	var localSerulian = make([]importInfo, 0)

	for _, importNode := range node.getChildrenOfType(parser.NodePredicateChild, parser.NodeTypeImport) {
		var source = importNode.getProperty(parser.NodeImportPredicateSource)
		if strings.HasPrefix(source, "`") || strings.HasPrefix(source, "\"") {
			source = source[1 : len(source)-1]
		}

		subsource := importNode.getProperty(parser.NodeImportPredicateSubsource)
		kind := importNode.getProperty(parser.NodeImportPredicateKind)
		name := importNode.getProperty(parser.NodeImportPredicateName)

		info := importInfo{
			node: importNode,

			source:    source,
			subsource: subsource,
			kind:      kind,
			name:      name,

			sortKey: kind + "://" + source + "/" + subsource + "/" + name,

			isVCS:      strings.Contains(importNode.getProperty(parser.NodeImportPredicateSource), "/"),
			isSerulian: kind == "",
		}

		if info.isSerulian {
			if info.isVCS {
				vcsSerulian = append(vcsSerulian, info)
			} else {
				localSerulian = append(localSerulian, info)
			}
		} else {
			if info.isVCS {
				vcsNonSerulian = append(vcsNonSerulian, info)
			} else {
				localNonSerulian = append(localNonSerulian, info)
			}
		}
	}

	// Sort the imports.
	sort.Sort(byImportSortKey(vcsNonSerulian))
	sort.Sort(byImportSortKey(vcsSerulian))
	sort.Sort(byImportSortKey(localNonSerulian))
	sort.Sort(byImportSortKey(localSerulian))

	// Emit the imports in the following order:
	// - VCS imports (non-Serulian)
	// - VCS imports (Serulian)
	// - Local imports (non-Serulian)
	// - Local imports (Serulian)
	sf.emitImportInfos(vcsNonSerulian)
	sf.emitImportInfos(vcsSerulian)
	sf.emitImportInfos(localNonSerulian)
	sf.emitImportInfos(localSerulian)
}

// emitImportInfos emits the importInfo structs as imports.
func (sf *sourceFormatter) emitImportInfos(infos []importInfo) {
	for _, info := range infos {
		sf.emitImport(info)
	}

	if len(infos) > 0 {
		sf.appendLine()
	}
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

	sf.append(info.source)

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
