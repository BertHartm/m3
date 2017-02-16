// Copyright (c) 2017 Uber Technologies, Inc.
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

// Code generated by protoc-gen-go.
// source: metadata.proto
// DO NOT EDIT!

/*
Package metadata is a generated protocol buffer package.

It is generated from these files:
	metadata.proto

It has these top-level messages:
	Metadata
*/
package metadata

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Metadata struct {
	Port              uint32 `protobuf:"varint,1,opt,name=port" json:"port,omitempty"`
	LivenessInterval  int64  `protobuf:"varint,2,opt,name=liveness_interval,json=livenessInterval" json:"liveness_interval,omitempty"`
	HeartbeatInterval int64  `protobuf:"varint,3,opt,name=heartbeat_interval,json=heartbeatInterval" json:"heartbeat_interval,omitempty"`
}

func (m *Metadata) Reset()                    { *m = Metadata{} }
func (m *Metadata) String() string            { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()               {}
func (*Metadata) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Metadata)(nil), "metadata.metadata")
}

func init() { proto.RegisterFile("metadata.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 127 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcb, 0x4d, 0x2d, 0x49,
	0x4c, 0x49, 0x2c, 0x49, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x80, 0xf1, 0x95, 0xaa,
	0xb8, 0xe0, 0x6c, 0x21, 0x21, 0x2e, 0x96, 0x82, 0xfc, 0xa2, 0x12, 0x09, 0x46, 0x05, 0x46, 0x0d,
	0xde, 0x20, 0x30, 0x5b, 0x48, 0x9b, 0x4b, 0x30, 0x27, 0xb3, 0x2c, 0x35, 0x2f, 0xb5, 0xb8, 0x38,
	0x3e, 0x33, 0xaf, 0x24, 0xb5, 0xa8, 0x2c, 0x31, 0x47, 0x82, 0x49, 0x81, 0x51, 0x83, 0x39, 0x48,
	0x00, 0x26, 0xe1, 0x09, 0x15, 0x17, 0xd2, 0xe5, 0x12, 0xca, 0x48, 0x4d, 0x2c, 0x2a, 0x49, 0x4a,
	0x4d, 0x2c, 0x41, 0xa8, 0x66, 0x06, 0xab, 0x16, 0x84, 0xcb, 0xc0, 0x94, 0x27, 0xb1, 0x81, 0x1d,
	0x63, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x0f, 0xe8, 0x23, 0x88, 0x9e, 0x00, 0x00, 0x00,
}
