// Code generated by mockery v2.26.1. DO NOT EDIT.

package github

import (
	context "context"

	v51github "github.com/google/go-github/v57/github"
	mock "github.com/stretchr/testify/mock"
)

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

type MockClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockClient) EXPECT() *MockClient_Expecter {
	return &MockClient_Expecter{mock: &_m.Mock}
}

// ArchiveIssue provides a mock function with given fields: ctx, repoOwner, repoName, issueNumber, title
func (_m *MockClient) ArchiveIssue(ctx context.Context, repoOwner string, repoName string, issueNumber int, title string) (*v51github.Issue, error) {
	ret := _m.Called(ctx, repoOwner, repoName, issueNumber, title)

	var r0 *v51github.Issue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, string) (*v51github.Issue, error)); ok {
		return rf(ctx, repoOwner, repoName, issueNumber, title)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, string) *v51github.Issue); ok {
		r0 = rf(ctx, repoOwner, repoName, issueNumber, title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v51github.Issue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, string) error); ok {
		r1 = rf(ctx, repoOwner, repoName, issueNumber, title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_ArchiveIssue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ArchiveIssue'
type MockClient_ArchiveIssue_Call struct {
	*mock.Call
}

// ArchiveIssue is a helper method to define mock.On call
//   - ctx context.Context
//   - repoOwner string
//   - repoName string
//   - issueNumber int
//   - title string
func (_e *MockClient_Expecter) ArchiveIssue(ctx interface{}, repoOwner interface{}, repoName interface{}, issueNumber interface{}, title interface{}) *MockClient_ArchiveIssue_Call {
	return &MockClient_ArchiveIssue_Call{Call: _e.mock.On("ArchiveIssue", ctx, repoOwner, repoName, issueNumber, title)}
}

func (_c *MockClient_ArchiveIssue_Call) Run(run func(ctx context.Context, repoOwner string, repoName string, issueNumber int, title string)) *MockClient_ArchiveIssue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int), args[4].(string))
	})
	return _c
}

func (_c *MockClient_ArchiveIssue_Call) Return(_a0 *v51github.Issue, _a1 error) *MockClient_ArchiveIssue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_ArchiveIssue_Call) RunAndReturn(run func(context.Context, string, string, int, string) (*v51github.Issue, error)) *MockClient_ArchiveIssue_Call {
	_c.Call.Return(run)
	return _c
}

// CloseIssue provides a mock function with given fields: ctx, repoOwner, repoName, issueNumber
func (_m *MockClient) CloseIssue(ctx context.Context, repoOwner string, repoName string, issueNumber int) (*v51github.Issue, error) {
	ret := _m.Called(ctx, repoOwner, repoName, issueNumber)

	var r0 *v51github.Issue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) (*v51github.Issue, error)); ok {
		return rf(ctx, repoOwner, repoName, issueNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) *v51github.Issue); ok {
		r0 = rf(ctx, repoOwner, repoName, issueNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v51github.Issue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int) error); ok {
		r1 = rf(ctx, repoOwner, repoName, issueNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_CloseIssue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CloseIssue'
type MockClient_CloseIssue_Call struct {
	*mock.Call
}

// CloseIssue is a helper method to define mock.On call
//   - ctx context.Context
//   - repoOwner string
//   - repoName string
//   - issueNumber int
func (_e *MockClient_Expecter) CloseIssue(ctx interface{}, repoOwner interface{}, repoName interface{}, issueNumber interface{}) *MockClient_CloseIssue_Call {
	return &MockClient_CloseIssue_Call{Call: _e.mock.On("CloseIssue", ctx, repoOwner, repoName, issueNumber)}
}

func (_c *MockClient_CloseIssue_Call) Run(run func(ctx context.Context, repoOwner string, repoName string, issueNumber int)) *MockClient_CloseIssue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int))
	})
	return _c
}

