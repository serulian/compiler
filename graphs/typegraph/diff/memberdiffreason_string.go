// Code generated by "stringer -type=MemberDiffReason"; DO NOT EDIT.

package diff

import "fmt"

const _MemberDiffReason_name = "MemberDiffReasonNotApplicable"

var _MemberDiffReason_index = [...]uint8{0, 29}

func (i MemberDiffReason) String() string {
	if i < 0 || i >= MemberDiffReason(len(_MemberDiffReason_index)-1) {
		return fmt.Sprintf("MemberDiffReason(%d)", i)
	}
	return _MemberDiffReason_name[_MemberDiffReason_index[i]:_MemberDiffReason_index[i+1]]
}
