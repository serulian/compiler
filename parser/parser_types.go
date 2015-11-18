// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

//go:generate stringer -type=NodeType

import (
	"strconv"

	"github.com/serulian/compiler/compilercommon"
)

type AstNode interface {
	// Connect connects this AstNode to another AstNode with the given predicate,
	// and returns the same AstNode.
	Connect(predicate string, other AstNode) AstNode

	// Decorate decorates this AstNode with the given property and string value,
	// and returns the same AstNode.
	Decorate(property string, value string) AstNode
}

// PackageImportType identifies the types of imports.
type PackageImportType int

const (
	ImportTypeLocal PackageImportType = iota
	ImportTypeVCS
)

// PackageImport defines the import of a package.
type PackageImport struct {
	Path           string
	ImportType     PackageImportType
	SourceLocation compilercommon.SourceAndLocation
}

// NodeType identifies the type of AST node.
type NodeType int

const (
	// Top-level
	NodeTypeError   NodeType = iota // error occurred; value is text of error
	NodeTypeFile                    // The file root node
	NodeTypeComment                 // A single or multiline comment

	// Decorator
	NodeTypeDecorator

	// Module-level
	NodeTypeImport    // An import
	NodeTypeClass     // A class
	NodeTypeInterface // An interface
	NodeTypeGeneric   // A generic definition on a type

	// Module and Type Members
	NodeTypeFunction // A function declaration or definition
	NodeTypeVariable // A variable definition at the module level.

	// Type Members
	NodeTypeConstructor // A constructor declaration or definition
	NodeTypeProperty    // A property declaration or definition
	NodeTypeOperator    // An operator declaration or definition
	NodeTypeField       // A field (var) definition

	// Type member blocks
	NodeTypePropertyBlock // A child block (get or set) of a property definition
	NodeTypeParameter     // A parameter under a type member (function, iterator, etc)

	// Statements
	NodeTypeStatementBlock       // A block of statements
	NodeTypeLoopStatement        // A for statement
	NodeTypeConditionalStatement // An if statement
	NodeTypeReturnStatement      // A return statement
	NodeTypeBreakStatement       // A break statement
	NodeTypeContinueStatement    // A continue statement
	NodeTypeVariableStatement    // A variable statement
	NodeTypeWithStatement        // A with statement
	NodeTypeMatchStatement       // A match statement
	NodeTypeAssignStatement      // An assignment state: a = b

	NodeTypeMatchStatementCase // A case of a match statement.

	NodeTypeNamedValue // A named value added to the scope of the parent statement.

	// Expressions
	NodeTypeAwaitExpression // An await expression: <- a
	NodeTypeArrowExpression // An arrow expression: a <- b

	NodeTypeLambdaExpression // A lambda expression

	NodeBitwiseXorExpression        // a ^ b
	NodeBitwiseOrExpression         // a | b
	NodeBitwiseAndExpression        // a & b
	NodeBitwiseShiftLeftExpression  // a << b
	NodeBitwiseShiftRightExpression // a >> b
	NodeBitwiseNotExpression        // ~a

	NodeBooleanOrExpression  // a || b
	NodeBooleanAndExpression // a && b
	NodeBooleanNotExpression // !a

	NodeComparisonEqualsExpression    // a == b
	NodeComparisonNotEqualsExpression // a != b

	NodeComparisonLTEExpression // a <= b
	NodeComparisonGTEExpression // a >= b
	NodeComparisonLTExpression  // a < b
	NodeComparisonGTExpression  // a > b

	NodeNullComparisonExpression // a ?? b

	NodeDefineRangeExpression // a .. b

	NodeBinaryAddExpression      // A plus expression: +
	NodeBinarySubtractExpression // A subtract expression: -
	NodeBinaryMultiplyExpression // A multiply expression: *
	NodeBinaryDivideExpression   // A divide expression: /
	NodeBinaryModuloExpression   // A modulo expression: %

	NodeMemberAccessExpression         // a.b
	NodeNullableMemberAccessExpression // a?.b
	NodeDynamicMemberAccessExpression  // a->b
	NodeStreamMemberAccessExpression   // a*.b
	NodeCastExpression                 // a.(b)
	NodeFunctionCallExpression         // a(b, c, d)
	NodeSliceExpression                // a[b:c]
	NodeGenericSpecifierExpression     // a<b>

	NodeNumericLiteralExpression        // 123
	NodeStringLiteralExpression         // 'hello'
	NodeBooleanLiteralExpression        // true
	NodeTemplateStringLiteralExpression // `foobar`
	NodeThisLiteralExpression           // this
	NodeNullLiteralExpression           // null

	NodeListExpression     // [1, 2, 3]
	NodeMapExpression      // {a: 1, b: 2}
	NodeMapExpressionEntry // a: 1

	NodeTypeIdentifierExpression // An identifier expression

	// Lambda expression members
	NodeTypeLambdaParameter

	// Type references
	NodeTypeTypeReference // A type reference
	NodeTypeStream
	NodeTypeNullable
	NodeTypeVoid

	// Misc
	NodeTypeIdentifierPath   // An identifier path
	NodeTypeIdentifierAccess // A named reference via an identifier or a dot access

	// NodeType is a tagged type.
	NodeTypeTagged
)

