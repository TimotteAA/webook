package entitymocks

import (
	context "context"
	reflect "reflect"
	entity "webook/internal/repository/entity"

	gomock "go.uber.org/mock/gomock"
)

// MockUserEntity is a mock of UserEntity interface.
type MockUserEntity struct {
        ctrl     *gomock.Controller
        recorder *MockUserEntityMockRecorder
}

// MockUserEntityMockRecorder is the mock recorder for MockUserEntity.
type MockUserEntityMockRecorder struct {
        mock *MockUserEntity
}

// NewMockUserEntity creates a new mock instance.
func NewMockUserEntity(ctrl *gomock.Controller) *MockUserEntity {
        mock := &MockUserEntity{ctrl: ctrl}
        mock.recorder = &MockUserEntityMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserEntity) EXPECT() *MockUserEntityMockRecorder {
        return m.recorder
}

// Create mocks base method.
func (m *MockUserEntity) Create(ctx context.Context, u entity.User) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Create", ctx, u)
        ret0, _ := ret[0].(error)
        return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserEntityMockRecorder) Create(ctx, u any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserEntity)(nil).Create), ctx, u)
}

// FindByEmail mocks base method.
func (m *MockUserEntity) FindByEmail(ctx context.Context, email string) (entity.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
        ret0, _ := ret[0].(entity.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserEntityMockRecorder) FindByEmail(ctx, email any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserEntity)(nil).FindByEmail), ctx, email)
}

// FindById mocks base method.
func (m *MockUserEntity) FindById(ctx context.Context, userId int64) (entity.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "FindById", ctx, userId)
        ret0, _ := ret[0].(entity.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockUserEntityMockRecorder) FindById(ctx, userId any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockUserEntity)(nil).FindById), ctx, userId)
}

// FindByPhone mocks base method.
func (m *MockUserEntity) FindByPhone(ctx context.Context, phone string) (entity.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "FindByPhone", ctx, phone)
        ret0, _ := ret[0].(entity.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// FindByPhone indicates an expected call of FindByPhone.
func (mr *MockUserEntityMockRecorder) FindByPhone(ctx, phone any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPhone", reflect.TypeOf((*MockUserEntity)(nil).FindByPhone), ctx, phone)
}

// Update mocks base method.
func (m *MockUserEntity) Update(ctx context.Context, userId int64, nickname, description string, birthday int64) (entity.User, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Update", ctx, userId, nickname, description, birthday)
        ret0, _ := ret[0].(entity.User)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockUserEntityMockRecorder) Update(ctx, userId, nickname, description, birthday any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserEntity)(nil).Update), ctx, userId, nickname, description, birthday)
}