// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/m3db/m3/src/query/generated/proto/admin/placement.proto

// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package admin

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import placementpb "github.com/m3db/m3cluster/generated/proto/placementpb"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type PlacementInitRequest struct {
	Instances         []*placementpb.Instance `protobuf:"bytes,1,rep,name=instances" json:"instances,omitempty"`
	NumShards         int32                   `protobuf:"varint,2,opt,name=num_shards,json=numShards,proto3" json:"num_shards,omitempty"`
	ReplicationFactor int32                   `protobuf:"varint,3,opt,name=replication_factor,json=replicationFactor,proto3" json:"replication_factor,omitempty"`
}

func (m *PlacementInitRequest) Reset()                    { *m = PlacementInitRequest{} }
func (m *PlacementInitRequest) String() string            { return proto.CompactTextString(m) }
func (*PlacementInitRequest) ProtoMessage()               {}
func (*PlacementInitRequest) Descriptor() ([]byte, []int) { return fileDescriptorPlacement, []int{0} }

func (m *PlacementInitRequest) GetInstances() []*placementpb.Instance {
	if m != nil {
		return m.Instances
	}
	return nil
}

func (m *PlacementInitRequest) GetNumShards() int32 {
	if m != nil {
		return m.NumShards
	}
	return 0
}

func (m *PlacementInitRequest) GetReplicationFactor() int32 {
	if m != nil {
		return m.ReplicationFactor
	}
	return 0
}

