// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"strconv"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
)

// noopScopeApplier is a scope applier that does nothing.
type noopScopeApplier struct{}

func (nsa noopScopeApplier) NodeScoped(node compilergraph.GraphNode, result proto.ScopeInfo) {}
func (nsa noopScopeApplier) DecorateWithSecondaryLabel(node compilergraph.GraphNode, label proto.ScopeLabel) {
}
func (nsa noopScopeApplier) AddErrorForSourceNode(node compilergraph.GraphNode, message string)   {}
func (nsa noopScopeApplier) AddWarningForSourceNode(node compilergraph.GraphNode, message string) {}

// concreteScopeApplier is a scope applier that writes changes back to the scope graph using
// the given modifier.
type concreteScopeApplier struct {
	modifier compilergraph.GraphLayerModifier
}

func (csa concreteScopeApplier) NodeScoped(node compilergraph.GraphNode, result proto.ScopeInfo) {
	scopeNode := csa.modifier.CreateNode(NodeTypeResolvedScope)
	scopeNode.DecorateWithTagged(NodePredicateScopeInfo, &result)
	scopeNode.Connect(NodePredicateSource, node)
}

func (csa concreteScopeApplier) DecorateWithSecondaryLabel(node compilergraph.GraphNode, label proto.ScopeLabel) {
	labelNode := csa.modifier.CreateNode(NodeTypeSecondaryLabel)
	labelNode.Decorate(NodePredicateSecondaryLabelValue, strconv.Itoa(int(label)))
	labelNode.Connect(NodePredicateLabelSource, node)
}

func (csa concreteScopeApplier) AddErrorForSourceNode(node compilergraph.GraphNode, message string) {
	errorNode := csa.modifier.CreateNode(NodeTypeError)
	errorNode.Decorate(NodePredicateNoticeMessage, message)
	errorNode.Connect(NodePredicateNoticeSource, node)
}

func (csa concreteScopeApplier) AddWarningForSourceNode(node compilergraph.GraphNode, message string) {
	warningNode := csa.modifier.CreateNode(NodeTypeWarning)
	warningNode.Decorate(NodePredicateNoticeMessage, message)
	warningNode.Connect(NodePredicateNoticeSource, node)
}
