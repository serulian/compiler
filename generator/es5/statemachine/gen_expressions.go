// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/generator/es5/codedom"
)

// generateFunctionDefinition generates the code for a function.
func (eg *expressionGenerator) generateFunctionDefinition(function *codedom.FunctionDefinitionNode) string {
	if function == nil || function.Body == nil {
		panic("Nil function")
	}

	templateStr := `
		({{ if .Item.Generics }}
		  function({{ range $index, $generic := .Item.Generics }}{{ if $index }}, {{ end }}{{ $generic }}{{ end }}) {
			{{ if .Item.RequiresThis }}var $this = this;{{ end }}
			var $f =
		{{ end }}
				function({{ range $index, $parameter := .Item.Parameters }}{{ if $index }}, {{ end }}{{ $parameter }}{{ end }}) {
					{{ if not .Item.Generics }}{{ if .Item.RequiresThis }}var $this = this;{{ end }}{{ end }}
					{{ $body := .Generator.GenerateMachine .Item.Body }}
					{{ if $body }}
						{{ $body }}
						return $promise.build($state);
					{{ else }}
						return $promise.empty();
					{{ end }}
				}
		{{ if .Item.Generics }}
			return $f;
		  }
		{{ end }})
	`

	stateGenerator := buildGenerator(eg.templater, eg.pather, eg.scopegraph)
	return eg.templater.Execute("functiondef", templateStr, generatingItem{function, stateGenerator})
}

// generateAwaitPromise generates the expression source for waiting for a promise.
func (eg *expressionGenerator) generateAwaitPromise(awaitPromise *codedom.AwaitPromiseNode) (string, expressionWrapper) {
	childExpr := eg.generateExpression(awaitPromise.ChildExpression)
	resultName := eg.generateUniqueName("$result")

	data := struct {
		ChildExpression string
		ResultName      string
	}{childExpr, resultName}

	templateStr := `
		({{ .Item.ChildExpression }}).then(function({{ .Item.ResultName }}) {
			{{ if .WrappedNested }}
				return ({{ .WrappedExpression }});
			{{ else }}
				{{ .WrappedExpression }}
			{{ end }}
		})
	`

	// An await expression is a reference to the "result" after waiting on the child expression
	// and then executing the parent expression inside the await callback.
	return resultName, expressionWrapper{
		data:        data,
		templateStr: templateStr,
	}
}

// generateUnaryOperation generates the expression source for a unary operator.
func (eg *expressionGenerator) generateUnaryOperation(unaryOp *codedom.UnaryOperationNode) string {
	childExpr := eg.generateExpression(unaryOp.ChildExpression)

	data := struct {
		ChildExpression string
		Operator        string
	}{childExpr, unaryOp.Operator}

	templateStr := `
		({{ .Operator }}({{ .ChildExpression }}))
	`

	return eg.templater.Execute("unaryop", templateStr, data)
}

// generateBinaryOperation generates the expression source for a binary operator.
func (eg *expressionGenerator) generateBinaryOperation(binaryOp *codedom.BinaryOperationNode) string {
	leftExpr := eg.generateExpression(binaryOp.LeftExpr)
	rightExpr := eg.generateExpression(binaryOp.RightExpr)

	data := struct {
		LeftExpression  string
		RightExpression string
		Operator        string
	}{leftExpr, rightExpr, binaryOp.Operator}

	templateStr := `
		(({{ .LeftExpression }}) {{ .Operator }} ({{ .RightExpression }}))
	`

	return eg.templater.Execute("binaryop", templateStr, data)
}

// generateFunctionCall generates the expression soruce for a function call.
func (eg *expressionGenerator) generateFunctionCall(functionCall *codedom.FunctionCallNode) string {
	childExpr := eg.generateExpression(functionCall.ChildExpression)
	arguments := eg.generateExpressions(functionCall.Arguments)

	data := struct {
		ChildExpression string
		Arguments       []string
	}{childExpr, arguments}

	templateStr := `
		({{ .ChildExpression }})({{ range $index, $arg := .Arguments }}{{ if $index }}, {{ end }}{{ $arg }}{{ end }})
	`

	return eg.templater.Execute("functioncall", templateStr, data)
}

