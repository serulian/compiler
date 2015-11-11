// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

//go:generate stringer -type=NodeType

import (
	"strconv"
)

// NodeType identifies the type of scope graph nodes.
type NodeType int

const (
	// Top-level
	NodeTypeError         NodeType = iota // A scope error
	NodeTypeWarning                       // A scope warning
	NodeTypeResolvedScope                 // Resolved scope for an SRG node

	// NodeType is a tagged type.
	NodeTypeTagged
)

const (
	// Connects a scope node to its SRG source.
	NodePredicateSource = "scope-source"

	// Decorates a scope node with its scope info.
	NodePredicateScopeInfo = "scope-info"

	// Connects an error or warning to its SRG source.
	NodePredicateNoticeSource = "scope-notice"

	// The error or warning message on a scope notice node.
	NodePredicateNoticeMessage = "notice-message"

	// Decorates a *SRG* node with its inferred type.
	NodePredicateInferredType = "scope-inferred-type"
)

func (t NodeType) Name() string {
	return "NodeType"
}

func (t NodeType) Value() string {
	return strconv.Itoa(int(t))
}

func (t NodeType) Build(value string) interface{} {
	i, err := strconv.Atoi(value)
	if err != nil {
		panic("Invalid value for NodeType: " + value)
	}
	return NodeType(i)
}
