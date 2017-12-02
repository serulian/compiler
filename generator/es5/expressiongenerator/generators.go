// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expressiongenerator

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/generator/es5/shared"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/sourceshape"
)

var _ = fmt.Printf

// wrapSynchronousExpression wraps the given synchronous expression and turns it into a promise.
func (eg *expressionGenerator) wrapSynchronousExpression(syncExpr esbuilder.ExpressionBuilder) esbuilder.ExpressionBuilder {
	// Wrap the expression in a resolve of a promise. We need the function wrapping to ensure
	// that if the expression raises an exception, we can handle that case as well.
	promiseExpr := esbuilder.Snippet(string(codedom.NewPromiseFunction)).Call(
		esbuilder.Function("",
			esbuilder.Identifier("$resolve").Call(syncExpr), "$resolve"))

	resultName := eg.generateUniqueName("$result")
	eg.addAsyncWrapper(promiseExpr, resultName)
	return esbuilder.Identifier(resultName)
}

// generateFunctionDefinition generates the code for a function.
func (eg *expressionGenerator) generateFunctionDefinition(function *codedom.FunctionDefinitionNode, context generationContext) esbuilder.ExpressionBuilder {
	templateStr := `
		{{ if .Item.WorkerExecute }}
			$t.workerwrap('{{ .Item.UniqueId }}',
		{{ end }}
		({{ if .Item.Generics }}
		  function({{ range $index, $generic := .Item.Generics }}{{ if $index }}, {{ end }}{{ $generic }}{{ end }}) {
			{{ if .Item.RequiresThis }}var $this = this;{{ end }}
			var $f =
		{{ end }}
			{{ if .Async }}
			$t.markpromising(
			{{ end }}
				function({{ range $index, $parameter := .Item.Parameters }}{{ if $index }}, {{ end }}{{ $parameter }}{{ end }}) {
					{{ if not .Item.Generics }}{{ if .Item.RequiresThis }}var $this = this;{{ end }}{{ end }}
					{{ $body := .GeneratedBody }}
					{{ if .Item.IsGenerator }}
						{{ if $body }}
							{{ emit $body }}
							return $generator.new($continue, {{ .BodyAsync }});
						{{ else }}
							return $generator.empty();
						{{ end }}						
					{{ else }}
						{{ if .Async }}
							{{ if $body }}
								{{ emit $body }}
								return $promise.new($continue);
							{{ else }}
								return $promise.empty();
							{{ end }}
						{{ else }}
							{{ if $body }}
								{{ emit $body }}
							{{ else }}
								return;
							{{ end }}
						{{ end }}
					{{ end }}
				}
			{{ if .Async }}
			)
			{{ end }}
		{{ if .Item.Generics }}
			return $f;
		  }
		{{ end }})
		{{ if .Item.WorkerExecute }}
			)
   	    {{ end }}
	`

	isAsync := function.IsAsynchronous(eg.scopegraph)
	bodyAsync := codedom.IsAsynchronous(function.Body, eg.scopegraph)
	isGenerator := function.IsGenerator()
	functionTraits := shared.FunctionTraits(isAsync, isGenerator, function.ManagesResources())

	data := struct {
		Item          *codedom.FunctionDefinitionNode
		GeneratedBody esbuilder.SourceBuilder
		Async         bool
		BodyAsync     bool
	}{function, eg.machineBuilder(function.Body, functionTraits), isAsync, bodyAsync}

	return esbuilder.Template("functiondef", templateStr, data).AsExpression()
}

// generateAwaitPromise generates the expression source for waiting for a promise.
func (eg *expressionGenerator) generateAwaitPromise(awaitPromise *codedom.AwaitPromiseNode, context generationContext) esbuilder.ExpressionBuilder {
	resultName := eg.generateUniqueName("$result")
	childExpr := eg.generateExpression(awaitPromise.ChildExpression, context)

	// Add an asynchronous wrapper for the await that executes the child expression as a promise and then
	// waits for it to return (via a call to then), at which point the wrapped expression is executed.
	if context.shortCircuiter != nil {
		eg.addAsyncWrapper(esbuilder.Binary(context.shortCircuiter, "||", childExpr), resultName)
	} else {
		eg.addAsyncWrapper(childExpr, resultName)
	}

	// An await expression is a reference to the "result" after waiting on the child expression
	// and then executing the parent expression inside the await callback.
	return esbuilder.Identifier(resultName)
}

