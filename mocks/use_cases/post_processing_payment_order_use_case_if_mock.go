// Code generated by mockery. DO NOT EDIT.

package use_cases

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	enums "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"

	use_cases "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/post_payment"
)

// PostProcessingPaymentOrderUseCaseIF is an autogenerated mock type for the PostProcessingPaymentOrderUseCaseIF type
type PostProcessingPaymentOrderUseCaseIF struct {
	mock.Mock
}

type PostProcessingPaymentOrderUseCaseIF_Expecter struct {
	mock *mock.Mock
}

func (_m *PostProcessingPaymentOrderUseCaseIF) EXPECT() *PostProcessingPaymentOrderUseCaseIF_Expecter {
	return &PostProcessingPaymentOrderUseCaseIF_Expecter{mock: &_m.Mock}
}

// PostProcessPaymentOrder provides a mock function with given fields: ctx, cmd
func (_m *PostProcessingPaymentOrderUseCaseIF) PostProcessPaymentOrder(ctx context.Context, cmd use_cases.PostProcessingPaymentOrderCommand) (enums.PaymentFlowEnum, error) {
	ret := _m.Called(ctx, cmd)

	if len(ret) == 0 {
		panic("no return value specified for PostProcessPaymentOrder")
	}

	var r0 enums.PaymentFlowEnum
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, use_cases.PostProcessingPaymentOrderCommand) (enums.PaymentFlowEnum, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, use_cases.PostProcessingPaymentOrderCommand) enums.PaymentFlowEnum); ok {
		r0 = rf(ctx, cmd)
	} else {
		r0 = ret.Get(0).(enums.PaymentFlowEnum)
	}

	if rf, ok := ret.Get(1).(func(context.Context, use_cases.PostProcessingPaymentOrderCommand) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PostProcessingPaymentOrderUseCaseIF_PostProcessPaymentOrder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PostProcessPaymentOrder'
type PostProcessingPaymentOrderUseCaseIF_PostProcessPaymentOrder_Call struct {
	*mock.Call
}

// PostProcessPaymentOrder is a helper method to define mock.On call
//   - ctx context.Context
//   - cmd use_cases.PostProcessingPaymentOrderCommand
func (_e *PostProcessingPaymentOrderUseCaseIF_Expecter) PostProcessPaymentOrder(ctx interface{}, cmd interface{}) *PostProcessingPaymentOrderUseCaseIF_PostProcessPaymentOrder_Call {
	return &PostProcessingPaymentOrderUseCaseIF_PostProcessPaymentOrder_Call{Call: _e.mock.On("PostProcessPaymentOrder", ctx, cmd)}
}

func (_c *PostProcessingPaymentOrderUseCaseIF_PostProcessPaymentOrder_Call) Run(run func(ctx context.Context, cmd use_cases.PostProcessingPaymentOrderCommand)) *PostProcessingPaymentOrderUseCaseIF_PostProcessPaymentOrder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(use_cases.PostProcessingPaymentOrderCommand))
	})
	return _c
}

func (_c *PostProcessingPaymentOrderUseCaseIF_PostProcessPaymentOrder_Call) Return(_a0 enums.PaymentFlowEnum, _a1 error) *PostProcessingPaymentOrderUseCaseIF_PostProcessPaymentOrder_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PostProcessingPaymentOrderUseCaseIF_PostProcessPaymentOrder_Call) RunAndReturn(run func(context.Context, use_cases.PostProcessingPaymentOrderCommand) (enums.PaymentFlowEnum, error)) *PostProcessingPaymentOrderUseCaseIF_PostProcessPaymentOrder_Call {
	_c.Call.Return(run)
	return _c
}

// NewPostProcessingPaymentOrderUseCaseIF creates a new instance of PostProcessingPaymentOrderUseCaseIF. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPostProcessingPaymentOrderUseCaseIF(t interface {
	mock.TestingT
	Cleanup(func())
}) *PostProcessingPaymentOrderUseCaseIF {
	mock := &PostProcessingPaymentOrderUseCaseIF{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
