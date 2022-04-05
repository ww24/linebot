// Code generated by MockGen. DO NOT EDIT.
// Source: reminder.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/ww24/linebot/domain/model"
)

// MockReminder is a mock of Reminder interface.
type MockReminder struct {
	ctrl     *gomock.Controller
	recorder *MockReminderMockRecorder
}

// MockReminderMockRecorder is the mock recorder for MockReminder.
type MockReminderMockRecorder struct {
	mock *MockReminder
}

// NewMockReminder creates a new mock instance.
func NewMockReminder(ctrl *gomock.Controller) *MockReminder {
	mock := &MockReminder{ctrl: ctrl}
	mock.recorder = &MockReminderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReminder) EXPECT() *MockReminderMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockReminder) Add(arg0 context.Context, arg1 *model.ReminderItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockReminderMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockReminder)(nil).Add), arg0, arg1)
}

// Delete mocks base method.
func (m *MockReminder) Delete(arg0 context.Context, arg1 model.ConversationID, arg2 model.ReminderItemID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockReminderMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockReminder)(nil).Delete), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockReminder) Get(arg0 context.Context, arg1 model.ConversationID, arg2 model.ReminderItemID) (*model.ReminderItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.ReminderItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockReminderMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockReminder)(nil).Get), arg0, arg1, arg2)
}

// List mocks base method.
func (m *MockReminder) List(arg0 context.Context, arg1 model.ConversationID) ([]*model.ReminderItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*model.ReminderItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockReminderMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockReminder)(nil).List), arg0, arg1)
}

// ListAll mocks base method.
func (m *MockReminder) ListAll(arg0 context.Context) ([]*model.ReminderItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAll", arg0)
	ret0, _ := ret[0].([]*model.ReminderItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAll indicates an expected call of ListAll.
func (mr *MockReminderMockRecorder) ListAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAll", reflect.TypeOf((*MockReminder)(nil).ListAll), arg0)
}
