// Code generated by MockGen. DO NOT EDIT.
// Source: conversation.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/ww24/linebot/domain/model"
)

// MockConversation is a mock of Conversation interface.
type MockConversation struct {
	ctrl     *gomock.Controller
	recorder *MockConversationMockRecorder
}

// MockConversationMockRecorder is the mock recorder for MockConversation.
type MockConversationMockRecorder struct {
	mock *MockConversation
}

// NewMockConversation creates a new mock instance.
func NewMockConversation(ctrl *gomock.Controller) *MockConversation {
	mock := &MockConversation{ctrl: ctrl}
	mock.recorder = &MockConversationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConversation) EXPECT() *MockConversationMockRecorder {
	return m.recorder
}

// AddShoppingItem mocks base method.
func (m *MockConversation) AddShoppingItem(arg0 context.Context, arg1 ...*model.ShoppingItem) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddShoppingItem", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddShoppingItem indicates an expected call of AddShoppingItem.
func (mr *MockConversationMockRecorder) AddShoppingItem(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddShoppingItem", reflect.TypeOf((*MockConversation)(nil).AddShoppingItem), varargs...)
}

// DeleteAllShoppingItem mocks base method.
func (m *MockConversation) DeleteAllShoppingItem(arg0 context.Context, arg1 model.ConversationID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllShoppingItem", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllShoppingItem indicates an expected call of DeleteAllShoppingItem.
func (mr *MockConversationMockRecorder) DeleteAllShoppingItem(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllShoppingItem", reflect.TypeOf((*MockConversation)(nil).DeleteAllShoppingItem), arg0, arg1)
}

// DeleteShoppingItems mocks base method.
func (m *MockConversation) DeleteShoppingItems(ctx context.Context, conversationID model.ConversationID, ids []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteShoppingItems", ctx, conversationID, ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteShoppingItems indicates an expected call of DeleteShoppingItems.
func (mr *MockConversationMockRecorder) DeleteShoppingItems(ctx, conversationID, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteShoppingItems", reflect.TypeOf((*MockConversation)(nil).DeleteShoppingItems), ctx, conversationID, ids)
}

// FindShoppingItem mocks base method.
func (m *MockConversation) FindShoppingItem(arg0 context.Context, arg1 model.ConversationID) ([]*model.ShoppingItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindShoppingItem", arg0, arg1)
	ret0, _ := ret[0].([]*model.ShoppingItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindShoppingItem indicates an expected call of FindShoppingItem.
func (mr *MockConversationMockRecorder) FindShoppingItem(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindShoppingItem", reflect.TypeOf((*MockConversation)(nil).FindShoppingItem), arg0, arg1)
}

// GetStatus mocks base method.
func (m *MockConversation) GetStatus(arg0 context.Context, arg1 model.ConversationID) (*model.ConversationStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", arg0, arg1)
	ret0, _ := ret[0].(*model.ConversationStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockConversationMockRecorder) GetStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockConversation)(nil).GetStatus), arg0, arg1)
}

// SetStatus mocks base method.
func (m *MockConversation) SetStatus(arg0 context.Context, arg1 *model.ConversationStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStatus", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStatus indicates an expected call of SetStatus.
func (mr *MockConversationMockRecorder) SetStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatus", reflect.TypeOf((*MockConversation)(nil).SetStatus), arg0, arg1)
}