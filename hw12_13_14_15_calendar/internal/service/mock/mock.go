// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"
	time "time"

	models "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockEvent is a mock of Event interface.
type MockEvent struct {
	ctrl     *gomock.Controller
	recorder *MockEventMockRecorder
}

// MockEventMockRecorder is the mock recorder for MockEvent.
type MockEventMockRecorder struct {
	mock *MockEvent
}

// NewMockEvent creates a new mock instance.
func NewMockEvent(ctrl *gomock.Controller) *MockEvent {
	mock := &MockEvent{ctrl: ctrl}
	mock.recorder = &MockEventMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEvent) EXPECT() *MockEventMockRecorder {
	return m.recorder
}

// CreateEvent mocks base method.
func (m *MockEvent) CreateEvent(ctx context.Context, event models.Event) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEvent", ctx, event)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEvent indicates an expected call of CreateEvent.
func (mr *MockEventMockRecorder) CreateEvent(ctx, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEvent", reflect.TypeOf((*MockEvent)(nil).CreateEvent), ctx, event)
}

// DeleteEvent mocks base method.
func (m *MockEvent) DeleteEvent(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEvent", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockEventMockRecorder) DeleteEvent(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockEvent)(nil).DeleteEvent), ctx, id)
}

// DeleteOutdatedEvents mocks base method.
func (m *MockEvent) DeleteOutdatedEvents(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOutdatedEvents", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOutdatedEvents indicates an expected call of DeleteOutdatedEvents.
func (mr *MockEventMockRecorder) DeleteOutdatedEvents(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOutdatedEvents", reflect.TypeOf((*MockEvent)(nil).DeleteOutdatedEvents), ctx)
}

// GetAllByDayEvents mocks base method.
func (m *MockEvent) GetAllByDayEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByDayEvents", ctx, date)
	ret0, _ := ret[0].([]models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByDayEvents indicates an expected call of GetAllByDayEvents.
func (mr *MockEventMockRecorder) GetAllByDayEvents(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByDayEvents", reflect.TypeOf((*MockEvent)(nil).GetAllByDayEvents), ctx, date)
}

// GetAllByMonthEvents mocks base method.
func (m *MockEvent) GetAllByMonthEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByMonthEvents", ctx, date)
	ret0, _ := ret[0].([]models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByMonthEvents indicates an expected call of GetAllByMonthEvents.
func (mr *MockEventMockRecorder) GetAllByMonthEvents(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByMonthEvents", reflect.TypeOf((*MockEvent)(nil).GetAllByMonthEvents), ctx, date)
}

// GetAllByWeekEvents mocks base method.
func (m *MockEvent) GetAllByWeekEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByWeekEvents", ctx, date)
	ret0, _ := ret[0].([]models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByWeekEvents indicates an expected call of GetAllByWeekEvents.
func (mr *MockEventMockRecorder) GetAllByWeekEvents(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByWeekEvents", reflect.TypeOf((*MockEvent)(nil).GetAllByWeekEvents), ctx, date)
}

// UpdateEvent mocks base method.
func (m *MockEvent) UpdateEvent(ctx context.Context, id string, event models.Event) (models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", ctx, id, event)
	ret0, _ := ret[0].(models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockEventMockRecorder) UpdateEvent(ctx, id, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockEvent)(nil).UpdateEvent), ctx, id, event)
}

// MockNotification is a mock of Notification interface.
type MockNotification struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationMockRecorder
}

// MockNotificationMockRecorder is the mock recorder for MockNotification.
type MockNotificationMockRecorder struct {
	mock *MockNotification
}

// NewMockNotification creates a new mock instance.
func NewMockNotification(ctrl *gomock.Controller) *MockNotification {
	mock := &MockNotification{ctrl: ctrl}
	mock.recorder = &MockNotificationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotification) EXPECT() *MockNotificationMockRecorder {
	return m.recorder
}

// GetNotificationInAdvance mocks base method.
func (m *MockNotification) GetNotificationInAdvance(ctx context.Context) ([]models.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotificationInAdvance", ctx)
	ret0, _ := ret[0].([]models.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotificationInAdvance indicates an expected call of GetNotificationInAdvance.
func (mr *MockNotificationMockRecorder) GetNotificationInAdvance(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotificationInAdvance", reflect.TypeOf((*MockNotification)(nil).GetNotificationInAdvance), ctx)
}

// UpdateScheduledNotification mocks base method.
func (m *MockNotification) UpdateScheduledNotification(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateScheduledNotification", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateScheduledNotification indicates an expected call of UpdateScheduledNotification.
func (mr *MockNotificationMockRecorder) UpdateScheduledNotification(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateScheduledNotification", reflect.TypeOf((*MockNotification)(nil).UpdateScheduledNotification), ctx, id)
}

// MockServices is a mock of Services interface.
type MockServices struct {
	ctrl     *gomock.Controller
	recorder *MockServicesMockRecorder
}

// MockServicesMockRecorder is the mock recorder for MockServices.
type MockServicesMockRecorder struct {
	mock *MockServices
}

// NewMockServices creates a new mock instance.
func NewMockServices(ctrl *gomock.Controller) *MockServices {
	mock := &MockServices{ctrl: ctrl}
	mock.recorder = &MockServicesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServices) EXPECT() *MockServicesMockRecorder {
	return m.recorder
}

// CreateEvent mocks base method.
func (m *MockServices) CreateEvent(ctx context.Context, event models.Event) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEvent", ctx, event)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEvent indicates an expected call of CreateEvent.
func (mr *MockServicesMockRecorder) CreateEvent(ctx, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEvent", reflect.TypeOf((*MockServices)(nil).CreateEvent), ctx, event)
}

// DeleteEvent mocks base method.
func (m *MockServices) DeleteEvent(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEvent", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockServicesMockRecorder) DeleteEvent(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockServices)(nil).DeleteEvent), ctx, id)
}

// DeleteOutdatedEvents mocks base method.
func (m *MockServices) DeleteOutdatedEvents(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOutdatedEvents", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOutdatedEvents indicates an expected call of DeleteOutdatedEvents.
func (mr *MockServicesMockRecorder) DeleteOutdatedEvents(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOutdatedEvents", reflect.TypeOf((*MockServices)(nil).DeleteOutdatedEvents), ctx)
}

// GetAllByDayEvents mocks base method.
func (m *MockServices) GetAllByDayEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByDayEvents", ctx, date)
	ret0, _ := ret[0].([]models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByDayEvents indicates an expected call of GetAllByDayEvents.
func (mr *MockServicesMockRecorder) GetAllByDayEvents(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByDayEvents", reflect.TypeOf((*MockServices)(nil).GetAllByDayEvents), ctx, date)
}

// GetAllByMonthEvents mocks base method.
func (m *MockServices) GetAllByMonthEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByMonthEvents", ctx, date)
	ret0, _ := ret[0].([]models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByMonthEvents indicates an expected call of GetAllByMonthEvents.
func (mr *MockServicesMockRecorder) GetAllByMonthEvents(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByMonthEvents", reflect.TypeOf((*MockServices)(nil).GetAllByMonthEvents), ctx, date)
}

// GetAllByWeekEvents mocks base method.
func (m *MockServices) GetAllByWeekEvents(ctx context.Context, date time.Time) ([]models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByWeekEvents", ctx, date)
	ret0, _ := ret[0].([]models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByWeekEvents indicates an expected call of GetAllByWeekEvents.
func (mr *MockServicesMockRecorder) GetAllByWeekEvents(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByWeekEvents", reflect.TypeOf((*MockServices)(nil).GetAllByWeekEvents), ctx, date)
}

// GetNotificationInAdvance mocks base method.
func (m *MockServices) GetNotificationInAdvance(ctx context.Context) ([]models.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotificationInAdvance", ctx)
	ret0, _ := ret[0].([]models.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotificationInAdvance indicates an expected call of GetNotificationInAdvance.
func (mr *MockServicesMockRecorder) GetNotificationInAdvance(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotificationInAdvance", reflect.TypeOf((*MockServices)(nil).GetNotificationInAdvance), ctx)
}

// UpdateEvent mocks base method.
func (m *MockServices) UpdateEvent(ctx context.Context, id string, event models.Event) (models.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", ctx, id, event)
	ret0, _ := ret[0].(models.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockServicesMockRecorder) UpdateEvent(ctx, id, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockServices)(nil).UpdateEvent), ctx, id, event)
}

// UpdateScheduledNotification mocks base method.
func (m *MockServices) UpdateScheduledNotification(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateScheduledNotification", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateScheduledNotification indicates an expected call of UpdateScheduledNotification.
func (mr *MockServicesMockRecorder) UpdateScheduledNotification(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateScheduledNotification", reflect.TypeOf((*MockServices)(nil).UpdateScheduledNotification), ctx, id)
}
