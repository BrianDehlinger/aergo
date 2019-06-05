// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aergoio/aergo/p2p/p2pcommon (interfaces: MsgSigner)

// Package p2pmock is a generated GoMock package.
package p2pmock

import (
	"github.com/aergoio/aergo/types"
	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core"
	"reflect"
)

// MockMsgSigner is a mock of MsgSigner interface
type MockMsgSigner struct {
	ctrl     *gomock.Controller
	recorder *MockMsgSignerMockRecorder
}

// MockMsgSignerMockRecorder is the mock recorder for MockMsgSigner
type MockMsgSignerMockRecorder struct {
	mock *MockMsgSigner
}

// NewMockMsgSigner creates a new mock instance
func NewMockMsgSigner(ctrl *gomock.Controller) *MockMsgSigner {
	mock := &MockMsgSigner{ctrl: ctrl}
	mock.recorder = &MockMsgSignerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMsgSigner) EXPECT() *MockMsgSignerMockRecorder {
	return m.recorder
}

// SignMsg mocks base method
func (m *MockMsgSigner) SignMsg(arg0 *types.P2PMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignMsg indicates an expected call of SignMsg
func (mr *MockMsgSignerMockRecorder) SignMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignMsg", reflect.TypeOf((*MockMsgSigner)(nil).SignMsg), arg0)
}

// VerifyMsg mocks base method
func (m *MockMsgSigner) VerifyMsg(arg0 *types.P2PMessage, arg1 core.PeerID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyMsg", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyMsg indicates an expected call of VerifyMsg
func (mr *MockMsgSignerMockRecorder) VerifyMsg(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyMsg", reflect.TypeOf((*MockMsgSigner)(nil).VerifyMsg), arg0, arg1)
}
