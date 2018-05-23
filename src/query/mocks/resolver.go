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

// Code generated by MockGen. DO NOT EDIT.
// Source: policy/resolver/interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	"context"
	"reflect"
	"time"

	"github.com/m3db/m3/src/query/models"
	"github.com/m3db/m3/src/query/tsdb"

	"github.com/golang/mock/gomock"
)

// MockPolicyResolver is a mock of PolicyResolver interface
type MockPolicyResolver struct {
	ctrl     *gomock.Controller
	recorder *MockPolicyResolverMockRecorder
}

// MockPolicyResolverMockRecorder is the mock recorder for MockPolicyResolver
type MockPolicyResolverMockRecorder struct {
	mock *MockPolicyResolver
}

// NewMockPolicyResolver creates a new mock instance
func NewMockPolicyResolver(ctrl *gomock.Controller) *MockPolicyResolver {
	mock := &MockPolicyResolver{ctrl: ctrl}
	mock.recorder = &MockPolicyResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPolicyResolver) EXPECT() *MockPolicyResolverMockRecorder {
	return m.recorder
}

// Resolve mocks base method
func (m *MockPolicyResolver) Resolve(ctx context.Context, tagMatchers models.Matchers, startTime, endTime time.Time) ([]tsdb.FetchRequest, error) {
	ret := m.ctrl.Call(m, "Resolve", ctx, tagMatchers, startTime, endTime)
	ret0, _ := ret[0].([]tsdb.FetchRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Resolve indicates an expected call of Resolve
func (mr *MockPolicyResolverMockRecorder) Resolve(ctx, tagMatchers, startTime, endTime interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resolve", reflect.TypeOf((*MockPolicyResolver)(nil).Resolve), ctx, tagMatchers, startTime, endTime)
}