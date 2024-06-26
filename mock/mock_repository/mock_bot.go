// Code generated by MockGen. DO NOT EDIT.
// Source: bot.go
//
// Generated by this command:
//
//	mockgen -source=bot.go -destination=../../mock/mock_repository/mock_bot.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	http "net/http"
	reflect "reflect"

	linebot "github.com/line/line-bot-sdk-go/v7/linebot"
	model "github.com/ww24/linebot/domain/model"
	repository "github.com/ww24/linebot/domain/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockBot is a mock of Bot interface.
type MockBot struct {
	ctrl     *gomock.Controller
	recorder *MockBotMockRecorder
}

// MockBotMockRecorder is the mock recorder for MockBot.
type MockBotMockRecorder struct {
	mock *MockBot
}

// NewMockBot creates a new mock instance.
func NewMockBot(ctrl *gomock.Controller) *MockBot {
	mock := &MockBot{ctrl: ctrl}
	mock.recorder = &MockBotMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBot) EXPECT() *MockBotMockRecorder {
	return m.recorder
}

// EventsFromRequest mocks base method.
func (m *MockBot) EventsFromRequest(r *http.Request) ([]*model.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EventsFromRequest", r)
	ret0, _ := ret[0].([]*model.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EventsFromRequest indicates an expected call of EventsFromRequest.
func (mr *MockBotMockRecorder) EventsFromRequest(r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EventsFromRequest", reflect.TypeOf((*MockBot)(nil).EventsFromRequest), r)
}

// PushMessage mocks base method.
func (m *MockBot) PushMessage(arg0 context.Context, arg1 model.ConversationID, arg2 repository.MessageProvider) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushMessage", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushMessage indicates an expected call of PushMessage.
func (mr *MockBotMockRecorder) PushMessage(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushMessage", reflect.TypeOf((*MockBot)(nil).PushMessage), arg0, arg1, arg2)
}

// ReplyMessage mocks base method.
func (m *MockBot) ReplyMessage(arg0 context.Context, arg1 *model.Event, arg2 repository.MessageProvider) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplyMessage", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplyMessage indicates an expected call of ReplyMessage.
func (mr *MockBotMockRecorder) ReplyMessage(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplyMessage", reflect.TypeOf((*MockBot)(nil).ReplyMessage), arg0, arg1, arg2)
}

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockHandler) Handle(arg0 context.Context, arg1 *model.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Handle indicates an expected call of Handle.
func (mr *MockHandlerMockRecorder) Handle(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockHandler)(nil).Handle), arg0, arg1)
}

// MockMessageProviderSet is a mock of MessageProviderSet interface.
type MockMessageProviderSet struct {
	ctrl     *gomock.Controller
	recorder *MockMessageProviderSetMockRecorder
}

// MockMessageProviderSetMockRecorder is the mock recorder for MockMessageProviderSet.
type MockMessageProviderSetMockRecorder struct {
	mock *MockMessageProviderSet
}

// NewMockMessageProviderSet creates a new mock instance.
func NewMockMessageProviderSet(ctrl *gomock.Controller) *MockMessageProviderSet {
	mock := &MockMessageProviderSet{ctrl: ctrl}
	mock.recorder = &MockMessageProviderSetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageProviderSet) EXPECT() *MockMessageProviderSetMockRecorder {
	return m.recorder
}

// Image mocks base method.
func (m *MockMessageProviderSet) Image(originalURL, previewURL string) repository.MessageProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Image", originalURL, previewURL)
	ret0, _ := ret[0].(repository.MessageProvider)
	return ret0
}

// Image indicates an expected call of Image.
func (mr *MockMessageProviderSetMockRecorder) Image(originalURL, previewURL any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Image", reflect.TypeOf((*MockMessageProviderSet)(nil).Image), originalURL, previewURL)
}

// ReminderChoices mocks base method.
func (m *MockMessageProviderSet) ReminderChoices(arg0 string, arg1 []string, arg2 []model.ExecutorType) repository.MessageProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReminderChoices", arg0, arg1, arg2)
	ret0, _ := ret[0].(repository.MessageProvider)
	return ret0
}

// ReminderChoices indicates an expected call of ReminderChoices.
func (mr *MockMessageProviderSetMockRecorder) ReminderChoices(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReminderChoices", reflect.TypeOf((*MockMessageProviderSet)(nil).ReminderChoices), arg0, arg1, arg2)
}

