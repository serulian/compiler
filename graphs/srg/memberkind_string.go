// Code generated by "stringer -type=MemberKind"; DO NOT EDIT.

package srg

import "fmt"

const _MemberKind_name = "ConstructorMemberVarMemberFunctionMemberPropertyMemberOperatorMember"

var _MemberKind_index = [...]uint8{0, 17, 26, 40, 54, 68}

func (i MemberKind) String() string {
	if i < 0 || i >= MemberKind(len(_MemberKind_index)-1) {
		return fmt.Sprintf("MemberKind(%d)", i)
	}
	return _MemberKind_name[_MemberKind_index[i]:_MemberKind_index[i+1]]
}
