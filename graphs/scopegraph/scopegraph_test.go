// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"
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

	/////////// Generator ///////////

	// Success test.
	scopegraphTest{"generator success test", "generator", "success", []expectedScopeEntry{},
		"", ""},

	// Nested success test.
	scopegraphTest{"generator nested success test", "generator", "nestedsuccess", []expectedScopeEntry{},
		"", ""},

	// Non-stream.
	scopegraphTest{"yield non-stream test", "generator", "nonstream", []expectedScopeEntry{},
		"'yield' statement must be under a function or property returning a Stream. Found: void", ""},

	// Value mismatch.
	scopegraphTest{"yield value mismatch test", "generator", "valuemismatch", []expectedScopeEntry{},
		"'yield' expression must have subtype of Boolean: 'Integer' cannot be used in place of non-interface 'Boolean'", ""},

	// Subgenerator mismatch.
	scopegraphTest{"yield in mismatch test", "generator", "inmismatch", []expectedScopeEntry{},
		"'yield in' expression must have subtype of Stream<Boolean>: member 'Next' under type 'Stream<Integer>' does not match that defined in type 'Stream<Boolean>'", ""},

	// Async not allowed.
	scopegraphTest{"async generator test", "generator", "async", []expectedScopeEntry{},
		"Asynchronous function DoSomethingAsync must return a structural type: Stream<Integer> is not structural nor serializable", ""},

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

	// Value returned by void function.
	scopegraphTest{"settlement void return value test", "settlement", "voidreturn", []expectedScopeEntry{},
		"No return value expected here, found value of type 'Integer'", ""},

	// Invalid value under conditional.
	scopegraphTest{"conditional invalid return test", "settlement", "conditionalinvalidreturn", []expectedScopeEntry{},
		"No return value expected here, found value of type 'Integer'", ""},

	/////////// Break ///////////

	// Normal break statement.
	scopegraphTest{"break statement test", "break", "normal", []expectedScopeEntry{},
		"", "Unreachable statement found"},

	// break statement not under a breakable node.
	scopegraphTest{"break statement error test", "break", "badparent", []expectedScopeEntry{},
		"'break' statement must be a under a loop or switch statement", ""},

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

	scopegraphTest{"conditional nullable assignment test", "conditional", "nullableassign",
		[]expectedScopeEntry{},
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
		"Expected return value of type 'Integer': 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	/////////// Loops ///////////

	// Empty loop test.
	scopegraphTest{"empty loop test", "loop", "empty", []expectedScopeEntry{
		expectedScopeEntry{"emptyloop", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
	}, "", "Unreachable statement found"},

	// Empty with return loop test.
	scopegraphTest{"empty with return loop test", "loop", "emptywithreturn", []expectedScopeEntry{
		expectedScopeEntry{"emptyloop", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
	}, "", ""},

	// Empty missing return test..
	scopegraphTest{"empty missing return loop test", "loop", "emptymissingreturn", []expectedScopeEntry{
		expectedScopeEntry{"emptyloop", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
	}, "", ""},

	// Empty with reject loop test.
	scopegraphTest{"empty with reject loop test", "loop", "emptywithreturn", []expectedScopeEntry{
		expectedScopeEntry{"emptyloop", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
	}, "", ""},

	// Empty with match that returns loop test.
	scopegraphTest{"empty with match return loop test", "loop", "emptywithmatch", []expectedScopeEntry{
		expectedScopeEntry{"emptyloop", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
	}, "", ""},

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

	scopegraphTest{"with as test", "with", "with_as", []expectedScopeEntry{
		expectedScopeEntry{"someInt", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
	}, "", ""},

	scopegraphTest{"invalid with expr test", "with", "nonreleasable", []expectedScopeEntry{},
		"With expression must implement the Releasable interface: Type 'Boolean' does not define or export member 'Release', which is required by type 'Releasable'", ""},

	/////////// Match ///////////

	scopegraphTest{"basic match switch test", "match", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"int", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"bool", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"string", expectedScope{true, proto.ScopeKind_VALUE, "String", "void"}},
			expectedScopeEntry{"default", expectedScope{true, proto.ScopeKind_VALUE, "any", "void"}},
		},
		"", ""},

	scopegraphTest{"match invalid branch test", "match", "invalidbranch",
		[]expectedScopeEntry{},
		"Match cases must be castable from type 'Integer': 'String' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"match nullable branch test", "match", "nullable",
		[]expectedScopeEntry{},
		"Match cases cannot be nullable. Found: String?", ""},

	/////////// Switch ///////////

	scopegraphTest{"basic bool switch test", "switch", "bool",
		[]expectedScopeEntry{
			expectedScopeEntry{"switch", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"switch no statement test", "switch", "nostatement",
		[]expectedScopeEntry{
			expectedScopeEntry{"switch", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"switch multi statement test", "switch", "multistatement",
		[]expectedScopeEntry{
			expectedScopeEntry{"switch", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
		},
		"", ""},

	scopegraphTest{"bool default switch test", "switch", "booldefault",
		[]expectedScopeEntry{
			expectedScopeEntry{"switch", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
		},
		"", ""},

	scopegraphTest{"basic non-bool switch test", "switch", "typed",
		[]expectedScopeEntry{
			expectedScopeEntry{"switch", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"non-bool default switch test", "switch", "typeddefault",
		[]expectedScopeEntry{
			expectedScopeEntry{"switch", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
		},
		"", ""},

	scopegraphTest{"invalid bool switch test", "switch", "boolinvalid", []expectedScopeEntry{},
		"Switch cases must have values matching type 'Boolean': 'Integer' cannot be used in place of non-interface 'Boolean'", ""},

	scopegraphTest{"invalid non-bool switch test", "switch", "typedinvalid", []expectedScopeEntry{},
		"Switch cases must have values matching type 'Integer': 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"switch no equals operator test", "switch", "nocompare", []expectedScopeEntry{},
		"Cannot switch over instance of type 'SomeClass', as it does not define or export an 'equals' operator", ""},

	/////////// Var ///////////

	scopegraphTest{"basic var test", "var", "basic",
		[]expectedScopeEntry{
			expectedScopeEntry{"var", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"assign any var test", "var", "assignany",
		[]expectedScopeEntry{},
		"", ""},

	scopegraphTest{"var known issue test", "var", "knownissue",
		[]expectedScopeEntry{},
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

	scopegraphTest{"bad expression var test", "var", "badexpr",
		[]expectedScopeEntry{},
		"The name 'bar' could not be found in this context", ""},

	/////////// SML expression ///////////

	scopegraphTest{"sml expression success test", "sml", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"simple", expectedScope{true, proto.ScopeKind_VALUE, "SimpleClass", "void"}},
			expectedScopeEntry{"classwithprops", expectedScope{true, proto.ScopeKind_VALUE, "ClassWithProps", "void"}},
			expectedScopeEntry{"funcwithprops", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},

			expectedScopeEntry{"propsstruct", expectedScope{true, proto.ScopeKind_VALUE, "SomePropsStruct", "void"}},

			expectedScopeEntry{"optionalchild", expectedScope{true, proto.ScopeKind_VALUE, "String", "void"}},
			expectedScopeEntry{"optionalchild2", expectedScope{true, proto.ScopeKind_VALUE, "String", "void"}},

			expectedScopeEntry{"requiredchild", expectedScope{true, proto.ScopeKind_VALUE, "String", "void"}},

			expectedScopeEntry{"childstream1", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"childstream2", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"childstream3", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"childstream4", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},

			expectedScopeEntry{"decorator", expectedScope{true, proto.ScopeKind_VALUE, "AnotherClass", "void"}},
			expectedScopeEntry{"decorator2", expectedScope{true, proto.ScopeKind_VALUE, "AnotherClass", "void"}},

			expectedScopeEntry{"chaineddecorator", expectedScope{true, proto.ScopeKind_VALUE, "ThirdClass", "void"}},

			expectedScopeEntry{"subtypedecorator", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"sml expression unknown ref test", "sml", "unknownref",
		[]expectedScopeEntry{},
		"The name 'something' could not be found in this context", ""},

	scopegraphTest{"sml expression unknown ref 2 test", "sml", "unknownref2",
		[]expectedScopeEntry{},
		"Could not find static name 'something' under class SomeClass", ""},

	scopegraphTest{"sml expression missing constructor test", "sml", "missingconstructor",
		[]expectedScopeEntry{},
		"A type used in a SML declaration tag must have a 'Declare' constructor: Could not find static name 'Declare' under class SomeClass", ""},

	scopegraphTest{"sml expression generic type test", "sml", "generic",
		[]expectedScopeEntry{},
		"A generic type cannot be used in a SML declaration tag", ""},

	scopegraphTest{"sml expression non-function type test", "sml", "nonfunction",
		[]expectedScopeEntry{},
		"Declared reference in an SML declaration tag must be a function. Found: Integer", ""},

	scopegraphTest{"sml expression invalid param count test", "sml", "invalidparamcount",
		[]expectedScopeEntry{},
		"Declarable function or constructor used in an SML declaration tag with attributes must have a 'props' parameter as parameter #1. Found: function<SomeClass>", ""},

	scopegraphTest{"sml expression void function test", "sml", "voidfunction",
		[]expectedScopeEntry{},
		"Declarable function used in an SML declaration tag cannot return void", ""},

	scopegraphTest{"sml expression invalid props type test", "sml", "invalidprops",
		[]expectedScopeEntry{},
		"Props parameter (parameter #1) of a declarable function or constructor used in an SML declaration tag must be a struct, a class with a ForProps constructor or a Mapping. Found: Integer", ""},

	scopegraphTest{"sml expression invalid props class test", "sml", "invalidpropsclass",
		[]expectedScopeEntry{},
		"Props parameter (parameter #1) of a declarable function or constructor used in an SML declaration tag has type SomeClass, which does not have any settable fields; use an empty `struct` instead if this is the intended behavior", ""},

	scopegraphTest{"sml expression invalid attribute name test", "sml", "invalidattributename",
		[]expectedScopeEntry{},
		"Could not find instance name 'unknownattr' under struct SomeType", ""},

	scopegraphTest{"sml expression invalid attribute value test", "sml", "invalidattributevalue",
		[]expectedScopeEntry{},
		"Cannot assign value of type String for attribute SomeAttribute: 'String' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"sml expression invalid attribute value 2 test", "sml", "invalidattributevalue2",
		[]expectedScopeEntry{},
		"Cannot assign value of type Boolean for attribute SomeAttribute: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"sml expression missing required attribute test", "sml", "missingrequiredattribute",
		[]expectedScopeEntry{},
		"Required attribute 'SomeAttribute' is missing for SML declaration props type SomeType", ""},

	scopegraphTest{"sml expression unknown decorator test", "sml", "unknowndecorator",
		[]expectedScopeEntry{},
		"The name 'SomeDecorator' could not be found in this context", ""},

	scopegraphTest{"sml expression non-function decorator test", "sml", "nonfunctiondecorator",
		[]expectedScopeEntry{},
		"SML declaration decorator 'sd' must refer to a function. Found: Integer", ""},

	scopegraphTest{"sml expression void function decorator test", "sml", "voidfunctiondecorator",
		[]expectedScopeEntry{},
		"SML declaration decorator 'sd' cannot return void", ""},

	scopegraphTest{"sml expression parameter count decorator test", "sml", "invalidparamcountdecorator",
		[]expectedScopeEntry{},
		"SML declaration decorator 'sd' must refer to a function with two parameters. Found: function<Integer>", ""},

	scopegraphTest{"sml expression decorator decorated mismatch test", "sml", "decoratedmismatch",
		[]expectedScopeEntry{},
		"SML declaration decorator 'sd' expects to decorate an instance of type String: 'Integer' cannot be used in place of non-interface 'String'", ""},

	scopegraphTest{"sml expression decorator value mismatch test", "sml", "decoratorvaluemismatch",
		[]expectedScopeEntry{},
		"Cannot assign value of type String for decorator 'sd': 'String' cannot be used in place of non-interface 'Boolean'", ""},

	scopegraphTest{"sml expression decorator value mismatch 2 test", "sml", "decoratorvaluemismatch2",
		[]expectedScopeEntry{},
		"Cannot assign value of type Boolean for decorator 'sd': 'Boolean' cannot be used in place of non-interface 'String'", ""},

	scopegraphTest{"sml expression decorator chained decorated mismatch test", "sml", "chaineddecoratedmismatch",
		[]expectedScopeEntry{},
		"SML declaration decorator 'second' expects to decorate an instance of type AnotherTypeEntirely: 'AnotherType' cannot be used in place of non-interface 'AnotherTypeEntirely'", ""},

	scopegraphTest{"sml expression child parameter count mismatch test", "sml", "childparamcountmismatch",
		[]expectedScopeEntry{},
		"Declarable function or constructor used in an SML declaration tag with children must have a 'children' parameter. Found: function<SomeClass>", ""},

	scopegraphTest{"sml expression child required test", "sml", "childrequired",
		[]expectedScopeEntry{},
		"SML declaration tag requires a single child. Found: 0", ""},

	scopegraphTest{"sml expression too many children test", "sml", "toomanychildren",
		[]expectedScopeEntry{},
		"SML declaration tag allows at most a single child. Found: 2", ""},

	scopegraphTest{"sml expression child type mismatch test", "sml", "childtypemismatch",
		[]expectedScopeEntry{},
		"SML declaration tag requires a child of type Integer: 'String' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"sml expression child stream mismatch test", "sml", "childstreammismatch",
		[]expectedScopeEntry{},
		"Child #1 under SML declaration must be subtype of Integer: 'String' cannot be used in place of non-interface 'Integer'", ""},

	/////////// Numeric literals ///////////

	scopegraphTest{"numeric literals test", "literals", "numeric",
		[]expectedScopeEntry{
			expectedScopeEntry{"int1", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"int2", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"int3", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"float1", expectedScope{true, proto.ScopeKind_VALUE, "Float64", "void"}},
			expectedScopeEntry{"float2", expectedScope{true, proto.ScopeKind_VALUE, "Float64", "void"}},
			expectedScopeEntry{"float3", expectedScope{true, proto.ScopeKind_VALUE, "Float64", "void"}},
		},
		"", ""},

	/////////// Resolve ///////////

	scopegraphTest{"resolve statement success test", "resolve", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"firstref", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"secondref", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},

			expectedScopeEntry{"thirdresolveref", expectedScope{true, proto.ScopeKind_VALUE, "Integer?", "void"}},
			expectedScopeEntry{"thirdrejectref", expectedScope{true, proto.ScopeKind_VALUE, "Error?", "void"}},

			expectedScopeEntry{"fourthresolveref", expectedScope{true, proto.ScopeKind_VALUE, "Integer?", "void"}},

			expectedScopeEntry{"fifthrejectref", expectedScope{true, proto.ScopeKind_VALUE, "Error?", "void"}},
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

	scopegraphTest{"assign interface property test", "assign", "interfaceprop", []expectedScopeEntry{},
		"", ""},

	scopegraphTest{"assign indexer test", "assign", "indexer",
		[]expectedScopeEntry{},
		"", ""},

	scopegraphTest{"assign indexer value type mismatch test", "assign", "indexermismatch",
		[]expectedScopeEntry{},
		"Cannot assign value to operator setindex: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	/////////// Loop expression ///////////

	scopegraphTest{"loop expression success test", "loopexpr", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"loopexpr", expectedScope{true, proto.ScopeKind_VALUE, "Stream<String>", "void"}},
		},
		"", ""},

	scopegraphTest{"loop expression non-stream test", "loopexpr", "nonstream",
		[]expectedScopeEntry{},
		"Loop iterable expression must implement type 'stream' or 'streamable': Type Integer cannot be used in place of type Stream as it does not implement member Next", ""},

	scopegraphTest{"loop expression invalid var test", "loopexpr", "invalidvar",
		[]expectedScopeEntry{},
		"The name 'les' could not be found in this context", ""},

	/////////// Conditional expression ///////////

	scopegraphTest{"conditional expression success test", "condexpr", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"condexpr", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"condexpr2", expectedScope{true, proto.ScopeKind_VALUE, "struct", "void"}},
		},
		"", ""},

	scopegraphTest{"conditional expression nullable success test", "condexpr", "nullable",
		[]expectedScopeEntry{
			expectedScopeEntry{"condexpr", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"condexpr2", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"conditional expression non-bool test", "condexpr", "nonbool",
		[]expectedScopeEntry{},
		"Conditional expression check must be of type 'bool', found: Integer", ""},

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

	scopegraphTest{"root type op nullable test", "rootop", "nullable",
		[]expectedScopeEntry{
			expectedScopeEntry{"someclass", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass?", "void"}},
			expectedScopeEntry{"generic", expectedScope{true, proto.ScopeKind_VALUE, "any", "void"}},
			expectedScopeEntry{"interface", expectedScope{true, proto.ScopeKind_VALUE, "any", "void"}},
		},
		"", ""},

	/////////// In operator expression ///////////

	scopegraphTest{"in op success test", "inop", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"in", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"in op nullable test", "inop", "nullable",
		[]expectedScopeEntry{},
		"Cannot invoke operator 'in' on nullable value of type 'SomeClass?'", ""},

	scopegraphTest{"in op no contains test", "inop", "nocontains",
		[]expectedScopeEntry{},
		"Operator 'contains' is not defined on type 'SomeClass'", ""},

	scopegraphTest{"in op invalid arg test", "inop", "invalidarg",
		[]expectedScopeEntry{},
		"Cannot invoke operator 'in' with value of type 'Boolean': 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	/////////// Is operator expression ///////////

	scopegraphTest{"is op success test", "isop", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"isresult", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},

			expectedScopeEntry{"aundercheck", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"aundercheck2", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"aundercheck3", expectedScope{true, proto.ScopeKind_VALUE, "null", "void"}},
		},
		"", ""},

	scopegraphTest{"is op not null failure test", "isop", "notnull",
		[]expectedScopeEntry{},
		"Right side of 'is' operator must be 'null' or 'not null'", ""},

	scopegraphTest{"is op not nullable failure test", "isop", "notnullable",
		[]expectedScopeEntry{},
		"Left side of 'is' operator must be a nullable type. Found: Integer", ""},

	/////////// Not keyword operator expression ///////////

	scopegraphTest{"not op success test", "notop", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"isnot", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"not", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"not op missing bool op failure test", "notop", "missingboolop",
		[]expectedScopeEntry{},
		"Operator 'bool' is not defined on type 'SomeClass'", ""},

	/////////// Nullable operator expression ///////////

	scopegraphTest{"nullable ops success test", "nullableops", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"cnullclass", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"cnullinter", expectedScope{true, proto.ScopeKind_VALUE, "ISomeInterface", "void"}},

			expectedScopeEntry{"inullclass", expectedScope{true, proto.ScopeKind_VALUE, "ISomeInterface", "void"}},
			expectedScopeEntry{"inullinter", expectedScope{true, proto.ScopeKind_VALUE, "ISomeInterface", "void"}},
		},
		"", ""},

	scopegraphTest{"nullable ops assert success test", "nullableops", "assert",
		[]expectedScopeEntry{
			expectedScopeEntry{"sc", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
		},
		"", ""},

	scopegraphTest{"nullable ops assert non-nullable child test", "nullableops", "assertnonnull",
		[]expectedScopeEntry{},
		"Child expression of an assert not nullable operator must be nullable. Found: SomeClass", ""},

	scopegraphTest{"nullable ops non-nullable left test", "nullableops", "nonnullable",
		[]expectedScopeEntry{},
		"Left hand side of a nullable operator must be nullable. Found: SomeClass", ""},

	scopegraphTest{"nullable ops nullable right test", "nullableops", "nullableright",
		[]expectedScopeEntry{},
		"Right hand side of a nullable operator cannot be nullable. Found: SomeClass?", ""},

	scopegraphTest{"nullable ops non-subtype test", "nullableops", "subtypemismatch",
		[]expectedScopeEntry{},
		"Left and right hand sides of a nullable operator must have common subtype. None found between 'Integer' and 'Boolean'", ""},

	/////////// Unary operator expressions ///////////

	scopegraphTest{"unary ops success test", "unaryops", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"not", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
		},
		"", ""},

	scopegraphTest{"unary ops generic test", "unaryops", "generic",
		[]expectedScopeEntry{
			expectedScopeEntry{"not", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass<Integer>", "void"}},
		},
		"", ""},

	scopegraphTest{"unary ops nullable fail test", "unaryops", "nullable",
		[]expectedScopeEntry{},
		"Cannot invoke operator 'not' on nullable type 'SomeClass?'", ""},

	/////////// Range operator expressions ///////////

	scopegraphTest{"range op simple test", "rangeop", "simple",
		[]expectedScopeEntry{
			expectedScopeEntry{"range", expectedScope{true, proto.ScopeKind_VALUE, "Stream<SomeClass>", "void"}},
		},
		"", ""},

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

	scopegraphTest{"binary ops generic test", "binaryops", "generic",
		[]expectedScopeEntry{
			expectedScopeEntry{"result", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass<Integer>", "void"}},
		},
		"", ""},

	scopegraphTest{"binary ops nullable fail test", "binaryops", "nullable",
		[]expectedScopeEntry{},
		"Cannot invoke operator 'plus' on nullable type 'SomeClass?'", ""},

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
			expectedScopeEntry{"memberref", expectedScope{true, proto.ScopeKind_VALUE, "function<void>", "void"}},
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
		"Right hand side of an arrow expression must be of type Awaitable: Type Integer cannot be used in place of type Awaitable as it does not implement member Then", ""},

	scopegraphTest{"arrow operator invalid destination test", "arrowops", "invaliddestination",
		[]expectedScopeEntry{},
		"Destination of arrow statement must accept type Boolean: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"arrow operator invalid rejection test", "arrowops", "invalidrejection",
		[]expectedScopeEntry{},
		"Rejection of arrow statement must accept type Error: 'Error' cannot be used in place of non-interface 'Boolean'", ""},

	/////////// List literal expression ///////////

	scopegraphTest{"list literal success test", "listliteral", "listliteral",
		[]expectedScopeEntry{
			expectedScopeEntry{"emptylist", expectedScope{true, proto.ScopeKind_VALUE, "Slice<any>", "void"}},
			expectedScopeEntry{"intlist", expectedScope{true, proto.ScopeKind_VALUE, "Slice<Integer>", "void"}},
			expectedScopeEntry{"mixedlist", expectedScope{true, proto.ScopeKind_VALUE, "Slice<struct>", "void"}},
		},
		"", ""},

	/////////// Slice literal expression ///////////

	scopegraphTest{"slice literal success test", "sliceliteral", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"slice", expectedScope{true, proto.ScopeKind_VALUE, "Slice<String>", "void"}},
		},
		"", ""},

	scopegraphTest{"slice literal invalid value", "sliceliteral", "invalidvalue",
		[]expectedScopeEntry{},
		"Invalid slice literal value: 'String' cannot be used in place of non-interface 'SomeClass'", ""},

	/////////// Mapping literal expression ///////////

	scopegraphTest{"mapping literal success test", "mappingliteral", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"mapping", expectedScope{true, proto.ScopeKind_VALUE, "Mapping<Integer>", "void"}},
		},
		"", ""},

	scopegraphTest{"mapping literal invalid key", "mappingliteral", "invalidkey",
		[]expectedScopeEntry{},
		"Mapping literal keys must be of type Stringable: Type 'SomeClass' does not define or export member 'String', which is required by type 'Stringable'", ""},

	scopegraphTest{"mapping literal invalid value", "mappingliteral", "invalidvalue",
		[]expectedScopeEntry{},
		"Expected mapping values of type Integer: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	/////////// Map literal expression ///////////

	scopegraphTest{"map literal success test", "mapliteral", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"emptymap", expectedScope{true, proto.ScopeKind_VALUE, "Mapping<any>", "void"}},
			expectedScopeEntry{"intmap", expectedScope{true, proto.ScopeKind_VALUE, "Mapping<Integer>", "void"}},
			expectedScopeEntry{"mixedmap", expectedScope{true, proto.ScopeKind_VALUE, "Mapping<struct>", "void"}},

			expectedScopeEntry{"intkeymap", expectedScope{true, proto.ScopeKind_VALUE, "Mapping<Integer>", "void"}},
			expectedScopeEntry{"mixedkeymap", expectedScope{true, proto.ScopeKind_VALUE, "Mapping<Integer>", "void"}},
			expectedScopeEntry{"mixedkeymap2", expectedScope{true, proto.ScopeKind_VALUE, "Mapping<Integer>", "void"}},
		},
		"", ""},

	scopegraphTest{"map literal nonmappable fail test", "mapliteral", "nonmappable",
		[]expectedScopeEntry{},
		"Map literal keys must be of type Stringable: Type 'SomeClass' does not define or export member 'String', which is required by type 'Stringable'", ""},

	scopegraphTest{"map literal null key fail test", "mapliteral", "nullkey",
		[]expectedScopeEntry{},
		"Map literal keys must be of type Stringable: null cannot be used in place of non-nullable type Stringable", ""},

	/////////// Slice expression ///////////

	scopegraphTest{"slice operator success test", "slice", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"slice", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"startslice", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"endslice", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"index", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"indexer nullable test", "slice", "indexernullable",
		[]expectedScopeEntry{},
		"Operator 'index' cannot be called on nullable type 'SomeClass?'", ""},

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

	scopegraphTest{"cast expr success test", "castexpr", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"cast", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"anycast", expectedScope{true, proto.ScopeKind_VALUE, "any", "void"}},
		},
		"", ""},

	scopegraphTest{"cast interfaces test", "castexpr", "interfaces",
		[]expectedScopeEntry{},
		"", ""},

	scopegraphTest{"cast struct success test", "castexpr", "structcast",
		[]expectedScopeEntry{
			expectedScopeEntry{"someField", expectedScope{true, proto.ScopeKind_VALUE, "AnotherStruct", "void"}},
			expectedScopeEntry{"bool", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"cast failure test", "castexpr", "failure",
		[]expectedScopeEntry{},
		"Cannot cast value of type 'ISomeInterface' to type 'SomeClass': Type 'SomeClass' does not define or export member 'DoSomething', which is required by type 'ISomeInterface'", ""},

	scopegraphTest{"cast nullable failure test", "castexpr", "castnull",
		[]expectedScopeEntry{},
		"Cannot cast value of type 'ISomeInterface?' to type 'SomeClass': Value may be null", ""},

	/////////// Function call expression ///////////

	scopegraphTest{"function call success test", "funccall", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"empty", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
			expectedScopeEntry{"somefunc", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"anotherfunc", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"function call missing argument test", "funccall", "missingarg",
		[]expectedScopeEntry{},
		"Function call on module member DoSomething expects 1 non-optional arguments, found 0", ""},

	scopegraphTest{"function call missing argument 2 test", "funccall", "missingarg2",
		[]expectedScopeEntry{},
		"Function call on module member someFunction expects 2 non-optional arguments, found 1", ""},

	scopegraphTest{"function call nullable access success test", "funccall", "nullaccess",
		[]expectedScopeEntry{
			expectedScopeEntry{"sm", expectedScope{true, proto.ScopeKind_VALUE, "Integer?", "void"}},
		},
		"", ""},

	scopegraphTest{"function call not function failure test", "funccall", "notfunc",
		[]expectedScopeEntry{},
		"Cannot invoke function call on non-function 'SomeClass'.", ""},

	scopegraphTest{"function call invalid count failure test", "funccall", "invalidcount",
		[]expectedScopeEntry{},
		"Function call on module member EmptyFunc expects 0 arguments, found 1", ""},

	scopegraphTest{"function call type mismatch failure test", "funccall", "typemismatch",
		[]expectedScopeEntry{},
		"Parameter #1 on module member EmptyFunc expects type Integer: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"function call nullable function failure test", "funccall", "nullablevar",
		[]expectedScopeEntry{},
		"Cannot invoke function call on non-function 'function<void>?'.", ""},

	/////////// Member access expression ///////////

	scopegraphTest{"member access success test", "memberaccess", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"varmember", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"constructor", expectedScope{true, proto.ScopeKind_VALUE, "function<SomeClass>(Integer)", "void"}},
			expectedScopeEntry{"modtype", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
			expectedScopeEntry{"modint", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"modfunc", expectedScope{true, proto.ScopeKind_VALUE, "function<void>", "void"}},
			expectedScopeEntry{"generictype", expectedScope{true, proto.ScopeKind_GENERIC, "void", "void"}},
			expectedScopeEntry{"prop", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},

			expectedScopeEntry{"genmembool", expectedScope{true, proto.ScopeKind_VALUE, "Boolean?", "void"}},
			expectedScopeEntry{"genmemint", expectedScope{true, proto.ScopeKind_VALUE, "Integer?", "void"}},
		},
		"", ""},

	scopegraphTest{"member access static under instance failure test", "memberaccess", "staticunderinstance",
		[]expectedScopeEntry{},
		"Could not find instance name 'Build' under class SomeClass", ""},

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
		"Could not find instance name 'UnknownProp' under nominal type Integer", ""},

	scopegraphTest{"member access unexported static test", "memberaccess", "unexported",
		[]expectedScopeEntry{},
		"Could not find static name 'someUnexportedThing' under import anotherpackage", ""},

	scopegraphTest{"member access unexported instance test", "memberaccess", "unexportedinstance",
		[]expectedScopeEntry{},
		"type member doSomething is not exported under class ThirdClass", ""},

	scopegraphTest{"member access miscapitalized test", "memberaccess", "miscapitalized",
		[]expectedScopeEntry{},
		"Could not find instance name 'doSomething' under class SomeClass; Did you mean 'DoSomething'?", ""},

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
		"", ""},

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

	/////////// principal literal expression ///////////

	scopegraphTest{"principal literal expression success test", "principalliteral", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"principal", expectedScope{true, proto.ScopeKind_VALUE, "SomeInterface", "void"}},
		},
		"", ""},

	scopegraphTest{"principal under class test", "principalliteral", "underclass",
		[]expectedScopeEntry{},
		"The 'principal' keyword cannot be used under non-agent class SomeClass", ""},

	scopegraphTest{"principal under module member test", "principalliteral", "modulemember",
		[]expectedScopeEntry{},
		"The 'principal' keyword cannot be used under module member DoSomething", ""},

	scopegraphTest{"principal under static member test", "principalliteral", "staticmember",
		[]expectedScopeEntry{},
		"The 'principal' keyword cannot be used under static type member Build", ""},

	scopegraphTest{"principal under property test", "principalliteral", "property",
		[]expectedScopeEntry{
			expectedScopeEntry{"scprincipal", expectedScope{true, proto.ScopeKind_VALUE, "SomeInterface", "void"}},
		},
		"", ""},

	/////////// generic specifier expression ///////////

	scopegraphTest{"generic specifier expression success test", "genericspecifier", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"someclassint", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
			expectedScopeEntry{"someclassbool", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},

			expectedScopeEntry{"somestructint", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
			expectedScopeEntry{"somestructbool", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},

			expectedScopeEntry{"someclassintbuild", expectedScope{true, proto.ScopeKind_VALUE, "function<SomeClass<Integer>>(Integer)", "void"}},
			expectedScopeEntry{"someclassboolbuild", expectedScope{true, proto.ScopeKind_VALUE, "function<SomeClass<Boolean>>(Boolean)", "void"}},

			expectedScopeEntry{"somefuncintbool", expectedScope{true, proto.ScopeKind_VALUE, "function<Integer?>(Boolean)", "void"}},
			expectedScopeEntry{"somefuncboolint", expectedScope{true, proto.ScopeKind_VALUE, "function<Boolean?>(Integer)", "void"}},
		},
		"", ""},

	scopegraphTest{"generic specifier static test", "genericspecifier", "static",
		[]expectedScopeEntry{
			expectedScopeEntry{"cons", expectedScope{true, proto.ScopeKind_VALUE, "function<SomeClass<Integer>>(Integer, Boolean)", "void"}},
			expectedScopeEntry{"sc", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass<Integer>", "void"}},
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

	scopegraphTest{"private interface member test", "types", "privateinterface",
		[]expectedScopeEntry{},
		"", ""},

	/////////// class field ///////////

	scopegraphTest{"class field uninitialized test", "var", "uninitializedfield",
		[]expectedScopeEntry{
			expectedScopeEntry{"new", expectedScope{true, proto.ScopeKind_VALUE, "function<SomeClass>(Integer, Boolean)", "void"}},
		},
		"", ""},

	/////////// lambda expression ///////////

	scopegraphTest{"lambda expression basic inference test", "lambda", "basicinference",
		[]expectedScopeEntry{
			expectedScopeEntry{"varref", expectedScope{true, proto.ScopeKind_VALUE, "function<Boolean>(Integer, String)", "void"}},
			expectedScopeEntry{"vardeclare", expectedScope{true, proto.ScopeKind_VALUE, "function<Boolean>(Boolean, Boolean)", "void"}},
			expectedScopeEntry{"callref", expectedScope{true, proto.ScopeKind_VALUE, "function<Boolean>(String, String)", "void"}},
			expectedScopeEntry{"nonref", expectedScope{true, proto.ScopeKind_VALUE, "function<Boolean>(any, any)", "void"}},
			expectedScopeEntry{"ripref", expectedScope{true, proto.ScopeKind_VALUE, "function<Integer>(Integer)", "void"}},
			expectedScopeEntry{"multiripref", expectedScope{true, proto.ScopeKind_VALUE, "function<struct>(struct)", "void"}},
		},
		"", ""},

	scopegraphTest{"lambda expression full definition test", "lambda", "full",
		[]expectedScopeEntry{
			expectedScopeEntry{"implicitreturn", expectedScope{true, proto.ScopeKind_VALUE, "function<Integer>(Integer, Boolean)", "void"}},
			expectedScopeEntry{"explicitreturn", expectedScope{true, proto.ScopeKind_VALUE, "function<String>(Integer, Boolean)", "void"}},
		},
		"", ""},

	/////////// chained inference test ///////////

	scopegraphTest{"chained inference test", "chained", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"someVar", expectedScope{true, proto.ScopeKind_VALUE, "function<void>(Integer, Boolean)", "void"}},
			expectedScopeEntry{"anotherVar", expectedScope{true, proto.ScopeKind_VALUE, "function<void>(Integer, Boolean)", "void"}},
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
		"Tagging expression for template string must be function with parameters ([]string, []stringable). Found: function<void>", ""},

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

	scopegraphTest{"webidl same name", "webidl", "samename",
		[]expectedScopeEntry{
			expectedScopeEntry{"global", expectedScope{true, proto.ScopeKind_VALUE, "Number", "void"}},
			expectedScopeEntry{"local", expectedScope{true, proto.ScopeKind_VALUE, "String", "void"}},
		},
		"", ""},

	/////////// nominal tests /////////////////

	scopegraphTest{"nominal type success", "nominal", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"this", expectedScope{true, proto.ScopeKind_VALUE, "MyType", "void"}},
			expectedScopeEntry{"sometype", expectedScope{true, proto.ScopeKind_VALUE, "SomeType", "void"}},
			expectedScopeEntry{"m", expectedScope{true, proto.ScopeKind_VALUE, "MyType", "void"}},
			expectedScopeEntry{"at", expectedScope{true, proto.ScopeKind_VALUE, "AnotherType", "void"}},
			expectedScopeEntry{"gt", expectedScope{true, proto.ScopeKind_VALUE, "GenericType<Integer>", "void"}},
			expectedScopeEntry{"nat", expectedScope{true, proto.ScopeKind_VALUE, "AnotherType?", "void"}},
		},
		"", ""},

	scopegraphTest{"nominal type over interface success", "nominal", "overinterface",
		[]expectedScopeEntry{
			expectedScopeEntry{"sn", expectedScope{true, proto.ScopeKind_VALUE, "SomeNominal", "void"}},
		},
		"", ""},

	scopegraphTest{"nominal conversion failure", "nominal", "cannotconvert",
		[]expectedScopeEntry{},
		"Cannot perform type conversion: Type 'MyType' cannot be converted to or from type 'SomeType'", ""},

	scopegraphTest{"nominal interface conversion failure", "nominal", "invalidinterface",
		[]expectedScopeEntry{},
		"Cannot perform type conversion: Type 'SomeClass' cannot be converted to or from type 'SomeNominal'", ""},

	scopegraphTest{"nominal conversion argument count mismatch", "nominal", "convertargcount",
		[]expectedScopeEntry{},
		"Type conversion requires a single argument", ""},

	scopegraphTest{"nominal conversion no argument mismatch", "nominal", "convertnoarg",
		[]expectedScopeEntry{},
		"Type conversion requires a single argument", ""},

	scopegraphTest{"nominal call base type member", "nominal", "basetypecall",
		[]expectedScopeEntry{},
		"Could not find instance name 'DoSomething' under nominal type MyType", ""},

	scopegraphTest{"nominal shortcut (wrapped in place of base)", "nominal", "wrappedinplaceofbase",
		[]expectedScopeEntry{},
		"", ""},

	/////////// structural tests /////////////////

	scopegraphTest{"structural equality test", "structnew", "equality",
		[]expectedScopeEntry{
			expectedScopeEntry{"eq", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

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
			expectedScopeEntry{"genericmodified", expectedScope{true, proto.ScopeKind_VALUE, "GenericStruct<Integer>", "void"}},
			expectedScopeEntry{"withdefaults", expectedScope{true, proto.ScopeKind_VALUE, "WithDefaults", "void"}},
			expectedScopeEntry{"withdefaults2", expectedScope{true, proto.ScopeKind_VALUE, "WithDefaults", "void"}},
		},
		"", ""},

	scopegraphTest{"structural new function test", "structnew", "function",
		[]expectedScopeEntry{
			expectedScopeEntry{"call1", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"call2", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"structural new invalid function test", "structnew", "invalidfunction",
		[]expectedScopeEntry{},
		"Structural mapping function must have 1 parameter. Found: function<Integer>", ""},

	scopegraphTest{"structural new function invalid map value test", "structnew", "functioninvalidvalue",
		[]expectedScopeEntry{},
		"Structural mapping function's parameter is Mapping with value type Boolean, but was given String: 'String' cannot be used in place of non-interface 'Boolean'", ""},

	scopegraphTest{"structural new function duplicate key test", "structnew", "functionduplicatekey",
		[]expectedScopeEntry{},
		"Structural mapping contains duplicate key: Foo", ""},

	scopegraphTest{"structural new invalid default test", "structnew", "invaliddefault",
		[]expectedScopeEntry{},
		"Field 'SomeField' has declared type 'Integer': 'String' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"structural new invalid generics test", "structnew", "invalidgenerics",
		[]expectedScopeEntry{},
		"SomeStruct<SomeClass> has non-structural generic type SomeClass: SomeClass is not structural nor serializable", ""},

	scopegraphTest{"structural new unconstructable type test", "structnew", "unconstructable",
		[]expectedScopeEntry{},
		"Cannot structurally construct type SomeInterface", ""},

	scopegraphTest{"structural new non-type test", "structnew", "nontype",
		[]expectedScopeEntry{},
		"Cannot clone and modify non-structural type Integer", ""},

	scopegraphTest{"structural new imported class test", "structnew", "importedclass",
		[]expectedScopeEntry{},
		"Cannot structurally construct type SomeClass, as it is imported from another module", ""},

	scopegraphTest{"structural new invalid name test", "structnew", "invalidname",
		[]expectedScopeEntry{},
		"Could not find instance name 'SomeField' under class SomeClass", ""},

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

	scopegraphTest{"async generic invalid test", "async", "genericinvalid",
		[]expectedScopeEntry{},
		"Asynchronous function DoSomethingAsync cannot have generics", ""},

	scopegraphTest{"async invalid return type test", "async", "invalidreturn",
		[]expectedScopeEntry{},
		"Asynchronous function DoSomethingAsync must return a structural type: SomeClass is not structural nor serializable", ""},

	scopegraphTest{"async invalid parameter test", "async", "invalidparam",
		[]expectedScopeEntry{},
		"Parameters of asynchronous function DoSomethingAsync must be structural: SomeClass is not structural nor serializable", ""},

	scopegraphTest{"async outside context test", "async", "outsidecontext",
		[]expectedScopeEntry{},
		"", "module member 'outside' is defined outside the async function and will therefore be unique for each call to this function"},

	scopegraphTest{"async non-void warning test", "async", "nonvoidwarn",
		[]expectedScopeEntry{},
		"", "Returned Awaitable resolves a value of type Integer which is not handled"},

	scopegraphTest{"async success test", "async", "success",
		[]expectedScopeEntry{},
		"", ""},

	/////////// module init tests /////////////////

	scopegraphTest{"module init basic test", "moduleinit", "basic",
		[]expectedScopeEntry{},
		"", ""},

	scopegraphTest{"module init cycle test", "moduleinit", "cycle",
		[]expectedScopeEntry{},
		"Initialization cycle found on module member foo: module member foo -> module member DoSomething -> module member DoSomethingElse -> module member foo", ""},

	/////////// agent tests /////////////////

	scopegraphTest{"agent constructor success test", "agent", "constructor",
		[]expectedScopeEntry{},
		"", ""},

	scopegraphTest{"agent constructor required field test", "agent", "requiredfield",
		[]expectedScopeEntry{},
		"", ""},

	scopegraphTest{"agent constructor fail test", "agent", "constructorfail",
		[]expectedScopeEntry{},
		"Cannot construct agent 'SomeAgent' outside its own constructor or without a composing type's context", ""},

	scopegraphTest{"agent struct constructor fail test", "agent", "structconstructfail",
		[]expectedScopeEntry{},
		"Cannot construct agent 'SomeAgent' outside its own constructor or without a composing type's context", ""},

	scopegraphTest{"agent constructor generic fail test", "agent", "genericfail",
		[]expectedScopeEntry{},
		"Cannot construct agent 'SomeAgent<Boolean>' outside its own constructor or without a composing type's context", ""},

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
		result := ParseAndBuildScopeGraph(entrypointFile, []string{}, packageloader.Library{TESTLIB_PATH, false, ""})

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
