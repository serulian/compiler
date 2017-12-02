// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Stack copied from https://gist.github.com/bemasher/1777766

package parser

import (
	"github.com/serulian/compiler/parser/shared"
)

type nodeStack struct {
	top  *nodeElement
	size int
}

type nodeElement struct {
	value shared.AstNode
	next  *nodeElement
}

func (s *nodeStack) topValue() shared.AstNode {
	if s.size == 0 {
		return nil
	}

	return s.top.value
}

// Push pushes a node onto the stack.
func (s *nodeStack) push(value shared.AstNode) {
	s.top = &nodeElement{value, s.top}
	s.size++
}

// Pop removes the node from the stack and returns it.
func (s *nodeStack) pop() (value shared.AstNode) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}
