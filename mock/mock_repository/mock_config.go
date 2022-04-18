// Code generated by MockGen. DO NOT EDIT.
// Source: config.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	url "net/url"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	model "github.com/ww24/linebot/domain/model"
	repository "github.com/ww24/linebot/domain/repository"
)

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// Addr mocks base method.
func (m *MockConfig) Addr() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Addr")
	ret0, _ := ret[0].(string)
	return ret0
}

// Addr indicates an expected call of Addr.
func (mr *MockConfigMockRecorder) Addr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Addr", reflect.TypeOf((*MockConfig)(nil).Addr))
}

// BrowserTimeout mocks base method.
func (m *MockConfig) BrowserTimeout() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BrowserTimeout")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// BrowserTimeout indicates an expected call of BrowserTimeout.
func (mr *MockConfigMockRecorder) BrowserTimeout() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BrowserTimeout", reflect.TypeOf((*MockConfig)(nil).BrowserTimeout))
}

// CloudTasksLocation mocks base method.
func (m *MockConfig) CloudTasksLocation() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudTasksLocation")
	ret0, _ := ret[0].(string)
	return ret0
}

// CloudTasksLocation indicates an expected call of CloudTasksLocation.
func (mr *MockConfigMockRecorder) CloudTasksLocation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudTasksLocation", reflect.TypeOf((*MockConfig)(nil).CloudTasksLocation))
}

// CloudTasksQueue mocks base method.
func (m *MockConfig) CloudTasksQueue() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudTasksQueue")
	ret0, _ := ret[0].(string)
	return ret0
}

// CloudTasksQueue indicates an expected call of CloudTasksQueue.
func (mr *MockConfigMockRecorder) CloudTasksQueue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudTasksQueue", reflect.TypeOf((*MockConfig)(nil).CloudTasksQueue))
}

// ConversationIDs mocks base method.
func (m *MockConfig) ConversationIDs() repository.ConversationIDs {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConversationIDs")
	ret0, _ := ret[0].(repository.ConversationIDs)
	return ret0
}

// ConversationIDs indicates an expected call of ConversationIDs.
func (mr *MockConfigMockRecorder) ConversationIDs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConversationIDs", reflect.TypeOf((*MockConfig)(nil).ConversationIDs))
}

// DefaultLocation mocks base method.
func (m *MockConfig) DefaultLocation() *time.Location {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultLocation")
	ret0, _ := ret[0].(*time.Location)
	return ret0
}

// DefaultLocation indicates an expected call of DefaultLocation.
func (mr *MockConfigMockRecorder) DefaultLocation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultLocation", reflect.TypeOf((*MockConfig)(nil).DefaultLocation))
}

// ImageBucket mocks base method.
func (m *MockConfig) ImageBucket() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageBucket")
	ret0, _ := ret[0].(string)
	return ret0
}

// ImageBucket indicates an expected call of ImageBucket.
func (mr *MockConfigMockRecorder) ImageBucket() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageBucket", reflect.TypeOf((*MockConfig)(nil).ImageBucket))
}

// InvokerServiceAccountEmail mocks base method.
func (m *MockConfig) InvokerServiceAccountEmail() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvokerServiceAccountEmail")
	ret0, _ := ret[0].(string)
	return ret0
}

// InvokerServiceAccountEmail indicates an expected call of InvokerServiceAccountEmail.
func (mr *MockConfigMockRecorder) InvokerServiceAccountEmail() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokerServiceAccountEmail", reflect.TypeOf((*MockConfig)(nil).InvokerServiceAccountEmail))
}

// InvokerServiceAccountID mocks base method.
func (m *MockConfig) InvokerServiceAccountID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvokerServiceAccountID")
	ret0, _ := ret[0].(string)
	return ret0
}

// InvokerServiceAccountID indicates an expected call of InvokerServiceAccountID.
func (mr *MockConfigMockRecorder) InvokerServiceAccountID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokerServiceAccountID", reflect.TypeOf((*MockConfig)(nil).InvokerServiceAccountID))
}

// LINEChannelSecret mocks base method.
func (m *MockConfig) LINEChannelSecret() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LINEChannelSecret")
	ret0, _ := ret[0].(string)
	return ret0
}

// LINEChannelSecret indicates an expected call of LINEChannelSecret.
func (mr *MockConfigMockRecorder) LINEChannelSecret() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LINEChannelSecret", reflect.TypeOf((*MockConfig)(nil).LINEChannelSecret))
}

// LINEChannelToken mocks base method.
func (m *MockConfig) LINEChannelToken() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LINEChannelToken")
	ret0, _ := ret[0].(string)
	return ret0
}

// LINEChannelToken indicates an expected call of LINEChannelToken.
func (mr *MockConfigMockRecorder) LINEChannelToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LINEChannelToken", reflect.TypeOf((*MockConfig)(nil).LINEChannelToken))
}

// ServiceEndpoint mocks base method.
func (m *MockConfig) ServiceEndpoint(path string) (*url.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServiceEndpoint", path)
	ret0, _ := ret[0].(*url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ServiceEndpoint indicates an expected call of ServiceEndpoint.
func (mr *MockConfigMockRecorder) ServiceEndpoint(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceEndpoint", reflect.TypeOf((*MockConfig)(nil).ServiceEndpoint), path)
}

// WeatherAPI mocks base method.
func (m *MockConfig) WeatherAPI() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WeatherAPI")
	ret0, _ := ret[0].(string)
	return ret0
}

// WeatherAPI indicates an expected call of WeatherAPI.
func (mr *MockConfigMockRecorder) WeatherAPI() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WeatherAPI", reflect.TypeOf((*MockConfig)(nil).WeatherAPI))
}

// WeatherAPITimeout mocks base method.
func (m *MockConfig) WeatherAPITimeout() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WeatherAPITimeout")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// WeatherAPITimeout indicates an expected call of WeatherAPITimeout.
func (mr *MockConfigMockRecorder) WeatherAPITimeout() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WeatherAPITimeout", reflect.TypeOf((*MockConfig)(nil).WeatherAPITimeout))
}

// MockConversationIDs is a mock of ConversationIDs interface.
type MockConversationIDs struct {
	ctrl     *gomock.Controller
	recorder *MockConversationIDsMockRecorder
}

// MockConversationIDsMockRecorder is the mock recorder for MockConversationIDs.
type MockConversationIDsMockRecorder struct {
	mock *MockConversationIDs
}

// NewMockConversationIDs creates a new mock instance.
func NewMockConversationIDs(ctrl *gomock.Controller) *MockConversationIDs {
	mock := &MockConversationIDs{ctrl: ctrl}
	mock.recorder = &MockConversationIDsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConversationIDs) EXPECT() *MockConversationIDsMockRecorder {
	return m.recorder
}

// Available mocks base method.
func (m *MockConversationIDs) Available(arg0 model.ConversationID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Available", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Available indicates an expected call of Available.
func (mr *MockConversationIDsMockRecorder) Available(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Available", reflect.TypeOf((*MockConversationIDs)(nil).Available), arg0)
}

// List mocks base method.
func (m *MockConversationIDs) List() []model.ConversationID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]model.ConversationID)
	return ret0
}

// List indicates an expected call of List.
func (mr *MockConversationIDsMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockConversationIDs)(nil).List))
}