// generateMemberAssignment generates the expression source for a member assignment.
func (eg *expressionGenerator) generateMemberAssignment(memberAssign *codedom.MemberAssignmentNode) string {
	basisNode := memberAssign.BasisNode()

	// If the member is an extension member, then we have to invoke it like a function call with the instance
	// being given as the first parameter and the value as the second.
	if memberAssign.Target.IsExtension() {
		childExpr := memberAssign.NameExpression.(*codedom.MemberReferenceNode).ChildExpression
		memberCall := codedom.MemberCall(
			childExpr,
			memberAssign.Target,
			[]codedom.Expression{memberAssign.Value},
			basisNode)

		return eg.generateExpression(codedom.WrapIfPromising(memberCall, memberAssign.Target, basisNode))
	}

	// If the target member is implicitly called, then this is a property that needs to be assigned via a call.
	if memberAssign.Target.IsImplicitlyCalled() {
		memberRef := memberAssign.NameExpression.(*codedom.MemberReferenceNode)

		memberCall := codedom.MemberCall(
			codedom.NativeAccess(memberRef.ChildExpression, memberRef.Member.Name(), memberRef.BasisNode()),
			memberAssign.Target,
			[]codedom.Expression{memberAssign.Value},
			basisNode)

		return eg.generateExpression(memberCall)
	}

	value := eg.generateExpression(memberAssign.Value)
	targetExpr := eg.generateExpression(memberAssign.NameExpression)

	data := struct {
		TargetExpression string
		ValueExpression  string
	}{targetExpr, value}

	templateStr := `
		({{ .TargetExpression }} = ({{ .ValueExpression }}))
	`

	return eg.templater.Execute("memberassignment", templateStr, data)
}

// generateLocalAssignment generates the expression source for a local assignment.
func (eg *expressionGenerator) generateLocalAssignment(localAssign *codedom.LocalAssignmentNode) string {
	value := eg.generateExpression(localAssign.Value)

	data := struct {
		TargetName      string
		ValueExpression string
	}{localAssign.Target, value}

	templateStr := `
		({{ .TargetName }} = ({{ .ValueExpression }}))
	`

	return eg.templater.Execute("localassignment", templateStr, data)
}

// generateLiteralValue generates the expression source for a literal value.
func (eg *expressionGenerator) generateLiteralValue(literalValue *codedom.LiteralValueNode) string {
	return literalValue.Value
}

// generateTypeLiteral generates the expression source for a type literal.
func (eg *expressionGenerator) generateTypeLiteral(typeLiteral *codedom.TypeLiteralNode) string {
	return eg.pather.TypeReferenceCall(typeLiteral.TypeRef)
}

// generateStaticTypeReference generates the expression source for a static type reference.
func (eg *expressionGenerator) generateStaticTypeReference(staticRef *codedom.StaticTypeReferenceNode) string {
	return eg.pather.GetTypePath(staticRef.Type)
}

// generateLocalReference generates the expression source for a local reference.
func (eg *expressionGenerator) generateLocalReference(localRef *codedom.LocalReferenceNode) string {
	return localRef.Name
}

// generateDynamicAccess generates the expression source for dynamic access.
func (eg *expressionGenerator) generateDynamicAccess(dynamicAccess *codedom.DynamicAccessNode) string {
	basisNode := dynamicAccess.BasisNode()

	funcCall := codedom.RuntimeFunctionCall(
		codedom.DynamicAccessFunction,
		[]codedom.Expression{
			dynamicAccess.ChildExpression,
			codedom.LiteralValue("'"+dynamicAccess.Name+"'", basisNode),
		},
		basisNode,
	)

	return eg.generateExpression(funcCall)
}

// generateNestedTypeAccess generates the expression source for a nested type access.
func (eg *expressionGenerator) generateNestedTypeAccess(nestedAccess *codedom.NestedTypeAccessNode) string {
	childExpr := eg.generateExpression(nestedAccess.ChildExpression)

	data := struct {
		ChildExpression string
		InnerTypePath   string
	}{childExpr, eg.pather.InnerInstanceName(nestedAccess.InnerType)}

	templateStr := `
		({{ .ChildExpression }}).{{ .InnerTypePath }}
	`

	return eg.templater.Execute("nestedtype", templateStr, data)
}

