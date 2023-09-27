package repomocks

import (
	context "context"
	reflect "reflect"
	domain "webook/internal/domain"

	gomock "go.uber.org/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
        ctrl     *gomock.Controller
        recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
        mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
        mock := &MockUserRepository{ctrl: ctrl}
        mock.recorder = &MockUserRepositoryMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
        return m.recorder
}

// Create mocks base method.
func (m *MockUserRepository) Create(ctx context.Context, user domain.User) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Create", ctx, user)
        ret0, _ := ret[0].(error)
        return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserRepositoryMockRecorder) Create(ctx, user any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepository)(nil).Create), ctx, user)
}

// FindByEmail mocks base method.
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
        ret0, _ := ret[0].(domain.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserRepositoryMockRecorder) FindByEmail(ctx, email any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindByEmail), ctx, email)
}

// FindById mocks base method.
func (m *MockUserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "FindById", ctx, userId)
        ret0, _ := ret[0].(domain.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockUserRepositoryMockRecorder) FindById(ctx, userId any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockUserRepository)(nil).FindById), ctx, userId)
}

// FindByPhone mocks base method.
func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "FindByPhone", ctx, phone)
        ret0, _ := ret[0].(domain.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// FindByPhone indicates an expected call of FindByPhone.
func (mr *MockUserRepositoryMockRecorder) FindByPhone(ctx, phone any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPhone", reflect.TypeOf((*MockUserRepository)(nil).FindByPhone), ctx, phone)
}

// Update mocks base method.
func (m *MockUserRepository) Update(ctx context.Context, userId int64, nickname, description string, birthday int64) (domain.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Update", ctx, userId, nickname, description, birthday)
        ret0, _ := ret[0].(domain.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockUserRepositoryMockRecorder) Update(ctx, userId, nickname, description, birthday any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserRepository)(nil).Update), ctx, userId, nickname, description, birthday)
}