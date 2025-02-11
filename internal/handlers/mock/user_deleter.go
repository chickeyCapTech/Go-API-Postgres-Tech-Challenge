// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// UserDeleter is an autogenerated mock type for the userDeleter type
type UserDeleter struct {
	mock.Mock
}

type UserDeleter_Expecter struct {
	mock *mock.Mock
}

func (_m *UserDeleter) EXPECT() *UserDeleter_Expecter {
	return &UserDeleter_Expecter{mock: &_m.Mock}
}

// DeleteUser provides a mock function with given fields: ctx, id
func (_m *UserDeleter) DeleteUser(ctx context.Context, id uint64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserDeleter_DeleteUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteUser'
type UserDeleter_DeleteUser_Call struct {
	*mock.Call
}

// DeleteUser is a helper method to define mock.On call
//   - ctx context.Context
//   - id uint64
func (_e *UserDeleter_Expecter) DeleteUser(ctx interface{}, id interface{}) *UserDeleter_DeleteUser_Call {
	return &UserDeleter_DeleteUser_Call{Call: _e.mock.On("DeleteUser", ctx, id)}
}

func (_c *UserDeleter_DeleteUser_Call) Run(run func(ctx context.Context, id uint64)) *UserDeleter_DeleteUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *UserDeleter_DeleteUser_Call) Return(_a0 error) *UserDeleter_DeleteUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserDeleter_DeleteUser_Call) RunAndReturn(run func(context.Context, uint64) error) *UserDeleter_DeleteUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserDeleter creates a new instance of UserDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserDeleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserDeleter {
	mock := &UserDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
