// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Stack copied from https://gist.github.com/bemasher/1777766

package parser

type nodeStack struct {
	top  *element
	size int
}

type element struct {
	value AstNode
	next  *element
}

func (s *nodeStack) topValue() AstNode {
	if s.size == 0 {
		return nil
	}

	return s.top.value
}

// Push pushes a node onto the stack.
func (s *nodeStack) push(value AstNode) {
	s.top = &element{value, s.top}
	s.size++
}

// Pop removes the node from the stack and returns it.
func (s *nodeStack) pop() (value AstNode) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}
