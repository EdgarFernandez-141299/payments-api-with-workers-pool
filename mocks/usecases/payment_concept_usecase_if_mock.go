// Code generated by mockery. DO NOT EDIT.

package usecases

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	request "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_concept/dto/request"
	response "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_concept/dto/response"
)

// PaymentConceptUsecaseIF is an autogenerated mock type for the PaymentConceptUsecaseIF type
type PaymentConceptUsecaseIF struct {
	mock.Mock
}

type PaymentConceptUsecaseIF_Expecter struct {
	mock *mock.Mock
}

func (_m *PaymentConceptUsecaseIF) EXPECT() *PaymentConceptUsecaseIF_Expecter {
	return &PaymentConceptUsecaseIF_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, payment, enterpriseId
func (_m *PaymentConceptUsecaseIF) Create(ctx context.Context, payment request.PaymentConceptRequest, enterpriseId string) (response.PaymentConceptResponse, error) {
	ret := _m.Called(ctx, payment, enterpriseId)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 response.PaymentConceptResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.PaymentConceptRequest, string) (response.PaymentConceptResponse, error)); ok {
		return rf(ctx, payment, enterpriseId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.PaymentConceptRequest, string) response.PaymentConceptResponse); ok {
		r0 = rf(ctx, payment, enterpriseId)
	} else {
		r0 = ret.Get(0).(response.PaymentConceptResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.PaymentConceptRequest, string) error); ok {
		r1 = rf(ctx, payment, enterpriseId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PaymentConceptUsecaseIF_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type PaymentConceptUsecaseIF_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - payment request.PaymentConceptRequest
//   - enterpriseId string
func (_e *PaymentConceptUsecaseIF_Expecter) Create(ctx interface{}, payment interface{}, enterpriseId interface{}) *PaymentConceptUsecaseIF_Create_Call {
	return &PaymentConceptUsecaseIF_Create_Call{Call: _e.mock.On("Create", ctx, payment, enterpriseId)}
}

func (_c *PaymentConceptUsecaseIF_Create_Call) Run(run func(ctx context.Context, payment request.PaymentConceptRequest, enterpriseId string)) *PaymentConceptUsecaseIF_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(request.PaymentConceptRequest), args[2].(string))
	})
	return _c
}

func (_c *PaymentConceptUsecaseIF_Create_Call) Return(_a0 response.PaymentConceptResponse, _a1 error) *PaymentConceptUsecaseIF_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PaymentConceptUsecaseIF_Create_Call) RunAndReturn(run func(context.Context, request.PaymentConceptRequest, string) (response.PaymentConceptResponse, error)) *PaymentConceptUsecaseIF_Create_Call {
	_c.Call.Return(run)
	return _c
}

// NewPaymentConceptUsecaseIF creates a new instance of PaymentConceptUsecaseIF. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPaymentConceptUsecaseIF(t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentConceptUsecaseIF {
	mock := &PaymentConceptUsecaseIF{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
