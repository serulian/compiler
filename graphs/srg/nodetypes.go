// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

var TYPE_KINDS = []parser.NodeType{
	parser.NodeTypeClass,
	parser.NodeTypeInterface,
	parser.NodeTypeNominal,
	parser.NodeTypeStruct,
	parser.NodeTypeAgent,
}

var TYPE_MEMBER_KINDS = []parser.NodeType{
	parser.NodeTypeFunction,
	parser.NodeTypeVariable,
	parser.NodeTypeConstructor,
	parser.NodeTypeProperty,
	parser.NodeTypeOperator,
}

var MEMBER_OR_TYPE_KINDS = append(TYPE_MEMBER_KINDS, TYPE_KINDS...)

var TYPE_KINDS_TAGGED = []compilergraph.TaggedValue{
	parser.NodeTypeClass,
	parser.NodeTypeInterface,
	parser.NodeTypeNominal,
	parser.NodeTypeStruct,
	parser.NodeTypeAgent,
}

var MODULE_MEMBER_KINDS_TAGGED = append(TYPE_KINDS_TAGGED,
	[]compilergraph.TaggedValue{
		parser.NodeTypeVariable,
		parser.NodeTypeFunction,
	}...)
