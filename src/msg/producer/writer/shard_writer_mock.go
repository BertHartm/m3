// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m3db/m3/src/msg/producer/writer/shard_go

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

// Package writer is a generated GoMock package.
package writer

import (
	"reflect"

	"github.com/m3db/m3/src/cluster/placement"
	"github.com/m3db/m3/src/msg/producer"

	"github.com/golang/mock/gomock"
)

// MockshardWriter is a mock of shardWriter interface
type MockshardWriter struct {
	ctrl     *gomock.Controller
	recorder *MockshardWriterMockRecorder
}

// MockshardWriterMockRecorder is the mock recorder for MockshardWriter
type MockshardWriterMockRecorder struct {
	mock *MockshardWriter
}

// NewMockshardWriter creates a new mock instance
func NewMockshardWriter(ctrl *gomock.Controller) *MockshardWriter {
	mock := &MockshardWriter{ctrl: ctrl}
	mock.recorder = &MockshardWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockshardWriter) EXPECT() *MockshardWriterMockRecorder {
	return m.recorder
}

// Write mocks base method
func (m *MockshardWriter) Write(rm *producer.RefCountedMessage) {
	m.ctrl.Call(m, "Write", rm)
}

// Write indicates an expected call of Write
func (mr *MockshardWriterMockRecorder) Write(rm interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockshardWriter)(nil).Write), rm)
}

// UpdateInstances mocks base method
func (m *MockshardWriter) UpdateInstances(instances []placement.Instance, cws map[string]consumerWriter) {
	m.ctrl.Call(m, "UpdateInstances", instances, cws)
}

// UpdateInstances indicates an expected call of UpdateInstances
func (mr *MockshardWriterMockRecorder) UpdateInstances(instances, cws interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInstances", reflect.TypeOf((*MockshardWriter)(nil).UpdateInstances), instances, cws)
}

// SetMessageTTLNanos mocks base method
func (m *MockshardWriter) SetMessageTTLNanos(value int64) {
	m.ctrl.Call(m, "SetMessageTTLNanos", value)
}

// SetMessageTTLNanos indicates an expected call of SetMessageTTLNanos
func (mr *MockshardWriterMockRecorder) SetMessageTTLNanos(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMessageTTLNanos", reflect.TypeOf((*MockshardWriter)(nil).SetMessageTTLNanos), value)
}

// Close mocks base method
func (m *MockshardWriter) Close() {
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockshardWriterMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockshardWriter)(nil).Close))
}

// QueueSize mocks base method
func (m *MockshardWriter) QueueSize() int {
	ret := m.ctrl.Call(m, "QueueSize")
	ret0, _ := ret[0].(int)
	return ret0
}

// QueueSize indicates an expected call of QueueSize
func (mr *MockshardWriterMockRecorder) QueueSize() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueSize", reflect.TypeOf((*MockshardWriter)(nil).QueueSize))
}
