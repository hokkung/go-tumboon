// Code generated by MockGen. DO NOT EDIT.
// Source: donation.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/hokkung/go-tumboon/internal/model"
	service "github.com/hokkung/go-tumboon/internal/service/donation"
)

// MockDonationService is a mock of DonationService interface.
type MockDonationService struct {
	ctrl     *gomock.Controller
	recorder *MockDonationServiceMockRecorder
}

// MockDonationServiceMockRecorder is the mock recorder for MockDonationService.
type MockDonationServiceMockRecorder struct {
	mock *MockDonationService
}

// NewMockDonationService creates a new mock instance.
func NewMockDonationService(ctrl *gomock.Controller) *MockDonationService {
	mock := &MockDonationService{ctrl: ctrl}
	mock.recorder = &MockDonationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDonationService) EXPECT() *MockDonationServiceMockRecorder {
	return m.recorder
}

// Donate mocks base method.
func (m *MockDonationService) Donate(detail model.DonationDetail) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Donate", detail)
	ret0, _ := ret[0].(error)
	return ret0
}

// Donate indicates an expected call of Donate.
func (mr *MockDonationServiceMockRecorder) Donate(detail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Donate", reflect.TypeOf((*MockDonationService)(nil).Donate), detail)
}

// Donates mocks base method.
func (m *MockDonationService) Donates(details []model.DonationDetail) (*service.SummaryDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Donates", details)
	ret0, _ := ret[0].(*service.SummaryDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Donates indicates an expected call of Donates.
func (mr *MockDonationServiceMockRecorder) Donates(details interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Donates", reflect.TypeOf((*MockDonationService)(nil).Donates), details)
}

// MakePermit mocks base method.
func (m *MockDonationService) MakePermit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakePermit")
	ret0, _ := ret[0].(error)
	return ret0
}

// MakePermit indicates an expected call of MakePermit.
func (mr *MockDonationServiceMockRecorder) MakePermit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakePermit", reflect.TypeOf((*MockDonationService)(nil).MakePermit))
}
