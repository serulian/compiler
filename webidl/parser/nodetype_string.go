// generated by stringer -type=NodeType; DO NOT EDIT

package parser

import "fmt"

const _NodeType_name = "NodeTypeErrorNodeTypeFileNodeTypeCommentNodeTypeAnnotationNodeTypeParameterNodeTypeDeclarationNodeTypeMemberNodeTypeImplementationNodeTypeTagged"

var _NodeType_index = [...]uint8{0, 13, 25, 40, 58, 75, 94, 108, 130, 144}

func (i NodeType) String() string {
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return fmt.Sprintf("NodeType(%d)", i)
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}
