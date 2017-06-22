// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/webidl/parser"
)

// SERIALIZABLE_OPS defines the WebIDL custom ops that mark a type as serializable.
var SERIALIZABLE_OPS = map[string]bool{
	"jsonifier":  true,
	"serializer": true,
}

// Declarations returns all the type declarations in the WebIDL IRG.
func (g *WebIRG) Declarations() []IRGDeclaration {
	dit := g.findAllNodes(parser.NodeTypeDeclaration).BuildNodeIterator()

	var declarations = make([]IRGDeclaration, 0)
	for dit.Next() {
		declaration := IRGDeclaration{dit.Node(), g}
		declarations = append(declarations, declaration)
	}

	return declarations
}

// FindDeclaration finds the declaration with the given name in the IRG, if any.
func (g *WebIRG) FindDeclaration(name string) (IRGDeclaration, bool) {
	declNode, hasDeclaration := g.layer.StartQuery(name).
		In(parser.NodePredicateDeclarationName).
		TryGetNode()

	if !hasDeclaration {
		return IRGDeclaration{}, false
	}

	return IRGDeclaration{declNode, g}, true
}

type DeclarationKind int

const (
	InterfaceDeclaration DeclarationKind = iota
)

// IRGDeclaration wraps a WebIDL declaration.
type IRGDeclaration struct {
	compilergraph.GraphNode
	irg *WebIRG // The parent IRG.
}

// Name returns the name of the declaration.
func (i *IRGDeclaration) Name() string {
	return i.GraphNode.Get(parser.NodePredicateDeclarationName)
}

// Kind returns the kind of declaration.
func (i *IRGDeclaration) Kind() DeclarationKind {
	kindStr := i.GraphNode.Get(parser.NodePredicateDeclarationKind)
	switch kindStr {
	case "interface":
		return InterfaceDeclaration

	default:
		panic("Unknown kind of WebIDL declaration")
	}
}

// IsSerializable returns whether the declaration contains one of the custom
// operations that makes the declared interface serializable.
func (i *IRGDeclaration) IsSerializable() bool {
	for _, customop := range i.CustomOperations() {
		if _, ok := SERIALIZABLE_OPS[customop]; ok {
			return true
		}
	}
	return false
}

// Module returns the parent module.
func (i *IRGDeclaration) Module() IRGModule {
	moduleNode := i.GraphNode.GetIncomingNode(parser.NodePredicateChild)
	return IRGModule{moduleNode, i.irg}
}

// ParentType returns the declared parent type of the declaration, if any.
func (i *IRGDeclaration) ParentType() (string, bool) {
	return i.GraphNode.TryGet(parser.NodePredicateDeclarationParentType)
}

// FindMember finds the member under this declaration with the given name, if any.
func (i *IRGDeclaration) FindMember(name string) (IRGMember, bool) {
	memberNode, hasMember := i.GraphNode.StartQuery().
		Out(parser.NodePredicateDeclarationMember).
		Has(parser.NodePredicateMemberName, name).
		TryGetNode()

	if !hasMember {
		return IRGMember{}, false
	}

	return IRGMember{memberNode, i.irg}, true
}

// Members returns all the members declared in the declaration.
func (i *IRGDeclaration) Members() []IRGMember {
	mit := i.GraphNode.StartQuery().
		Out(parser.NodePredicateDeclarationMember).
		BuildNodeIterator()

	var members = make([]IRGMember, 0)
	for mit.Next() {
		member := IRGMember{mit.Node(), i.irg}
		members = append(members, member)
	}

	return members
}

// CustomOperations returns all the custom operations defined on the declaration.
func (i *IRGDeclaration) CustomOperations() []string {
	mit := i.GraphNode.StartQuery().
		Out(parser.NodePredicateDeclarationCustomOperation).
		BuildNodeIterator(parser.NodePredicateCustomOpName)

	var customOps = make([]string, 0)
	for mit.Next() {
		customOps = append(customOps, mit.GetPredicate(parser.NodePredicateCustomOpName).String())
	}

	return customOps
}

// HasAnnotation returns true if the declaration is decorated with the given annotation.
func (i *IRGDeclaration) HasAnnotation(name string) bool {
	_, hasNode := i.GraphNode.StartQuery().
		Out(parser.NodePredicateDeclarationAnnotation).
		Has(parser.NodePredicateAnnotationName, name).
		TryGetNode()

	return hasNode
}

// HasOneAnnotation returns true if the declaration is decorated with one of the given annotations.
func (i *IRGDeclaration) HasOneAnnotation(names ...interface{}) bool {
	_, hasNode := i.GraphNode.StartQuery().
		Out(parser.NodePredicateDeclarationAnnotation).
		Has(parser.NodePredicateAnnotationName, names...).
		TryGetNode()

	return hasNode
}

// GetAnnotations returns all the annotations with the given name declared on the declaration.
func (i *IRGDeclaration) GetAnnotations(name string) []IRGAnnotation {
	ait := i.GraphNode.StartQuery().
		Out(parser.NodePredicateDeclarationAnnotation).
		Has(parser.NodePredicateAnnotationName, name).
		BuildNodeIterator()

	var annotations = make([]IRGAnnotation, 0)
	for ait.Next() {
		annotation := IRGAnnotation{ait.Node(), i.irg}
		annotations = append(annotations, annotation)
	}

	return annotations
}

// Annotations returns all the annotations declared on the declaration.
func (i *IRGDeclaration) Annotations() []IRGAnnotation {
	ait := i.GraphNode.StartQuery().
		Out(parser.NodePredicateDeclarationAnnotation).
		BuildNodeIterator()

	var annotations = make([]IRGAnnotation, 0)
	for ait.Next() {
		annotation := IRGAnnotation{ait.Node(), i.irg}
		annotations = append(annotations, annotation)
	}

	return annotations
}

// SourceRange returns the source range of the declaration in source.
func (i *IRGDeclaration) SourceRange() (compilercommon.SourceRange, bool) {
	return i.irg.SourceRangeOf(i.GraphNode)
}
