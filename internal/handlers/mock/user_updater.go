// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/chickey/blog/internal/models"
)

// UserUpdater is an autogenerated mock type for the userUpdater type
type UserUpdater struct {
	mock.Mock
}

type UserUpdater_Expecter struct {
	mock *mock.Mock
}

func (_m *UserUpdater) EXPECT() *UserUpdater_Expecter {
	return &UserUpdater_Expecter{mock: &_m.Mock}
}

// UpdateUser provides a mock function with given fields: ctx, id, patch
func (_m *UserUpdater) UpdateUser(ctx context.Context, id uint64, patch models.User) (models.User, error) {
	ret := _m.Called(ctx, id, patch)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, models.User) (models.User, error)); ok {
		return rf(ctx, id, patch)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, models.User) models.User); ok {
		r0 = rf(ctx, id, patch)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, models.User) error); ok {
		r1 = rf(ctx, id, patch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserUpdater_UpdateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUser'
type UserUpdater_UpdateUser_Call struct {
	*mock.Call
}

// UpdateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - id uint64
//   - patch models.User
func (_e *UserUpdater_Expecter) UpdateUser(ctx interface{}, id interface{}, patch interface{}) *UserUpdater_UpdateUser_Call {
	return &UserUpdater_UpdateUser_Call{Call: _e.mock.On("UpdateUser", ctx, id, patch)}
}

func (_c *UserUpdater_UpdateUser_Call) Run(run func(ctx context.Context, id uint64, patch models.User)) *UserUpdater_UpdateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(models.User))
	})
	return _c
}

func (_c *UserUpdater_UpdateUser_Call) Return(_a0 models.User, _a1 error) *UserUpdater_UpdateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserUpdater_UpdateUser_Call) RunAndReturn(run func(context.Context, uint64, models.User) (models.User, error)) *UserUpdater_UpdateUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserUpdater creates a new instance of UserUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserUpdater {
	mock := &UserUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
