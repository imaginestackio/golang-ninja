// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ImagineDevOps/Hands-On-Software-Engineering-with-Golang/Chapter12/dbspgraph/proto (interfaces: JobQueue_JobStreamServer)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	proto "github.com/ImagineDevOps/Hands-On-Software-Engineering-with-Golang/Chapter12/dbspgraph/proto"
	gomock "github.com/golang/mock/gomock"
	metadata "google.golang.org/grpc/metadata"
	reflect "reflect"
)

// MockJobQueue_JobStreamServer is a mock of JobQueue_JobStreamServer interface
type MockJobQueue_JobStreamServer struct {
	ctrl     *gomock.Controller
	recorder *MockJobQueue_JobStreamServerMockRecorder
}

// MockJobQueue_JobStreamServerMockRecorder is the mock recorder for MockJobQueue_JobStreamServer
type MockJobQueue_JobStreamServerMockRecorder struct {
	mock *MockJobQueue_JobStreamServer
}

// NewMockJobQueue_JobStreamServer creates a new mock instance
func NewMockJobQueue_JobStreamServer(ctrl *gomock.Controller) *MockJobQueue_JobStreamServer {
	mock := &MockJobQueue_JobStreamServer{ctrl: ctrl}
	mock.recorder = &MockJobQueue_JobStreamServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockJobQueue_JobStreamServer) EXPECT() *MockJobQueue_JobStreamServerMockRecorder {
	return m.recorder
}

// Context mocks base method
func (m *MockJobQueue_JobStreamServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockJobQueue_JobStreamServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockJobQueue_JobStreamServer)(nil).Context))
}

// Recv mocks base method
func (m *MockJobQueue_JobStreamServer) Recv() (*proto.WorkerPayload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*proto.WorkerPayload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockJobQueue_JobStreamServerMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockJobQueue_JobStreamServer)(nil).Recv))
}

// RecvMsg mocks base method
func (m *MockJobQueue_JobStreamServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockJobQueue_JobStreamServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockJobQueue_JobStreamServer)(nil).RecvMsg), arg0)
}

// Send mocks base method
func (m *MockJobQueue_JobStreamServer) Send(arg0 *proto.MasterPayload) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockJobQueue_JobStreamServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockJobQueue_JobStreamServer)(nil).Send), arg0)
}

// SendHeader mocks base method
func (m *MockJobQueue_JobStreamServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader
func (mr *MockJobQueue_JobStreamServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockJobQueue_JobStreamServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method
func (m *MockJobQueue_JobStreamServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockJobQueue_JobStreamServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockJobQueue_JobStreamServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method
func (m *MockJobQueue_JobStreamServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader
func (mr *MockJobQueue_JobStreamServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockJobQueue_JobStreamServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method
func (m *MockJobQueue_JobStreamServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer
func (mr *MockJobQueue_JobStreamServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockJobQueue_JobStreamServer)(nil).SetTrailer), arg0)
}