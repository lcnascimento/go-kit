// Code generated by MockGen. DO NOT EDIT.
// Source: ./messaging/contracts.go
//
// Generated by this command:
//
//	mockgen -source=./messaging/contracts.go -destination=./messaging/mocks/mocks.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	cqrs "github.com/ThreeDotsLabs/watermill/components/cqrs"
	message "github.com/ThreeDotsLabs/watermill/message"
	gomock "go.uber.org/mock/gomock"
)

// MockPubSub is a mock of PubSub interface.
type MockPubSub struct {
	ctrl     *gomock.Controller
	recorder *MockPubSubMockRecorder
}

// MockPubSubMockRecorder is the mock recorder for MockPubSub.
type MockPubSubMockRecorder struct {
	mock *MockPubSub
}

// NewMockPubSub creates a new mock instance.
func NewMockPubSub(ctrl *gomock.Controller) *MockPubSub {
	mock := &MockPubSub{ctrl: ctrl}
	mock.recorder = &MockPubSubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPubSub) EXPECT() *MockPubSubMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockPubSub) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockPubSubMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPubSub)(nil).Close))
}

// Publish mocks base method.
func (m *MockPubSub) Publish(topic string, messages ...*message.Message) error {
	m.ctrl.T.Helper()
	varargs := []any{topic}
	for _, a := range messages {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Publish", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockPubSubMockRecorder) Publish(topic any, messages ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{topic}, messages...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockPubSub)(nil).Publish), varargs...)
}

// Subscribe mocks base method.
func (m *MockPubSub) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx, topic)
	ret0, _ := ret[0].(<-chan *message.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockPubSubMockRecorder) Subscribe(ctx, topic any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockPubSub)(nil).Subscribe), ctx, topic)
}

// MockBrokerCQRS is a mock of BrokerCQRS interface.
type MockBrokerCQRS struct {
	ctrl     *gomock.Controller
	recorder *MockBrokerCQRSMockRecorder
}

// MockBrokerCQRSMockRecorder is the mock recorder for MockBrokerCQRS.
type MockBrokerCQRSMockRecorder struct {
	mock *MockBrokerCQRS
}

// NewMockBrokerCQRS creates a new mock instance.
func NewMockBrokerCQRS(ctrl *gomock.Controller) *MockBrokerCQRS {
	mock := &MockBrokerCQRS{ctrl: ctrl}
	mock.recorder = &MockBrokerCQRSMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBrokerCQRS) EXPECT() *MockBrokerCQRSMockRecorder {
	return m.recorder
}

// AddCommandHandlers mocks base method.
func (m *MockBrokerCQRS) AddCommandHandlers(ctx context.Context, handlers ...cqrs.CommandHandler) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range handlers {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddCommandHandlers", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCommandHandlers indicates an expected call of AddCommandHandlers.
func (mr *MockBrokerCQRSMockRecorder) AddCommandHandlers(ctx any, handlers ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, handlers...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCommandHandlers", reflect.TypeOf((*MockBrokerCQRS)(nil).AddCommandHandlers), varargs...)
}

// AddEventHandlers mocks base method.
func (m *MockBrokerCQRS) AddEventHandlers(ctx context.Context, handlers ...cqrs.EventHandler) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range handlers {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddEventHandlers", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEventHandlers indicates an expected call of AddEventHandlers.
func (mr *MockBrokerCQRSMockRecorder) AddEventHandlers(ctx any, handlers ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, handlers...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEventHandlers", reflect.TypeOf((*MockBrokerCQRS)(nil).AddEventHandlers), varargs...)
}

// Running mocks base method.
func (m *MockBrokerCQRS) Running(ctx context.Context) chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Running", ctx)
	ret0, _ := ret[0].(chan struct{})
	return ret0
}

// Running indicates an expected call of Running.
func (mr *MockBrokerCQRSMockRecorder) Running(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Running", reflect.TypeOf((*MockBrokerCQRS)(nil).Running), ctx)
}

// SendCommand mocks base method.
func (m *MockBrokerCQRS) SendCommand(ctx context.Context, command any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCommand", ctx, command)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCommand indicates an expected call of SendCommand.
func (mr *MockBrokerCQRSMockRecorder) SendCommand(ctx, command any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCommand", reflect.TypeOf((*MockBrokerCQRS)(nil).SendCommand), ctx, command)
}

// SendEvent mocks base method.
func (m *MockBrokerCQRS) SendEvent(ctx context.Context, event any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEvent", ctx, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEvent indicates an expected call of SendEvent.
func (mr *MockBrokerCQRSMockRecorder) SendEvent(ctx, event any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEvent", reflect.TypeOf((*MockBrokerCQRS)(nil).SendEvent), ctx, event)
}

// Start mocks base method.
func (m *MockBrokerCQRS) Start(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockBrokerCQRSMockRecorder) Start(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockBrokerCQRS)(nil).Start), ctx)
}

// Stop mocks base method.
func (m *MockBrokerCQRS) Stop(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockBrokerCQRSMockRecorder) Stop(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockBrokerCQRS)(nil).Stop), ctx)
}