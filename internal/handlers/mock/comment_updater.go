// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/chickey/blog/internal/models"
)

// CommentUpdater is an autogenerated mock type for the commentUpdater type
type CommentUpdater struct {
	mock.Mock
}

type CommentUpdater_Expecter struct {
	mock *mock.Mock
}

func (_m *CommentUpdater) EXPECT() *CommentUpdater_Expecter {
	return &CommentUpdater_Expecter{mock: &_m.Mock}
}

// UpdateComment provides a mock function with given fields: ctx, patch
func (_m *CommentUpdater) UpdateComment(ctx context.Context, patch models.Comment) (models.Comment, error) {
	ret := _m.Called(ctx, patch)

	if len(ret) == 0 {
		panic("no return value specified for UpdateComment")
	}

	var r0 models.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Comment) (models.Comment, error)); ok {
		return rf(ctx, patch)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Comment) models.Comment); ok {
		r0 = rf(ctx, patch)
	} else {
		r0 = ret.Get(0).(models.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Comment) error); ok {
		r1 = rf(ctx, patch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CommentUpdater_UpdateComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateComment'
type CommentUpdater_UpdateComment_Call struct {
	*mock.Call
}

// UpdateComment is a helper method to define mock.On call
//   - ctx context.Context
//   - patch models.Comment
func (_e *CommentUpdater_Expecter) UpdateComment(ctx interface{}, patch interface{}) *CommentUpdater_UpdateComment_Call {
	return &CommentUpdater_UpdateComment_Call{Call: _e.mock.On("UpdateComment", ctx, patch)}
}

func (_c *CommentUpdater_UpdateComment_Call) Run(run func(ctx context.Context, patch models.Comment)) *CommentUpdater_UpdateComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.Comment))
	})
	return _c
}

func (_c *CommentUpdater_UpdateComment_Call) Return(_a0 models.Comment, _a1 error) *CommentUpdater_UpdateComment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CommentUpdater_UpdateComment_Call) RunAndReturn(run func(context.Context, models.Comment) (models.Comment, error)) *CommentUpdater_UpdateComment_Call {
	_c.Call.Return(run)
	return _c
}

// NewCommentUpdater creates a new instance of CommentUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCommentUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *CommentUpdater {
	mock := &CommentUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