// generateAnonymousClosureCall generates the call to an anonymous closure.
func (eg *expressionGenerator) generateAnonymousClosureCall(closureCall *codedom.AnonymousClosureCallNode, context generationContext) esbuilder.ExpressionBuilder {
	functionCall := codedom.FunctionCall(closureCall.Closure, closureCall.Arguments, closureCall.BasisNode())
	if functionCall.IsAsynchronous(eg.scopegraph) {
		functionCall = codedom.AwaitPromise(functionCall, closureCall.BasisNode())
	}

	return eg.generateExpression(
		functionCall,
		context)
}

// generateCompoundExpression generates the expression source for a compound expression.
func (eg *expressionGenerator) generateCompoundExpression(compound *codedom.CompoundExpressionNode, context generationContext) esbuilder.ExpressionBuilder {
	// Create a slice of expressions that includes an assignment to the value var
	// of the input value. This ensures that the input var has its value before the rest of
	// expressions are executed.
	inputAssignment := codedom.LocalAssignment(compound.InputVarName, compound.InputValue, compound.InputValue.BasisNode())
	fullExpressions := append([]codedom.Expression{inputAssignment}, compound.Expressions...)

	allExpressions := eg.generateExpressions(fullExpressions, context)
	outputValueExpression := eg.generateExpression(compound.OutputValue, context)

	// Add the variable to the generator's result.
	eg.variables = append(eg.variables, compound.InputVarName)

	// Return an expression list for the expression.
	return esbuilder.ExpressionList(outputValueExpression, allExpressions...)
}

// generateUnaryOperation generates the expression source for a unary operator.
func (eg *expressionGenerator) generateUnaryOperation(unaryOp *codedom.UnaryOperationNode, context generationContext) esbuilder.ExpressionBuilder {
	childExpr := eg.generateExpression(unaryOp.ChildExpression, context)
	return esbuilder.Prefix(unaryOp.Operator, childExpr)
}

// generateBinaryOperation generates the expression source for a binary operator.
func (eg *expressionGenerator) generateBinaryOperation(binaryOp *codedom.BinaryOperationNode, context generationContext) esbuilder.ExpressionBuilder {
	if binaryOp.Operator == "??" {
		return eg.generateNullComparisonOperator(binaryOp, context)
	}

	// If the binary expression's operator short circuits, then we need to generate
	// specialized wrappers and expressions to ensure that only the necessary functions get called.
	if (binaryOp.Operator == "&&" || binaryOp.Operator == "||") && binaryOp.IsAsynchronous(eg.scopegraph) {
		return eg.generateShortCircuitedBinaryOperator(binaryOp, context)
	} else {
		return eg.generateNormalBinaryOperator(binaryOp, context)
	}
}

// generateNormalBinaryOperator generates the expression source for a non-short circuiting binary operator.
func (eg *expressionGenerator) generateNormalBinaryOperator(binaryOp *codedom.BinaryOperationNode, context generationContext) esbuilder.ExpressionBuilder {
	binaryLeftExpr := eg.generateExpression(binaryOp.LeftExpr, context)
	binaryRightExpr := eg.generateExpression(binaryOp.RightExpr, context)

	return esbuilder.Binary(binaryLeftExpr, binaryOp.Operator, binaryRightExpr)
}

