// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/chickey/blog/internal/models"
)

// CommentsLister is an autogenerated mock type for the commentsLister type
type CommentsLister struct {
	mock.Mock
}

type CommentsLister_Expecter struct {
	mock *mock.Mock
}

func (_m *CommentsLister) EXPECT() *CommentsLister_Expecter {
	return &CommentsLister_Expecter{mock: &_m.Mock}
}

// ListComments provides a mock function with given fields: ctx, authorId, blogId
func (_m *CommentsLister) ListComments(ctx context.Context, authorId uint, blogId uint) ([]models.Comment, error) {
	ret := _m.Called(ctx, authorId, blogId)

	if len(ret) == 0 {
		panic("no return value specified for ListComments")
	}

	var r0 []models.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) ([]models.Comment, error)); ok {
		return rf(ctx, authorId, blogId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) []models.Comment); ok {
		r0 = rf(ctx, authorId, blogId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Comment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, uint) error); ok {
		r1 = rf(ctx, authorId, blogId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CommentsLister_ListComments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListComments'
type CommentsLister_ListComments_Call struct {
	*mock.Call
}

// ListComments is a helper method to define mock.On call
//   - ctx context.Context
//   - authorId uint
//   - blogId uint
func (_e *CommentsLister_Expecter) ListComments(ctx interface{}, authorId interface{}, blogId interface{}) *CommentsLister_ListComments_Call {
	return &CommentsLister_ListComments_Call{Call: _e.mock.On("ListComments", ctx, authorId, blogId)}
}

func (_c *CommentsLister_ListComments_Call) Run(run func(ctx context.Context, authorId uint, blogId uint)) *CommentsLister_ListComments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint), args[2].(uint))
	})
	return _c
}

func (_c *CommentsLister_ListComments_Call) Return(_a0 []models.Comment, _a1 error) *CommentsLister_ListComments_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CommentsLister_ListComments_Call) RunAndReturn(run func(context.Context, uint, uint) ([]models.Comment, error)) *CommentsLister_ListComments_Call {
	_c.Call.Return(run)
	return _c
}

// NewCommentsLister creates a new instance of CommentsLister. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCommentsLister(t interface {
	mock.TestingT
	Cleanup(func())
}) *CommentsLister {
	mock := &CommentsLister{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
