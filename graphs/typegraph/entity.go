// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import "github.com/serulian/compiler/compilergraph"

// EntityKind defines the various kinds of entities in which types or members can
// be found
type EntityKind string

const (
	// EntityKindModule indicates the entity is a module.
	EntityKindModule EntityKind = "module"

	// EntityKindType indicates the entity is a type.
	EntityKindType EntityKind = "type"

	// EntityKindMember indicates the entity is a member.
	EntityKindMember EntityKind = "member"
)

// Entity represents an entity: module, type, or member.
type Entity struct {
	// Kind is it the kind of the entity.
	Kind EntityKind

	// NameOrPath represents the name or path of the entity. If a module,
	// this is the module's path. If a type, this is the type's name.
	// If a member, this is the member's *child name*.
	NameOrPath string

	// SourceGraphId is the ID of the source graph for the type node. If none,
	// this will be `typegraph`.
	SourceGraphId string
}

// TGEntity represents any form of entity: module, member or type.
type TGEntity interface {
	Name() string
	Node() compilergraph.GraphNode
	IsType() bool
	AsType() (TGTypeDecl, bool)
	EntityPath() []Entity
}
