// Code generated by mockery. DO NOT EDIT.

package resources

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	request "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/request"

	response "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/response"

	utils "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
)

// DeunaPaymentResourceIF is an autogenerated mock type for the DeunaPaymentResourceIF type
type DeunaPaymentResourceIF struct {
	mock.Mock
}

type DeunaPaymentResourceIF_Expecter struct {
	mock *mock.Mock
}

func (_m *DeunaPaymentResourceIF) EXPECT() *DeunaPaymentResourceIF_Expecter {
	return &DeunaPaymentResourceIF_Expecter{mock: &_m.Mock}
}

// MakeOrderPayment provides a mock function with given fields: ctx, body, token
func (_m *DeunaPaymentResourceIF) MakeOrderPayment(ctx context.Context, body request.DeunaOrderPaymentRequest, token string) error {
	ret := _m.Called(ctx, body, token)

	if len(ret) == 0 {
		panic("no return value specified for MakeOrderPayment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, request.DeunaOrderPaymentRequest, string) error); ok {
		r0 = rf(ctx, body, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeunaPaymentResourceIF_MakeOrderPayment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MakeOrderPayment'
type DeunaPaymentResourceIF_MakeOrderPayment_Call struct {
	*mock.Call
}

// MakeOrderPayment is a helper method to define mock.On call
//   - ctx context.Context
//   - body request.DeunaOrderPaymentRequest
//   - token string
func (_e *DeunaPaymentResourceIF_Expecter) MakeOrderPayment(ctx interface{}, body interface{}, token interface{}) *DeunaPaymentResourceIF_MakeOrderPayment_Call {
	return &DeunaPaymentResourceIF_MakeOrderPayment_Call{Call: _e.mock.On("MakeOrderPayment", ctx, body, token)}
}

func (_c *DeunaPaymentResourceIF_MakeOrderPayment_Call) Run(run func(ctx context.Context, body request.DeunaOrderPaymentRequest, token string)) *DeunaPaymentResourceIF_MakeOrderPayment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(request.DeunaOrderPaymentRequest), args[2].(string))
	})
	return _c
}

func (_c *DeunaPaymentResourceIF_MakeOrderPayment_Call) Return(_a0 error) *DeunaPaymentResourceIF_MakeOrderPayment_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DeunaPaymentResourceIF_MakeOrderPayment_Call) RunAndReturn(run func(context.Context, request.DeunaOrderPaymentRequest, string) error) *DeunaPaymentResourceIF_MakeOrderPayment_Call {
	_c.Call.Return(run)
	return _c
}

// MakeOrderPaymentV2 provides a mock function with given fields: ctx, body
func (_m *DeunaPaymentResourceIF) MakeOrderPaymentV2(ctx context.Context, body request.DeunaOrderPaymentRequestV2) (response.DeunaOrderPaymentResponseV2, error) {
	ret := _m.Called(ctx, body)

	if len(ret) == 0 {
		panic("no return value specified for MakeOrderPaymentV2")
	}

	var r0 response.DeunaOrderPaymentResponseV2
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, request.DeunaOrderPaymentRequestV2) (response.DeunaOrderPaymentResponseV2, error)); ok {
		return rf(ctx, body)
	}
	if rf, ok := ret.Get(0).(func(context.Context, request.DeunaOrderPaymentRequestV2) response.DeunaOrderPaymentResponseV2); ok {
		r0 = rf(ctx, body)
	} else {
		r0 = ret.Get(0).(response.DeunaOrderPaymentResponseV2)
	}

	if rf, ok := ret.Get(1).(func(context.Context, request.DeunaOrderPaymentRequestV2) error); ok {
		r1 = rf(ctx, body)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeunaPaymentResourceIF_MakeOrderPaymentV2_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MakeOrderPaymentV2'
type DeunaPaymentResourceIF_MakeOrderPaymentV2_Call struct {
	*mock.Call
}

// MakeOrderPaymentV2 is a helper method to define mock.On call
//   - ctx context.Context
//   - body request.DeunaOrderPaymentRequestV2
func (_e *DeunaPaymentResourceIF_Expecter) MakeOrderPaymentV2(ctx interface{}, body interface{}) *DeunaPaymentResourceIF_MakeOrderPaymentV2_Call {
	return &DeunaPaymentResourceIF_MakeOrderPaymentV2_Call{Call: _e.mock.On("MakeOrderPaymentV2", ctx, body)}
}

func (_c *DeunaPaymentResourceIF_MakeOrderPaymentV2_Call) Run(run func(ctx context.Context, body request.DeunaOrderPaymentRequestV2)) *DeunaPaymentResourceIF_MakeOrderPaymentV2_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(request.DeunaOrderPaymentRequestV2))
	})
	return _c
}

func (_c *DeunaPaymentResourceIF_MakeOrderPaymentV2_Call) Return(_a0 response.DeunaOrderPaymentResponseV2, _a1 error) *DeunaPaymentResourceIF_MakeOrderPaymentV2_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DeunaPaymentResourceIF_MakeOrderPaymentV2_Call) RunAndReturn(run func(context.Context, request.DeunaOrderPaymentRequestV2) (response.DeunaOrderPaymentResponseV2, error)) *DeunaPaymentResourceIF_MakeOrderPaymentV2_Call {
	_c.Call.Return(run)
	return _c
}

