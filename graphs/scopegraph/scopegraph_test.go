// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/packageloader"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

const TESTLIB_PATH = "../../testlib"

type expectedScope struct {
	IsValid      bool
	ScopeKind    proto.ScopeKind
	ResolvedType string
	ReturnedType string
}

type expectedScopeEntry struct {
	name  string
	scope expectedScope
}

type scopegraphTest struct {
	name            string
	input           string
	entrypoint      string
	expectedScope   []expectedScopeEntry
	expectedError   string
	expectedWarning string
}

func (sgt *scopegraphTest) json() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s/%s.json", sgt.input, sgt.entrypoint))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (sgt *scopegraphTest) writeJson(value string) {
	err := ioutil.WriteFile(fmt.Sprintf("tests/%s/%s.json", sgt.input, sgt.entrypoint), []byte(value), 0644)
	if err != nil {
		panic(err)
	}
}

var scopeGraphTests = []scopegraphTest{
	/////////// Empty block ///////////

	// Empty block on void function.
	scopegraphTest{"empty block void test", "empty", "block", []expectedScopeEntry{
		expectedScopeEntry{"emptyblock", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
	}, "", ""},

	// Unreachable statement.
	scopegraphTest{"unreachable statement test", "unreachable", "return", []expectedScopeEntry{},
		"", "Unreachable statement found"},

	// Empty block on int function.
	scopegraphTest{"empty block int test", "empty", "missingreturn", []expectedScopeEntry{},
		"Expected return value of type 'Integer' but not all paths return a value", ""},

	/////////// Settling (return and reject) ///////////

	// Success test.
	scopegraphTest{"settlement statements success test", "settlement", "success", []expectedScopeEntry{},
		"", ""},

	// Missing branch test.
	scopegraphTest{"settlement missing branch test", "settlement", "missingbranch", []expectedScopeEntry{},
		"Expected return value of type 'Integer' but not all paths return a value", ""},

	// Invalid reject value test.
	scopegraphTest{"settlement invalid reject test", "settlement", "invalidreject", []expectedScopeEntry{},
		"'reject' statement value must be an Error: Type 'String' does not define or export member 'Message', which is required by type 'Error'", ""},

	/////////// Break ///////////

	// Normal break statement.
	scopegraphTest{"break statement test", "break", "normal", []expectedScopeEntry{},
		"", "Unreachable statement found"},

	// break statement not under a breakable node.
	scopegraphTest{"break statement error test", "break", "badparent", []expectedScopeEntry{},
		"'break' statement must be a under a loop or match statement", ""},

	/////////// Continue ///////////

	// Normal continue statement.
	scopegraphTest{"continue statement test", "continue", "normal", []expectedScopeEntry{},
		"", "Unreachable statement found"},

	// continue statement not under a breakable node.
	scopegraphTest{"continue statement error test", "continue", "badparent", []expectedScopeEntry{},
		"'continue' statement must be a under a loop statement", ""},

	/////////// Conditionals ///////////

	scopegraphTest{"basic conditional test", "conditional", "basic",
		[]expectedScopeEntry{
			expectedScopeEntry{"conditional", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"conditional with else test", "conditional", "else",
		[]expectedScopeEntry{
			expectedScopeEntry{"conditional", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
			expectedScopeEntry{"trueblock", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
			expectedScopeEntry{"falseblock", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
		},
		"", ""},

	scopegraphTest{"chained conditional test", "conditional", "chained",
		[]expectedScopeEntry{
			expectedScopeEntry{"conditional", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
			expectedScopeEntry{"true1block", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
			expectedScopeEntry{"true2block", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
			expectedScopeEntry{"falseblock", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
		},
		"", ""},

	scopegraphTest{"conditional invalid expr test", "conditional", "nonboolexpr", []expectedScopeEntry{},
		"Conditional expression must be of type 'bool', found: Integer", ""},

	scopegraphTest{"conditional return test", "conditional", "return", []expectedScopeEntry{},
		"Expected return value of type 'Integer' but not all paths return a value", ""},

	scopegraphTest{"conditional chained return intersect test", "conditional", "chainedintersect", []expectedScopeEntry{},
		"Expected return value of type 'Integer': Cannot use type 'any' in place of type 'Integer'", ""},

	/////////// Loops ///////////

	// Empty loop test.
	scopegraphTest{"empty loop test", "loop", "empty", []expectedScopeEntry{
		expectedScopeEntry{"emptyloop", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
	}, "", "Unreachable statement found"},

	// Boolean loop test.
	scopegraphTest{"bool loop test", "loop", "boolloop", []expectedScopeEntry{
		expectedScopeEntry{"loopexpr", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		expectedScopeEntry{"loop", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
	}, "", ""},

	// Stream loop test.
	scopegraphTest{"stream loop test", "loop", "streamloop", []expectedScopeEntry{
		expectedScopeEntry{"loop", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		expectedScopeEntry{"a", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
	}, "", ""},

	// Streamable loop test.
	scopegraphTest{"streamable loop test", "loop", "streamableloop", []expectedScopeEntry{
		expectedScopeEntry{"loop", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		expectedScopeEntry{"a", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
	}, "", ""},

	// Expected bool loop test.
	scopegraphTest{"expected bool loop test", "loop", "expectedboolloop", []expectedScopeEntry{},
		"Loop conditional expression must be of type 'bool', found: Integer", ""},

	// Expected stream loop test.
	scopegraphTest{"expected stream loop test", "loop", "expectedstreamloop", []expectedScopeEntry{},
		"Loop iterable expression must implement type 'stream' or 'streamable': Type Integer cannot be used in place of type Stream as it does not implement member Next", ""},

	/////////// With ///////////

	// Basic with test.
	scopegraphTest{"basic with test", "with", "basic", []expectedScopeEntry{
		expectedScopeEntry{"with", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
	}, "", ""},

	scopegraphTest{"invalid with expr test", "with", "nonreleasable", []expectedScopeEntry{},
		"With expression must implement the Releasable interface: Type 'Boolean' does not define or export member 'Release', which is required by type 'Releasable'", ""},

	/////////// Match ///////////

	scopegraphTest{"basic bool match test", "match", "bool",
		[]expectedScopeEntry{
			expectedScopeEntry{"match", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"bool default match test", "match", "booldefault",
		[]expectedScopeEntry{
			expectedScopeEntry{"match", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
		},
		"", ""},

	scopegraphTest{"basic non-bool match test", "match", "typed",
		[]expectedScopeEntry{
			expectedScopeEntry{"match", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"non-bool default match test", "match", "typeddefault",
		[]expectedScopeEntry{
			expectedScopeEntry{"match", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
		},
		"", ""},

	scopegraphTest{"invalid bool match test", "match", "boolinvalid", []expectedScopeEntry{},
		"Match cases must have values matching type 'Boolean': 'Integer' cannot be used in place of non-interface 'Boolean'", ""},

	scopegraphTest{"invalid non-bool match test", "match", "typedinvalid", []expectedScopeEntry{},
		"Match cases must have values matching type 'Integer': 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"match no equals operator test", "match", "nocompare", []expectedScopeEntry{},
		"Cannot match over instance of type 'SomeClass', as it does not define or export an 'equals' operator", ""},

	/////////// Var ///////////

	scopegraphTest{"basic var test", "var", "basic",
		[]expectedScopeEntry{
			expectedScopeEntry{"var", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"no declared type var test", "var", "notype",
		[]expectedScopeEntry{
			expectedScopeEntry{"var", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
			expectedScopeEntry{"number", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"initialized var test", "var", "initialized",
		[]expectedScopeEntry{
			expectedScopeEntry{"var", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
			expectedScopeEntry{"number", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	/////////// Assign ///////////

	scopegraphTest{"assign unknown name test", "assign", "unknown", []expectedScopeEntry{},
		"The name 'something' could not be found in this context", ""},

	scopegraphTest{"assign readonly test", "assign", "readonly", []expectedScopeEntry{},
		"Cannot assign to non-assignable module member DoSomething", ""},

	scopegraphTest{"assign import test", "assign", "import", []expectedScopeEntry{},
		"Cannot assign to non-assignable import success", ""},

	scopegraphTest{"assign type mismatch test", "assign", "typemismatch", []expectedScopeEntry{},
		"Cannot assign value to variable something: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"assign success test", "assign", "success", []expectedScopeEntry{},
		"", ""},

	scopegraphTest{"assign indexer test", "assign", "indexer",
		[]expectedScopeEntry{},
		"", ""},

	scopegraphTest{"assign indexer value type mismatch test", "assign", "indexermismatch",
		[]expectedScopeEntry{},
		"Cannot assign value to operator setindex: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	/////////// Comparison operator expressions ///////////

	scopegraphTest{"comparison op success test", "compareops", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"equals", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"notequals", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},

			expectedScopeEntry{"lt", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"lte", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"gt", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"gte", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	/////////// Root type operator expression ///////////

	scopegraphTest{"root type op success test", "rootop", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"someclass", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"generic", expectedScope{true, proto.ScopeKind_VALUE, "any", "void"}},
			expectedScopeEntry{"interface", expectedScope{true, proto.ScopeKind_VALUE, "any", "void"}},
		},
		"", ""},

	scopegraphTest{"root type op class failure test", "rootop", "classfail",
		[]expectedScopeEntry{},
		"Root type operator (&) cannot be applied to value of type SomeClass", ""},

	/////////// Is operator expression ///////////

	scopegraphTest{"is op success test", "isop", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"isresult", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"is op not null failure test", "isop", "notnull",
		[]expectedScopeEntry{},
		"Right side of 'is' operator must be 'null'", ""},

	scopegraphTest{"is op not nullable failure test", "isop", "notnullable",
		[]expectedScopeEntry{},
		"Left side of 'is' operator must be a nullable type. Found: Integer", ""},

	/////////// Nullable operator expression ///////////

	scopegraphTest{"nullable ops success test", "nullableops", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"cnullclass", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"cnullinter", expectedScope{true, proto.ScopeKind_VALUE, "ISomeInterface", "void"}},

			expectedScopeEntry{"inullclass", expectedScope{true, proto.ScopeKind_VALUE, "ISomeInterface", "void"}},
			expectedScopeEntry{"inullinter", expectedScope{true, proto.ScopeKind_VALUE, "ISomeInterface", "void"}},
		},
		"", ""},

	scopegraphTest{"nullable ops non-nullable left test", "nullableops", "nonnullable",
		[]expectedScopeEntry{},
		"Left hand side of a nullable operator must be nullable. Found: SomeClass", ""},

	scopegraphTest{"nullable ops nullable right test", "nullableops", "nullableright",
		[]expectedScopeEntry{},
		"Right hand side of a nullable operator cannot be nullable. Found: SomeClass?", ""},

	scopegraphTest{"nullable ops non-subtype test", "nullableops", "subtypemismatch",
		[]expectedScopeEntry{},
		"Left and right hand sides of a nullable operator must have common subtype. None found between 'Integer' and 'Boolean'", ""},

	/////////// Binary operator expressions ///////////

	scopegraphTest{"binary ops success test", "binaryops", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"xor", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"or", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"and", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"shiftleft", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},

			expectedScopeEntry{"plus", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"minus", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"times", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"div", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"mod", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},

			expectedScopeEntry{"bor", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"band", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"boolean ops fail test", "binaryops", "boolfail",
		[]expectedScopeEntry{},
		"Boolean operator requires type Boolean for operands. Left hand operand has type: Integer", ""},

	/////////// Identifier expression ///////////

	scopegraphTest{"identifier expr invalid anonymous test", "identexpr", "anonymous", []expectedScopeEntry{},
		"Anonymous identifier '_' cannot be used as a value", ""},

	scopegraphTest{"identifier expr unknown name test", "identexpr", "unknown", []expectedScopeEntry{},
		"The name 'unknown' could not be found in this context", ""},

	scopegraphTest{"identifier expr parameter test", "identexpr", "parameter",
		[]expectedScopeEntry{
			expectedScopeEntry{"paramref", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr variable test", "identexpr", "var",
		[]expectedScopeEntry{
			expectedScopeEntry{"varref", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr implicit variable test", "identexpr", "implicitvar",
		[]expectedScopeEntry{
			expectedScopeEntry{"varref", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr import test", "identexpr", "import",
		[]expectedScopeEntry{
			expectedScopeEntry{"importref", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr type test", "identexpr", "type",
		[]expectedScopeEntry{
			expectedScopeEntry{"typeref", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr aliased type test", "identexpr", "alias",
		[]expectedScopeEntry{
			expectedScopeEntry{"typeref", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr generic type test", "identexpr", "generictype",
		[]expectedScopeEntry{
			expectedScopeEntry{"typeref", expectedScope{true, proto.ScopeKind_GENERIC, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr member var test", "identexpr", "membervar",
		[]expectedScopeEntry{
			expectedScopeEntry{"memberref", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr member func test", "identexpr", "memberfunc",
		[]expectedScopeEntry{
			expectedScopeEntry{"memberref", expectedScope{true, proto.ScopeKind_VALUE, "Function<void>", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr generic member test", "identexpr", "genericmember",
		[]expectedScopeEntry{
			expectedScopeEntry{"memberref", expectedScope{true, proto.ScopeKind_GENERIC, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr value test", "identexpr", "value",
		[]expectedScopeEntry{
			expectedScopeEntry{"valueref", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr nested test", "identexpr", "nested",
		[]expectedScopeEntry{
			expectedScopeEntry{"first", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"second", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
		},
		"", ""},

	/////////// Arrow operator ///////////

	scopegraphTest{"arrow operator success test", "arrowops", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"await", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"arrow operator invalid source test", "arrowops", "invalidsource",
		[]expectedScopeEntry{},
		"Right hand side of an arrow expression must be of type Promise: Type Integer cannot be used in place of type Promise as it does not implement member Then", ""},

	scopegraphTest{"arrow operator invalid destination test", "arrowops", "invaliddestination",
		[]expectedScopeEntry{},
		"Destination of arrow statement must accept type Boolean: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"arrow operator invalid rejection test", "arrowops", "invalidrejection",
		[]expectedScopeEntry{},
		"Rejection of arrow statement must accept type Error: 'Error' cannot be used in place of non-interface 'Boolean'", ""},

	/////////// List literal expression ///////////

	scopegraphTest{"list literal success test", "listliteral", "listliteral",
		[]expectedScopeEntry{
			expectedScopeEntry{"emptylist", expectedScope{true, proto.ScopeKind_VALUE, "List<any>", "void"}},
			expectedScopeEntry{"intlist", expectedScope{true, proto.ScopeKind_VALUE, "List<Integer>", "void"}},
			expectedScopeEntry{"mixedlist", expectedScope{true, proto.ScopeKind_VALUE, "List<any>", "void"}},
		},
		"", ""},

	/////////// Map literal expression ///////////

	scopegraphTest{"map literal success test", "mapliteral", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"emptymap", expectedScope{true, proto.ScopeKind_VALUE, "Map<Mappable, any>", "void"}},
			expectedScopeEntry{"intmap", expectedScope{true, proto.ScopeKind_VALUE, "Map<String, Integer>", "void"}},
			expectedScopeEntry{"mixedmap", expectedScope{true, proto.ScopeKind_VALUE, "Map<String, any>", "void"}},

			expectedScopeEntry{"intkeymap", expectedScope{true, proto.ScopeKind_VALUE, "Map<Integer, Integer>", "void"}},
			expectedScopeEntry{"mixedkeymap", expectedScope{true, proto.ScopeKind_VALUE, "Map<Mappable, Integer>", "void"}},
		},
		"", ""},

	scopegraphTest{"map literal nonmappable fail test", "mapliteral", "nonmappable",
		[]expectedScopeEntry{},
		"Map literal keys must be of type Mappable: Type 'SomeClass' does not define or export member 'MapKey', which is required by type 'Mappable'", ""},

	scopegraphTest{"map literal null key fail test", "mapliteral", "nullkey",
		[]expectedScopeEntry{},
		"Map literal keys must be of type Mappable: null cannot be used in place of non-nullable type Mappable", ""},

	/////////// Slice expression ///////////

	scopegraphTest{"slice operator success test", "slice", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"slice", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"startslice", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"endslice", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"index", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"indexer invalid param test", "slice", "invalidindex",
		[]expectedScopeEntry{},
		"Indexer parameter must be type Boolean: 'Integer' cannot be used in place of non-interface 'Boolean'", ""},

	scopegraphTest{"slice invalid param test", "slice", "invalidslice",
		[]expectedScopeEntry{},
		"Slice index must be of type Integer, found: Boolean", ""},

	scopegraphTest{"generic indexer success tests", "slice", "genericindexer",
		[]expectedScopeEntry{
			expectedScopeEntry{"getter", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"generic slice success tests", "slice", "genericslice",
		[]expectedScopeEntry{
			expectedScopeEntry{"getter", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"indexer success tests", "slice", "indexer",
		[]expectedScopeEntry{},
		"", ""},

	scopegraphTest{"invalid generic indexer tests", "slice", "invalidgenericindexer",
		[]expectedScopeEntry{},
		"Cannot assign value to operator setindex: Cannot use type 'T' in place of type 'Boolean'", ""},

	/////////// Cast expression ///////////

	scopegraphTest{"cast interface success test", "castexpr", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"cast", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
		},
		"", ""},

	scopegraphTest{"cast subclass success test", "castexpr", "subclass",
		[]expectedScopeEntry{
			expectedScopeEntry{"bc", expectedScope{true, proto.ScopeKind_VALUE, "BaseClass", "void"}},
			expectedScopeEntry{"abc", expectedScope{true, proto.ScopeKind_VALUE, "AnotherBaseClass", "void"}},
		},
		"", ""},

	scopegraphTest{"cast failure test", "castexpr", "failure",
		[]expectedScopeEntry{},
		"Cannot cast value of type 'ISomeInterface' to type 'SomeClass': Type 'SomeClass' does not define or export member 'DoSomething', which is required by type 'ISomeInterface'", ""},

	/////////// Function call expression ///////////

	scopegraphTest{"function call success test", "funccall", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"empty", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
			expectedScopeEntry{"somefunc", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"anotherfunc", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"function call not function failure test", "funccall", "notfunc",
		[]expectedScopeEntry{},
		"Cannot invoke function call on non-function 'SomeClass'.", ""},

	scopegraphTest{"function call invalid count failure test", "funccall", "invalidcount",
		[]expectedScopeEntry{},
		"Function call expects 0 arguments, found 1", ""},

	scopegraphTest{"function call type mismatch failure test", "funccall", "typemismatch",
		[]expectedScopeEntry{},
		"Parameter #1 expects type Integer: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	/////////// Member access expression ///////////

	scopegraphTest{"member access success test", "memberaccess", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"varmember", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"constructor", expectedScope{true, proto.ScopeKind_VALUE, "Function<SomeClass>(Integer)", "void"}},
			expectedScopeEntry{"modtype", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
			expectedScopeEntry{"modint", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"modfunc", expectedScope{true, proto.ScopeKind_VALUE, "Function<void>", "void"}},
			expectedScopeEntry{"generictype", expectedScope{true, proto.ScopeKind_GENERIC, "void", "void"}},
			expectedScopeEntry{"prop", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},

			expectedScopeEntry{"genmembool", expectedScope{true, proto.ScopeKind_VALUE, "Boolean?", "void"}},
			expectedScopeEntry{"genmemint", expectedScope{true, proto.ScopeKind_VALUE, "Integer?", "void"}},
		},
		"", ""},

	scopegraphTest{"member access static under instance failure test", "memberaccess", "staticunderinstance",
		[]expectedScopeEntry{},
		"Could not find instance name 'Build' under type SomeClass", ""},

	scopegraphTest{"member access instance under static failure test", "memberaccess", "instanceunderstatic",
		[]expectedScopeEntry{},
		"Could not find static name 'SomeInt' under class SomeClass", ""},

	scopegraphTest{"member access nullable failure test", "memberaccess", "nullable",
		[]expectedScopeEntry{},
		"Cannot access name 'someInt' under nullable type 'SomeClass?'. Please use the ?. operator to ensure type safety.", ""},

	scopegraphTest{"member access generic failure test", "memberaccess", "generic",
		[]expectedScopeEntry{},
		"Cannot attempt member access of 'Build' under class SomeClass, as it is generic without specification", ""},

	scopegraphTest{"member access generic func failure test", "memberaccess", "genericfunc",
		[]expectedScopeEntry{},
		"Cannot attempt member access of 'someMember' under module member GenericFunc, as it is generic without specification", ""},

	scopegraphTest{"member access literal test", "memberaccess", "literalaccess",
		[]expectedScopeEntry{},
		"Could not find instance name 'UnknownProp' under type Integer", ""},

	/////////// Nullable member access expression ///////////

	scopegraphTest{"nullable member access success test", "nullablememberaccess", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"sci", expectedScope{true, proto.ScopeKind_VALUE, "Integer?", "void"}},
			expectedScopeEntry{"scb", expectedScope{true, proto.ScopeKind_VALUE, "Boolean?", "void"}},
		},
		"", ""},

	scopegraphTest{"nullable member access non-nullable failure test", "nullablememberaccess", "nonnullable",
		[]expectedScopeEntry{},
		"Cannot access name 'SomeInt' under non-nullable type 'SomeClass'. Please use the . operator to ensure type safety.", ""},

	scopegraphTest{"nullable member access static failure test", "nullablememberaccess", "static",
		[]expectedScopeEntry{},
		"Cannot attempt nullable member access of 'Build' under class SomeClass, as it is a static type", ""},

	scopegraphTest{"nullable member access generic failure test", "nullablememberaccess", "generic",
		[]expectedScopeEntry{},
		"Cannot attempt nullable member access of 'something' under class SomeClass, as it is generic without specification", ""},

	/////////// Dynamic member access expression ///////////

	scopegraphTest{"dynamic member access success unknown test", "dynamicmemberaccess", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"dynaccess", expectedScope{true, proto.ScopeKind_VALUE, "any", "void"}},
		},
		"", "Member 'someValue' is unknown under known type Integer. This call will return null."},

	scopegraphTest{"dynamic member access success known test", "dynamicmemberaccess", "successknown",
		[]expectedScopeEntry{
			expectedScopeEntry{"dynaccess", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", "Dynamic access of known member 'someInt' under type SomeClass. The . operator is suggested."},

	scopegraphTest{"dynamic member access success nullable known test", "dynamicmemberaccess", "successnullableknown",
		[]expectedScopeEntry{
			expectedScopeEntry{"dynaccess", expectedScope{true, proto.ScopeKind_VALUE, "Integer?", "void"}},
		},
		"", "Dynamic access of known member 'someInt' under type SomeClass?. The ?. operator is suggested."},

	scopegraphTest{"dynamic member access static under non-static failure test", "dynamicmemberaccess", "staticunderinstance",
		[]expectedScopeEntry{},
		"Member 'Build' is static but accessed under an instance value", ""},

	scopegraphTest{"dynamic member access non-static under static failure test", "dynamicmemberaccess", "instanceunderstatic",
		[]expectedScopeEntry{},
		"Member 'SomeInt' is non-static but accessed under a static value", ""},

	/////////// Stream member access expression ///////////

	scopegraphTest{"stream member access success test", "streammemberaccess", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"intstream", expectedScope{true, proto.ScopeKind_VALUE, "Stream<Integer?>", "void"}},
			expectedScopeEntry{"boolstream", expectedScope{true, proto.ScopeKind_VALUE, "Stream<Boolean?>", "void"}},
		},
		"", ""},

	scopegraphTest{"stream member access non-stream failure test", "streammemberaccess", "nonstream",
		[]expectedScopeEntry{},
		"Cannot attempt stream access of name 'something' under non-stream type 'Integer': Type Integer cannot be used in place of type Stream as it does not implement member Next", ""},

	scopegraphTest{"stream member access generic failure test", "streammemberaccess", "generic",
		[]expectedScopeEntry{},
		"Cannot attempt stream member access of 'something' under class SomeClass, as it is generic without specification", ""},

	scopegraphTest{"stream member access static failure test", "streammemberaccess", "static",
		[]expectedScopeEntry{},
		"Cannot attempt stream member access of 'something' under class SomeClass, as it is a static type", ""},

	/////////// Null literal expression ///////////

	scopegraphTest{"null literal expression success test", "nullliteral", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"null", expectedScope{true, proto.ScopeKind_VALUE, "null", "void"}},
		},
		"", ""},

	scopegraphTest{"non-nullable null assign failure test", "nullliteral", "nonnullableassign",
		[]expectedScopeEntry{},
		"Variable 'foo' has declared type 'Integer': null cannot be used in place of non-nullable type Integer", ""},

	/////////// this literal expression ///////////

	scopegraphTest{"this literal expression success test", "thisliteral", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"scthis", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"acthis", expectedScope{true, proto.ScopeKind_VALUE, "AnotherClass<T>", "void"}},
		},
		"", ""},

	scopegraphTest{"this under module member test", "thisliteral", "modulemember",
		[]expectedScopeEntry{},
		"The 'this' keyword cannot be used under module member DoSomething", ""},

	scopegraphTest{"this under static member test", "thisliteral", "staticmember",
		[]expectedScopeEntry{},
		"The 'this' keyword cannot be used under static type member Build", ""},

	scopegraphTest{"this under property test", "thisliteral", "property",
		[]expectedScopeEntry{
			expectedScopeEntry{"scthis", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
		},
		"", ""},

	/////////// generic specifier expression ///////////

	scopegraphTest{"generic specifier expression success test", "genericspecifier", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"someclassint", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
			expectedScopeEntry{"someclassbool", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},

			expectedScopeEntry{"somestructint", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
			expectedScopeEntry{"somestructbool", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},

			expectedScopeEntry{"someclassintbuild", expectedScope{true, proto.ScopeKind_VALUE, "Function<SomeClass<Integer>>(Integer)", "void"}},
			expectedScopeEntry{"someclassboolbuild", expectedScope{true, proto.ScopeKind_VALUE, "Function<SomeClass<Boolean>>(Boolean)", "void"}},

			expectedScopeEntry{"somefuncintbool", expectedScope{true, proto.ScopeKind_VALUE, "Function<Integer?>(Boolean)", "void"}},
			expectedScopeEntry{"somefuncboolint", expectedScope{true, proto.ScopeKind_VALUE, "Function<Boolean?>(Integer)", "void"}},
		},
		"", ""},

	scopegraphTest{"generic specifier non-generic test", "genericspecifier", "nongeneric",
		[]expectedScopeEntry{},
		"Cannot apply generics to non-generic scope", ""},

	scopegraphTest{"generic specifier not enough generics test", "genericspecifier", "countmismatch",
		[]expectedScopeEntry{},
		"Generic count must match. Found: 1, expected: 2 on class SomeClass", ""},

	scopegraphTest{"generic specifier too many generics test", "genericspecifier", "countmismatch2",
		[]expectedScopeEntry{},
		"Generic count must match. Found: 2, expected: 1 on class SomeClass", ""},

	scopegraphTest{"generic specifier constraint failure test", "genericspecifier", "constraintfail",
		[]expectedScopeEntry{},
		"Cannot use type Boolean as generic T (#1) over class SomeClass: Type 'Boolean' does not define or export member 'DoSomething', which is required by type 'ISomeInterface'", ""},

	scopegraphTest{"generic specifier non-structural test", "genericspecifier", "structfail",
		[]expectedScopeEntry{},
		"Cannot use type SomeClass as generic T (#1) over struct SomeStruct: SomeClass is not structural nor serializable", ""},

	/////////// constructable types ///////////

	scopegraphTest{"constructable interface test", "types", "constructableinterface",
		[]expectedScopeEntry{},
		"", ""},

	scopegraphTest{"constructable generic test", "types", "constructablegeneric",
		[]expectedScopeEntry{
			expectedScopeEntry{"get", expectedScope{true, proto.ScopeKind_VALUE, "SomeInterface", "void"}},
		},
		"", ""},

	/////////// class field ///////////

	scopegraphTest{"class field uninitialized test", "var", "uninitializedfield",
		[]expectedScopeEntry{
			expectedScopeEntry{"new", expectedScope{true, proto.ScopeKind_VALUE, "Function<SomeClass>(Integer, Boolean)", "void"}},
		},
		"", ""},

	/////////// lambda expression ///////////

	scopegraphTest{"lambda expression basic inference test", "lambda", "basicinference",
		[]expectedScopeEntry{
			expectedScopeEntry{"varref", expectedScope{true, proto.ScopeKind_VALUE, "Function<Boolean>(Integer, String)", "void"}},
			expectedScopeEntry{"vardeclare", expectedScope{true, proto.ScopeKind_VALUE, "Function<Boolean>(Boolean, Boolean)", "void"}},
			expectedScopeEntry{"callref", expectedScope{true, proto.ScopeKind_VALUE, "Function<Boolean>(String, String)", "void"}},
			expectedScopeEntry{"nonref", expectedScope{true, proto.ScopeKind_VALUE, "Function<Boolean>(any, any)", "void"}},
			expectedScopeEntry{"ripref", expectedScope{true, proto.ScopeKind_VALUE, "Function<Integer>(Integer)", "void"}},
			expectedScopeEntry{"multiripref", expectedScope{true, proto.ScopeKind_VALUE, "Function<any>(any)", "void"}},
		},
		"", ""},

	scopegraphTest{"lambda expression full definition test", "lambda", "full",
		[]expectedScopeEntry{
			expectedScopeEntry{"implicitreturn", expectedScope{true, proto.ScopeKind_VALUE, "Function<Integer>(Integer, Boolean)", "void"}},
			expectedScopeEntry{"explicitreturn", expectedScope{true, proto.ScopeKind_VALUE, "Function<String>(Integer, Boolean)", "void"}},
		},
		"", ""},

	/////////// chained inference test ///////////

	scopegraphTest{"chained inference test", "chained", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"someVar", expectedScope{true, proto.ScopeKind_VALUE, "Function<void>(Integer, Boolean)", "void"}},
			expectedScopeEntry{"anotherVar", expectedScope{true, proto.ScopeKind_VALUE, "Function<void>(Integer, Boolean)", "void"}},
		},
		"", ""},

	/////////// property value test /////////////

	scopegraphTest{"property val test", "property", "getterval",
		[]expectedScopeEntry{},
		"The 'val' keyword can only be used under property setters", ""},

	scopegraphTest{"property val test", "property", "value",
		[]expectedScopeEntry{
			expectedScopeEntry{"val", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	/////////// template strings /////////////

	scopegraphTest{"untagged template string success", "templatestr", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"templatestr", expectedScope{true, proto.ScopeKind_VALUE, "String", "void"}},
		},
		"", ""},

	scopegraphTest{"tagged template string success", "templatestr", "taggedsuccess",
		[]expectedScopeEntry{
			expectedScopeEntry{"templatestr", expectedScope{true, proto.ScopeKind_VALUE, "String", "void"}},
			expectedScopeEntry{"nonstrtemplatestr", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"template string non-stringable failure", "templatestr", "nonstringable",
		[]expectedScopeEntry{},
		"All expressions in a template string must be of type Stringable: Type 'SomeClass' does not define or export member 'String', which is required by type 'Stringable'", ""},

	scopegraphTest{"tagged template failure success", "templatestr", "taggedfailure",
		[]expectedScopeEntry{},
		"Tagging expression for template string must be function with parameters ([]string, []stringable). Found: Function<void>", ""},

	/////////// webidl tests /////////////////

	scopegraphTest{"webidl success", "webidl", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"sometype", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
			expectedScopeEntry{"someparam", expectedScope{true, proto.ScopeKind_VALUE, "SomeType", "void"}},
			expectedScopeEntry{"somefunc", expectedScope{true, proto.ScopeKind_VALUE, "any", "void"}},
			expectedScopeEntry{"con1", expectedScope{true, proto.ScopeKind_VALUE, "SomeType", "void"}},
			expectedScopeEntry{"con2", expectedScope{true, proto.ScopeKind_VALUE, "SomeType", "void"}},
			expectedScopeEntry{"addition", expectedScope{true, proto.ScopeKind_VALUE, "SomeType", "void"}},
			expectedScopeEntry{"getindex", expectedScope{true, proto.ScopeKind_VALUE, "Second", "void"}},
		},
		"", ""},

	scopegraphTest{"webidl unknown type failure", "webidl", "unknowntype",
		[]expectedScopeEntry{},
		"Type 'global.UnknownType' could not be found", ""},

	scopegraphTest{"webidl non-static function failure", "webidl", "nonstatic",
		[]expectedScopeEntry{},
		"Could not find static name 'SomeFunction' under external interface SomeType", ""},

	scopegraphTest{"webidl static function failure", "webidl", "static",
		[]expectedScopeEntry{},
		"Could not find static name 'SomeFunction' under external interface SomeType", ""},

	/////////// nominal tests /////////////////

	scopegraphTest{"nominal type success", "nominal", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"this", expectedScope{true, proto.ScopeKind_VALUE, "MyType", "void"}},
			expectedScopeEntry{"sometype", expectedScope{true, proto.ScopeKind_VALUE, "SomeType", "void"}},
			expectedScopeEntry{"m", expectedScope{true, proto.ScopeKind_VALUE, "MyType", "void"}},
			expectedScopeEntry{"at", expectedScope{true, proto.ScopeKind_VALUE, "AnotherType", "void"}},
			expectedScopeEntry{"gt", expectedScope{true, proto.ScopeKind_VALUE, "GenericType<Integer>", "void"}},
		},
		"", ""},

	scopegraphTest{"nominal conversion failure", "nominal", "cannotconvert",
		[]expectedScopeEntry{},
		"Cannot perform type conversion: Type 'MyType' cannot be converted to or from type 'SomeType'", ""},

	scopegraphTest{"nominal conversion argument count mismatch", "nominal", "convertargcount",
		[]expectedScopeEntry{},
		"Type conversion requires a single argument", ""},

	scopegraphTest{"nominal conversion no argument mismatch", "nominal", "convertnoarg",
		[]expectedScopeEntry{},
		"Type conversion requires a single argument", ""},

	scopegraphTest{"nominal call base type member", "nominal", "basetypecall",
		[]expectedScopeEntry{},
		"Could not find instance name 'DoSomething' under type MyType", ""},

	/////////// structural new expression tests /////////////////

	scopegraphTest{"structural new success tests", "structnew", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"empty", expectedScope{true, proto.ScopeKind_VALUE, "EmptyClass", "void"}},
			expectedScopeEntry{"nonempty", expectedScope{true, proto.ScopeKind_VALUE, "NonRequiredClass", "void"}},
			expectedScopeEntry{"nonesome", expectedScope{true, proto.ScopeKind_VALUE, "NonRequiredClass", "void"}},
			expectedScopeEntry{"noneanother", expectedScope{true, proto.ScopeKind_VALUE, "NonRequiredClass", "void"}},
			expectedScopeEntry{"somestruct", expectedScope{true, proto.ScopeKind_VALUE, "SomeStruct", "void"}},
			expectedScopeEntry{"required", expectedScope{true, proto.ScopeKind_VALUE, "RequiredClass", "void"}},
			expectedScopeEntry{"generic", expectedScope{true, proto.ScopeKind_VALUE, "GenericStruct<Integer>", "void"}},
		},
		"", ""},

	scopegraphTest{"structural new unconstructable type test", "structnew", "unconstructable",
		[]expectedScopeEntry{},
		"Cannot structurally construct type SomeInterface", ""},

	scopegraphTest{"structural new non-type test", "structnew", "nontype",
		[]expectedScopeEntry{},
		"Cannot construct non-type expression", ""},

	scopegraphTest{"structural new imported class test", "structnew", "importedclass",
		[]expectedScopeEntry{},
		"Cannot structurally construct type SomeClass, as it is imported from another module", ""},

	scopegraphTest{"structural new invalid name test", "structnew", "invalidname",
		[]expectedScopeEntry{},
		"Name 'SomeField' could not be found under type SomeClass", ""},

	scopegraphTest{"structural new invalid value test", "structnew", "invalidvalue",
		[]expectedScopeEntry{},
		"The name 'foo' could not be found in this context", ""},

	scopegraphTest{"structural new read-only field test", "structnew", "readonlyname",
		[]expectedScopeEntry{},
		"type member SomeField under type SomeClass is read-only", ""},

	scopegraphTest{"structural new value type mismatch test", "structnew", "valuetypemismatch",
		[]expectedScopeEntry{},
		"Cannot assign value of type Integer to type member SomeField: 'Integer' cannot be used in place of non-interface 'Boolean'", ""},

	scopegraphTest{"structural new value missing required field test", "structnew", "missingrequired",
		[]expectedScopeEntry{},
		"Non-nullable type member 'SomeField' is required to construct type SomeClass", ""},

	/////////// async function tests /////////////////

	scopegraphTest{"async under class test", "async", "underclass",
		[]expectedScopeEntry{},
		"Asynchronous functions must be declared under modules: 'DoSomethingAsync' defined under class SomeClass", ""},

	scopegraphTest{"async invalid return type test", "async", "invalidreturn",
		[]expectedScopeEntry{},
		"Asynchronous function DoSomethingAsync must return a structural type: SomeClass is not structural nor serializable", ""},

	scopegraphTest{"async invalid parameter test", "async", "invalidparam",
		[]expectedScopeEntry{},
		"Parameters of asynchronous function DoSomethingAsync must be structural: SomeClass is not structural nor serializable", ""},

	scopegraphTest{"async outside context test", "async", "outsidecontext",
		[]expectedScopeEntry{},
		"", "module member 'outside' is defined outside the async function and will therefore be unique for each call to this function"},

	scopegraphTest{"async success test", "async", "success",
		[]expectedScopeEntry{},
		"", ""},

	/////////// known issue tests /////////////////

	scopegraphTest{"known issue panic test", "knownissues", "knownissue1",
		[]expectedScopeEntry{},
		"Operator 'equals' is not defined on type 'T'", ""},
}

func TestGraphs(t *testing.T) {
	for _, test := range scopeGraphTests {
		if os.Getenv("FILTER") != "" && !strings.Contains(test.name, os.Getenv("FILTER")) {
			continue
		}

		fmt.Printf("Running test: %v\n", test.name)

		entrypointFile := "tests/" + test.input + "/" + test.entrypoint + ".seru"
		result := ParseAndBuildScopeGraph(entrypointFile, packageloader.Library{TESTLIB_PATH, false, ""})

		if test.expectedError != "" {
			if !assert.False(t, result.Status, "Expected failure in scoping on test : %v", test.name) {
				continue
			}

			assert.Equal(t, 1, len(result.Errors), "Expected 1 error on test %v, found: %v", test.name, result.Errors)
			assert.Equal(t, test.expectedError, result.Errors[0].Error(), "Error mismatch on test %v", test.name)
			continue
		} else {
			if !assert.True(t, result.Status, "Expected success in scoping on test: %v\n%v\n%v", test.name, result.Errors, result.Warnings) {
				continue
			}
		}

		if test.expectedWarning != "" {
			if assert.Equal(t, 1, len(result.Warnings), "Expected 1 warning on test %v, found: %v", test.name, len(result.Warnings)) {
				assert.Equal(t, test.expectedWarning, result.Warnings[0].Warning(), "Warning mismatch on test %v", test.name)
			}
		}

		// Check each of the scopes.
		for _, expected := range test.expectedScope {
			// Collect the scopes for all requested nodes and compare.
			node, found := result.Graph.SourceGraph().FindCommentedNode(fmt.Sprintf("/* %s */", expected.name))
			if !assert.True(t, found, "Missing commented node %s in test: %v", expected.name, test.name) {
				continue
			}

			scope, valid := result.Graph.GetScope(node)
			if !assert.True(t, valid, "Could not get scope for commented node %s (%v) in test: %v (%v)", expected.name, node, test.name, result.Errors) {
				continue
			}

			// Compare the scope found to that expected.
			if !assert.Equal(t, expected.scope.IsValid, scope.GetIsValid(), "Scope valid mismatch for commented node %s in test: %v", expected.name, test.name) {
				continue
			}

			if !expected.scope.IsValid {
				continue
			}

			if !assert.Equal(t, expected.scope.ScopeKind, scope.GetKind(), "Scope kind mismatch for commented node %s in test: %v", expected.name, test.name) {
				continue
			}

			if !assert.Equal(t, expected.scope.ResolvedType, scope.ResolvedTypeRef(result.Graph.TypeGraph()).String(), "Resolved type mismatch for commented node %s in test: %v", expected.name, test.name) {
				continue
			}

			if !assert.Equal(t, expected.scope.ReturnedType, scope.ReturnedTypeRef(result.Graph.TypeGraph()).String(), "Returned type mismatch for commented node %s in test: %v", expected.name, test.name) {
				continue
			}
		}
	}
}