// generateTernary generates the expression for a ternary expression.
func (eg *expressionGenerator) generateTernary(ternary *codedom.TernaryNode, context generationContext) esbuilder.ExpressionBuilder {
	unwrappedCheckExpr := codedom.NominalUnwrapping(ternary.CheckExpr, eg.scopegraph.TypeGraph().BoolTypeReference(), ternary.CheckExpr.BasisNode())

	// If not async, generate as a direct ternary expression.
	if !ternary.IsAsynchronous(eg.scopegraph) {
		return esbuilder.Ternary(
			eg.generateExpression(unwrappedCheckExpr, context),
			eg.generateExpression(ternary.ThenExpr, context),
			eg.generateExpression(ternary.ElseExpr, context),
		)
	}

	// Generate a specialized wrapper which resolves the conditional value of the ternary and
	// places it into the result.
	resultName := eg.generateUniqueName("$result")
	resolveConditionalValue := codedom.RuntimeFunctionCall(codedom.ResolvePromiseFunction,
		[]codedom.Expression{unwrappedCheckExpr},
		ternary.BasisNode())

	eg.addAsyncWrapper(eg.generateExpression(resolveConditionalValue, context), resultName)

	// Generate the then and else expressions as short circuited by the value.
	thenExpr := eg.generateWithShortCircuiting(ternary.ThenExpr, resultName, codedom.LiteralValue("true", ternary.BasisNode()), context)
	elseExpr := eg.generateWithShortCircuiting(ternary.ElseExpr, resultName, codedom.LiteralValue("false", ternary.BasisNode()), context)

	// Return an expression which compares the value and either return the then or else values.
	return esbuilder.Ternary(esbuilder.Identifier(resultName), thenExpr, elseExpr)
}

// generateWithShortCircuiting generates an expression with automatic short circuiting based on the value found
// in the resultName variable and compared to the given compare value.
func (eg *expressionGenerator) generateWithShortCircuiting(expr codedom.Expression, resultName string, compareValue codedom.Expression, context generationContext) esbuilder.ExpressionBuilder {
	shortCircuiter := eg.generateExpression(
		codedom.RuntimeFunctionCall(codedom.ShortCircuitPromiseFunction,
			[]codedom.Expression{codedom.LocalReference(resultName, expr.BasisNode()), compareValue},
			expr.BasisNode()),
		context)

	return eg.generateExpression(expr, generationContext{shortCircuiter})
}

// generateNullComparisonOperator generates the expression for a null comparison operator.
func (eg *expressionGenerator) generateNullComparisonOperator(compareOp *codedom.BinaryOperationNode, context generationContext) esbuilder.ExpressionBuilder {
	// If not asynchronous, generate as a call to the sync nullable comparison function,
	// with the left expression and a *closure* of the right expression. The closure is necessary
	// to ensure the right hand expression is not executed unless the left hand expression is null.
	if !compareOp.IsAsynchronous(eg.scopegraph) {
		return eg.generateExpression(
			codedom.RuntimeFunctionCall(codedom.SyncNullableComparisonFunction,
				[]codedom.Expression{compareOp.LeftExpr,
					codedom.FunctionDefinition([]string{}, []string{},
						compareOp.RightExpr, false, codedom.NormalFunction, compareOp.BasisNode())},
				compareOp.BasisNode()),
			context)
	}

	return eg.generateShortCircuiter(
		compareOp.LeftExpr,
		codedom.LiteralValue("null", compareOp.BasisNode()),
		compareOp.RightExpr,
		context,

		func(resultName string, rightSide esbuilder.ExpressionBuilder) esbuilder.ExpressionBuilder {
			return esbuilder.Call(esbuilder.Snippet(string(codedom.AsyncNullableComparisonFunction)), esbuilder.Identifier(resultName), rightSide)
		})
}

// generateShortCircuitedBinaryOperator generates the expression source for a short circuiting binary operator.
func (eg *expressionGenerator) generateShortCircuitedBinaryOperator(binaryOp *codedom.BinaryOperationNode, context generationContext) esbuilder.ExpressionBuilder {
	compareValue := codedom.LiteralValue("false", binaryOp.BasisNode())
	if binaryOp.Operator == "&&" {
		compareValue = codedom.LiteralValue("true", binaryOp.BasisNode())
	}

	return eg.generateShortCircuiter(
		binaryOp.LeftExpr,
		compareValue,
		binaryOp.RightExpr,
		context,

		func(resultName string, rightSide esbuilder.ExpressionBuilder) esbuilder.ExpressionBuilder {
			return esbuilder.Binary(esbuilder.Identifier(resultName), binaryOp.Operator, rightSide)
		})
}

