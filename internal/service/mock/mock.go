package mock

//
//import (
//	"github.com/golang/mock/gomock"
//)
//
//// MockAuthorization is a mock of Authorization interface.
//type MockAuthorization struct {
//	ctrl     *gomock.Controller
//	recorder *MockAuthorizationMockRecorder
//}
//
//// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
//type MockAuthorizationMockRecorder struct {
//	mock *MockAuthorization
//}
//
//// NewMockAuthorization creates a new mock instance.
//func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
//	mock := &MockAuthorization{ctrl: ctrl}
//	mock.recorder = &MockAuthorizationMockRecorder{mock}
//	return mock
//}
//
//// EXPECT returns an object that allows the caller to indicate expected use.
//func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
//	return m.recorder
//}
//
//// CreateUser mocks base method.
//func (m *MockAuthorization) CreateUser() (int, error) {
//	m.ctrl.T.Helper()
//	ret := m.ctrl.Call(m, "CreateUser", user)
//	ret0, _ := ret[0].(int)
//	ret1, _ := ret[1].(error)
//	return ret0, ret1
//}
