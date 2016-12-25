// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/compilergraph"
)

// Based on stack found here: https://gist.github.com/bemasher/1777766

// resource wraps information about a managed resource  added into the scope.
type resource struct {
	name       string
	basis      compilergraph.GraphNode
	startState *state
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

// OutOfScope returns the names of any resource that will be out of scope when context changes to the given state.
func (s *ResourceStack) OutOfScope(stateId stateId) []string {
	targetStateId := int(stateId)
	resourceNames := make([]string, 0, s.size)

	var current = s.top
	for index := 0; index < s.size; index++ {
		currentResource := current.value
		startStateId := int(currentResource.startState.ID)
		if targetStateId <= startStateId {
			resourceNames = append(resourceNames, currentResource.name)
		}

		current = current.next
	}

	return resourceNames
}
