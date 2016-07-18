// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// Based on stack found here: https://gist.github.com/bemasher/1777766

// resource wraps information about a resource ('with') added into the scope.
type resource struct {
	name  string
	basis compilergraph.GraphNode
}

func (r resource) Name() string {
	return r.name
}

type ResourceStack struct {
	top  *ResourceElement
	size int
}

type ResourceElement struct {
	value resource
	next  *ResourceElement
}

// Return the stack's length
func (s *ResourceStack) Len() int {
	return s.size
}

// Push a new element onto the stack
func (s *ResourceStack) Push(value resource) {
	s.top = &ResourceElement{value, s.top}
	s.size++
}

// Remove the top element from the stack and return it's value
// If the stack is empty, return nil
func (s *ResourceStack) Pop() (value resource) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return resource{}
}

// OutOfScope returns any resources that will be out of scope when context changes to the given reference
// node. Used to determine which resources need to be popped off of the stack when jumps occurr. Note that
// the reference node *must be an SRG node*.
func (s *ResourceStack) OutOfScope(referenceNode compilergraph.GraphNode) []resource {
	resources := make([]resource, 0)

	// Find the start and end location for the reference node.
	referenceStart := referenceNode.GetValue(parser.NodePredicateStartRune).Int()
	referenceEnd := referenceNode.GetValue(parser.NodePredicateEndRune).Int()

	var current = s.top
	for index := 0; index < s.size; index++ {
		currentResource := current.value

		// Find the start and end location for the resource node.
		resourceStart := currentResource.basis.GetValue(parser.NodePredicateStartRune).Int()
		resourceEnd := currentResource.basis.GetValue(parser.NodePredicateEndRune).Int()

		if !(referenceStart >= resourceStart && referenceEnd <= resourceEnd) {
			resources = append(resources, currentResource)
		}

		current = current.next
	}

	return resources
}