func (t NodeType) Name() string {
	return "NodeType"
}

func (t NodeType) Value() string {
	return strconv.Itoa(int(t))
}

func (t NodeType) Build(value string) interface{} {
	i, err := strconv.Atoi(value)
	if err != nil {
		panic("Invalid value for NodeType: " + value)
	}
	return NodeType(i)
}

const (
	//
	// All nodes
	//
	// The source of this node.
	NodePredicateSource = "input-source"

	// The rune position in the input string at which this node begins.
	NodePredicateStartRune = "start-rune"

	// The rune position in the input string at which this node ends.
	NodePredicateEndRune = "end-rune"

	// A direct child of this node. Implementations should handle the ordering
	// automatically for this predicate.
	NodePredicateChild = "child-node"

	//
	// NodeTypeError
	//

	// The message for the parsing error.
	NodePredicateErrorMessage = "error-message"

	//
	// NodeTypeComment
	//

	// The value of the comment, including its delimeter(s)
	NodeCommentPredicateValue = "comment-value"

	//
	// NodeTypeDecorator
	//
	NodeDecoratorPredicateInternal  = "decorator-internal"
	NodeDecoratorPredicateParameter = "decorator-parameter"

	//
	// NodeTypeImport
	//
	NodeImportPredicateKind        = "import-kind"
	NodeImportPredicateSource      = "import-source"
	NodeImportPredicateSubsource   = "import-subsource"
	NodeImportPredicateName        = "named"
	NodeImportPredicatePackageName = "import-package-name"

	NodeImportPredicateLocation = "import-location-ref"

	//
	// NodeTypeClass + NodeTypeInterface
	//
	NodeTypeDefinitionGeneric   = "type-generic"
	NodeTypeDefinitionMember    = "type-member"
	NodeTypeDefinitionName      = "named"
	NodeTypeDefinitionDecorator = "decorator"

	//
	// NodeTypeClass
	//
	NodeClassPredicateBaseType = "class-basetypepath"

	//
	// NodeTypeGeneric
	//
	NodeGenericPredicateName = "named"
	NodeGenericSubtype       = "generic-subtype"

	//
	// Type members and properties
	//
	NodePredicateBody = "definition-body"

	//
	// Type members: NodeTypeProperty, NodeTypeFunction, NodeTypeField, NodeTypeConstructor
	//
	NodePredicateTypeMemberName         = "named"
	NodePredicateTypeMemberDeclaredType = "typemember-declared-type"
	NodePredicateTypeMemberReturnType   = "typemember-return-type"
	NodePredicateTypeMemberGeneric      = "typemember-generic"
	NodePredicateTypeMemberParameter    = "typemember-parameter"

	//
	// NodeTypeOperator
	//
	NodeOperatorName = "operator-named"

	//
	// NodeTypeProperty
	//
	NodePropertyReadOnly = "typemember-readonly"
	NodePropertyGetter   = "property-getter"
	NodePropertySetter   = "property-setter"

	//
	// NodeTypePropertyBlock
	//
	NodePropertyBlockType = "propertyblock-type"

	//
	// NodeTypeParameter
	//
	NodeParameterType = "parameter-type"
	NodeParameterName = "named"

	//
	// NodeTypeTypeReference
	//
	NodeTypeReferencePath      = "typereference-path"
	NodeTypeReferenceGeneric   = "typereference-generic"
	NodeTypeReferenceParameter = "typereference-parameter"
	NodeTypeReferenceInnerType = "typereference-inner-type"

	//
	// NodeTypeIdentifierPath
	//
	NodeIdentifierPathRoot = "identifierpath-root"

	//
	// NodeTypeIdentifierAccess
	//
	NodeIdentifierAccessName   = "identifieraccess-name"
	NodeIdentifierAccessSource = "identifieraccess-source"

	//
	// NodeTypeStatementBlock
	//
	NodeStatementBlockStatement = "block-child"
	NodeStatementLabel          = "statement-label"

	//
	// NodeTypeLoopStatement + NodeTypeWithStatement
	//
	NodeStatementNamedValue = "named-value"

	//
	// NodeTypeLoopStatement
	//
	NodeLoopStatementExpression = "loop-expression"
	NodeLoopStatementBlock      = "loop-block"

	//
	// NodeTypeAssignStatement
	//
	NodeAssignStatementName  = "assign-statement-name"
	NodeAssignStatementValue = "assign-statement-expr"

	//
	// NodeTypeField/NodeTypeVariable/NodeTypeVariableStatement
	//
	NodeVariableStatementDeclaredType = "var-declared-type"
	NodeVariableStatementName         = "named"
	NodeVariableStatementExpression   = "var-expr"

	//
	// NodeTypeConditionalStatement
	//
	NodeConditionalStatementConditional = "conditional-expr"
	NodeConditionalStatementBlock       = "conditional-block"
	NodeConditionalStatementElseClause  = "conditional-else"

	//
	// NodeTypeReturnStatement
	//
	NodeReturnStatementValue = "return-expr"

	//
	// NodeTypeBreakStatement
	//
	NodeBreakStatementLabel = "statement-label-destination"

	//
	// NodeTypeContinueStatement
	//
	NodeContinueStatementLabel = "statement-label-destination"

	//
	// NodeTypeWithStatement
	//
	NodeWithStatementExpression = "with-expression"
	NodeWithStatementBlock      = "with-block"

	//
	// NodeTypeMatchStatement
	//
	NodeMatchStatementExpression = "match-expression"
	NodeMatchStatementCase       = "match-case"

	//
	// NodeTypeMatchStatementCase
	//
	NodeMatchStatementCaseExpression = "match-case-expression"
	NodeMatchStatementCaseStatement  = "match-case-statement"

	//
	// NodeTypeAwaitExpression
	//
	NodeAwaitExpressionSource = "await-expression-source"

	//
	// NodeTypeArrowExpression
	//
	NodeArrowExpressionDestination = "arrow-expression-left"
	NodeArrowExpressionSource      = "arrow-expression-right"

	//
	// NodeTypeLambdaExpression
	//
	NodeLambdaExpressionReturnType        = "lambda-expression-return-type"
	NodeLambdaExpressionParameter         = "lambda-expression-parameter"
	NodeLambdaExpressionInferredParameter = "lambda-expression-inferred-parameter"
	NodeLambdaExpressionBlock             = "lambda-expression-block"
	NodeLambdaExpressionChildExpr         = "lambda-expression-child-expr"

	//
	//	NodeTypeLambdaParameter
	//
	NodeLambdaExpressionParameterName = "named"

	//
	// Binary expressions.
	//
	NodeBinaryExpressionLeftExpr  = "binary-expression-left"
	NodeBinaryExpressionRightExpr = "binary-expression-right"

	//
	// Unary expressions.
	//
	NodeUnaryExpressionChildExpr = "unary-expression-child"

	//
	// Member Access expressions.
	//
	NodeMemberAccessChildExpr  = "member-access-expr"
	NodeMemberAccessIdentifier = "member-access-identifier"

	//
	// NodeCastExpression
	//
	NodeCastExpressionType      = "cast-expr-type"
	NodeCastExpressionChildExpr = "cast-expr-expr"

	//
	// NodeSliceExpression
	//
	NodeSliceExpressionChildExpr  = "slice-expr-expr"
	NodeSliceExpressionIndex      = "slice-expr-index"
	NodeSliceExpressionLeftIndex  = "slice-expr-left-index"
	NodeSliceExpressionRightIndex = "slice-expr-right-index"

	//
	// NodeFunctionCallExpression
	//
	NodeFunctionCallArgument            = "function-call-argument"
	NodeFunctionCallExpressionChildExpr = "function-call-expr"

	//
	// NodeListExpression
	//
	NodeListExpressionValue = "list-expr-value"

	//
	// NodeMapExpression
	//
	NodeMapExpressionChildEntry = "map-expr-entry"

	NodeMapExpressionEntryKey   = "map-entry-key"
	NodeMapExpressionEntryValue = "map-entry-value"

	//
	// NodeGenericSpecifierExpression
	//

	NodeGenericSpecifierChildExpr = "generic-specifier-expr"
	NodeGenericSpecifierType      = "generic-specifier-type"

	//
	// Literals.
	//
	NodeNumericLiteralExpressionValue        = "literal-value"
	NodeStringLiteralExpressionValue         = "literal-value"
	NodeBooleanLiteralExpressionValue        = "literal-value"
	NodeTemplateStringLiteralExpressionValue = "literal-value"

	//
	// NodeTypeIdentifierExpression
	//
	NodeIdentifierExpressionName = "identexpr-name"

	//
	// NodeTypeNamedValue
	//
	NodeNamedValueName = "named"
)
