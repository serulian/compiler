// Code generated by protoc-gen-gogo.
// source: membersig.proto
// DO NOT EDIT!

/*
	Package proto is a generated protocol buffer package.

	It is generated from these files:
		membersig.proto

	It has these top-level messages:
		MemberSig
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MemberSig struct {
	MemberName         *string  `protobuf:"bytes,1,opt,name=MemberName" json:"MemberName,omitempty"`
	MemberKind         *uint64  `protobuf:"varint,2,opt,name=MemberKind" json:"MemberKind,omitempty"`
	IsWritable         *bool    `protobuf:"varint,3,opt,name=IsWritable" json:"IsWritable,omitempty"`
	IsExported         *bool    `protobuf:"varint,4,opt,name=IsExported" json:"IsExported,omitempty"`
	MemberType         *string  `protobuf:"bytes,5,opt,name=MemberType" json:"MemberType,omitempty"`
	GenericConstraints []string `protobuf:"bytes,6,rep,name=GenericConstraints" json:"GenericConstraints,omitempty"`
	XXX_unrecognized   []byte   `json:"-"`
}

func (m *MemberSig) Reset()         { *m = MemberSig{} }
func (m *MemberSig) String() string { return proto1.CompactTextString(m) }
func (*MemberSig) ProtoMessage()    {}

func (m *MemberSig) GetMemberName() string {
	if m != nil && m.MemberName != nil {
		return *m.MemberName
	}
	return ""
}

func (m *MemberSig) GetMemberKind() uint64 {
	if m != nil && m.MemberKind != nil {
		return *m.MemberKind
	}
	return 0
}

func (m *MemberSig) GetIsWritable() bool {
	if m != nil && m.IsWritable != nil {
		return *m.IsWritable
	}
	return false
}

func (m *MemberSig) GetIsExported() bool {
	if m != nil && m.IsExported != nil {
		return *m.IsExported
	}
	return false
}

func (m *MemberSig) GetMemberType() string {
	if m != nil && m.MemberType != nil {
		return *m.MemberType
	}
	return ""
}

func (m *MemberSig) GetGenericConstraints() []string {
	if m != nil {
		return m.GenericConstraints
	}
	return nil
}

func (m *MemberSig) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *MemberSig) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.MemberName != nil {
		data[i] = 0xa
		i++
		i = encodeVarintMembersig(data, i, uint64(len(*m.MemberName)))
		i += copy(data[i:], *m.MemberName)
	}
	if m.MemberKind != nil {
		data[i] = 0x10
		i++
		i = encodeVarintMembersig(data, i, uint64(*m.MemberKind))
	}
	if m.IsWritable != nil {
		data[i] = 0x18
		i++
		if *m.IsWritable {
			data[i] = 1
		} else {
			data[i] = 0
		}
		i++
	}
	if m.IsExported != nil {
		data[i] = 0x20
		i++
		if *m.IsExported {
			data[i] = 1
		} else {
			data[i] = 0
		}
		i++
	}
	if m.MemberType != nil {
		data[i] = 0x2a
		i++
		i = encodeVarintMembersig(data, i, uint64(len(*m.MemberType)))
		i += copy(data[i:], *m.MemberType)
	}
	if len(m.GenericConstraints) > 0 {
		for _, s := range m.GenericConstraints {
			data[i] = 0x32
			i++
			l = len(s)
			for l >= 1<<7 {
				data[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			data[i] = uint8(l)
			i++
			i += copy(data[i:], s)
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeFixed64Membersig(data []byte, offset int, v uint64) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	data[offset+4] = uint8(v >> 32)
	data[offset+5] = uint8(v >> 40)
	data[offset+6] = uint8(v >> 48)
	data[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Membersig(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintMembersig(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}
func (m *MemberSig) Size() (n int) {
	var l int
	_ = l
	if m.MemberName != nil {
		l = len(*m.MemberName)
		n += 1 + l + sovMembersig(uint64(l))
	}
	if m.MemberKind != nil {
		n += 1 + sovMembersig(uint64(*m.MemberKind))
	}
	if m.IsWritable != nil {
		n += 2
	}
	if m.IsExported != nil {
		n += 2
	}
	if m.MemberType != nil {
		l = len(*m.MemberType)
		n += 1 + l + sovMembersig(uint64(l))
	}
	if len(m.GenericConstraints) > 0 {
		for _, s := range m.GenericConstraints {
			l = len(s)
			n += 1 + l + sovMembersig(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovMembersig(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozMembersig(x uint64) (n int) {
	return sovMembersig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MemberSig) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMembersig
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MemberSig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MemberSig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MemberName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMembersig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMembersig
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[iNdEx:postIndex])
			m.MemberName = &s
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MemberKind", wireType)
			}
			var v uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMembersig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				v |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.MemberKind = &v
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsWritable", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMembersig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			b := bool(v != 0)
			m.IsWritable = &b
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsExported", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMembersig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			b := bool(v != 0)
			m.IsExported = &b
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MemberType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMembersig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMembersig
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[iNdEx:postIndex])
			m.MemberType = &s
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenericConstraints", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMembersig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMembersig
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GenericConstraints = append(m.GenericConstraints, string(data[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMembersig(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMembersig
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMembersig(data []byte) (n int, err error) {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMembersig
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMembersig
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if data[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMembersig
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthMembersig
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowMembersig
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := data[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipMembersig(data[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthMembersig = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMembersig   = fmt.Errorf("proto: integer overflow")
)
