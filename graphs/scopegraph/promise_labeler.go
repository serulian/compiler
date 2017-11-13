// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"container/list"
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
)

var _ = fmt.Printf

type promiseLabeler struct {
	builder     *scopeBuilder
	entrypoints *list.List
	workset     *list.List
}

func newPromiseLabeler(builder *scopeBuilder) *promiseLabeler {
	return &promiseLabeler{
		builder:     builder,
		entrypoints: list.New(),
		workset:     list.New(),
	}
}

func (pl *promiseLabeler) addEntrypoint(entrypoint srg.SRGImplementable) {
	pl.workset.PushBack(entrypoint)
	pl.entrypoints.PushBack(entrypoint)
}

func (pl *promiseLabeler) labelEntrypoints() map[string]bool {
	currentLabel := map[compilergraph.GraphNodeId]proto.ScopeLabel{}
	dynamicNames := map[string]bool{}

	var hasChanges = false

	labelEntrypoint := func(e *list.Element, entrypoint srg.SRGImplementable, entrypointScope *proto.ScopeInfo, label proto.ScopeLabel) {
		// Special handling: If the entrypoint is a generator, then Next is also dynamic.
		if entrypointScope.HasLabel(proto.ScopeLabel_GENERATOR_STATEMENT) {
			dynamicNames["Next"] = true
		}

		currentLabel[entrypoint.Node().NodeId] = label
		if !entrypoint.IsMember() {
			currentLabel[entrypoint.ContainingMember().NodeId] = label
		}

		pl.workset.Remove(e)
		memberName, _ := entrypoint.ContainingMember().Name()
		dynamicNames[memberName] = true
		hasChanges = true
	}

outerloop:
	for {
		hasChanges = false
		var next *list.Element
		for e := pl.workset.Front(); e != nil; e = next {
			next = e.Next()
			entrypoint := e.Value.(srg.SRGImplementable)

			// If the entrypoint awaits on another node, then it is a YES.
			entrypointScope := pl.builder.getScopeForRootNode(entrypoint.Node())
			if entrypointScope.HasLabel(proto.ScopeLabel_AWAITS) {
				labelEntrypoint(e, entrypoint, entrypointScope, proto.ScopeLabel_SML_PROMISING_YES)
				continue outerloop
			}

			// If the entrypoint calls an anonymous closure, it might be promising (since the function)
			// may promise.
			if entrypointScope.HasLabel(proto.ScopeLabel_CALLS_ANONYMOUS_CLOSURE) {
				labelEntrypoint(e, entrypoint, entrypointScope, proto.ScopeLabel_SML_PROMISING_MAYBE)
				continue outerloop
			}

			// For each of the dynamic dependencies, check if they are promising.
			for _, name := range entrypointScope.GetDynamicDependencies() {
				if _, exists := dynamicNames[name]; exists {
					labelEntrypoint(e, entrypoint, entrypointScope, proto.ScopeLabel_SML_PROMISING_MAYBE)
					continue outerloop
				}
			}

			// For each of the static dependencies, check if they are promising.
			var found proto.ScopeLabel = proto.ScopeLabel_SML_PROMISING_NO
			for _, staticDep := range entrypointScope.GetStaticDependencies() {
				memberNodeId := compilergraph.GraphNodeId(staticDep.GetReferencedNode())
				member := pl.builder.sg.tdg.GetTypeOrMember(memberNodeId)

				// Skip fields, as they are always accessed after their value has been computed
				if member.IsField() {
					continue
				}

				// If the member is under an interface, check it dynamically.
				parent := member.Parent()
				if asType, isType := parent.AsType(); isType {
					if asType.TypeKind() == typegraph.ImplicitInterfaceType {
						if _, exists := dynamicNames[member.Name()]; exists {
							labelEntrypoint(e, entrypoint, entrypointScope, proto.ScopeLabel_SML_PROMISING_MAYBE)
							continue outerloop
						}

						continue
					}
				}

				// Otherwise, check the member directly.
				sourceNodeId, hasSourceNode := member.SourceNodeId()
				if !hasSourceNode {
					continue
				}

				if label, exists := currentLabel[sourceNodeId]; exists && label > found {
					found = label
				}
			}

			if found != proto.ScopeLabel_SML_PROMISING_NO {
				labelEntrypoint(e, entrypoint, entrypointScope, found)
				continue outerloop
			}
		}

		if !hasChanges {
			break
		}
	}

	// Any remaining members do not need to promise.
	for e := pl.workset.Front(); e != nil; e = e.Next() {
		entrypoint := e.Value.(srg.SRGImplementable)
		pl.builder.decorateWithSecondaryLabel(entrypoint.GraphNode, proto.ScopeLabel_SML_PROMISING_NO)
	}

	for e := pl.entrypoints.Front(); e != nil; e = e.Next() {
		entrypoint := e.Value.(srg.SRGImplementable)
		if value, exists := currentLabel[entrypoint.Node().NodeId]; exists {
			pl.builder.decorateWithSecondaryLabel(entrypoint.GraphNode, value)
		}
	}

	return dynamicNames
}
