// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Stack copied from https://gist.github.com/bemasher/1777766

package grok

type nestingStack struct {
	top  *nestingElement
	size int
}

type nestingElement struct {
	openRune   rune
	startIndex int
	next       *nestingElement
}

func (s *nestingStack) topValue() (rune, int, bool) {
	if s.size == 0 {
		return ' ', -1, false
	}

	return s.top.openRune, s.top.startIndex, true
}

// Push pushes an opening rune onto the stack.
func (s *nestingStack) push(openRune rune, startIndex int) {
	s.top = &nestingElement{openRune, startIndex, s.top}
	s.size++
}

// Pop removes the nesting rune from the stack and returns it.
func (s *nestingStack) pop() (openRune rune, startIndex int) {
	if s.size > 0 {
		openRune, startIndex, s.top = s.top.openRune, s.top.startIndex, s.top.next
		s.size--
		return
	}
	return ' ', -1
}