type shortCircuitHandler func(string, esbuilder.ExpressionBuilder) esbuilder.ExpressionBuilder

func (eg *expressionGenerator) generateShortCircuiter(compareExpr codedom.Expression,
	compareValue codedom.Expression,
	childExpr codedom.Expression,
	context generationContext,
	handler shortCircuitHandler) esbuilder.ExpressionBuilder {

	// Generate a specialized wrapper which resolves the left side value and places it into the result.
	resultName := eg.generateUniqueName("$result")
	resolveCompareValue := codedom.RuntimeFunctionCall(codedom.ResolvePromiseFunction,
		[]codedom.Expression{compareExpr},
		compareExpr.BasisNode())

	eg.addAsyncWrapper(eg.generateExpression(resolveCompareValue, context), resultName)

	shortedChildExpr := eg.generateWithShortCircuiting(childExpr, resultName, compareValue, context)
	return handler(resultName, shortedChildExpr)
}

// generateGenericSpecification generates the expression source for a generic specification.
func (eg *expressionGenerator) generateGenericSpecification(genericSpecification *codedom.GenericSpecificationNode, context generationContext) esbuilder.ExpressionBuilder {
	childExpr := eg.generateExpression(genericSpecification.ChildExpression, context)
	arguments := eg.generateExpressions(genericSpecification.TypeArguments, context)
	return esbuilder.Call(childExpr, arguments...)
}

// generateFunctionCall generates the expression source for a function call.
func (eg *expressionGenerator) generateFunctionCall(functionCall *codedom.FunctionCallNode, context generationContext) esbuilder.ExpressionBuilder {
	childExpr := eg.generateExpression(functionCall.ChildExpression, context)
	arguments := eg.generateExpressions(functionCall.Arguments, context)
	return esbuilder.Call(childExpr, arguments...)
}

// generateLocalAssignment generates the expression source for a local assignment.
func (eg *expressionGenerator) generateLocalAssignment(localAssign *codedom.LocalAssignmentNode, context generationContext) esbuilder.ExpressionBuilder {
	value := eg.generateExpression(localAssign.Value, context)
	assignment := esbuilder.Assignment(esbuilder.Identifier(localAssign.Target), value)

	// If this assignment is under an async expression wrapper, then we add it to the wrapper itself,
	// rather than doing the assignment inline. This ensures that the variable's value is updated when
	// expected in the async flow, rather than once all the promises have returned.
	if wrapper, hasWrapper := eg.currentAsyncWrapper(); hasWrapper {
		wrapper.addIntermediateExpression(assignment)
		return esbuilder.Identifier(localAssign.Target)
	}

	return assignment
}

// ByKey is a helper to sort ObjectLiteralEntryNode's.
type ByKey []codedom.ObjectLiteralEntryNode

func (s ByKey) Len() int {
	return len(s)
}

func (s ByKey) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByKey) Less(i, j int) bool {
	iKey := s[i].KeyExpression
	jKey := s[j].KeyExpression

	var iKeyString = ""
	var jKeyString = ""

	// First we try the literal value of the key. If not found, we use the node's location
	// in the source file as a stable fallback.
	if lvn, ok := iKey.(*codedom.LiteralValueNode); ok {
		iKeyString = lvn.Value
	} else {
		iKeyString = strconv.Itoa(iKey.BasisNode().GetValue(sourceshape.NodePredicateStartRune).Int())
	}

	if lvn, ok := jKey.(*codedom.LiteralValueNode); ok {
		jKeyString = lvn.Value
	} else {
		jKeyString = strconv.Itoa(jKey.BasisNode().GetValue(sourceshape.NodePredicateStartRune).Int())
	}

	return iKeyString < jKeyString
}

