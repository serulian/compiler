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
	"strconv"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser/shared"
	"github.com/serulian/compiler/sourceshape"

	"github.com/stretchr/testify/assert"
)

type testNode struct {
	nodeType   sourceshape.NodeType
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

func createAstNode(source compilercommon.InputSource, kind sourceshape.NodeType) shared.AstNode {
	return &testNode{
		nodeType:   kind,
		properties: make(map[string]string),
		children:   make(map[string]*list.List),
	}
}

func (tn *testNode) GetType() sourceshape.NodeType {
	return tn.nodeType
}

func (tn *testNode) Connect(predicate string, other shared.AstNode) shared.AstNode {
	if tn.children[predicate] == nil {
		tn.children[predicate] = list.New()
	}

	tn.children[predicate].PushBack(other)
	return tn
}

func (tn *testNode) Decorate(property string, value string) shared.AstNode {
	if _, ok := tn.properties[property]; ok {
		panic(fmt.Sprintf("Existing key for property %s\n\tNode: %v", property, tn.properties))
	}

	tn.properties[property] = value
	return tn
}

func (tn *testNode) DecorateWithInt(property string, value int) shared.AstNode {
	return tn.Decorate(property, strconv.Itoa(value))
}

var parserTests = []parserTest{
	// Import success tests.
	{"basic import test", "import/basic"},
	{"from import test", "import/from"},
	{"string import test", "import/string_import"},
	{"multiple imports test", "import/multiple"},
	{"complex imports test", "import/complex"},
	{"multiple imports under package test", "import/multiple_package"},
	{"relative import test", "import/relative"},
	{"relative import test 2", "import/relative2"},
	{"absolute import test", "import/absolute"},
	{"empty import test", "import/empty"},
	{"alias import test", "import/aliases"},

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
	{"struct default test", "struct/default"},

	// Class success tests.
	{"empty class test", "class/empty"},
	{"basic formatted class test", "class/basic_formatted"},

	{"basic agents test", "class/basic_agents"},
	{"aliased agents test", "class/aliased_agents"},
	{"missing agent test", "class/missing_agent"},

	{"generic class test", "class/generic"},
	{"complex generic class test", "class/complex_generic"},

	// Class failure tests.
	{"missing name class test", "class/missing_name"},
	{"missing generic test", "class/missing_generic"},

	// Interface success tests.
	{"basic interface test", "interface/basic"},
	{"generic interface test", "interface/generic"},

	// Nominal type tests.
	{"basic nominal type test", "nominal/basic"},
	{"generic nominal type test", "nominal/generic"},
	{"missing subtype nominal type test", "nominal/missingsubtype"},

	// Agent tests.
	{"basic agent test", "agent/basic"},
	{"generic agent test", "agent/generic"},
	{"missing principal type agent test", "agent/missingprincipal"},

	// Class member success tests.
	{"basic class function test", "class/basic_function"},
	{"generic function test", "class/generic_function"},
	{"basic class constructor test", "class/basic_constructor"},
	{"basic class property test", "class/basic_property"},
	{"readonly class property test", "class/readonly_property"},
	{"class function no parameters test", "class/function_noparams"},
	{"basic class operator test", "class/basic_operator"},
	{"keyword class operator test", "class/keyword_operator"},
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
	{"compressed statement test", "statement/compressed"},

	{"labeled statement test", "statement/labeled"},

	{"assignment statement test", "statement/assign"},
	{"resolve statement test", "statement/resolve"},

	{"no expr loop statement test", "statement/loop_noexpr"},
	{"expr loop statement test", "statement/loop_expr"},
	{"in expr loop statement test", "statement/loop_inexpr"},

	{"basic conditional statement test", "statement/conditional"},
	{"else conditional statement test", "statement/else_conditional"},
	{"chained conditional statement test", "statement/chained_conditional"},

	{"yield value statement test", "statement/yield_value"},
	{"yield in statement test", "statement/yield_in"},
	{"yield break statement test", "statement/yield_break"},

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

	{"switch statement basic test", "statement/switch"},
	{"switch statement no expr test", "statement/switch_noexpr"},
	{"switch statement multi statement test", "statement/switch_multi"},

	{"match statement basic test", "statement/match"},
	{"match statement as test", "statement/match_as"},

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
	{"map literal expr test", "expression/mapliteral"},
	{"conditional expr test", "expression/conditional"},
	{"complex conditional expr test", "expression/complexconditional"},
	{"loop expr test", "expression/loop"},
	{"slice literal slice expr test", "expression/sliceliteralslice"},
	{"numeric literal expr test", "expression/numericliteral"},
	{"not precedence expr test", "expression/notprecedence"},

	// Serulian Markup Language expression tests.
	{"tag only sml test", "sml/tagonly"},
	{"open and close sml test", "sml/openclose"},
	{"missing close sml test", "sml/missingclose"},
	{"multiline sml test", "sml/multiline"},
	{"mismatch sml test", "sml/mismatch"},
	{"attributes sml test", "sml/attributes"},
	{"keyword attribute sml test", "sml/keywordattribute"},
	{"invalid attribute value sml test", "sml/invalidattributevalue"},
	{"children sml test", "sml/children"},
	{"multiple children sml test", "sml/multichild"},
	{"whitespace text sml test", "sml/whitespacetext"},
	{"nested property sml test", "sml/nestedprop"},
	{"loop sml test", "sml/loop"},
	{"multiline loop sml test", "sml/multilineloop"},
	{"inline loop sml test", "sml/inlineloop"},

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
	{"known issue 14 test", "knownissue/knownissue14"},

	// Comment tests.
	{"comment function test", "comment/function"},
	{"comment block test", "comment/block"},
	{"comment loop test", "comment/loop"},
	{"comment literal test", "comment/literal"},
	{"comment line test", "comment/line"},
	{"comment parens test", "comment/parens"},
	{"comment assign statement test", "comment/assignstatement"},
	{"comment keywords test", "comment/keywords"},
	{"comment keyword op test", "comment/keywordop"},

	// Partial tests.
	{"partial function call test", "partial/funccall"},
	{"partial member access test", "partial/memberaccess"},
	{"partial type definition test", "partial/typedef"},
	{"partial sml definition test", "partial/sml"},
	{"partial type member test", "partial/typemember"},
	{"partial type member test 2", "partial/typemember2"},
	{"partial type ref test", "partial/typeref"},
}

func reportImport(sourceKind string, importPath string, importType packageloader.PackageImportType, importSource compilercommon.InputSource, runePosition int) string {
	return "location:" + importPath
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

		rootNode, _ := Parse(createAstNode, reportImport, compilercommon.InputSource(test.name), test.input())
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

type parseExpressionTest struct {
	expressionString string
	expectedNodeType sourceshape.NodeType
	isOkay           bool
}

var parseExpressionTests = []parseExpressionTest{
	// Success tests.
	parseExpressionTest{"this", sourceshape.NodeThisLiteralExpression, true},
	parseExpressionTest{"principal", sourceshape.NodePrincipalLiteralExpression, true},
	parseExpressionTest{"a.b", sourceshape.NodeMemberAccessExpression, true},
	parseExpressionTest{"a.b.c", sourceshape.NodeMemberAccessExpression, true},
	parseExpressionTest{"a.b.c.d", sourceshape.NodeMemberAccessExpression, true},
	parseExpressionTest{"a[1]", sourceshape.NodeSliceExpression, true},
	parseExpressionTest{"a == b", sourceshape.NodeComparisonEqualsExpression, true},

	// Failure tests.
	parseExpressionTest{"a +", sourceshape.NodeTypeError, false},
	parseExpressionTest{"a ++", sourceshape.NodeTypeError, false},
	parseExpressionTest{"a....b", sourceshape.NodeTypeError, false},
}

func TestParseExpression(t *testing.T) {
	for index, test := range parseExpressionTests {
		parsed, ok := ParseExpression(createAstNode, compilercommon.InputSource(test.expressionString), index+10, test.expressionString)
		if !assert.Equal(t, test.isOkay, ok, "Mismatch in success expected for parsing test: %s", test.expressionString) {
			continue
		}

		if !test.isOkay {
			continue
		}

		testNode := parsed.(*testNode)
		if !assert.Equal(t, test.expectedNodeType, testNode.nodeType, "Mismatch in node type found for parsing test: %s", test.expressionString) {
			continue
		}

		if !assert.Equal(t, fmt.Sprintf("%v", index+10), testNode.properties[sourceshape.NodePredicateStartRune], "Mismatch in start rune") {
			continue
		}
	}
}
