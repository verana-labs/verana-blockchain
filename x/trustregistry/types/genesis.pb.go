// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: verana/tr/v1/genesis.proto

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

// Counter defines an entity type and its current counter value
type Counter struct {
	EntityType string `protobuf:"bytes,1,opt,name=entity_type,json=entityType,proto3" json:"entity_type,omitempty"`
	Value      uint64 `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *Counter) Reset()         { *m = Counter{} }
func (m *Counter) String() string { return proto.CompactTextString(m) }
func (*Counter) ProtoMessage()    {}
func (*Counter) Descriptor() ([]byte, []int) {
	return fileDescriptor_adf26977d00c7756, []int{0}
}
func (m *Counter) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Counter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Counter.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Counter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Counter.Merge(m, src)
}
func (m *Counter) XXX_Size() int {
	return m.Size()
}
func (m *Counter) XXX_DiscardUnknown() {
	xxx_messageInfo_Counter.DiscardUnknown(m)
}

var xxx_messageInfo_Counter proto.InternalMessageInfo

func (m *Counter) GetEntityType() string {
	if m != nil {
		return m.EntityType
	}
	return ""
}

func (m *Counter) GetValue() uint64 {
	if m != nil {
		return m.Value
	}
	return 0
}

// GenesisState defines the trustregistry module's genesis state.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	// Collection of all trust registries
	TrustRegistries []TrustRegistry `protobuf:"bytes,2,rep,name=trust_registries,json=trustRegistries,proto3" json:"trust_registries"`
	// Collection of all governance framework versions
	GovernanceFrameworkVersions []GovernanceFrameworkVersion `protobuf:"bytes,3,rep,name=governance_framework_versions,json=governanceFrameworkVersions,proto3" json:"governance_framework_versions"`
	// Collection of all governance framework documents
	GovernanceFrameworkDocuments []GovernanceFrameworkDocument `protobuf:"bytes,4,rep,name=governance_framework_documents,json=governanceFrameworkDocuments,proto3" json:"governance_framework_documents"`
	// List of counters by entity type (tr, gfv, gfd)
	Counters []Counter `protobuf:"bytes,5,rep,name=counters,proto3" json:"counters"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_adf26977d00c7756, []int{1}
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

func (m *GenesisState) GetTrustRegistries() []TrustRegistry {
	if m != nil {
		return m.TrustRegistries
	}
	return nil
}

func (m *GenesisState) GetGovernanceFrameworkVersions() []GovernanceFrameworkVersion {
	if m != nil {
		return m.GovernanceFrameworkVersions
	}
	return nil
}

func (m *GenesisState) GetGovernanceFrameworkDocuments() []GovernanceFrameworkDocument {
	if m != nil {
		return m.GovernanceFrameworkDocuments
	}
	return nil
}

func (m *GenesisState) GetCounters() []Counter {
	if m != nil {
		return m.Counters
	}
	return nil
}

func init() {
	proto.RegisterType((*Counter)(nil), "verana.tr.v1.Counter")
	proto.RegisterType((*GenesisState)(nil), "verana.tr.v1.GenesisState")
}

func init() { proto.RegisterFile("verana/tr/v1/genesis.proto", fileDescriptor_adf26977d00c7756) }

var fileDescriptor_adf26977d00c7756 = []byte{
	// 433 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0x31, 0x6f, 0xd4, 0x30,
	0x14, 0xc7, 0x2f, 0xbd, 0x6b, 0xa1, 0xbe, 0x4a, 0x80, 0x75, 0x48, 0xe1, 0x0a, 0xe9, 0xa9, 0x53,
	0x40, 0x22, 0x56, 0xcb, 0xd0, 0x0d, 0xa1, 0x03, 0xd1, 0x85, 0x01, 0x85, 0x0a, 0x21, 0x96, 0xc8,
	0x09, 0x0f, 0xd7, 0xea, 0xc5, 0x8e, 0x6c, 0x27, 0x90, 0x6f, 0xc1, 0xc7, 0x60, 0x64, 0xe5, 0x1b,
	0x74, 0xec, 0xc8, 0x84, 0xd0, 0xdd, 0xc0, 0xd7, 0x40, 0xb1, 0xdd, 0xaa, 0x91, 0x38, 0x75, 0x89,
	0x5e, 0xde, 0xfb, 0xff, 0x7f, 0xff, 0xc4, 0xcf, 0x68, 0xda, 0x80, 0xa2, 0x82, 0x12, 0xa3, 0x48,
	0x73, 0x40, 0x18, 0x08, 0xd0, 0x5c, 0x27, 0x95, 0x92, 0x46, 0xe2, 0x1d, 0x37, 0x4b, 0x8c, 0x4a,
	0x9a, 0x83, 0xe9, 0x3d, 0x5a, 0x72, 0x21, 0x89, 0x7d, 0x3a, 0xc1, 0x74, 0xc2, 0x24, 0x93, 0xb6,
	0x24, 0x5d, 0xe5, 0xbb, 0x0f, 0x7a, 0xc8, 0x8a, 0x2a, 0x5a, 0x7a, 0xe2, 0x34, 0xec, 0x8d, 0x4c,
	0x5b, 0x81, 0x9f, 0xec, 0xbf, 0x40, 0xb7, 0x5e, 0xca, 0x5a, 0x18, 0x50, 0x78, 0x0f, 0x8d, 0x41,
	0x18, 0x6e, 0xda, 0xac, 0x13, 0x84, 0xc1, 0x2c, 0x88, 0xb7, 0x53, 0xe4, 0x5a, 0x27, 0x6d, 0x05,
	0x78, 0x82, 0x36, 0x1b, 0xba, 0xa8, 0x21, 0xdc, 0x98, 0x05, 0xf1, 0x28, 0x75, 0x2f, 0xfb, 0x3f,
	0x87, 0x68, 0xe7, 0xd8, 0x7d, 0xff, 0x3b, 0x43, 0x0d, 0xe0, 0x23, 0xb4, 0xe5, 0xc2, 0x2d, 0x62,
	0x7c, 0x38, 0x49, 0xae, 0xff, 0x4f, 0xf2, 0xd6, 0xce, 0xe6, 0xdb, 0xe7, 0xbf, 0xf7, 0x06, 0xdf,
	0xff, 0xfe, 0x78, 0x12, 0xa4, 0x5e, 0x8e, 0xdf, 0xa0, 0xbb, 0x46, 0xd5, 0xda, 0x64, 0x0a, 0x18,
	0xd7, 0x46, 0x71, 0xd0, 0xe1, 0xc6, 0x6c, 0x18, 0x8f, 0x0f, 0x77, 0xfb, 0x88, 0x93, 0x4e, 0x95,
	0x3a, 0x51, 0x3b, 0x1f, 0x75, 0xa4, 0xf4, 0x8e, 0xb9, 0xd6, 0xe4, 0xa0, 0xb1, 0x42, 0x8f, 0x98,
	0x6c, 0x40, 0x09, 0x2a, 0x0a, 0xc8, 0x3e, 0x2b, 0x5a, 0xc2, 0x17, 0xa9, 0xce, 0xb2, 0x06, 0x94,
	0xe6, 0x52, 0xe8, 0x70, 0x68, 0xd1, 0x71, 0x1f, 0x7d, 0x7c, 0x65, 0x79, 0x7d, 0xe9, 0x78, 0xef,
	0x0c, 0x3e, 0x67, 0x97, 0xad, 0x55, 0x68, 0x5c, 0xa3, 0xe8, 0xbf, 0x99, 0x9f, 0x64, 0x51, 0x97,
	0x20, 0x8c, 0x0e, 0x47, 0x36, 0xf4, 0xf1, 0x8d, 0xa1, 0xaf, 0xbc, 0xc3, 0xa7, 0x3e, 0x64, 0xeb,
	0x25, 0x1a, 0x1f, 0xa1, 0xdb, 0x85, 0x5b, 0xa2, 0x0e, 0x37, 0x6d, 0xc0, 0xfd, 0x7e, 0x80, 0x5f,
	0xb1, 0x87, 0x5d, 0x89, 0xe7, 0x1f, 0xce, 0x97, 0x51, 0x70, 0xb1, 0x8c, 0x82, 0x3f, 0xcb, 0x28,
	0xf8, 0xb6, 0x8a, 0x06, 0x17, 0xab, 0x68, 0xf0, 0x6b, 0x15, 0x0d, 0x3e, 0x3e, 0x67, 0xdc, 0x9c,
	0xd6, 0x79, 0x52, 0xc8, 0x92, 0x38, 0xd4, 0xd3, 0x05, 0xcd, 0xf5, 0x65, 0x9d, 0x2f, 0x64, 0x71,
	0x56, 0x9c, 0x52, 0x2e, 0xc8, 0x57, 0x62, 0x4f, 0xde, 0xef, 0xac, 0x75, 0xb7, 0x2b, 0xdf, 0xb2,
	0xd7, 0xeb, 0xd9, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x59, 0x29, 0x3a, 0xbe, 0xe8, 0x02, 0x00,
	0x00,
}

func (m *Counter) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Counter) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Counter) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Value))
		i--
		dAtA[i] = 0x10
	}
	if len(m.EntityType) > 0 {
		i -= len(m.EntityType)
		copy(dAtA[i:], m.EntityType)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.EntityType)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
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
	if len(m.Counters) > 0 {
		for iNdEx := len(m.Counters) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Counters[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.GovernanceFrameworkDocuments) > 0 {
		for iNdEx := len(m.GovernanceFrameworkDocuments) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.GovernanceFrameworkDocuments[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.GovernanceFrameworkVersions) > 0 {
		for iNdEx := len(m.GovernanceFrameworkVersions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.GovernanceFrameworkVersions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.TrustRegistries) > 0 {
		for iNdEx := len(m.TrustRegistries) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TrustRegistries[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
func (m *Counter) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.EntityType)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Value != 0 {
		n += 1 + sovGenesis(uint64(m.Value))
	}
	return n
}

func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.TrustRegistries) > 0 {
		for _, e := range m.TrustRegistries {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.GovernanceFrameworkVersions) > 0 {
		for _, e := range m.GovernanceFrameworkVersions {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.GovernanceFrameworkDocuments) > 0 {
		for _, e := range m.GovernanceFrameworkDocuments {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Counters) > 0 {
		for _, e := range m.Counters {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Counter) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Counter: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Counter: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EntityType", wireType)
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
			m.EntityType = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			m.Value = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Value |= uint64(b&0x7F) << shift
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
				return fmt.Errorf("proto: wrong wireType = %d for field TrustRegistries", wireType)
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
			m.TrustRegistries = append(m.TrustRegistries, TrustRegistry{})
			if err := m.TrustRegistries[len(m.TrustRegistries)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GovernanceFrameworkVersions", wireType)
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
			m.GovernanceFrameworkVersions = append(m.GovernanceFrameworkVersions, GovernanceFrameworkVersion{})
			if err := m.GovernanceFrameworkVersions[len(m.GovernanceFrameworkVersions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GovernanceFrameworkDocuments", wireType)
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
			m.GovernanceFrameworkDocuments = append(m.GovernanceFrameworkDocuments, GovernanceFrameworkDocument{})
			if err := m.GovernanceFrameworkDocuments[len(m.GovernanceFrameworkDocuments)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Counters", wireType)
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
			m.Counters = append(m.Counters, Counter{})
			if err := m.Counters[len(m.Counters)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