// generateObjectLiteral generates the expression source for a literal object value.
func (eg *expressionGenerator) generateObjectLiteral(objectLiteral *codedom.ObjectLiteralNode, context generationContext) esbuilder.ExpressionBuilder {
	// Sort the entries by key to ensure a consistent ordering.
	sortedEntries := objectLiteral.Entries
	sort.Sort(ByKey(sortedEntries))

	// Determine whether we can use the compact form of object literals. The compact form is only
	// possible if all the keys are string literals.
	var compactFormAllowed = true
	entries := make([]interface{}, len(objectLiteral.Entries))
	for index, entry := range sortedEntries {
		if _, ok := entry.KeyExpression.(*codedom.LiteralValueNode); !ok {
			compactFormAllowed = false
		}

		entries[index] = struct {
			Key   esbuilder.ExpressionBuilder
			Value esbuilder.ExpressionBuilder
		}{eg.generateExpression(entry.KeyExpression, context),
			eg.generateExpression(entry.ValueExpression, context)}
	}

	data := struct {
		Entries []interface{}
	}{entries}

	if compactFormAllowed {
		templateStr := `
			({
				{{ range $idx, $entry := .Entries }}
					{{ emit $entry.Key }}: {{ emit $entry.Value }},
				{{ end }}
			})
		`

		return esbuilder.Template("compactobjectliteral", templateStr, data).AsExpression()
	} else {
		templateStr := `
			((function() {
				var obj = {};
				{{ range $idx, $entry := .Entries }}
					obj[{{ emit $entry.Key }}] = {{ emit $entry.Value }};
				{{ end }}
				return obj;
			})())
		`

		return esbuilder.Template("expandedobjectliteral", templateStr, data).AsExpression()
	}
}

// generateArrayLiteral generates the expression source for a literal array value.
func (eg *expressionGenerator) generateArrayLiteral(arrayLiteral *codedom.ArrayLiteralNode, context generationContext) esbuilder.ExpressionBuilder {
	values := eg.generateExpressions(arrayLiteral.Values, context)
	return esbuilder.Array(values...)
}

// generateLiteralValue generates the expression source for a literal value.
func (eg *expressionGenerator) generateLiteralValue(literalValue *codedom.LiteralValueNode, context generationContext) esbuilder.ExpressionBuilder {
	return esbuilder.LiteralValue(literalValue.Value)
}

// generateTypeLiteral generates the expression source for a type literal.
func (eg *expressionGenerator) generateTypeLiteral(typeLiteral *codedom.TypeLiteralNode, context generationContext) esbuilder.ExpressionBuilder {
	return esbuilder.Snippet(eg.pather.TypeReferenceCall(typeLiteral.TypeRef))
}

// generateStaticTypeReference generates the expression source for a static type reference.
func (eg *expressionGenerator) generateStaticTypeReference(staticRef *codedom.StaticTypeReferenceNode, context generationContext) esbuilder.ExpressionBuilder {
	return esbuilder.Snippet(eg.pather.GetTypePath(staticRef.Type))
}

// generateLocalReference generates the expression source for a local reference.
func (eg *expressionGenerator) generateLocalReference(localRef *codedom.LocalReferenceNode, context generationContext) esbuilder.ExpressionBuilder {
	return esbuilder.Identifier(localRef.Name)
}

// generateDynamicAccess generates the expression source for dynamic access.
func (eg *expressionGenerator) generateDynamicAccess(dynamicAccess *codedom.DynamicAccessNode, context generationContext) esbuilder.ExpressionBuilder {
	basisNode := dynamicAccess.BasisNode()

	isPromisingLiteral := "false"
	if dynamicAccess.IsPossiblyPromising {
		isPromisingLiteral = "true"
	}

	funcCall := codedom.RuntimeFunctionCall(
		codedom.DynamicAccessFunction,
		[]codedom.Expression{
			dynamicAccess.ChildExpression,
			codedom.LiteralValue("'"+dynamicAccess.Name+"'", basisNode),
			codedom.LiteralValue(isPromisingLiteral, basisNode),
		},
		basisNode,
	)

	if dynamicAccess.IsPossiblyPromising {
		funcCall = codedom.AwaitPromise(funcCall, basisNode)
	}

	return eg.generateExpression(funcCall, context)
}