// MakePartialRefund provides a mock function with given fields: ctx, body, orderToken
func (_m *DeunaPaymentResourceIF) MakePartialRefund(ctx context.Context, body utils.DeunaPartialRefundRequest, orderToken string) (response.DeunaRefundPaymentResponse, error) {
	ret := _m.Called(ctx, body, orderToken)

	if len(ret) == 0 {
		panic("no return value specified for MakePartialRefund")
	}

	var r0 response.DeunaRefundPaymentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, utils.DeunaPartialRefundRequest, string) (response.DeunaRefundPaymentResponse, error)); ok {
		return rf(ctx, body, orderToken)
	}
	if rf, ok := ret.Get(0).(func(context.Context, utils.DeunaPartialRefundRequest, string) response.DeunaRefundPaymentResponse); ok {
		r0 = rf(ctx, body, orderToken)
	} else {
		r0 = ret.Get(0).(response.DeunaRefundPaymentResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, utils.DeunaPartialRefundRequest, string) error); ok {
		r1 = rf(ctx, body, orderToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeunaPaymentResourceIF_MakePartialRefund_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MakePartialRefund'
type DeunaPaymentResourceIF_MakePartialRefund_Call struct {
	*mock.Call
}

// MakePartialRefund is a helper method to define mock.On call
//   - ctx context.Context
//   - body utils.DeunaPartialRefundRequest
//   - orderToken string
func (_e *DeunaPaymentResourceIF_Expecter) MakePartialRefund(ctx interface{}, body interface{}, orderToken interface{}) *DeunaPaymentResourceIF_MakePartialRefund_Call {
	return &DeunaPaymentResourceIF_MakePartialRefund_Call{Call: _e.mock.On("MakePartialRefund", ctx, body, orderToken)}
}

func (_c *DeunaPaymentResourceIF_MakePartialRefund_Call) Run(run func(ctx context.Context, body utils.DeunaPartialRefundRequest, orderToken string)) *DeunaPaymentResourceIF_MakePartialRefund_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(utils.DeunaPartialRefundRequest), args[2].(string))
	})
	return _c
}

func (_c *DeunaPaymentResourceIF_MakePartialRefund_Call) Return(_a0 response.DeunaRefundPaymentResponse, _a1 error) *DeunaPaymentResourceIF_MakePartialRefund_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DeunaPaymentResourceIF_MakePartialRefund_Call) RunAndReturn(run func(context.Context, utils.DeunaPartialRefundRequest, string) (response.DeunaRefundPaymentResponse, error)) *DeunaPaymentResourceIF_MakePartialRefund_Call {
	_c.Call.Return(run)
	return _c
}

// MakeTotalRefund provides a mock function with given fields: ctx, body, orderToken
func (_m *DeunaPaymentResourceIF) MakeTotalRefund(ctx context.Context, body utils.DeunaTotalRefundRequest, orderToken string) (response.DeunaRefundPaymentResponse, error) {
	ret := _m.Called(ctx, body, orderToken)

	if len(ret) == 0 {
		panic("no return value specified for MakeTotalRefund")
	}

	var r0 response.DeunaRefundPaymentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, utils.DeunaTotalRefundRequest, string) (response.DeunaRefundPaymentResponse, error)); ok {
		return rf(ctx, body, orderToken)
	}
	if rf, ok := ret.Get(0).(func(context.Context, utils.DeunaTotalRefundRequest, string) response.DeunaRefundPaymentResponse); ok {
		r0 = rf(ctx, body, orderToken)
	} else {
		r0 = ret.Get(0).(response.DeunaRefundPaymentResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, utils.DeunaTotalRefundRequest, string) error); ok {
		r1 = rf(ctx, body, orderToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeunaPaymentResourceIF_MakeTotalRefund_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MakeTotalRefund'
type DeunaPaymentResourceIF_MakeTotalRefund_Call struct {
	*mock.Call
}

// MakeTotalRefund is a helper method to define mock.On call
//   - ctx context.Context
//   - body utils.DeunaTotalRefundRequest
//   - orderToken string
func (_e *DeunaPaymentResourceIF_Expecter) MakeTotalRefund(ctx interface{}, body interface{}, orderToken interface{}) *DeunaPaymentResourceIF_MakeTotalRefund_Call {
	return &DeunaPaymentResourceIF_MakeTotalRefund_Call{Call: _e.mock.On("MakeTotalRefund", ctx, body, orderToken)}
}

func (_c *DeunaPaymentResourceIF_MakeTotalRefund_Call) Run(run func(ctx context.Context, body utils.DeunaTotalRefundRequest, orderToken string)) *DeunaPaymentResourceIF_MakeTotalRefund_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(utils.DeunaTotalRefundRequest), args[2].(string))
	})
	return _c
}

func (_c *DeunaPaymentResourceIF_MakeTotalRefund_Call) Return(_a0 response.DeunaRefundPaymentResponse, _a1 error) *DeunaPaymentResourceIF_MakeTotalRefund_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DeunaPaymentResourceIF_MakeTotalRefund_Call) RunAndReturn(run func(context.Context, utils.DeunaTotalRefundRequest, string) (response.DeunaRefundPaymentResponse, error)) *DeunaPaymentResourceIF_MakeTotalRefund_Call {
	_c.Call.Return(run)
	return _c
}

// NewDeunaPaymentResourceIF creates a new instance of DeunaPaymentResourceIF. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDeunaPaymentResourceIF(t interface {
	mock.TestingT
	Cleanup(func())
}) *DeunaPaymentResourceIF {
	mock := &DeunaPaymentResourceIF{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
