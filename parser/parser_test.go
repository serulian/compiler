// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/packageloader"
	"github.com/stretchr/testify/assert"
)

type testNode struct {
	nodeType   NodeType
	properties map[string]string
	children   map[string]*list.List
}

type parserTest struct {
	name     string
	filename string
}

func (pt *parserTest) input() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s.seru", pt.filename))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (pt *parserTest) tree() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s.tree", pt.filename))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (pt *parserTest) writeTree(value string) {
	err := ioutil.WriteFile(fmt.Sprintf("tests/%s.tree", pt.filename), []byte(value), 0644)
	if err != nil {
		panic(err)
	}
}

func createAstNode(source compilercommon.InputSource, kind NodeType) AstNode {
	return &testNode{
		nodeType:   kind,
		properties: make(map[string]string),
		children:   make(map[string]*list.List),
	}
}

func (tn *testNode) GetType() NodeType {
	return tn.nodeType
}

func (tn *testNode) Connect(predicate string, other AstNode) AstNode {
	if tn.children[predicate] == nil {
		tn.children[predicate] = list.New()
	}

	tn.children[predicate].PushBack(other)
	return tn
}

func (tn *testNode) Decorate(property string, value string) AstNode {
	if _, ok := tn.properties[property]; ok {
		panic(fmt.Sprintf("Existing key for property %s\n\tNode: %v", property, tn.properties))
	}

	tn.properties[property] = value
	return tn
}

var parserTests = []parserTest{
	// Import success tests.
	{"basic import test", "import/basic"},
	{"from import test", "import/from"},
	{"string import test", "import/string_import"},
	{"multiple imports test", "import/multiple"},
	{"complex imports test", "import/complex"},

	// Import failure tests.
	{"missing source import test", "import/missing_source"},
	{"invalid source", "import/invalid_source"},
	{"invalid subsource", "import/invalid_subsource"},
	{"missing as for SCM test", "import/missing_as_scm"},

	// Module success tests.
	{"module variable test", "module/module_var"},
	{"module variable no declared type test", "module/module_var_nodeclare"},
	{"module function test", "module/module_function"},

	// Struct success tests.
	{"basic struct test", "struct/basic"},
	{"generic struct test", "struct/generic"},
	{"construct struct test", "struct/construct"},
	{"referenced struct test", "struct/referenced"},
	{"nullable field struct test", "struct/nullable"},
	{"tagged struct test", "struct/tagged"},
	{"slice struct test", "struct/slice"},

	// Class success tests.
	{"empty class test", "class/empty"},
	{"basic formatted class test", "class/basic_formatted"},
	{"basic inheritance test", "class/inherits"},
	{"complex inheritance test", "class/complex_inherits"},
	{"generic class test", "class/generic"},
	{"complex generic class test", "class/complex_generic"},
	{"complex class def test", "class/complex_def"},

	// Class failure tests.
	{"missing name class test", "class/missing_name"},
	{"missing subtype test", "class/missing_subtype"},
	{"missing another subtype test", "class/missing_another_subtype"},
	{"missing generic test", "class/missing_generic"},
	{"invalid inherits test", "class/invalidinherits"},

	// Interface success tests.
	{"basic interface test", "interface/basic"},
	{"generic interface test", "interface/generic"},

	// Nominal type tests.
	{"basic nominal type test", "nominal/basic"},
	{"generic nominal type test", "nominal/generic"},
	{"missing subtype nominal type test", "nominal/missingsubtype"},

	// Class member success tests.
	{"basic class function test", "class/basic_function"},
	{"generic function test", "class/generic_function"},
	{"basic class constructor test", "class/basic_constructor"},
	{"basic class property test", "class/basic_property"},
	{"readonly class property test", "class/readonly_property"},
	{"class function no parameters test", "class/function_noparams"},
	{"basic class operator test", "class/basic_operator"},
	{"basic class field test", "class/basic_field"},
	{"basic class operator with return type test", "class/operator_returns"},
	{"class constructor no params stest", "class/constructor_noparams"},

	// Interface member success tests.
	{"basic interface function test", "interface/basic_function"},
	{"basic interface constructor test", "interface/basic_constructor"},

	// Interface member failure tests.
	{"basic interface function failure test", "interface/basic_function_fail"},
	{"basic interface property test", "interface/basic_property"},
	{"basic interface operator test", "interface/basic_operator"},
	{"readonly interface property test", "interface/readonly_property"},

	// Statement tests.
	{"labeled statement test", "statement/labeled"},

	{"assignment statement test", "statement/assign"},

	{"no expr loop statement test", "statement/loop_noexpr"},
	{"expr loop statement test", "statement/loop_expr"},
	{"in expr loop statement test", "statement/loop_inexpr"},

	{"basic conditional statement test", "statement/conditional"},
	{"else conditional statement test", "statement/else_conditional"},
	{"chained conditional statement test", "statement/chained_conditional"},

	{"return no value statement test", "statement/return"},
	{"return value statement test", "statement/return_value"},
	{"reject value statement test", "statement/reject"},

	{"jump statements test", "statement/jump"},

	{"basic variable test", "statement/var"},
	{"variable no type success test", "statement/var_notype"},
	{"variable no type failure test", "statement/var_notype_failure"},

	{"basic with test", "statement/with"},
	{"with as test", "statement/with_as"},

	{"expression statement test", "statement/expression"},

	{"match statement basic test", "statement/match"},
	{"match statement no expr test", "statement/match_noexpr"},
	{"match statement multi statement test", "statement/match_multi"},

	// Expression tests.
	{"arrow expr test", "expression/arrow"},
	{"await expr test", "expression/await"},
	{"basic binary expr test", "expression/binary_basic"},
	{"binary unary expr test", "expression/binary_unary"},
	{"subtract binary expr test", "expression/binary_subtract"},
	{"multiple binary expr test", "expression/binary_multiple"},
	{"basic access expr test", "expression/access_basic"},
	{"multiple access expr test", "expression/access_multiple"},
	{"parens expr test", "expression/parens"},
	{"list expr test", "expression/list"},
	{"list call expr test", "expression/list_call"},
	{"lambda expr test", "expression/lambda"},
	{"map missing comma expr test", "expression/map_missingcomma"},
	{"generic specifier expr test", "expression/generic_specifier"},
	{"less than and generic specifier expr test", "expression/ltandgeneric"},
	{"template string expr test", "expression/templatestring"},
	{"range expr test", "expression/range"},
	{"assert not null expr test", "expression/assertnotnull"},
	{"is null expr test", "expression/isnull"},
	{"is null conditional expr test", "expression/isnullconditional"},
	{"slice literal expr test", "expression/sliceliteral"},
	{"mapping literal expr test", "expression/mappingliteral"},

	{"all expr test", "expression/all"},
	{"complex expr test", "expression/complex"},

	// Type reference tests.
	{"all type reference test", "typeref/all"},

	// Full example tests.
	{"basic full example test", "full/basic"},

	// Decorator tests.
	{"basic decorator test", "decorator/basic"},

	// Known issue tests.
	{"known issue test", "class/knownissue"},
	{"known issue 1 test", "knownissue/knownissue1"},
	{"known issue 2 test", "knownissue/knownissue2"},
	{"known issue 3 test", "knownissue/knownissue3"},
	{"known issue 4 test", "knownissue/knownissue4"},
	{"known issue 5 test", "knownissue/knownissue5"},
	{"known issue 6 test", "knownissue/knownissue6"},
	{"known issue 7 test", "knownissue/knownissue7"},
	{"known issue 8 test", "knownissue/knownissue8"},
	{"known issue 9 test", "knownissue/knownissue9"},
	{"known issue 10 test", "knownissue/knownissue10"},
	{"known issue 11 test", "knownissue/knownissue11"},
	{"known issue 12 test", "knownissue/knownissue12"},
	{"known issue 13 test", "knownissue/knownissue13"},

	// Comment tests.
	{"comment function test", "comment/function"},
	{"comment block test", "comment/block"},
	{"comment loop test", "comment/loop"},
	{"comment literal test", "comment/literal"},
	{"comment line test", "comment/line"},
	{"comment parens test", "comment/parens"},
	{"comment assign statement test", "comment/assignstatement"},
}

