// Code generated by mockery. DO NOT EDIT.

package usecases

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	request "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account_route/dto/request"
	response "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account_route/dto/response"
)

// CollectionAccountRouteUsecaseIF is an autogenerated mock type for the CollectionAccountRouteUsecaseIF type
type CollectionAccountRouteUsecaseIF struct {
	mock.Mock
}

type CollectionAccountRouteUsecaseIF_Expecter struct {
	mock *mock.Mock
}

func (_m *CollectionAccountRouteUsecaseIF) EXPECT() *CollectionAccountRouteUsecaseIF_Expecter {
	return &CollectionAccountRouteUsecaseIF_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, collection, enterpriseId
func (_m *CollectionAccountRouteUsecaseIF) Create(ctx context.Context, collection request.CollectionAccountRouteRequest, enterpriseId string) (response.CollectionAccountRouteResponse, error) {
	ret := _m.Called(ctx, collection, enterpriseId)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 response.CollectionAccountRouteResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.CollectionAccountRouteRequest, string) (response.CollectionAccountRouteResponse, error)); ok {
		return rf(ctx, collection, enterpriseId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.CollectionAccountRouteRequest, string) response.CollectionAccountRouteResponse); ok {
		r0 = rf(ctx, collection, enterpriseId)
	} else {
		r0 = ret.Get(0).(response.CollectionAccountRouteResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.CollectionAccountRouteRequest, string) error); ok {
		r1 = rf(ctx, collection, enterpriseId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionAccountRouteUsecaseIF_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type CollectionAccountRouteUsecaseIF_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - collection request.CollectionAccountRouteRequest
//   - enterpriseId string
func (_e *CollectionAccountRouteUsecaseIF_Expecter) Create(ctx interface{}, collection interface{}, enterpriseId interface{}) *CollectionAccountRouteUsecaseIF_Create_Call {
	return &CollectionAccountRouteUsecaseIF_Create_Call{Call: _e.mock.On("Create", ctx, collection, enterpriseId)}
}

func (_c *CollectionAccountRouteUsecaseIF_Create_Call) Run(run func(ctx context.Context, collection request.CollectionAccountRouteRequest, enterpriseId string)) *CollectionAccountRouteUsecaseIF_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(request.CollectionAccountRouteRequest), args[2].(string))
	})
	return _c
}

func (_c *CollectionAccountRouteUsecaseIF_Create_Call) Return(_a0 response.CollectionAccountRouteResponse, _a1 error) *CollectionAccountRouteUsecaseIF_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectionAccountRouteUsecaseIF_Create_Call) RunAndReturn(run func(context.Context, request.CollectionAccountRouteRequest, string) (response.CollectionAccountRouteResponse, error)) *CollectionAccountRouteUsecaseIF_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Disable provides a mock function with given fields: ctx, id, enterpriseId
func (_m *CollectionAccountRouteUsecaseIF) Disable(ctx context.Context, id string, enterpriseId string) (response.CollectionAccountRouteDisableResponse, error) {
	ret := _m.Called(ctx, id, enterpriseId)

	if len(ret) == 0 {
		panic("no return value specified for Disable")
	}

	var r0 response.CollectionAccountRouteDisableResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (response.CollectionAccountRouteDisableResponse, error)); ok {
		return rf(ctx, id, enterpriseId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) response.CollectionAccountRouteDisableResponse); ok {
		r0 = rf(ctx, id, enterpriseId)
	} else {
		r0 = ret.Get(0).(response.CollectionAccountRouteDisableResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, id, enterpriseId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionAccountRouteUsecaseIF_Disable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Disable'
type CollectionAccountRouteUsecaseIF_Disable_Call struct {
	*mock.Call
}

// Disable is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - enterpriseId string
func (_e *CollectionAccountRouteUsecaseIF_Expecter) Disable(ctx interface{}, id interface{}, enterpriseId interface{}) *CollectionAccountRouteUsecaseIF_Disable_Call {
	return &CollectionAccountRouteUsecaseIF_Disable_Call{Call: _e.mock.On("Disable", ctx, id, enterpriseId)}
}

func (_c *CollectionAccountRouteUsecaseIF_Disable_Call) Run(run func(ctx context.Context, id string, enterpriseId string)) *CollectionAccountRouteUsecaseIF_Disable_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *CollectionAccountRouteUsecaseIF_Disable_Call) Return(_a0 response.CollectionAccountRouteDisableResponse, _a1 error) *CollectionAccountRouteUsecaseIF_Disable_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectionAccountRouteUsecaseIF_Disable_Call) RunAndReturn(run func(context.Context, string, string) (response.CollectionAccountRouteDisableResponse, error)) *CollectionAccountRouteUsecaseIF_Disable_Call {
	_c.Call.Return(run)
	return _c
}

// NewCollectionAccountRouteUsecaseIF creates a new instance of CollectionAccountRouteUsecaseIF. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCollectionAccountRouteUsecaseIF(t interface {
	mock.TestingT
	Cleanup(func())
}) *CollectionAccountRouteUsecaseIF {
	mock := &CollectionAccountRouteUsecaseIF{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
