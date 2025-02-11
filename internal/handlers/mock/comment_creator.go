// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/chickey/blog/internal/models"
)

// CommentCreator is an autogenerated mock type for the commentCreator type
type CommentCreator struct {
	mock.Mock
}

type CommentCreator_Expecter struct {
	mock *mock.Mock
}

func (_m *CommentCreator) EXPECT() *CommentCreator_Expecter {
	return &CommentCreator_Expecter{mock: &_m.Mock}
}

// CreateComment provides a mock function with given fields: ctx, comment
func (_m *CommentCreator) CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error) {
	ret := _m.Called(ctx, comment)

	if len(ret) == 0 {
		panic("no return value specified for CreateComment")
	}

	var r0 models.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Comment) (models.Comment, error)); ok {
		return rf(ctx, comment)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Comment) models.Comment); ok {
		r0 = rf(ctx, comment)
	} else {
		r0 = ret.Get(0).(models.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Comment) error); ok {
		r1 = rf(ctx, comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CommentCreator_CreateComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateComment'
type CommentCreator_CreateComment_Call struct {
	*mock.Call
}

// CreateComment is a helper method to define mock.On call
//   - ctx context.Context
//   - comment models.Comment
func (_e *CommentCreator_Expecter) CreateComment(ctx interface{}, comment interface{}) *CommentCreator_CreateComment_Call {
	return &CommentCreator_CreateComment_Call{Call: _e.mock.On("CreateComment", ctx, comment)}
}

func (_c *CommentCreator_CreateComment_Call) Run(run func(ctx context.Context, comment models.Comment)) *CommentCreator_CreateComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.Comment))
	})
	return _c
}

func (_c *CommentCreator_CreateComment_Call) Return(_a0 models.Comment, _a1 error) *CommentCreator_CreateComment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CommentCreator_CreateComment_Call) RunAndReturn(run func(context.Context, models.Comment) (models.Comment, error)) *CommentCreator_CreateComment_Call {
	_c.Call.Return(run)
	return _c
}

// NewCommentCreator creates a new instance of CommentCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCommentCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *CommentCreator {
	mock := &CommentCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
