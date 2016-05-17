// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esbuilder

import (
	"testing"

	"github.com/serulian/compiler/sourcemap"
	"github.com/stretchr/testify/assert"
)

type expectedMapping struct {
	lineNumber  int
	colPosition int
	mapping     sourcemap.SourceMapping
}

type builderTest struct {
	name             string
	node             ExpressionOrStatementBuilder
	expectedCode     string
	expectedMappings []expectedMapping
}

var generationTests = []builderTest{
	builderTest{"identfier", Identifier("something"), "something", []expectedMapping{}},
	builderTest{"access",
		Identifier("something").
			Member("somemember"),
		"(something).somemember", []expectedMapping{}},

	builderTest{"expr access",
		Identifier("something").
			Member(Value("somemember")),
		"(something)[\"somemember\"]", []expectedMapping{}},

	builderTest{"basic call", Identifier("something").Call(), "(something)()", []expectedMapping{}},
	builderTest{"basic call with args", Identifier("something").Call(1, 2, 3), "(something)(1,2,3)", []expectedMapping{}},
	builderTest{"basic call with expr args", Identifier("something").Call(Identifier("foo")), "(something)(foo)", []expectedMapping{}},

	builderTest{"prefix expr",
		Prefix("!", Identifier("something")),
		"!(something)", []expectedMapping{}},

	builderTest{"postfix expr",
		Postfix(Identifier("something"), "!"),
		"(something)!", []expectedMapping{}},

	builderTest{"binary expr",
		Binary(Identifier("something"), "+", Identifier("somethingElse")),
		"((something)+(somethingElse))", []expectedMapping{}},

	builderTest{"basic function",
		Function("somefunction", Identifier("someexpr"), "foo", "bar", "baz"),

		`function somefunction(foo,bar,baz){
  someexpr
}`, []expectedMapping{}},

	builderTest{"basic closure",
		Closure(Identifier("someexpr"), "foo", "bar", "baz"),

		`(function (foo,bar,baz){
  someexpr
})`, []expectedMapping{}},

	builderTest{"function with return",
		Closure(Return()),

		`(function (){
  return;
})`, []expectedMapping{}},

	builderTest{"function with returns",
		Closure(Returns(Identifier("something"))),

		`(function (){
  return (something);
})`, []expectedMapping{}},

	builderTest{"function with conditional",
		Closure(
			IfElse(Value(true),
				Statements(Returns(Value(1))),
				Statements(Returns(Value(2)))),
		),

		`(function (){
  if (true) {
    return (1);
  } else {
    return (2);
  }
})`, []expectedMapping{}},

	builderTest{"simple template",
		Template("simple", "{{ .CoolThing }}", map[string]interface{}{
			"CoolThing": "coolthing",
		}),

		"coolthing", []expectedMapping{}},

	builderTest{"expr template",
		Template("expr", "some {{ emit .CoolThing }}", map[string]interface{}{
			"CoolThing": Identifier("coolthing"),
		}),

		"some coolthing", []expectedMapping{}},

	builderTest{"identfier with mapping",
		Identifier("something").WithMapping(sourcemap.SourceMapping{"foo.js", 10, 5, ""}),

		"something",

		[]expectedMapping{
			expectedMapping{0, 0, sourcemap.SourceMapping{"foo.js", 10, 5, ""}},
		}},

	builderTest{"prefix expr with mapping",
		Prefix("!",
			Identifier("something").WithMapping(sourcemap.SourceMapping{"foo.js", 10, 5, ""})).
			WithMapping(sourcemap.SourceMapping{"bar.js", 10, 2, ""}),

		"!(something)",

		[]expectedMapping{
			expectedMapping{0, 0, sourcemap.SourceMapping{"bar.js", 10, 2, ""}},
			expectedMapping{0, 2, sourcemap.SourceMapping{"foo.js", 10, 5, ""}},
		}},

	builderTest{"expr template with mapping",
		Template("expr", "some {{ emit .CoolThing }}", map[string]interface{}{
			"CoolThing": Identifier("coolthing").WithMapping(sourcemap.SourceMapping{"foo.js", 10, 5, ""}),
		}),

		"some coolthing",

		[]expectedMapping{
			expectedMapping{0, 5, sourcemap.SourceMapping{"foo.js", 10, 5, ""}},
		}},

	builderTest{"expr template with mapping and newlines",
		Template("expr", "some \ncooler {{ emit .CoolThing }}", map[string]interface{}{
			"CoolThing": Identifier("coolthing").WithMapping(sourcemap.SourceMapping{"foo.js", 10, 5, ""}),
		}),

		"some \ncooler coolthing",

		[]expectedMapping{
			expectedMapping{1, 7, sourcemap.SourceMapping{"foo.js", 10, 5, ""}},
		}},

	builderTest{"multiple expr template with mapping",
		Template("expr",
			`function () {
  if ({{ emit .FirstThing }}) {
    {{ emit .SecondThing }};
  }
}`,

			map[string]interface{}{
				"FirstThing":  Value(true).WithMapping(sourcemap.SourceMapping{"foo.js", 10, 5, ""}),
				"SecondThing": Value("somevalue").WithMapping(sourcemap.SourceMapping{"bar.js", 40, 2, ""}),
			}),

		`function () {
  if (true) {
    "somevalue";
  }
}`,

		[]expectedMapping{
			expectedMapping{1, 6, sourcemap.SourceMapping{"foo.js", 10, 5, ""}},
			expectedMapping{2, 4, sourcemap.SourceMapping{"bar.js", 40, 2, ""}},
		}},
}

func TestGeneration(t *testing.T) {
	for _, test := range generationTests {
		generated := buildSource(test.node)
		if !assert.Equal(t, generated.buf.String(), test.expectedCode, "Mismatch on generation test %s", test.name) {
			continue
		}

		built := generated.sourcemap.Build()
		for _, expectedMapping := range test.expectedMappings {
			mapping, ok := built.LookupMapping(expectedMapping.lineNumber, expectedMapping.colPosition)
			if !assert.True(t, ok, "Expected mapping %v on test %s", expectedMapping, test.name) {
				continue
			}

			if !assert.Equal(t, expectedMapping.mapping, mapping, "Mapping mismatch on test %s", test.name) {
				continue
			}
		}
	}
}
