// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/chickey/blog/internal/models"
)

// UsersLister is an autogenerated mock type for the usersLister type
type UsersLister struct {
	mock.Mock
}

type UsersLister_Expecter struct {
	mock *mock.Mock
}

func (_m *UsersLister) EXPECT() *UsersLister_Expecter {
	return &UsersLister_Expecter{mock: &_m.Mock}
}

// ListUsers provides a mock function with given fields: ctx, name
func (_m *UsersLister) ListUsers(ctx context.Context, name string) ([]models.User, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for ListUsers")
	}

	var r0 []models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]models.User, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []models.User); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UsersLister_ListUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListUsers'
type UsersLister_ListUsers_Call struct {
	*mock.Call
}

// ListUsers is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
func (_e *UsersLister_Expecter) ListUsers(ctx interface{}, name interface{}) *UsersLister_ListUsers_Call {
	return &UsersLister_ListUsers_Call{Call: _e.mock.On("ListUsers", ctx, name)}
}

func (_c *UsersLister_ListUsers_Call) Run(run func(ctx context.Context, name string)) *UsersLister_ListUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UsersLister_ListUsers_Call) Return(_a0 []models.User, _a1 error) *UsersLister_ListUsers_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UsersLister_ListUsers_Call) RunAndReturn(run func(context.Context, string) ([]models.User, error)) *UsersLister_ListUsers_Call {
	_c.Call.Return(run)
	return _c
}

// NewUsersLister creates a new instance of UsersLister. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUsersLister(t interface {
	mock.TestingT
	Cleanup(func())
}) *UsersLister {
	mock := &UsersLister{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
