// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

// ExpressionResult represents the result of generating an expression.
type ExpressionResult struct {
	resultExpression string               // The expression source representing the final value.
	generator        *expressionGenerator // The underlying generator.
	isPromise        bool                 // Whether the result is a promise.
}

// IsPromise returns true if the generated expression is a promise.
func (er ExpressionResult) IsPromise() bool {
	return er.isPromise
}

// IsAsync returns true if the generated expression is asynchronous.
func (er ExpressionResult) IsAsync() bool {
	return len(er.generator.wrappers) > 0
}

// ExprSource returns the source for this expression. If specified, the expression will be wrapped
// via the given template string. The original expression value reference will be placed into
// a template field with name "ResultExpr", while any data passed into this method will be placed
// into "Data".
func (er ExpressionResult) ExprSource(wrappingTemplateStr string, data interface{}) string {
	var expressionResult = er.resultExpression

	// If specified, wrap the expression.
	if wrappingTemplateStr != "" {
		fullData := struct {
			ResultExpr string
			Data       interface{}
		}{expressionResult, data}

		expressionResult = er.generator.templater.Execute("resolutionwrap", wrappingTemplateStr, fullData)
	}

	// For each expression wrapper (in *reverse order*), wrap the result expression.
	for rindex, _ := range er.generator.wrappers {
		index := len(er.generator.wrappers) - rindex - 1
		wrapper := er.generator.wrappers[index]

		data := struct {
			Item                    interface{}
			WrappedExpression       string
			WrappedNested           bool
			IntermediateExpressions []string
		}{wrapper.data, expressionResult, rindex > 0, wrapper.intermediateExpressions}

		expressionResult = er.generator.templater.Execute("expressionwrap", wrapper.templateStr, data)
	}

	return expressionResult
}