func reportImport(path packageloader.PackageImport) string {
	return "location:" + path.Path
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		if os.Getenv("FILTER") != "" {
			if !strings.Contains(test.name, os.Getenv("FILTER")) {
				continue
			} else {
				fmt.Printf("Matched Test: %v\n", test.name)
			}
		}

		rootNode := Parse(createAstNode, reportImport, compilercommon.InputSource(test.name), test.input())
		parseTree := getParseTree(t, (rootNode).(*testNode), 0)
		assert := assert.New(t)

		expected := strings.TrimSpace(test.tree())
		found := strings.TrimSpace(parseTree)

		if os.Getenv("REGEN") == "true" {
			test.writeTree(found)
		} else {
			if !assert.Equal(expected, found, test.name) {
				t.Log(parseTree)
			}
		}
	}
}

func getParseTree(t *testing.T, currentNode *testNode, indentation int) string {
	parseTree := ""
	parseTree = parseTree + strings.Repeat(" ", indentation)
	parseTree = parseTree + fmt.Sprintf("%v", currentNode.nodeType)
	parseTree = parseTree + "\n"

	keys := make([]string, 0)

	for key, _ := range currentNode.properties {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		parseTree = parseTree + strings.Repeat(" ", indentation+2)
		parseTree = parseTree + fmt.Sprintf("%s = %s", key, currentNode.properties[key])
		parseTree = parseTree + "\n"
	}

	keys = make([]string, 0)

	for key, _ := range currentNode.children {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		value := currentNode.children[key]
		parseTree = parseTree + fmt.Sprintf("%s%v =>", strings.Repeat(" ", indentation+2), key)
		parseTree = parseTree + "\n"

		for e := value.Front(); e != nil; e = e.Next() {
			if !assert.NotNil(t, e.Value, "Got nil value for predicate %v", key) {
				continue
			}

			parseTree = parseTree + getParseTree(t, e.Value.(*testNode), indentation+4)
		}
	}

	return parseTree
}
