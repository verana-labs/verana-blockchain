// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: verana/td/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the trustdeposit module's genesis state.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params        Params               `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	TrustDeposits []TrustDepositRecord `protobuf:"bytes,2,rep,name=trust_deposits,json=trustDeposits,proto3" json:"trust_deposits"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_daff73aaff90939d, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetTrustDeposits() []TrustDepositRecord {
	if m != nil {
		return m.TrustDeposits
	}
	return nil
}

// TrustDepositRecord defines a trust deposit entry for genesis state
type TrustDepositRecord struct {
	Account   string `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	Share     uint64 `protobuf:"varint,2,opt,name=share,proto3" json:"share,omitempty"`
	Amount    uint64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Claimable uint64 `protobuf:"varint,4,opt,name=claimable,proto3" json:"claimable,omitempty"`
}

func (m *TrustDepositRecord) Reset()         { *m = TrustDepositRecord{} }
func (m *TrustDepositRecord) String() string { return proto.CompactTextString(m) }
func (*TrustDepositRecord) ProtoMessage()    {}
func (*TrustDepositRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_daff73aaff90939d, []int{1}
}
func (m *TrustDepositRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TrustDepositRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TrustDepositRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TrustDepositRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TrustDepositRecord.Merge(m, src)
}
func (m *TrustDepositRecord) XXX_Size() int {
	return m.Size()
}
func (m *TrustDepositRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_TrustDepositRecord.DiscardUnknown(m)
}

var xxx_messageInfo_TrustDepositRecord proto.InternalMessageInfo

func (m *TrustDepositRecord) GetAccount() string {
	if m != nil {
		return m.Account
	}
	return ""
}

func (m *TrustDepositRecord) GetShare() uint64 {
	if m != nil {
		return m.Share
	}
	return 0
}

func (m *TrustDepositRecord) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *TrustDepositRecord) GetClaimable() uint64 {
	if m != nil {
		return m.Claimable
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "verana.td.v1.GenesisState")
	proto.RegisterType((*TrustDepositRecord)(nil), "verana.td.v1.TrustDepositRecord")
}

func init() { proto.RegisterFile("verana/td/v1/genesis.proto", fileDescriptor_daff73aaff90939d) }

var fileDescriptor_daff73aaff90939d = []byte{
	// 342 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0xb3, 0x6d, 0xad, 0x74, 0x5b, 0x05, 0x97, 0x22, 0xb1, 0x48, 0x0c, 0x3d, 0x15, 0xc1,
	0x2c, 0xad, 0x07, 0x4f, 0x5e, 0x8a, 0xe0, 0x49, 0x90, 0x28, 0x08, 0x5e, 0x64, 0xb2, 0x59, 0xd2,
	0x60, 0x92, 0x0d, 0xd9, 0x6d, 0x51, 0x9f, 0xc2, 0x93, 0xcf, 0xe0, 0xd1, 0xc7, 0xe8, 0xb1, 0x47,
	0x4f, 0x22, 0xed, 0xc1, 0xd7, 0x90, 0xee, 0xa6, 0xd8, 0xe2, 0x65, 0x99, 0xf9, 0xff, 0x8f, 0x7f,
	0x76, 0x06, 0x77, 0x26, 0xbc, 0x80, 0x0c, 0xa8, 0x0a, 0xe9, 0xa4, 0x4f, 0x23, 0x9e, 0x71, 0x19,
	0x4b, 0x2f, 0x2f, 0x84, 0x12, 0xa4, 0x65, 0x3c, 0x4f, 0x85, 0xde, 0xa4, 0xdf, 0xd9, 0x83, 0x34,
	0xce, 0x04, 0xd5, 0xaf, 0x01, 0x3a, 0xed, 0x48, 0x44, 0x42, 0x97, 0x74, 0x59, 0x95, 0xea, 0xc1,
	0x46, 0x64, 0x0e, 0x05, 0xa4, 0x65, 0x62, 0xf7, 0x0d, 0xe1, 0xd6, 0xa5, 0x99, 0x71, 0xa3, 0x40,
	0x71, 0x72, 0x86, 0xeb, 0x06, 0xb0, 0x91, 0x8b, 0x7a, 0xcd, 0x41, 0xdb, 0x5b, 0x9f, 0xe9, 0x5d,
	0x6b, 0x6f, 0xd8, 0x98, 0x7e, 0x1d, 0x59, 0xef, 0x3f, 0x1f, 0xc7, 0xc8, 0x2f, 0x71, 0x72, 0x85,
	0x77, 0x55, 0x31, 0x96, 0xea, 0x21, 0xe4, 0xb9, 0x90, 0xb1, 0x92, 0x76, 0xc5, 0xad, 0xf6, 0x9a,
	0x03, 0x77, 0x33, 0xe0, 0x76, 0xc9, 0x5c, 0x18, 0xc4, 0xe7, 0x4c, 0x14, 0xe1, 0xb0, 0xb6, 0x0c,
	0xf3, 0x77, 0xd4, 0x9a, 0x23, 0xbb, 0x2f, 0x98, 0xfc, 0x47, 0x89, 0x8d, 0xb7, 0x81, 0x31, 0x31,
	0xce, 0x94, 0xfe, 0x5e, 0xc3, 0x5f, 0xb5, 0xa4, 0x8d, 0xb7, 0xe4, 0x08, 0x0a, 0x6e, 0x57, 0x5c,
	0xd4, 0xab, 0xf9, 0xa6, 0x21, 0xfb, 0xb8, 0x0e, 0xa9, 0xc6, 0xab, 0x5a, 0x2e, 0x3b, 0x72, 0x88,
	0x1b, 0x2c, 0x81, 0x38, 0x85, 0x20, 0xe1, 0x76, 0x4d, 0x5b, 0x7f, 0xc2, 0xf0, 0x6e, 0x3a, 0x77,
	0xd0, 0x6c, 0xee, 0xa0, 0xef, 0xb9, 0x83, 0x5e, 0x17, 0x8e, 0x35, 0x5b, 0x38, 0xd6, 0xe7, 0xc2,
	0xb1, 0xee, 0xcf, 0xa3, 0x58, 0x8d, 0xc6, 0x81, 0xc7, 0x44, 0x4a, 0xcd, 0x5a, 0x27, 0x09, 0x04,
	0x72, 0x55, 0x07, 0x89, 0x60, 0x8f, 0x6c, 0x04, 0x71, 0x46, 0x9f, 0xa8, 0xde, 0xa7, 0x3c, 0x06,
	0x55, 0xcf, 0x39, 0x97, 0x41, 0x5d, 0x1f, 0xfd, 0xf4, 0x37, 0x00, 0x00, 0xff, 0xff, 0x85, 0xeb,
	0x3f, 0xf3, 0xe4, 0x01, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TrustDeposits) > 0 {
		for iNdEx := len(m.TrustDeposits) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TrustDeposits[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *TrustDepositRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TrustDepositRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TrustDepositRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Claimable != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Claimable))
		i--
		dAtA[i] = 0x20
	}
	if m.Amount != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x18
	}
	if m.Share != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Share))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Account) > 0 {
		i -= len(m.Account)
		copy(dAtA[i:], m.Account)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Account)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.TrustDeposits) > 0 {
		for _, e := range m.TrustDeposits {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *TrustDepositRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Account)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Share != 0 {
		n += 1 + sovGenesis(uint64(m.Share))
	}
	if m.Amount != 0 {
		n += 1 + sovGenesis(uint64(m.Amount))
	}
	if m.Claimable != 0 {
		n += 1 + sovGenesis(uint64(m.Claimable))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TrustDeposits", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TrustDeposits = append(m.TrustDeposits, TrustDepositRecord{})
			if err := m.TrustDeposits[len(m.TrustDeposits)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TrustDepositRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TrustDepositRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TrustDepositRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Account = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Share", wireType)
			}
			m.Share = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Share |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Claimable", wireType)
			}
			m.Claimable = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Claimable |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
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
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
