// Code generated by MockGen. DO NOT EDIT.
// Source: nl.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/ww24/linebot/domain/model"
)

// MockNLParser is a mock of NLParser interface.
type MockNLParser struct {
	ctrl     *gomock.Controller
	recorder *MockNLParserMockRecorder
}

// MockNLParserMockRecorder is the mock recorder for MockNLParser.
type MockNLParserMockRecorder struct {
	mock *MockNLParser
}

// NewMockNLParser creates a new mock instance.
func NewMockNLParser(ctrl *gomock.Controller) *MockNLParser {
	mock := &MockNLParser{ctrl: ctrl}
	mock.recorder = &MockNLParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNLParser) EXPECT() *MockNLParserMockRecorder {
	return m.recorder
}

// Parse mocks base method.
func (m *MockNLParser) Parse(arg0 string) *model.Item {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", arg0)
	ret0, _ := ret[0].(*model.Item)
	return ret0
}

// Parse indicates an expected call of Parse.
func (mr *MockNLParserMockRecorder) Parse(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockNLParser)(nil).Parse), arg0)
}
