// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/sourceshape"
)

var TYPE_KINDS = []sourceshape.NodeType{
	sourceshape.NodeTypeClass,
	sourceshape.NodeTypeInterface,
	sourceshape.NodeTypeNominal,
	sourceshape.NodeTypeStruct,
	sourceshape.NodeTypeAgent,
}

var TYPE_MEMBER_KINDS = []sourceshape.NodeType{
	sourceshape.NodeTypeFunction,
	sourceshape.NodeTypeVariable,
	sourceshape.NodeTypeConstructor,
	sourceshape.NodeTypeProperty,
	sourceshape.NodeTypeOperator,
}

var MEMBER_OR_TYPE_KINDS = append(TYPE_MEMBER_KINDS, TYPE_KINDS...)

var TYPE_KINDS_TAGGED = []compilergraph.TaggedValue{
	sourceshape.NodeTypeClass,
	sourceshape.NodeTypeInterface,
	sourceshape.NodeTypeNominal,
	sourceshape.NodeTypeStruct,
	sourceshape.NodeTypeAgent,
}

var MODULE_MEMBER_KINDS_TAGGED = append(TYPE_KINDS_TAGGED,
	[]compilergraph.TaggedValue{
		sourceshape.NodeTypeVariable,
		sourceshape.NodeTypeFunction,
	}...)
