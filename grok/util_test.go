// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type highlightParameterTest struct {
	documentation string
	parameterName string
	expected      string
}

var parameterTests = []highlightParameterTest{
	highlightParameterTest{
		"This functions is invoked with `param1`, `param2` and `param3`",
		"param1",
		"This functions is invoked with `param1`, **param2** and **param3**",
	},

	highlightParameterTest{
		"This functions is invoked with `param1`, `param2` and `param3`",
		"param2",
		"This functions is invoked with **param1**, `param2` and **param3**",
	},

	highlightParameterTest{
		"This functions is invoked with `param1`, `param2` and `param3`",
		"param3",
		"This functions is invoked with **param1**, **param2** and `param3`",
	},

	highlightParameterTest{
		"DoSomething does something really cool with `p`. If `p` is not available, then we use `q`!",
		"p",
		"DoSomething does something really cool with `p`. If `p` is not available, then we use **q**!",
	},

	highlightParameterTest{
		"DoSomething does something really cool with `p`. If `p` is not available, then we use `q`!",
		"q",
		"DoSomething does something really cool with **p**. If **p** is not available, then we use `q`!",
	},
}

func TestHighlightParameter(t *testing.T) {
	for _, parameterTest := range parameterTests {
		assert.Equal(t, parameterTest.expected, highlightParameter(parameterTest.documentation, parameterTest.parameterName))
	}
}