func (_c *MockClient_CloseIssue_Call) Return(_a0 *v51github.Issue, _a1 error) *MockClient_CloseIssue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_CloseIssue_Call) RunAndReturn(run func(context.Context, string, string, int) (*v51github.Issue, error)) *MockClient_CloseIssue_Call {
	_c.Call.Return(run)
	return _c
}

// CreateIssue provides a mock function with given fields: ctx, repoOwner, repoName, param
func (_m *MockClient) CreateIssue(ctx context.Context, repoOwner string, repoName string, param *v51github.IssueRequest) (*v51github.Issue, error) {
	ret := _m.Called(ctx, repoOwner, repoName, param)

	var r0 *v51github.Issue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *v51github.IssueRequest) (*v51github.Issue, error)); ok {
		return rf(ctx, repoOwner, repoName, param)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *v51github.IssueRequest) *v51github.Issue); ok {
		r0 = rf(ctx, repoOwner, repoName, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v51github.Issue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *v51github.IssueRequest) error); ok {
		r1 = rf(ctx, repoOwner, repoName, param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_CreateIssue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateIssue'
type MockClient_CreateIssue_Call struct {
	*mock.Call
}

// CreateIssue is a helper method to define mock.On call
//   - ctx context.Context
//   - repoOwner string
//   - repoName string
//   - param *v51github.IssueRequest
func (_e *MockClient_Expecter) CreateIssue(ctx interface{}, repoOwner interface{}, repoName interface{}, param interface{}) *MockClient_CreateIssue_Call {
	return &MockClient_CreateIssue_Call{Call: _e.mock.On("CreateIssue", ctx, repoOwner, repoName, param)}
}

func (_c *MockClient_CreateIssue_Call) Run(run func(ctx context.Context, repoOwner string, repoName string, param *v51github.IssueRequest)) *MockClient_CreateIssue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*v51github.IssueRequest))
	})
	return _c
}

func (_c *MockClient_CreateIssue_Call) Return(_a0 *v51github.Issue, _a1 error) *MockClient_CreateIssue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_CreateIssue_Call) RunAndReturn(run func(context.Context, string, string, *v51github.IssueRequest) (*v51github.Issue, error)) *MockClient_CreateIssue_Call {
	_c.Call.Return(run)
	return _c
}

// GetIssue provides a mock function with given fields: ctx, repoOwner, repoName, title
func (_m *MockClient) GetIssue(ctx context.Context, repoOwner string, repoName string, title string) (*Issue, error) {
	ret := _m.Called(ctx, repoOwner, repoName, title)

	var r0 *Issue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) (*Issue, error)); ok {
		return rf(ctx, repoOwner, repoName, title)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *Issue); ok {
		r0 = rf(ctx, repoOwner, repoName, title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Issue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, repoOwner, repoName, title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_GetIssue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIssue'
type MockClient_GetIssue_Call struct {
	*mock.Call
}

// GetIssue is a helper method to define mock.On call
//   - ctx context.Context
//   - repoOwner string
//   - repoName string
//   - title string
func (_e *MockClient_Expecter) GetIssue(ctx interface{}, repoOwner interface{}, repoName interface{}, title interface{}) *MockClient_GetIssue_Call {
	return &MockClient_GetIssue_Call{Call: _e.mock.On("GetIssue", ctx, repoOwner, repoName, title)}
}

func (_c *MockClient_GetIssue_Call) Run(run func(ctx context.Context, repoOwner string, repoName string, title string)) *MockClient_GetIssue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockClient_GetIssue_Call) Return(_a0 *Issue, _a1 error) *MockClient_GetIssue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_GetIssue_Call) RunAndReturn(run func(context.Context, string, string, string) (*Issue, error)) *MockClient_GetIssue_Call {
	_c.Call.Return(run)
	return _c
}