// generateMemberReference generates the expression for a reference to a module or type member.
func (eg *expressionGenerator) generateMemberReference(memberReference *codedom.MemberReferenceNode) string {
	// If the target member is implicitly called, then this is a property that needs to be accessed via a call.
	if memberReference.Member.IsImplicitlyCalled() {
		basisNode := memberReference.BasisNode()

		// Note: If the member is an extension member, then we don't add the property name, as it'll be done
		// for us in the MemberCall transformation.
		var memberAccess = codedom.NativeAccess(memberReference.ChildExpression, memberReference.Member.Name(), basisNode)
		if memberReference.Member.IsExtension() {
			memberAccess = memberReference.ChildExpression
		}

		memberCall := codedom.MemberCall(
			memberAccess,
			memberReference.Member,
			[]codedom.Expression{},
			basisNode)

		return eg.generateExpression(memberCall)
	}

	// This handles the native new case for WebIDL. We should probably handle this directly.
	if memberReference.Member.IsStatic() && !memberReference.Member.IsPromising() {
		return eg.generateExpression(codedom.StaticMemberReference(memberReference.Member, memberReference.BasisNode()))
	}

	childExpr := eg.generateExpression(memberReference.ChildExpression)

	data := struct {
		ChildExpression string
		MemberName      string
	}{childExpr, eg.pather.GetMemberName(memberReference.Member)}

	templateStr := `
		({{ .ChildExpression }}).{{ .MemberName }}
	`

	return eg.templater.Execute("memberref", templateStr, data)
}

// generateStaticMemberReference generates the expression for a static reference to a module or type member.
func (eg *expressionGenerator) generateStaticMemberReference(memberReference *codedom.StaticMemberReferenceNode) string {
	data := struct {
		StaticPath string
	}{eg.pather.GetStaticMemberPath(memberReference.Member, eg.scopegraph.TypeGraph().AnyTypeReference())}

	templateStr := `
		{{ .StaticPath }}
	`

	return eg.templater.Execute("staticmemberref", templateStr, data)
}

// generateRuntineFunctionCall generates the expression source for a call to a runtime function.
func (eg *expressionGenerator) generateRuntineFunctionCall(runtimeCall *codedom.RuntimeFunctionCallNode) string {
	arguments := eg.generateExpressions(runtimeCall.Arguments)

	data := struct {
		RuntimeFunction string
		Arguments       []string
	}{string(runtimeCall.Function), arguments}

	templateStr := `
		{{ .RuntimeFunction }}({{ range $index, $arg := .Arguments }}{{ if $index }}, {{ end }}{{ $arg }}{{ end }})
	`

	return eg.templater.Execute("runtimecall", templateStr, data)
}

// generateNativeAccess generates the expression source for a native assign.
func (eg *expressionGenerator) generateNativeAssign(nativeAssign *codedom.NativeAssignNode) string {
	target := eg.generateExpression(nativeAssign.TargetExpression)
	value := eg.generateExpression(nativeAssign.ValueExpression)

	data := struct {
		TargetExpression string
		ValueExpression  string
	}{target, value}

	templateStr := `
		{{ .TargetExpression }} = {{ .ValueExpression }}
	`

	return eg.templater.Execute("nativeassign", templateStr, data)
}

// generateNativeAccess generates the expression source for a native access to a member.
func (eg *expressionGenerator) generateNativeAccess(nativeAccess *codedom.NativeAccessNode) string {
	childExpr := eg.generateExpression(nativeAccess.ChildExpression)

	data := struct {
		ChildExpression string
		Name            string
	}{childExpr, nativeAccess.Name}

	templateStr := `
		({{ .ChildExpression }}).{{ .Name }}
	`

	return eg.templater.Execute("nativeaccess", templateStr, data)
}

// generateMemberCall generates the expression source for a call to a module or type member.
func (eg *expressionGenerator) generateMemberCall(memberCall *codedom.MemberCallNode) string {
	basisNode := memberCall.BasisNode()
	staticMemberPath := eg.pather.GetStaticMemberPath(memberCall.Member, eg.scopegraph.TypeGraph().AnyTypeReference())

	var callPath codedom.Expression = nil
	var arguments []codedom.Expression = []codedom.Expression{}

	if memberCall.Member.IsExtension() {
		// If the member is an extension member, we call it statically, with the instance as the first param.
		var childExpr = memberCall.ChildExpression
		if refNode, ok := childExpr.(*codedom.MemberReferenceNode); ok {
			childExpr = refNode.ChildExpression
		}

		callPath = codedom.LocalReference(staticMemberPath, basisNode)
		arguments = append([]codedom.Expression{childExpr}, memberCall.Arguments...)
	} else {
		// Otherwise this is a normal function call over a child expression.
		callPath = memberCall.ChildExpression
		arguments = memberCall.Arguments
	}

	functionCall := codedom.FunctionCall(callPath, arguments, basisNode)
	return eg.generateExpression(codedom.WrapIfPromising(functionCall, memberCall.Member, basisNode))
}
