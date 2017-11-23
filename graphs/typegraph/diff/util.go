// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"sort"

	"github.com/serulian/compiler/graphs/typegraph"
)

func getTypesAndMembers(modules []typegraph.TGModule) (map[string]typegraph.TGTypeDecl, map[string]typegraph.TGMember) {
	types := map[string]typegraph.TGTypeDecl{}
	members := map[string]typegraph.TGMember{}
	for _, module := range modules {
		for _, typedecl := range module.Types() {
			// Aliases are treated as if their refer to their underlying type,
			// just with the aliased name.
			aliasedType, hasAliasedType := typedecl.AliasedType()
			if hasAliasedType {
				types[typedecl.Name()] = aliasedType
				continue
			}

			types[typedecl.Name()] = typedecl
		}
		for _, member := range module.MembersAndOperators() {
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

func (mhw memberHolderWrap) GetMemberOrOperator(name string) (typegraph.TGMember, bool) {
	member, ok := mhw[name]
	return member, ok
}

func (mhw memberHolderWrap) MembersAndOperators() []typegraph.TGMember {
	members := make([]typegraph.TGMember, 0, len(mhw))
	for _, member := range mhw {
		members = append(members, member)
	}
	return members
}

type sortableMemberDiff []MemberDiff

func (s sortableMemberDiff) Len() int           { return len(s) }
func (s sortableMemberDiff) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortableMemberDiff) Less(i, j int) bool { return s[i].Name < s[j].Name }

type sortableTypeDiff []TypeDiff

func (s sortableTypeDiff) Len() int           { return len(s) }
func (s sortableTypeDiff) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortableTypeDiff) Less(i, j int) bool { return s[i].Name < s[j].Name }

func sortedMemberDiffs(diffs []MemberDiff) []MemberDiff {
	sort.Sort(sortableMemberDiff(diffs))
	return diffs
}

func sortedTypeDiffs(diffs []TypeDiff) []TypeDiff {
	sort.Sort(sortableTypeDiff(diffs))
	return diffs
}
