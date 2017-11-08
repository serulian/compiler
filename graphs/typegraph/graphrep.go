// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/typegraph/proto"
)

// GetJSONForm returns this type graph serialized to JSON. *For testing purposes only*.
func (tg *TypeGraph) GetJSONForm() string {
	return tg.GetFilteredJSONForm([]string{}, []compilergraph.TaggedValue{})
}

// GetFilteredJSONForm returns the filtered type graph serialized to JSON. *For testing purposes only*.
func (tg *TypeGraph) GetFilteredJSONForm(pathFilters []string, skipNodeKinds []compilergraph.TaggedValue) string {
	repMap := map[compilergraph.GraphNodeId]*graphNodeRep{}
	filterMap := map[string]bool{}
	for _, path := range pathFilters {
		filterMap[path] = true
	}

	// Start the walk at the type declarations.
	var startingNodes = make([]compilergraph.GraphNode, 0)
	for _, typeDecl := range tg.TypeDecls() {
		if len(pathFilters) == 0 || filterMap[typeDecl.ParentModule().Path()] {
			startingNodes = append(startingNodes, typeDecl.Node())
		}
	}

	for _, typeAlias := range tg.TypeAliases() {
		if len(pathFilters) == 0 || filterMap[typeAlias.ParentModule().Path()] {
			startingNodes = append(startingNodes, typeAlias.Node())
		}
	}

	for _, module := range tg.ModulesWithMembers() {
		if len(pathFilters) == 0 || filterMap[module.Path()] {
			startingNodes = append(startingNodes, module.Node())
		}
	}

	// Walk the graph outward from the type declaration nodes, building an in-memory tree
	// representation along the way.
	tg.layer.WalkOutward(startingNodes, func(result *compilergraph.WalkResult) bool {
		for _, nodeKind := range skipNodeKinds {
			if result.Node.Kind() == nodeKind {
				return false
			}
		}

		// Filter any predicates that match UUIDs, as they attach to other graph layers
		// and will have rotating IDs.
		normalizedPredicates := map[string]string{}

		var keys []string
		for name, value := range result.Predicates {
			if compilerutil.IsId(value) {
				normalizedPredicates[name] = "(NodeRef)"
			} else if strings.Contains(value, "|TypeReference") {
				// Convert type references into human-readable strings so that they don't change constantly
				// due to the underlying IDs.
				normalizedPredicates[name] = tg.AnyTypeReference().Build(value).(TypeReference).String()
			} else if name == "tdg-"+NodePredicateMemberSignature {
				esig := &proto.MemberSig{}
				sig := esig.Build(value[:len(value)-len("|MemberSig|tdg")]).(*proto.MemberSig)

				// Normalize the member type and constraints into human-readable strings.
				memberType := tg.AnyTypeReference().Build(sig.GetMemberType()).(TypeReference).String()
				sig.MemberType = memberType

				genericTypes := make([]string, len(sig.GetGenericConstraints()))
				for index, constraint := range sig.GetGenericConstraints() {
					genericTypes[index] = tg.AnyTypeReference().Build(constraint).(TypeReference).String()
				}

				sig.GenericConstraints = genericTypes
				marshalled, _ := sig.Marshal()
				normalizedPredicates[name] = string(marshalled)
			} else {
				normalizedPredicates[name] = value
			}

			keys = append(keys, name)
		}

		// Build a hash of all predicates and values.
		sort.Strings(keys)
		h := md5.New()
		for _, key := range keys {
			io.WriteString(h, key+":"+normalizedPredicates[key])
		}

		// Build the representation of the node.
		repKey := fmt.Sprintf("%x", h.Sum(nil))
		repMap[result.Node.NodeId] = &graphNodeRep{
			Key:        repKey,
			Kind:       result.Node.Kind(),
			Children:   map[string]graphChildRep{},
			Predicates: normalizedPredicates,
		}

		if result.ParentNode != nil {
			parentRep := repMap[result.ParentNode.NodeId]
			childRep := repMap[result.Node.NodeId]

			parentRep.Children[repKey] = graphChildRep{
				Predicate: result.IncomingPredicate,
				Child:     childRep,
			}
		}

		return result.Node.Kind() != NodeTypeAlias
	})

	rootReps := map[string]*graphNodeRep{}
	for _, typeDecl := range tg.TypeDecls() {
		if len(pathFilters) == 0 || filterMap[typeDecl.ParentModule().Path()] {
			rootReps[repMap[typeDecl.Node().NodeId].Key] = repMap[typeDecl.Node().NodeId]
		}
	}

	for _, typeAlias := range tg.TypeAliases() {
		if len(pathFilters) == 0 || filterMap[typeAlias.ParentModule().Path()] {
			rootReps[repMap[typeAlias.Node().NodeId].Key] = repMap[typeAlias.Node().NodeId]
		}
	}

	for _, module := range tg.ModulesWithMembers() {
		if len(pathFilters) == 0 || filterMap[module.Path()] {
			if _, exists := repMap[module.Node().NodeId]; exists {
				rootReps[repMap[module.Node().NodeId].Key] = repMap[module.Node().NodeId]
			}
		}
	}

	// Marshal the tree to JSON.
	b, _ := json.MarshalIndent(rootReps, "", "    ")
	return string(b)
}

type graphChildRep struct {
	Predicate string
	Child     *graphNodeRep
}

type graphNodeRep struct {
	Key        string
	Kind       interface{}
	Children   map[string]graphChildRep
	Predicates map[string]string
}
