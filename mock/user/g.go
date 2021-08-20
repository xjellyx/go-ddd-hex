// Code generated by MockGen. DO NOT EDIT.
// Source: internal/adapter/xhttp/xgin/handler/user.go

// Package mock_handler is a generated GoMock package.
package mock_handler

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockUserHandler is a mock of UserHandler interface.
type MockUserHandler struct {
	ctrl     *gomock.Controller
	recorder *MockUserHandlerMockRecorder
}

// MockUserHandlerMockRecorder is the mock recorder for MockUserHandler.
type MockUserHandlerMockRecorder struct {
	mock *MockUserHandler
}

// NewMockUserHandler creates a new mock instance.
func NewMockUserHandler(ctrl *gomock.Controller) *MockUserHandler {
	mock := &MockUserHandler{ctrl: ctrl}
	mock.recorder = &MockUserHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserHandler) EXPECT() *MockUserHandlerMockRecorder {
	return m.recorder
}

// ChangePasswd mocks base method.
func (m *MockUserHandler) ChangePasswd(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ChangePasswd", ctx)
}

// ChangePasswd indicates an expected call of ChangePasswd.
func (mr *MockUserHandlerMockRecorder) ChangePasswd(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePasswd", reflect.TypeOf((*MockUserHandler)(nil).ChangePasswd), ctx)
}

// Create mocks base method.
func (m *MockUserHandler) Create(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Create", ctx)
}

// Create indicates an expected call of Create.
func (mr *MockUserHandlerMockRecorder) Create(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserHandler)(nil).Create), ctx)
}

// Get mocks base method.
func (m *MockUserHandler) Get(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Get", ctx)
}

// Get indicates an expected call of Get.
func (mr *MockUserHandlerMockRecorder) Get(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserHandler)(nil).Get), ctx)
}
