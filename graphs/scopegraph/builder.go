// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/parser"

	"github.com/streamrail/concurrent-map"
)

// scopeHandler is a handler function for scoping an SRG node of a particular kind.
type scopeHandler func(node compilergraph.GraphNode) proto.ScopeInfo

// scopeBuilder defines a type for easy scoping of the SRG.
type scopeBuilder struct {
	sg      *ScopeGraph
	nodeMap cmap.ConcurrentMap

	Status bool
}

// newScopeBuilder returns a new scope builder for the given scope graph.
func newScopeBuilder(sg *ScopeGraph) *scopeBuilder {
	return &scopeBuilder{sg, cmap.New(), true}
}

// getScope returns the scope for the given node, building (and waiting) if necessary.
func (sb *scopeBuilder) getScope(node compilergraph.GraphNode) *proto.ScopeInfo {
	// Check the map cache for the scope.
	found, ok := sb.nodeMap.Get(string(node.NodeId))
	if ok {
		result := found.(proto.ScopeInfo)
		return &result
	}

	built := <-sb.buildScope(node)
	return &built
}

// buildScope builds the scope for the given node, returning a channel
// that can be watched for the result.
func (sb *scopeBuilder) buildScope(node compilergraph.GraphNode) chan proto.ScopeInfo {
	var handler scopeHandler

	switch node.Kind {
	case parser.NodeTypeStatementBlock:
		handler = sb.scopeStatementBlock

	default:
		panic(fmt.Sprintf("Unknown SRG node in scoping: %v", node.Kind))
	}

	// Execute the handler in a gorountine and return the result channel.
	resultChan := make(chan proto.ScopeInfo)
	go (func() {
		result := handler(node)
		if !result.GetIsValid() {
			sb.Status = false
		}

		// Add the scope to the map and on the node.
		sb.nodeMap.Set(string(node.NodeId), result)

		scopeNode := sb.sg.layer.CreateNode(NodeTypeResolvedScope)
		scopeNode.DecorateWithTagged(NodePredicateScopeInfo, &result)
		scopeNode.Connect(NodePredicateSource, node)

		resultChan <- result
	})()

	return resultChan
}

// decorateWithError decorates an *SRG* node with the specified scope error.
func (sb *scopeBuilder) decorateWithError(node compilergraph.GraphNode, message string, args ...interface{}) {
	errorNode := sb.sg.layer.CreateNode(NodeTypeError)
	errorNode.Decorate(NodePredicateErrorMessage, fmt.Sprintf(message, args...))
	errorNode.Connect(NodePredicateErrorSource, node)
}

// GetWarnings returns the warnings created during the build pass.
func (sb *scopeBuilder) GetWarnings() []*compilercommon.SourceWarning {
	return make([]*compilercommon.SourceWarning, 0)
}

// GetErrors returns the errors created during the build pass.
func (sb *scopeBuilder) GetErrors() []*compilercommon.SourceError {
	var errors = make([]*compilercommon.SourceError, 0)

	it := sb.sg.layer.StartQuery().
		IsKind(NodeTypeError).
		BuildNodeIterator()

	for it.Next() {
		errNode := it.Node()

		// Lookup the location of the SRG source node.
		errorSource := sb.sg.srg.GetNode(compilergraph.GraphNodeId(errNode.Get(NodePredicateErrorSource)))
		location := sb.sg.srg.NodeLocation(errorSource)

		// Add the error.
		msg := errNode.Get(NodePredicateErrorMessage)
		errors = append(errors, compilercommon.NewSourceError(location, msg))
	}

	return errors
}
