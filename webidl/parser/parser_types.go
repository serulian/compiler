// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package parser

import (
	"strconv"
)

// NodeType identifies the type of AST node.
type NodeType int

const (
	// Top-level
	NodeTypeError   NodeType = iota // error occurred; value is text of error
	NodeTypeFile                    // The file root node
	NodeTypeComment                 // A single or multiline comment

	NodeTypeTagged
)

func (t NodeType) Name() string {
	return "IDLNodeType"
}

func (t NodeType) Value() string {
	return strconv.Itoa(int(t))
}

func (t NodeType) Build(value string) interface{} {
	i, err := strconv.Atoi(value)
	if err != nil {
		panic("Invalid value for IDLNodeType: " + value)
	}
	return NodeType(i)
}
