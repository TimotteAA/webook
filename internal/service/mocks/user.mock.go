package svcmocks

import (
	context "context"
	reflect "reflect"
	domain "webook/internal/domain"

	gomock "go.uber.org/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
        ctrl     *gomock.Controller
        recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
        mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
        mock := &MockUserService{ctrl: ctrl}
        mock.recorder = &MockUserServiceMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
        return m.recorder
}

// Edit mocks base method.
func (m *MockUserService) Edit(ctx context.Context, userId int64, nickname, description string, birthday int64) (domain.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Edit", ctx, userId, nickname, description, birthday)
        ret0, _ := ret[0].(domain.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// Edit indicates an expected call of Edit.
func (mr *MockUserServiceMockRecorder) Edit(ctx, userId, nickname, description, birthday any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockUserService)(nil).Edit), ctx, userId, nickname, description, birthday)
}

// FindOne mocks base method.
func (m *MockUserService) FindOne(ctx context.Context, userId int64) (domain.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "FindOne", ctx, userId)
        ret0, _ := ret[0].(domain.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// FindOne indicates an expected call of FindOne.
func (mr *MockUserServiceMockRecorder) FindOne(ctx, userId any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockUserService)(nil).FindOne), ctx, userId)
}

// FindOrCreate mocks base method.
func (m *MockUserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "FindOrCreate", ctx, phone)
        ret0, _ := ret[0].(domain.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// FindOrCreate indicates an expected call of FindOrCreate.
func (mr *MockUserServiceMockRecorder) FindOrCreate(ctx, phone any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrCreate", reflect.TypeOf((*MockUserService)(nil).FindOrCreate), ctx, phone)
}

// Login mocks base method.
func (m *MockUserService) Login(ctx context.Context, user domain.User) (domain.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Login", ctx, user)
        ret0, _ := ret[0].(domain.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserServiceMockRecorder) Login(ctx, user any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserService)(nil).Login), ctx, user)
}

// SignUp mocks base method.
func (m *MockUserService) SignUp(ctx context.Context, user domain.User) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "SignUp", ctx, user)
        ret0, _ := ret[0].(error)
        return ret0
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUserServiceMockRecorder) SignUp(ctx, user any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUserService)(nil).SignUp), ctx, user)
}