// generateNestedTypeAccess generates the expression source for a nested type access.
func (eg *expressionGenerator) generateNestedTypeAccess(nestedAccess *codedom.NestedTypeAccessNode, context generationContext) esbuilder.ExpressionBuilder {
	childExpr := eg.generateExpression(nestedAccess.ChildExpression, context)
	return childExpr.Member(eg.pather.InnerInstanceName(nestedAccess.InnerType))
}

// generateStaticMemberReference generates the expression for a static reference to a module or type member.
func (eg *expressionGenerator) generateStaticMemberReference(memberReference *codedom.StaticMemberReferenceNode, context generationContext) esbuilder.ExpressionBuilder {
	staticPath := eg.pather.GetStaticMemberPath(memberReference.Member, memberReference.ParentType)
	return esbuilder.Snippet(staticPath)
}

// generateRuntimeFunctionCall generates the expression source for a call to a runtime function.
func (eg *expressionGenerator) generateRuntimeFunctionCall(runtimeCall *codedom.RuntimeFunctionCallNode, context generationContext) esbuilder.ExpressionBuilder {
	arguments := eg.generateExpressions(runtimeCall.Arguments, context)
	return esbuilder.Call(esbuilder.Snippet(string(runtimeCall.Function)), arguments...)
}

// generateNativeAccess generates the expression source for a native assign.
func (eg *expressionGenerator) generateNativeAssign(nativeAssign *codedom.NativeAssignNode, context generationContext) esbuilder.ExpressionBuilder {
	target := eg.generateExpression(nativeAssign.TargetExpression, context)
	value := eg.generateExpression(nativeAssign.ValueExpression, context)
	return esbuilder.Assignment(target, value)
}

// generateNativeAccess generates the expression source for a native access to a name under an expression.
func (eg *expressionGenerator) generateNativeAccess(nativeAccess *codedom.NativeAccessNode, context generationContext) esbuilder.ExpressionBuilder {
	childExpr := eg.generateExpression(nativeAccess.ChildExpression, context)
	return childExpr.Member(nativeAccess.Name)
}

// generateNativeMemberAccess generates the expression source for a native access to a defined member.
func (eg *expressionGenerator) generateNativeMemberAccess(nativeAccess *codedom.NativeMemberAccessNode, context generationContext) esbuilder.ExpressionBuilder {
	childExpr := eg.generateExpression(nativeAccess.ChildExpression, context)
	return childExpr.Member(nativeAccess.NativeName)
}

// generateNativeIndexing generates the expression source for a native index on an expression.
func (eg *expressionGenerator) generateNativeIndexing(nativeIndex *codedom.NativeIndexingNode, context generationContext) esbuilder.ExpressionBuilder {
	childExpr := eg.generateExpression(nativeIndex.ChildExpression, context)
	indexExpr := eg.generateExpression(nativeIndex.IndexExpression, context)
	return childExpr.Member(indexExpr)
}

// generateNominalWrapping generates the expression source for the nominal wrapping of an instance of a base type.
func (eg *expressionGenerator) generateNominalWrapping(nominalWrapping *codedom.NominalWrappingNode, context generationContext) esbuilder.ExpressionBuilder {
	// If this is a wrap is of an unwrap, then cancel both operations if we are simply rewrapping
	// the unwrapped type.
	if unwrapping, ok := nominalWrapping.ChildExpression.(*codedom.NominalUnwrappingNode); ok {
		if unwrapping.ChildExpressionType.NominalDataType() == nominalWrapping.NominalTypeRef {
			// Skip entirely.
			return eg.generateExpression(unwrapping.ChildExpression, context)
		}
	}

	// If the child expression being wrapped is non-nullable and not a nominal or struct,
	// then we can use a fast-path box call without the extra unboxing.
	boxFunction := codedom.BoxFunction
	if !nominalWrapping.ChildExpressionType.NullValueAllowed() &&
		(!nominalWrapping.ChildExpressionType.IsNominalOrStruct() || nominalWrapping.IsLiteralWrap) {
		boxFunction = codedom.FastBoxFunction
	}

	call := codedom.RuntimeFunctionCall(
		boxFunction,
		[]codedom.Expression{
			nominalWrapping.ChildExpression,
			codedom.TypeLiteral(nominalWrapping.NominalTypeRef, nominalWrapping.BasisNode())},
		nominalWrapping.BasisNode())
	return eg.generateExpression(call, context)
}

