// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"runtime"
	"testing"

	"github.com/serulian/compiler/compilerutil"
	"github.com/stretchr/testify/assert"
)

type typePrefixTest struct {
	snippet  string
	isPrefix bool
}

var typePrefixTests = []typePrefixTest{
	typePrefixTest{"class DoSomething<T: ", true},
	typePrefixTest{"interface DoSomething<T: ", true},

	typePrefixTest{"type SomeType : ", true},

	typePrefixTest{"agent SomeAgent for ", true},
	typePrefixTest{"agent SomeAgent for a", false},

	typePrefixTest{"class DoSomething with ", true},
	typePrefixTest{"class DoSomething with a + ", true},
	typePrefixTest{"class DoSomething with a + b", false},

	typePrefixTest{"function DoSomething(something int)", true},
	typePrefixTest{"function DoSomething(something int, ", false},
	typePrefixTest{"function DoSomething(something int, a ", true},

	typePrefixTest{"function DoSomething<T>(something int)", true},
	typePrefixTest{"function DoSomething<T>(something int, ", false},
	typePrefixTest{"function DoSomething<T>(something int, a ", true},

	typePrefixTest{"operator DoSomething(something int)", true},
	typePrefixTest{"operator DoSomething(something int, ", false},
	typePrefixTest{"operator DoSomething(something int, a ", true},

	typePrefixTest{"function DoSomething<T:", true},
	typePrefixTest{"function DoSomething<T: ", true},
	typePrefixTest{"function DoSomething<Q, T: ", true},
	typePrefixTest{"function DoSomething<Q, T: s", false},

	typePrefixTest{"property Something ", true},
	typePrefixTest{"var Something ", true},

	typePrefixTest{"foo.(", true},
	typePrefixTest{"foo.(a", false},

	typePrefixTest{"[]", true},
	typePrefixTest{"[]{", true},
}

func TestIsTypePrefix(t *testing.T) {
	for _, test := range typePrefixTests {
		t.Run(test.snippet, func(t *testing.T) {
			defer compilerutil.DetectGoroutineLeak(t, runtime.NumGoroutine())
			assert.Equal(t, test.isPrefix, IsTypePrefix(test.snippet))
		})
	}
}
