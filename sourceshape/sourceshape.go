// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate stringer -type=NodeType

// Package sourceshape defines the types representing the structure of source code.
package sourceshape

import "strconv"

// SerulianFileExtension defines the file extension for Serulian files.
const SerulianFileExtension = ".seru"

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
	NodeTypeImport        // An import
	NodeTypeImportPackage // A package imported under an import statement.

	NodeTypeClass     // A class
	NodeTypeInterface // An interface
	NodeTypeNominal   // A nominal type
	NodeTypeStruct    // A structural type
	NodeTypeAgent     // An agent type

	NodeTypeGeneric        // A generic definition on a type
	NodeTypeAgentReference // A single agent included in a class

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
	NodeTypeMemberTag     // A tag on a type member.

	// Statements
	NodeTypeArrowStatement       // An arrow statement: a <- b
	NodeTypeStatementBlock       // A block of statements
	NodeTypeLoopStatement        // A for statement
	NodeTypeConditionalStatement // An if statement
	NodeTypeReturnStatement      // A return statement
	NodeTypeYieldStatement       // A yield statement
	NodeTypeRejectStatement      // A reject statement
	NodeTypeBreakStatement       // A break statement
	NodeTypeContinueStatement    // A continue statement
	NodeTypeVariableStatement    // A variable statement
	NodeTypeWithStatement        // A with statement
	NodeTypeSwitchStatement      // A switch statement
	NodeTypeMatchStatement       // A match statement
	NodeTypeAssignStatement      // An assignment statement: a = b
	NodeTypeResolveStatement     // A resolve assignment statement: a := b
	NodeTypeExpressionStatement  // A statement containing a single expression

	NodeTypeSwitchStatementCase // A case of a switch statement.
	NodeTypeMatchStatementCase  // A case of a match statement.

	NodeTypeNamedValue    // A named value added to the scope of the parent statement.
	NodeTypeAssignedValue // A named value assigned to the scope by a parent statement.

	// Expressions
	NodeTypeAwaitExpression // An await expression: <- a

	NodeTypeLambdaExpression // A lambda expression

	NodeTypeSmlExpression // <sometag />
	NodeTypeSmlAttribute  // a="somevalue" or <.a>
	NodeTypeSmlDecorator  // @a="somevalue"
	NodeTypeSmlText       // some text

	NodeTypeConditionalExpression // a if b else c
	NodeTypeLoopExpression        // a for a in c

	NodeBitwiseXorExpression        // a ^ b
	NodeBitwiseOrExpression         // a | b
	NodeBitwiseAndExpression        // a & b
	NodeBitwiseShiftLeftExpression  // a << b
	NodeBitwiseShiftRightExpression // a >> b
	NodeBitwiseNotExpression        // ~a

	NodeBooleanOrExpression  // a || b
	NodeBooleanAndExpression // a && b
	NodeBooleanNotExpression // !a
	NodeKeywordNotExpression // not a

	NodeRootTypeExpression // &a

	NodeComparisonEqualsExpression    // a == b
	NodeComparisonNotEqualsExpression // a != b

	NodeComparisonLTEExpression // a <= b
	NodeComparisonGTEExpression // a >= b
	NodeComparisonLTExpression  // a < b
	NodeComparisonGTExpression  // a > b

	NodeNullComparisonExpression // a ?? b
	NodeIsComparisonExpression   // a is b
	NodeAssertNotNullExpression  // a!

	NodeInCollectionExpression // a in b

	NodeDefineRangeExpression          // a .. b
	NodeDefineExclusiveRangeExpression // a ..< b

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

	NodeTaggedTemplateLiteralString // someexpr`foo`
	NodeTypeTemplateString          // `foo`

	NodeNumericLiteralExpression   // 123
	NodeStringLiteralExpression    // 'hello'
	NodeBooleanLiteralExpression   // true
	NodeThisLiteralExpression      // this
	NodePrincipalLiteralExpression // principal
	NodeNullLiteralExpression      // null
	NodeValLiteralExpression       // val

	NodeListLiteralExpression         // [1, 2, 3]
	NodeSliceLiteralExpression        // []int{1, 2, 3}
	NodeMappingLiteralExpression      // []{string}{a: 1, b: 2}
	NodeMappingLiteralExpressionEntry // a: 1

	NodeStructuralNewExpression      // SomeName{a: 1, b: 2}
	NodeStructuralNewExpressionEntry // a: 1

	NodeMapLiteralExpression      // {a: 1, b: 2}
	NodeMapLiteralExpressionEntry // a: 1

	NodeTypeIdentifierExpression // An identifier expression

	// Lambda expression members
	NodeTypeLambdaParameter

	// Type references
	NodeTypeTypeReference // A type reference
	NodeTypeStream
	NodeTypeSlice
	NodeTypeMapping
	NodeTypeNullable
	NodeTypeVoid
	NodeTypeAny
	NodeTypeStructReference

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
	NodeImportPredicateKind       = "import-kind"
	NodeImportPredicateSource     = "import-source"
	NodeImportPredicatePackageRef = "import-package"

	NodeImportPredicateSubsource   = "import-subsource"
	NodeImportPredicateName        = "named"
	NodeImportPredicatePackageName = "import-package-name"
	NodeImportPredicateLocation    = "import-location-ref"

	//
	// NodeTypeClass + NodeTypeInterface
	//
	NodeTypeDefinitionGeneric   = "type-generic"
	NodeTypeDefinitionMember    = "type-member"
	NodeTypeDefinitionName      = "named"
	NodeTypeDefinitionDecorator = "decorator"

	//
	// NodeTypeClass + NodeTypeAgent
	//
	NodePredicateComposedAgent = "composed-agent"

	//
	// NodeTypeNominal
	//
	NodeNominalPredicateBaseType = "nominal-basetypepath"

	//
	// NodeTypeAgent
	//
	NodeAgentPredicatePrincipalType = "agent-principaltypepath"

	//
	// NodeTypeAgentReference
	//
	NodeAgentReferencePredicateReferenceType = "agent-reference-reftypepath"
	NodeAgentReferencePredicateAlias         = "agent-alias"

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
	NodePredicateTypeMemberTag          = "typemember-tag"

	//
	// NodeTypeField
	//
	NodePredicateTypeFieldDefaultValue = "typemember-field-default-value"

	//
	// NodeTypeMemberTag
	//
	NodePredicateTypeMemberTagName  = "typemembertag-name"
	NodePredicateTypeMemberTagValue = "typemembertag-value"

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
	// NodeTypeResolveStatement
	//
	NodeResolveStatementSource = "resolve-statement-expr"

	NodeAssignedDestination = "assigned-value-destination"
	NodeAssignedRejection   = "assigned-value-rejection"

	//
	// NodeTypeField/NodeTypeVariable/NodeTypeVariableStatement
	//
	NodeVariableStatementDeclaredType = "var-declared-type"
	NodeVariableStatementName         = "named"
	NodeVariableStatementExpression   = "var-expr"
	NodeVariableStatementConstant     = "var-const"

	//
	// NodeTypeSmlExpression
	//
	NodeSmlExpressionTypeOrFunction = "sml-expression-typefunc"
	NodeSmlExpressionAttribute      = "sml-expression-attribute"
	NodeSmlExpressionDecorator      = "sml-expression-decorator"
	NodeSmlExpressionChild          = "sml-expression-child"
	NodeSmlExpressionNestedProperty = "sml-expression-nested-prop"

	//
	// NodeTypeSmlAttribute
	//
	NodeSmlAttributeName   = "sml-attribute-name"
	NodeSmlAttributeValue  = "sml-attribute-value"
	NodeSmlAttributeNested = "sml-attribute-nested"

	//
	// NodeTypeSmlDecorator
	//
	NodeSmlDecoratorPath  = "sml-decorator-path"
	NodeSmlDecoratorValue = "sml-decorator-value"

	//
	// NodeTypeSmlText
	//
	NodeSmlTextValue = "sml-text-value"

	//
	// NodeTypeConditionalStatement
	//
	NodeConditionalStatementConditional = "conditional-expr"
	NodeConditionalStatementBlock       = "conditional-block"
	NodeConditionalStatementElseClause  = "conditional-else"

	//
	// NodeTypeYieldStatement
	//
	NodeYieldStatementValue       = "yield-value"
	NodeYieldStatementStreamValue = "yield-stream-value"
	NodeYieldStatementBreak       = "yield-break"

	//
	// NodeTypeReturnStatement
	//
	NodeReturnStatementValue = "return-expr"

	//
	// NodeTypeRejectStatement
	//
	NodeRejectStatementValue = "reject-expr"

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
	// NodeTypeSwitchStatement
	//
	NodeSwitchStatementExpression = "switch-expression"
	NodeSwitchStatementCase       = "switch-case"

	//
	// NodeTypeSwitchStatementCase
	//
	NodeSwitchStatementCaseExpression = "switch-case-expression"
	NodeSwitchStatementCaseStatement  = "switch-case-statement"

	//
	// NodeTypeMatchStatement
	//
	NodeMatchStatementExpression = "match-expression"
	NodeMatchStatementCase       = "match-case"

	//
	// NodeTypeMatchStatementCase
	//
	NodeMatchStatementCaseTypeReference = "match-case-typeref"
	NodeMatchStatementCaseStatement     = "match-case-statement"

	//
	// NodeTypeExpressionStatement
	//
	NodeExpressionStatementExpression = "expr-statement-expr"

	//
	// NodeTypeArrowStatement
	//
	NodeArrowStatementDestination = "arrow-statement-destination"
	NodeArrowStatementRejection   = "arrow-statement-rejection"
	NodeArrowStatementSource      = "arrow-statement-right"

	//
	// NodeTypeLoopExpression
	//
	NodeLoopExpressionStreamExpression = "loop-expr-stream-expr"
	NodeLoopExpressionNamedValue       = "named-value"
	NodeLoopExpressionMapExpression    = "loop-expr-map-expr"

	//
	// NodeTypeConditionalExpression
	//
	NodeConditionalExpressionCheckExpression = "comparison-expr-check"
	NodeConditionalExpressionThenExpression  = "comparison-expr-then"
	NodeConditionalExpressionElseExpression  = "comparison-expr-else"

	//
	// NodeTypeAwaitExpression
	//
	NodeAwaitExpressionSource = "await-expression-source"

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
	// NodeListLiteralExpression
	//
	NodeListLiteralExpressionValue = "list-expr-value"

	//
	// NodeSliceLiteralExpression
	//
	NodeSliceLiteralExpressionValue = "slice-literal-expr-value"
	NodeSliceLiteralExpressionType  = "slice-literal-expr-type"

	//
	// NodeMappingLiteralExpression
	//
	NodeMappingLiteralExpressionEntryRef = "mapping-literal-expr-entry"
	NodeMappingLiteralExpressionType     = "mapping-literal-expr-type"

	//
	// NodeMappingLiteralExpressionEntry
	//
	NodeMappingLiteralExpressionEntryKey   = "mapping-literal-entry-key"
	NodeMappingLiteralExpressionEntryValue = "mapping-literal-entry-value"

	//
	// NodeStructuralNewExpression
	//
	NodeStructuralNewTypeExpression       = "structural-new-type-expr"
	NodeStructuralNewExpressionChildEntry = "structural-new-entry"

	//
	// NodeStructuralNewExpressionEntry
	//
	NodeStructuralNewEntryKey   = "structural-new-entry-key"
	NodeStructuralNewEntryValue = "structural-new-entry-value"

	//
	// NodeMapLiteralExpression
	//
	NodeMapLiteralExpressionChildEntry = "map-expr-entry"

	NodeMapLiteralExpressionEntryKey   = "map-entry-key"
	NodeMapLiteralExpressionEntryValue = "map-entry-value"

	//
	// NodeGenericSpecifierExpression
	//

	NodeGenericSpecifierChildExpr = "generic-specifier-expr"
	NodeGenericSpecifierType      = "generic-specifier-type"

	//
	// NodeTaggedTemplateLiteralString
	//
	NodeTaggedTemplateCallExpression = "tagged-template-callexpr"
	NodeTaggedTemplateParsed         = "tagged-template-parsed"

	//
	// NodeTypeTemplateString
	//
	NodeTemplateStringPiece = "template-string-piece"

	//
	// Literals.
	//
	NodeNumericLiteralExpressionValue = "literal-value"
	NodeStringLiteralExpressionValue  = "literal-value"
	NodeBooleanLiteralExpressionValue = "literal-value"

	//
	// NodeTypeIdentifierExpression
	//
	NodeIdentifierExpressionName = "identexpr-name"

	//
	// NodeTypeNamedValue/NodeTypeAssignedValue
	//
	NodeNamedValueName = "named"
)