// generateNominalUnwrapping generates the expression source for the unwrapping of a nominal instance of a base type.
func (eg *expressionGenerator) generateNominalUnwrapping(nominalUnwrapping *codedom.NominalUnwrappingNode, context generationContext) esbuilder.ExpressionBuilder {
	// If this is an unwrap is of a wrap, then try to either cancel both operations (if we can),
	// or simplify the unwrapping operation.
	childExpression := nominalUnwrapping.ChildExpression
	if wrapping, ok := nominalUnwrapping.ChildExpression.(*codedom.NominalWrappingNode); ok {
		// If the value that was being wrapped is, itself, nominal or structural, then we still
		// need to unwrap it, but we can fast-path by ignoring the wrap and directly unwrapping
		// the nominal/struct. Otherwise, we can simply just collapse the unwrap+wrap into a NOOP.
		if !wrapping.IsLiteralWrap && wrapping.ChildExpressionType.IsNominalOrStruct() {
			childExpression = wrapping.ChildExpression
		} else {
			// Otherwise, we can just collapse both operations into nothing.
			return eg.generateExpression(wrapping.ChildExpression, context)
		}
	}

	// If the child expression being unwrapped is non-nullable and we know it is boxed, then
	// we can use a fast-path unbox call without the extra checks.
	if !nominalUnwrapping.ChildExpressionType.NullValueAllowed() &&
		nominalUnwrapping.ChildExpressionType.IsNominalOrStruct() {
		access := codedom.NativeAccess(
			childExpression,
			codedom.BoxedDataProperty,
			nominalUnwrapping.BasisNode())
		return eg.generateExpression(access, context)
	}

	call := codedom.RuntimeFunctionCall(
		codedom.UnboxFunction,
		[]codedom.Expression{
			childExpression,
		},
		nominalUnwrapping.BasisNode())
	return eg.generateExpression(call, context)
}

// generateMemberAssignment generates the expression source for a member assignment.
func (eg *expressionGenerator) generateMemberAssignment(memberAssign *codedom.MemberAssignmentNode, context generationContext) esbuilder.ExpressionBuilder {
	basisNode := memberAssign.BasisNode()

	// If the target member is an operator, then we need to invoke it as a function call, with the first
	// argument being the argument to the child call, and the second argument being the assigned child
	// expression.
	if memberAssign.Target.IsOperator() {
		childCall := memberAssign.NameExpression.(*codedom.MemberCallNode)
		memberRef := childCall.ChildExpression.(*codedom.MemberReferenceNode)

		// If this is a native operator, change it into a native indexing and assignment.
		if memberAssign.Target.IsNative() {
			nativeAssign := codedom.NativeAssign(
				codedom.NativeIndexing(memberRef.ChildExpression,
					childCall.Arguments[0], basisNode),
				memberAssign.Value,
				basisNode)

			return eg.generateExpression(nativeAssign, context)
		} else {
			memberCall := codedom.MemberCall(
				codedom.NativeMemberAccess(memberRef.ChildExpression, eg.pather.GetMemberName(memberAssign.Target), memberAssign.Target, memberRef.BasisNode()),
				memberAssign.Target,
				[]codedom.Expression{childCall.Arguments[0], memberAssign.Value},
				basisNode)

			return eg.generateExpression(memberCall, context)
		}
	}

	// If the target member is implicitly called, then this is a property that needs to be assigned via a call.
	if memberAssign.Target.IsImplicitlyCalled() {
		memberRef := memberAssign.NameExpression.(*codedom.MemberReferenceNode)

		memberCall := codedom.MemberCallWithCallType(
			codedom.NativeMemberAccess(memberRef.ChildExpression, eg.pather.GetSetterName(memberRef.Member), memberRef.Member, memberRef.BasisNode()),
			memberAssign.Target,
			[]codedom.Expression{memberAssign.Value},
			scopegraph.PromisingAccessImplicitSet,
			basisNode)

		return eg.generateExpression(memberCall, context)
	}

	value := eg.generateExpression(memberAssign.Value, context)
	targetExpr := eg.generateExpression(memberAssign.NameExpression, context)
	return esbuilder.Assignment(targetExpr, value)
}

