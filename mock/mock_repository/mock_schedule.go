// Code generated by MockGen. DO NOT EDIT.
// Source: schedule.go
//
// Generated by this command:
//
//	mockgen -source=schedule.go -destination=../../mock/mock_repository/mock_schedule.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"
	time "time"

	model "github.com/ww24/linebot/domain/model"
	gomock "go.uber.org/mock/gomock"
)

// MockScheduleHandler is a mock of ScheduleHandler interface.
type MockScheduleHandler struct {
	ctrl     *gomock.Controller
	recorder *MockScheduleHandlerMockRecorder
}

// MockScheduleHandlerMockRecorder is the mock recorder for MockScheduleHandler.
type MockScheduleHandlerMockRecorder struct {
	mock *MockScheduleHandler
}

// NewMockScheduleHandler creates a new mock instance.
func NewMockScheduleHandler(ctrl *gomock.Controller) *MockScheduleHandler {
	mock := &MockScheduleHandler{ctrl: ctrl}
	mock.recorder = &MockScheduleHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduleHandler) EXPECT() *MockScheduleHandlerMockRecorder {
	return m.recorder
}

// HandleSchedule mocks base method.
func (m *MockScheduleHandler) HandleSchedule(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleSchedule", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleSchedule indicates an expected call of HandleSchedule.
func (mr *MockScheduleHandlerMockRecorder) HandleSchedule(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleSchedule", reflect.TypeOf((*MockScheduleHandler)(nil).HandleSchedule), arg0)
}

// MockScheduleSynchronizer is a mock of ScheduleSynchronizer interface.
type MockScheduleSynchronizer struct {
	ctrl     *gomock.Controller
	recorder *MockScheduleSynchronizerMockRecorder
}

// MockScheduleSynchronizerMockRecorder is the mock recorder for MockScheduleSynchronizer.
type MockScheduleSynchronizerMockRecorder struct {
	mock *MockScheduleSynchronizer
}

// NewMockScheduleSynchronizer creates a new mock instance.
func NewMockScheduleSynchronizer(ctrl *gomock.Controller) *MockScheduleSynchronizer {
	mock := &MockScheduleSynchronizer{ctrl: ctrl}
	mock.recorder = &MockScheduleSynchronizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduleSynchronizer) EXPECT() *MockScheduleSynchronizerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockScheduleSynchronizer) Create(arg0 context.Context, arg1 model.ConversationID, arg2 *model.ReminderItem, arg3 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockScheduleSynchronizerMockRecorder) Create(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockScheduleSynchronizer)(nil).Create), arg0, arg1, arg2, arg3)
}

// Delete mocks base method.
func (m *MockScheduleSynchronizer) Delete(arg0 context.Context, arg1 model.ConversationID, arg2 *model.ReminderItem, arg3 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockScheduleSynchronizerMockRecorder) Delete(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockScheduleSynchronizer)(nil).Delete), arg0, arg1, arg2, arg3)
}

// Sync mocks base method.
func (m *MockScheduleSynchronizer) Sync(arg0 context.Context, arg1 model.ConversationID, arg2 model.ReminderItems, arg3 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Sync indicates an expected call of Sync.
func (mr *MockScheduleSynchronizerMockRecorder) Sync(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockScheduleSynchronizer)(nil).Sync), arg0, arg1, arg2, arg3)
}

// MockRemindHandler is a mock of RemindHandler interface.
type MockRemindHandler struct {
	ctrl     *gomock.Controller
	recorder *MockRemindHandlerMockRecorder
}

// MockRemindHandlerMockRecorder is the mock recorder for MockRemindHandler.
type MockRemindHandlerMockRecorder struct {
	mock *MockRemindHandler
}

// NewMockRemindHandler creates a new mock instance.
func NewMockRemindHandler(ctrl *gomock.Controller) *MockRemindHandler {
	mock := &MockRemindHandler{ctrl: ctrl}
	mock.recorder = &MockRemindHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRemindHandler) EXPECT() *MockRemindHandlerMockRecorder {
	return m.recorder
}

// HandleReminder mocks base method.
func (m *MockRemindHandler) HandleReminder(arg0 context.Context, arg1 *model.ReminderItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleReminder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleReminder indicates an expected call of HandleReminder.
func (mr *MockRemindHandlerMockRecorder) HandleReminder(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleReminder", reflect.TypeOf((*MockRemindHandler)(nil).HandleReminder), arg0, arg1)
}
