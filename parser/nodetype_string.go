// Code generated by "stringer -type=NodeType"; DO NOT EDIT

package parser

import "fmt"

const _NodeType_name = "NodeTypeErrorNodeTypeFileNodeTypeCommentNodeTypeDecoratorNodeTypeImportNodeTypeImportPackageNodeTypeClassNodeTypeInterfaceNodeTypeNominalNodeTypeStructNodeTypeGenericNodeTypeFunctionNodeTypeVariableNodeTypeConstructorNodeTypePropertyNodeTypeOperatorNodeTypeFieldNodeTypePropertyBlockNodeTypeParameterNodeTypeMemberTagNodeTypeArrowStatementNodeTypeStatementBlockNodeTypeLoopStatementNodeTypeConditionalStatementNodeTypeReturnStatementNodeTypeYieldStatementNodeTypeRejectStatementNodeTypeBreakStatementNodeTypeContinueStatementNodeTypeVariableStatementNodeTypeWithStatementNodeTypeSwitchStatementNodeTypeMatchStatementNodeTypeAssignStatementNodeTypeResolveStatementNodeTypeExpressionStatementNodeTypeSwitchStatementCaseNodeTypeMatchStatementCaseNodeTypeNamedValueNodeTypeAssignedValueNodeTypeAwaitExpressionNodeTypeLambdaExpressionNodeTypeSmlExpressionNodeTypeSmlAttributeNodeTypeSmlDecoratorNodeTypeSmlTextNodeTypeConditionalExpressionNodeTypeLoopExpressionNodeBitwiseXorExpressionNodeBitwiseOrExpressionNodeBitwiseAndExpressionNodeBitwiseShiftLeftExpressionNodeBitwiseShiftRightExpressionNodeBitwiseNotExpressionNodeBooleanOrExpressionNodeBooleanAndExpressionNodeBooleanNotExpressionNodeKeywordNotExpressionNodeRootTypeExpressionNodeComparisonEqualsExpressionNodeComparisonNotEqualsExpressionNodeComparisonLTEExpressionNodeComparisonGTEExpressionNodeComparisonLTExpressionNodeComparisonGTExpressionNodeNullComparisonExpressionNodeIsComparisonExpressionNodeAssertNotNullExpressionNodeInCollectionExpressionNodeDefineRangeExpressionNodeBinaryAddExpressionNodeBinarySubtractExpressionNodeBinaryMultiplyExpressionNodeBinaryDivideExpressionNodeBinaryModuloExpressionNodeMemberAccessExpressionNodeNullableMemberAccessExpressionNodeDynamicMemberAccessExpressionNodeStreamMemberAccessExpressionNodeCastExpressionNodeFunctionCallExpressionNodeSliceExpressionNodeGenericSpecifierExpressionNodeTaggedTemplateLiteralStringNodeTypeTemplateStringNodeNumericLiteralExpressionNodeStringLiteralExpressionNodeBooleanLiteralExpressionNodeThisLiteralExpressionNodeNullLiteralExpressionNodeValLiteralExpressionNodeListExpressionNodeSliceLiteralExpressionNodeMappingLiteralExpressionNodeMappingLiteralExpressionEntryNodeStructuralNewExpressionNodeStructuralNewExpressionEntryNodeMapExpressionNodeMapExpressionEntryNodeTypeIdentifierExpressionNodeTypeLambdaParameterNodeTypeTypeReferenceNodeTypeStreamNodeTypeSliceNodeTypeMappingNodeTypeNullableNodeTypeVoidNodeTypeAnyNodeTypeStructReferenceNodeTypeIdentifierPathNodeTypeIdentifierAccessNodeTypeTagged"

var _NodeType_index = [...]uint16{0, 13, 25, 40, 57, 71, 92, 105, 122, 137, 151, 166, 182, 198, 217, 233, 249, 262, 283, 300, 317, 339, 361, 382, 410, 433, 455, 478, 500, 525, 550, 571, 594, 616, 639, 663, 690, 717, 743, 761, 782, 805, 829, 850, 870, 890, 905, 934, 956, 980, 1003, 1027, 1057, 1088, 1112, 1135, 1159, 1183, 1207, 1229, 1259, 1292, 1319, 1346, 1372, 1398, 1426, 1452, 1479, 1505, 1530, 1553, 1581, 1609, 1635, 1661, 1687, 1721, 1754, 1786, 1804, 1830, 1849, 1879, 1910, 1932, 1960, 1987, 2015, 2040, 2065, 2089, 2107, 2133, 2161, 2194, 2221, 2253, 2270, 2292, 2320, 2343, 2364, 2378, 2391, 2406, 2422, 2434, 2445, 2468, 2490, 2514, 2528}

func (i NodeType) String() string {
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return fmt.Sprintf("NodeType(%d)", i)
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}