// generateMemberReference generates the expression for a reference to a module or type member.
func (eg *expressionGenerator) generateMemberReference(memberReference *codedom.MemberReferenceNode, context generationContext) esbuilder.ExpressionBuilder {
	// If the access is nullable, generate as a dynamic access, since it handles nullables for us.
	if memberReference.Nullable {
		return eg.generateExpression(
			codedom.DynamicAccess(
				memberReference.ChildExpression,
				memberReference.ExprName(),
				memberReference.IsAsynchronous(eg.scopegraph),
				memberReference.BasisNode()),
			context)
	}

	// If the target member is implicitly called, then this is a property that needs to be accessed via a call.
	if memberReference.Member.IsImplicitlyCalled() {
		basisNode := memberReference.BasisNode()
		memberCall := codedom.MemberCallWithCallType(
			codedom.NativeMemberAccess(memberReference.ChildExpression, memberReference.Member.Name(), memberReference.Member, basisNode),
			memberReference.Member,
			[]codedom.Expression{},
			scopegraph.PromisingAccessImplicitGet,
			basisNode)

		return eg.generateExpression(memberCall, context)
	}

	// This handles the native new case for WebIDL. We should probably handle this directly.
	// TODO: generalize this.
	if memberReference.Member.IsStatic() {
		if memberReference.Member.SourceGraphId() == "webidl" {
			return eg.generateExpression(codedom.StaticMemberReference(memberReference.Member, eg.scopegraph.TypeGraph().AnyTypeReference(), memberReference.BasisNode()), context)
		}
	}

	childExpr := eg.generateExpression(memberReference.ChildExpression, context)
	return childExpr.Member(eg.pather.GetMemberName(memberReference.Member))
}

// generateMemberCall generates the expression source for a call to a module or type member.
func (eg *expressionGenerator) generateMemberCall(memberCall *codedom.MemberCallNode, context generationContext) esbuilder.ExpressionBuilder {
	if memberCall.Member.IsOperator() && memberCall.Member.IsNative() {
		// This is a call to a native operator.
		if memberCall.Member.Name() != "index" {
			panic("Native call to non-index operator")
		}

		refExpr := memberCall.ChildExpression.(*codedom.MemberReferenceNode).ChildExpression
		return eg.generateExpression(codedom.NativeIndexing(refExpr, memberCall.Arguments[0], memberCall.BasisNode()), context)
	}

	callPath := memberCall.ChildExpression
	arguments := memberCall.Arguments

	var functionCall = codedom.InvokeFunction(callPath, arguments, memberCall.CallType, eg.scopegraph, memberCall.BasisNode())

	// Special case: A function call invoke on an otherwise nullable function.
	if memberCall.Nullable {
		// Invoke the function with a specialized nullable-invoke.
		refExpr := callPath.(*codedom.MemberReferenceNode).ChildExpression

		var isPromising = "false"
		if eg.scopegraph.IsPromisingMember(memberCall.Member, memberCall.CallType) {
			isPromising = "true"
		}

		localArguments := []codedom.Expression{
			refExpr,
			codedom.LiteralValue("'"+memberCall.Member.Name()+"'", refExpr.BasisNode()),
			codedom.LiteralValue(isPromising, memberCall.BasisNode()),
			codedom.ArrayLiteral(arguments, memberCall.BasisNode()),
		}

		functionCall = codedom.RuntimeFunctionCall(codedom.NullableInvokeFunction, localArguments, memberCall.BasisNode())
		if eg.scopegraph.IsPromisingMember(memberCall.Member, memberCall.CallType) {
			functionCall = codedom.AwaitPromise(functionCall, memberCall.BasisNode())
		}
	}

	return eg.generateExpression(functionCall, context)
}
