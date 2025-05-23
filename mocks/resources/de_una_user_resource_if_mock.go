// Code generated by mockery. DO NOT EDIT.

package resources

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	request "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/request"

	response "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
)

// DeUnaUserResourceIF is an autogenerated mock type for the DeUnaUserResourceIF type
type DeUnaUserResourceIF struct {
	mock.Mock
}

type DeUnaUserResourceIF_Expecter struct {
	mock *mock.Mock
}

func (_m *DeUnaUserResourceIF) EXPECT() *DeUnaUserResourceIF_Expecter {
	return &DeUnaUserResourceIF_Expecter{mock: &_m.Mock}
}

// CreateUser provides a mock function with given fields: ctx, _a1
func (_m *DeUnaUserResourceIF) CreateUser(ctx context.Context, _a1 request.CreateUserRequestDTO) (response.CreatedUserResponse, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 response.CreatedUserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.CreateUserRequestDTO) (response.CreatedUserResponse, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.CreateUserRequestDTO) response.CreatedUserResponse); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(response.CreatedUserResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.CreateUserRequestDTO) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeUnaUserResourceIF_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type DeUnaUserResourceIF_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 request.CreateUserRequestDTO
func (_e *DeUnaUserResourceIF_Expecter) CreateUser(ctx interface{}, _a1 interface{}) *DeUnaUserResourceIF_CreateUser_Call {
	return &DeUnaUserResourceIF_CreateUser_Call{Call: _e.mock.On("CreateUser", ctx, _a1)}
}

func (_c *DeUnaUserResourceIF_CreateUser_Call) Run(run func(ctx context.Context, _a1 request.CreateUserRequestDTO)) *DeUnaUserResourceIF_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(request.CreateUserRequestDTO))
	})
	return _c
}

func (_c *DeUnaUserResourceIF_CreateUser_Call) Return(_a0 response.CreatedUserResponse, _a1 error) *DeUnaUserResourceIF_CreateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DeUnaUserResourceIF_CreateUser_Call) RunAndReturn(run func(context.Context, request.CreateUserRequestDTO) (response.CreatedUserResponse, error)) *DeUnaUserResourceIF_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewDeUnaUserResourceIF creates a new instance of DeUnaUserResourceIF. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDeUnaUserResourceIF(t interface {
	mock.TestingT
	Cleanup(func())
}) *DeUnaUserResourceIF {
	mock := &DeUnaUserResourceIF{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