// ListIssues provides a mock function with given fields: ctx, repoOwner, repoName
func (_m *MockClient) ListIssues(ctx context.Context, repoOwner string, repoName string) ([]*Issue, error) {
	ret := _m.Called(ctx, repoOwner, repoName)

	var r0 []*Issue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]*Issue, error)); ok {
		return rf(ctx, repoOwner, repoName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*Issue); ok {
		r0 = rf(ctx, repoOwner, repoName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Issue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, repoOwner, repoName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_ListIssues_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListIssues'
type MockClient_ListIssues_Call struct {
	*mock.Call
}

// ListIssues is a helper method to define mock.On call
//   - ctx context.Context
//   - repoOwner string
//   - repoName string
func (_e *MockClient_Expecter) ListIssues(ctx interface{}, repoOwner interface{}, repoName interface{}) *MockClient_ListIssues_Call {
	return &MockClient_ListIssues_Call{Call: _e.mock.On("ListIssues", ctx, repoOwner, repoName)}
}

func (_c *MockClient_ListIssues_Call) Run(run func(ctx context.Context, repoOwner string, repoName string)) *MockClient_ListIssues_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockClient_ListIssues_Call) Return(_a0 []*Issue, _a1 error) *MockClient_ListIssues_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_ListIssues_Call) RunAndReturn(run func(context.Context, string, string) ([]*Issue, error)) *MockClient_ListIssues_Call {
	_c.Call.Return(run)
	return _c
}

// ListLeastRecentlyUpdatedIssues provides a mock function with given fields: ctx, repoOwner, repoName, numOfIssues, deadline
func (_m *MockClient) ListLeastRecentlyUpdatedIssues(ctx context.Context, repoOwner string, repoName string, numOfIssues int, deadline string) ([]*Issue, error) {
	ret := _m.Called(ctx, repoOwner, repoName, numOfIssues, deadline)

	var r0 []*Issue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, string) ([]*Issue, error)); ok {
		return rf(ctx, repoOwner, repoName, numOfIssues, deadline)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, string) []*Issue); ok {
		r0 = rf(ctx, repoOwner, repoName, numOfIssues, deadline)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Issue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, string) error); ok {
		r1 = rf(ctx, repoOwner, repoName, numOfIssues, deadline)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_ListLeastRecentlyUpdatedIssues_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListLeastRecentlyUpdatedIssues'
type MockClient_ListLeastRecentlyUpdatedIssues_Call struct {
	*mock.Call
}

// ListLeastRecentlyUpdatedIssues is a helper method to define mock.On call
//   - ctx context.Context
//   - repoOwner string
//   - repoName string
//   - numOfIssues int
//   - deadline string
func (_e *MockClient_Expecter) ListLeastRecentlyUpdatedIssues(ctx interface{}, repoOwner interface{}, repoName interface{}, numOfIssues interface{}, deadline interface{}) *MockClient_ListLeastRecentlyUpdatedIssues_Call {
	return &MockClient_ListLeastRecentlyUpdatedIssues_Call{Call: _e.mock.On("ListLeastRecentlyUpdatedIssues", ctx, repoOwner, repoName, numOfIssues, deadline)}
}

func (_c *MockClient_ListLeastRecentlyUpdatedIssues_Call) Run(run func(ctx context.Context, repoOwner string, repoName string, numOfIssues int, deadline string)) *MockClient_ListLeastRecentlyUpdatedIssues_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int), args[4].(string))
	})
	return _c
}

func (_c *MockClient_ListLeastRecentlyUpdatedIssues_Call) Return(_a0 []*Issue, _a1 error) *MockClient_ListLeastRecentlyUpdatedIssues_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_ListLeastRecentlyUpdatedIssues_Call) RunAndReturn(run func(context.Context, string, string, int, string) ([]*Issue, error)) *MockClient_ListLeastRecentlyUpdatedIssues_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockClient creates a new instance of MockClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockClient(t mockConstructorTestingTNewMockClient) *MockClient {
	mock := &MockClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
