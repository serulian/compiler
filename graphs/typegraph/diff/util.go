// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"github.com/serulian/compiler/graphs/typegraph"
)

func getTypesAndMembers(modules []typegraph.TGModule) (map[string]typegraph.TGTypeDecl, map[string]typegraph.TGMember) {
	types := map[string]typegraph.TGTypeDecl{}
	members := map[string]typegraph.TGMember{}
	for _, module := range modules {
		for _, typedecl := range module.Types() {
			// Filter aliases out.
			if typedecl.TypeKind() == typegraph.AliasType {
				continue
			}

			types[typedecl.Name()] = typedecl
		}
		for _, member := range module.Members() {
			members[member.Name()] = member
		}
	}
	return types, members
}

type typeHolderWrap map[string]typegraph.TGTypeDecl

func (thw typeHolderWrap) GetType(name string) (typegraph.TGTypeDecl, bool) {
	typeDecl, ok := thw[name]
	return typeDecl, ok
}

func (thw typeHolderWrap) Types() []typegraph.TGTypeDecl {
	types := make([]typegraph.TGTypeDecl, 0, len(thw))
	for _, typeDecl := range thw {
		types = append(types, typeDecl)
	}
	return types
}

type memberHolderWrap map[string]typegraph.TGMember

func (mhw memberHolderWrap) GetMember(name string) (typegraph.TGMember, bool) {
	member, ok := mhw[name]
	return member, ok
}

func (mhw memberHolderWrap) Members() []typegraph.TGMember {
	members := make([]typegraph.TGMember, 0, len(mhw))
	for _, member := range mhw {
		members = append(members, member)
	}
	return members
}