type PlacementGetResponse struct {
	Placement *placementpb.Placement `protobuf:"bytes,1,opt,name=placement" json:"placement,omitempty"`
	Version   int32                  `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (m *PlacementGetResponse) Reset()                    { *m = PlacementGetResponse{} }
func (m *PlacementGetResponse) String() string            { return proto.CompactTextString(m) }
func (*PlacementGetResponse) ProtoMessage()               {}
func (*PlacementGetResponse) Descriptor() ([]byte, []int) { return fileDescriptorPlacement, []int{1} }

func (m *PlacementGetResponse) GetPlacement() *placementpb.Placement {
	if m != nil {
		return m.Placement
	}
	return nil
}

func (m *PlacementGetResponse) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

type PlacementAddRequest struct {
	Instances []*placementpb.Instance `protobuf:"bytes,1,rep,name=instances" json:"instances,omitempty"`
	// By default add requests will only succeed if all instances in the placement
	// are AVAILABLE for all their shards. force overrides that.
	Force bool `protobuf:"varint,2,opt,name=force,proto3" json:"force,omitempty"`
}

func (m *PlacementAddRequest) Reset()                    { *m = PlacementAddRequest{} }
func (m *PlacementAddRequest) String() string            { return proto.CompactTextString(m) }
func (*PlacementAddRequest) ProtoMessage()               {}
func (*PlacementAddRequest) Descriptor() ([]byte, []int) { return fileDescriptorPlacement, []int{2} }

func (m *PlacementAddRequest) GetInstances() []*placementpb.Instance {
	if m != nil {
		return m.Instances
	}
	return nil
}

func (m *PlacementAddRequest) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
}

func init() {
	proto.RegisterType((*PlacementInitRequest)(nil), "admin.PlacementInitRequest")
	proto.RegisterType((*PlacementGetResponse)(nil), "admin.PlacementGetResponse")
	proto.RegisterType((*PlacementAddRequest)(nil), "admin.PlacementAddRequest")
}
func (m *PlacementInitRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PlacementInitRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Instances) > 0 {
		for _, msg := range m.Instances {
			dAtA[i] = 0xa
			i++
			i = encodeVarintPlacement(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.NumShards != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintPlacement(dAtA, i, uint64(m.NumShards))
	}
	if m.ReplicationFactor != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintPlacement(dAtA, i, uint64(m.ReplicationFactor))
	}
	return i, nil
}

func (m *PlacementGetResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PlacementGetResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Placement != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintPlacement(dAtA, i, uint64(m.Placement.Size()))
		n1, err := m.Placement.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.Version != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintPlacement(dAtA, i, uint64(m.Version))
	}
	return i, nil
}

func (m *PlacementAddRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PlacementAddRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Instances) > 0 {
		for _, msg := range m.Instances {
			dAtA[i] = 0xa
			i++
			i = encodeVarintPlacement(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Force {
		dAtA[i] = 0x10
		i++
		if m.Force {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func encodeVarintPlacement(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *PlacementInitRequest) Size() (n int) {
	var l int
	_ = l
	if len(m.Instances) > 0 {
		for _, e := range m.Instances {
			l = e.Size()
			n += 1 + l + sovPlacement(uint64(l))
		}
	}
	if m.NumShards != 0 {
		n += 1 + sovPlacement(uint64(m.NumShards))
	}
	if m.ReplicationFactor != 0 {
		n += 1 + sovPlacement(uint64(m.ReplicationFactor))
	}
	return n
}

func (m *PlacementGetResponse) Size() (n int) {
	var l int
	_ = l
	if m.Placement != nil {
		l = m.Placement.Size()
		n += 1 + l + sovPlacement(uint64(l))
	}
	if m.Version != 0 {
		n += 1 + sovPlacement(uint64(m.Version))
	}
	return n
}

func (m *PlacementAddRequest) Size() (n int) {
	var l int
	_ = l
	if len(m.Instances) > 0 {
		for _, e := range m.Instances {
			l = e.Size()
			n += 1 + l + sovPlacement(uint64(l))
		}
	}
	if m.Force {
		n += 2
	}
	return n
}

func sovPlacement(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozPlacement(x uint64) (n int) {
	return sovPlacement(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PlacementInitRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPlacement
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PlacementInitRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PlacementInitRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Instances", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlacement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPlacement
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Instances = append(m.Instances, &placementpb.Instance{})
			if err := m.Instances[len(m.Instances)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumShards", wireType)
			}
			m.NumShards = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlacement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumShards |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReplicationFactor", wireType)
			}
			m.ReplicationFactor = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlacement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReplicationFactor |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPlacement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPlacement
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
func (m *PlacementGetResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPlacement
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PlacementGetResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PlacementGetResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Placement", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlacement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPlacement
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Placement == nil {
				m.Placement = &placementpb.Placement{}
			}
			if err := m.Placement.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			m.Version = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlacement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Version |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPlacement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPlacement
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
func (m *PlacementAddRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPlacement
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PlacementAddRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PlacementAddRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Instances", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlacement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPlacement
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Instances = append(m.Instances, &placementpb.Instance{})
			if err := m.Instances[len(m.Instances)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Force", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlacement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Force = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipPlacement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPlacement
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
func skipPlacement(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPlacement
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
					return 0, ErrIntOverflowPlacement
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
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
					return 0, ErrIntOverflowPlacement
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthPlacement
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowPlacement
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
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
				next, err := skipPlacement(dAtA[start:])
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
	ErrInvalidLengthPlacement = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPlacement   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/m3db/m3/src/query/generated/proto/admin/placement.proto", fileDescriptorPlacement)
}

var fileDescriptorPlacement = []byte{
	// 310 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x90, 0x4f, 0x4e, 0x02, 0x31,
	0x14, 0xc6, 0xad, 0x04, 0x95, 0xb2, 0xd1, 0x8a, 0x66, 0x62, 0xe2, 0x84, 0xb0, 0x62, 0xe3, 0x34,
	0x61, 0xbc, 0x80, 0x24, 0x6a, 0xd8, 0x99, 0xf1, 0x00, 0xd8, 0x69, 0x1f, 0xd0, 0x84, 0xb6, 0x43,
	0xff, 0x98, 0x78, 0x0b, 0xb7, 0xde, 0xc8, 0xa5, 0x47, 0x30, 0x78, 0x11, 0x93, 0x0a, 0x03, 0xca,
	0xd2, 0xe5, 0xfb, 0xbe, 0xef, 0xfd, 0xde, 0x97, 0x87, 0x87, 0x53, 0xe9, 0x67, 0xa1, 0xcc, 0xb8,
	0x51, 0x54, 0xe5, 0xa2, 0xa4, 0x2a, 0xa7, 0xce, 0x72, 0xba, 0x08, 0x60, 0x5f, 0xe8, 0x14, 0x34,
	0x58, 0xe6, 0x41, 0xd0, 0xca, 0x1a, 0x6f, 0x28, 0x13, 0x4a, 0x6a, 0x5a, 0xcd, 0x19, 0x07, 0x05,
	0xda, 0x67, 0x51, 0x25, 0xcd, 0x28, 0x5f, 0xdc, 0xee, 0xa2, 0xf8, 0x3c, 0x38, 0x0f, 0x76, 0x87,
	0x53, 0x13, 0xaa, 0xf2, 0x2f, 0xad, 0xf7, 0x86, 0x70, 0xe7, 0x61, 0xad, 0x8d, 0xb4, 0xf4, 0x05,
	0x2c, 0x02, 0x38, 0x4f, 0x72, 0xdc, 0x92, 0xda, 0x79, 0xa6, 0x39, 0xb8, 0x04, 0x75, 0x1b, 0xfd,
	0xf6, 0xe0, 0x2c, 0xdb, 0x22, 0x65, 0xa3, 0x95, 0x5b, 0x6c, 0x72, 0xe4, 0x12, 0x63, 0x1d, 0xd4,
	0xd8, 0xcd, 0x98, 0x15, 0x2e, 0xd9, 0xef, 0xa2, 0x7e, 0xb3, 0x68, 0xe9, 0xa0, 0x1e, 0xa3, 0x40,
	0xae, 0x30, 0xb1, 0x50, 0xcd, 0x25, 0x67, 0x5e, 0x1a, 0x3d, 0x9e, 0x30, 0xee, 0x8d, 0x4d, 0x1a,
	0x31, 0x76, 0xb2, 0xe5, 0xdc, 0x45, 0xa3, 0x37, 0xd9, 0xaa, 0x76, 0x0f, 0xbe, 0x00, 0x57, 0x19,
	0xed, 0x80, 0x5c, 0xe3, 0x56, 0x5d, 0x24, 0x41, 0x5d, 0xd4, 0x6f, 0x0f, 0xce, 0x7f, 0x55, 0xab,
	0xb7, 0x8a, 0x4d, 0x90, 0x24, 0xf8, 0xf0, 0x19, 0xac, 0x93, 0x46, 0xaf, 0x8a, 0xad, 0xc7, 0xde,
	0x13, 0x3e, 0xad, 0x37, 0x6e, 0x84, 0xf8, 0xd7, 0x07, 0x3a, 0xb8, 0x39, 0x31, 0x96, 0x43, 0xbc,
	0x71, 0x54, 0xfc, 0x0c, 0xc3, 0xe3, 0xf7, 0x65, 0x8a, 0x3e, 0x96, 0x29, 0xfa, 0x5c, 0xa6, 0xe8,
	0xf5, 0x2b, 0xdd, 0x2b, 0x0f, 0xe2, 0xfb, 0xf3, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xbe, 0x40,
	0x2b, 0x71, 0x12, 0x02, 0x00, 0x00,
}
