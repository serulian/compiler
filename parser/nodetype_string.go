// generated by stringer -type=NodeType; DO NOT EDIT

package parser

import "fmt"

const _NodeType_name = "NodeTypeErrorNodeTypeFileNodeTypeCommentNodeTypeDecoratorNodeTypeImportNodeTypeImportPackageNodeTypeClassNodeTypeInterfaceNodeTypeNominalNodeTypeStructNodeTypeGenericNodeTypeFunctionNodeTypeVariableNodeTypeConstructorNodeTypePropertyNodeTypeOperatorNodeTypeFieldNodeTypePropertyBlockNodeTypeParameterNodeTypeMemberTagNodeTypeArrowStatementNodeTypeStatementBlockNodeTypeLoopStatementNodeTypeConditionalStatementNodeTypeReturnStatementNodeTypeYieldStatementNodeTypeRejectStatementNodeTypeBreakStatementNodeTypeContinueStatementNodeTypeVariableStatementNodeTypeWithStatementNodeTypeMatchStatementNodeTypeAssignStatementNodeTypeExpressionStatementNodeTypeMatchStatementCaseNodeTypeNamedValueNodeTypeAwaitExpressionNodeTypeLambdaExpressionNodeTypeConditionalExpressionNodeTypeLoopExpressionNodeBitwiseXorExpressionNodeBitwiseOrExpressionNodeBitwiseAndExpressionNodeBitwiseShiftLeftExpressionNodeBitwiseShiftRightExpressionNodeBitwiseNotExpressionNodeBooleanOrExpressionNodeBooleanAndExpressionNodeBooleanNotExpressionNodeRootTypeExpressionNodeComparisonEqualsExpressionNodeComparisonNotEqualsExpressionNodeComparisonLTEExpressionNodeComparisonGTEExpressionNodeComparisonLTExpressionNodeComparisonGTExpressionNodeNullComparisonExpressionNodeIsComparisonExpressionNodeAssertNotNullExpressionNodeInCollectionExpressionNodeDefineRangeExpressionNodeBinaryAddExpressionNodeBinarySubtractExpressionNodeBinaryMultiplyExpressionNodeBinaryDivideExpressionNodeBinaryModuloExpressionNodeMemberAccessExpressionNodeNullableMemberAccessExpressionNodeDynamicMemberAccessExpressionNodeStreamMemberAccessExpressionNodeCastExpressionNodeFunctionCallExpressionNodeSliceExpressionNodeGenericSpecifierExpressionNodeTaggedTemplateLiteralStringNodeTypeTemplateStringNodeNumericLiteralExpressionNodeStringLiteralExpressionNodeBooleanLiteralExpressionNodeThisLiteralExpressionNodeNullLiteralExpressionNodeValLiteralExpressionNodeListExpressionNodeSliceLiteralExpressionNodeMappingLiteralExpressionNodeMappingLiteralExpressionEntryNodeStructuralNewExpressionNodeStructuralNewExpressionEntryNodeMapExpressionNodeMapExpressionEntryNodeTypeIdentifierExpressionNodeTypeLambdaParameterNodeTypeTypeReferenceNodeTypeStreamNodeTypeSliceNodeTypeMappingNodeTypeNullableNodeTypeVoidNodeTypeAnyNodeTypeIdentifierPathNodeTypeIdentifierAccessNodeTypeTagged"

var _NodeType_index = [...]uint16{0, 13, 25, 40, 57, 71, 92, 105, 122, 137, 151, 166, 182, 198, 217, 233, 249, 262, 283, 300, 317, 339, 361, 382, 410, 433, 455, 478, 500, 525, 550, 571, 593, 616, 643, 669, 687, 710, 734, 763, 785, 809, 832, 856, 886, 917, 941, 964, 988, 1012, 1034, 1064, 1097, 1124, 1151, 1177, 1203, 1231, 1257, 1284, 1310, 1335, 1358, 1386, 1414, 1440, 1466, 1492, 1526, 1559, 1591, 1609, 1635, 1654, 1684, 1715, 1737, 1765, 1792, 1820, 1845, 1870, 1894, 1912, 1938, 1966, 1999, 2026, 2058, 2075, 2097, 2125, 2148, 2169, 2183, 2196, 2211, 2227, 2239, 2250, 2272, 2296, 2310}

func (i NodeType) String() string {
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return fmt.Sprintf("NodeType(%d)", i)
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}
