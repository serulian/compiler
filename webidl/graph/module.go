// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"path"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/webidl/parser"
)

// IRGModule wraps a module defined in the IRG.
type IRGModule struct {
	compilergraph.GraphNode
	irg *WebIRG // The parent IRG.
}

// GetModules returns all the modules defined in the IRG.
func (g *WebIRG) GetModules() []IRGModule {
	it := g.findAllNodes(parser.NodeTypeFile).BuildNodeIterator()
	var modules []IRGModule

	for it.Next() {
		modules = append(modules, IRGModule{it.Node(), g})
	}

	return modules
}

// InputSource returns the input source for this module.
func (m IRGModule) InputSource() compilercommon.InputSource {
	return compilercommon.InputSource(m.GraphNode.Get(parser.NodePredicateSource))
}

// Name returns the name of the module.
func (m IRGModule) Name() string {
	return path.Base(string(m.InputSource()))
}

// Node returns the underlying node.
func (m IRGModule) Node() compilergraph.GraphNode {
	return m.GraphNode
}

// Declarations returns the declarations directly under the module.
func (m IRGModule) Declarations() []IRGDeclaration {
	it := m.GraphNode.StartQuery().
		Out(parser.NodePredicateChild).
		IsKind(parser.NodeTypeDeclaration).
		BuildNodeIterator()

	var decls []IRGDeclaration
	for it.Next() {
		decls = append(decls, IRGDeclaration{it.Node(), m.irg})
	}

	return decls
}