// ReminderDeleteConfirmation mocks base method.
func (m *MockMessageProviderSet) ReminderDeleteConfirmation(text, data string) repository.MessageProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReminderDeleteConfirmation", text, data)
	ret0, _ := ret[0].(repository.MessageProvider)
	return ret0
}

// ReminderDeleteConfirmation indicates an expected call of ReminderDeleteConfirmation.
func (mr *MockMessageProviderSetMockRecorder) ReminderDeleteConfirmation(text, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReminderDeleteConfirmation", reflect.TypeOf((*MockMessageProviderSet)(nil).ReminderDeleteConfirmation), text, data)
}

// ReminderMenu mocks base method.
func (m *MockMessageProviderSet) ReminderMenu(arg0 string, arg1 model.ReminderReplyType, arg2 []*model.ReminderItem) repository.MessageProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReminderMenu", arg0, arg1, arg2)
	ret0, _ := ret[0].(repository.MessageProvider)
	return ret0
}

// ReminderMenu indicates an expected call of ReminderMenu.
func (mr *MockMessageProviderSetMockRecorder) ReminderMenu(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReminderMenu", reflect.TypeOf((*MockMessageProviderSet)(nil).ReminderMenu), arg0, arg1, arg2)
}

// ShoppingDeleteConfirmation mocks base method.
func (m *MockMessageProviderSet) ShoppingDeleteConfirmation(arg0 string) repository.MessageProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShoppingDeleteConfirmation", arg0)
	ret0, _ := ret[0].(repository.MessageProvider)
	return ret0
}

// ShoppingDeleteConfirmation indicates an expected call of ShoppingDeleteConfirmation.
func (mr *MockMessageProviderSetMockRecorder) ShoppingDeleteConfirmation(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShoppingDeleteConfirmation", reflect.TypeOf((*MockMessageProviderSet)(nil).ShoppingDeleteConfirmation), arg0)
}

// ShoppingMenu mocks base method.
func (m *MockMessageProviderSet) ShoppingMenu(arg0 string, arg1 model.ShoppingReplyType) repository.MessageProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShoppingMenu", arg0, arg1)
	ret0, _ := ret[0].(repository.MessageProvider)
	return ret0
}

// ShoppingMenu indicates an expected call of ShoppingMenu.
func (mr *MockMessageProviderSetMockRecorder) ShoppingMenu(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShoppingMenu", reflect.TypeOf((*MockMessageProviderSet)(nil).ShoppingMenu), arg0, arg1)
}

// Text mocks base method.
func (m *MockMessageProviderSet) Text(arg0 string) repository.MessageProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Text", arg0)
	ret0, _ := ret[0].(repository.MessageProvider)
	return ret0
}

// Text indicates an expected call of Text.
func (mr *MockMessageProviderSetMockRecorder) Text(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Text", reflect.TypeOf((*MockMessageProviderSet)(nil).Text), arg0)
}

// TimePicker mocks base method.
func (m *MockMessageProviderSet) TimePicker(text, data string) repository.MessageProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TimePicker", text, data)
	ret0, _ := ret[0].(repository.MessageProvider)
	return ret0
}

// TimePicker indicates an expected call of TimePicker.
func (mr *MockMessageProviderSetMockRecorder) TimePicker(text, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimePicker", reflect.TypeOf((*MockMessageProviderSet)(nil).TimePicker), text, data)
}

// MockMessageProvider is a mock of MessageProvider interface.
type MockMessageProvider struct {
	ctrl     *gomock.Controller
	recorder *MockMessageProviderMockRecorder
}

// MockMessageProviderMockRecorder is the mock recorder for MockMessageProvider.
type MockMessageProviderMockRecorder struct {
	mock *MockMessageProvider
}

// NewMockMessageProvider creates a new mock instance.
func NewMockMessageProvider(ctrl *gomock.Controller) *MockMessageProvider {
	mock := &MockMessageProvider{ctrl: ctrl}
	mock.recorder = &MockMessageProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageProvider) EXPECT() *MockMessageProviderMockRecorder {
	return m.recorder
}

// ToMessage mocks base method.
func (m *MockMessageProvider) ToMessage() linebot.SendingMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToMessage")
	ret0, _ := ret[0].(linebot.SendingMessage)
	return ret0
}

// ToMessage indicates an expected call of ToMessage.
func (mr *MockMessageProviderMockRecorder) ToMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToMessage", reflect.TypeOf((*MockMessageProvider)(nil).ToMessage))
}
