// Code generated by mockery. DO NOT EDIT.

package repository

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// DeunaOrderRepository is an autogenerated mock type for the DeunaOrderRepository type
type DeunaOrderRepository struct {
	mock.Mock
}

type DeunaOrderRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *DeunaOrderRepository) EXPECT() *DeunaOrderRepository_Expecter {
	return &DeunaOrderRepository_Expecter{mock: &_m.Mock}
}

// CreatePaymentOrderDeuna provides a mock function with given fields: ctx, paymentID, orderID, deunaOrderToken
func (_m *DeunaOrderRepository) CreatePaymentOrderDeuna(ctx context.Context, paymentID string, orderID string, deunaOrderToken string) error {
	ret := _m.Called(ctx, paymentID, orderID, deunaOrderToken)

	if len(ret) == 0 {
		panic("no return value specified for CreatePaymentOrderDeuna")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, paymentID, orderID, deunaOrderToken)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeunaOrderRepository_CreatePaymentOrderDeuna_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePaymentOrderDeuna'
type DeunaOrderRepository_CreatePaymentOrderDeuna_Call struct {
	*mock.Call
}

// CreatePaymentOrderDeuna is a helper method to define mock.On call
//   - ctx context.Context
//   - paymentID string
//   - orderID string
//   - deunaOrderToken string
func (_e *DeunaOrderRepository_Expecter) CreatePaymentOrderDeuna(ctx interface{}, paymentID interface{}, orderID interface{}, deunaOrderToken interface{}) *DeunaOrderRepository_CreatePaymentOrderDeuna_Call {
	return &DeunaOrderRepository_CreatePaymentOrderDeuna_Call{Call: _e.mock.On("CreatePaymentOrderDeuna", ctx, paymentID, orderID, deunaOrderToken)}
}

func (_c *DeunaOrderRepository_CreatePaymentOrderDeuna_Call) Run(run func(ctx context.Context, paymentID string, orderID string, deunaOrderToken string)) *DeunaOrderRepository_CreatePaymentOrderDeuna_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *DeunaOrderRepository_CreatePaymentOrderDeuna_Call) Return(_a0 error) *DeunaOrderRepository_CreatePaymentOrderDeuna_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DeunaOrderRepository_CreatePaymentOrderDeuna_Call) RunAndReturn(run func(context.Context, string, string, string) error) *DeunaOrderRepository_CreatePaymentOrderDeuna_Call {
	_c.Call.Return(run)
	return _c
}

// GetTokenByOrderAndPaymentID provides a mock function with given fields: ctx, orderID, paymentID
func (_m *DeunaOrderRepository) GetTokenByOrderAndPaymentID(ctx context.Context, orderID string, paymentID string) (string, error) {
	ret := _m.Called(ctx, orderID, paymentID)

	if len(ret) == 0 {
		panic("no return value specified for GetTokenByOrderAndPaymentID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, orderID, paymentID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, orderID, paymentID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, orderID, paymentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeunaOrderRepository_GetTokenByOrderAndPaymentID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTokenByOrderAndPaymentID'
type DeunaOrderRepository_GetTokenByOrderAndPaymentID_Call struct {
	*mock.Call
}

// GetTokenByOrderAndPaymentID is a helper method to define mock.On call
//   - ctx context.Context
//   - orderID string
//   - paymentID string
func (_e *DeunaOrderRepository_Expecter) GetTokenByOrderAndPaymentID(ctx interface{}, orderID interface{}, paymentID interface{}) *DeunaOrderRepository_GetTokenByOrderAndPaymentID_Call {
	return &DeunaOrderRepository_GetTokenByOrderAndPaymentID_Call{Call: _e.mock.On("GetTokenByOrderAndPaymentID", ctx, orderID, paymentID)}
}

func (_c *DeunaOrderRepository_GetTokenByOrderAndPaymentID_Call) Run(run func(ctx context.Context, orderID string, paymentID string)) *DeunaOrderRepository_GetTokenByOrderAndPaymentID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *DeunaOrderRepository_GetTokenByOrderAndPaymentID_Call) Return(token string, err error) *DeunaOrderRepository_GetTokenByOrderAndPaymentID_Call {
	_c.Call.Return(token, err)
	return _c
}

func (_c *DeunaOrderRepository_GetTokenByOrderAndPaymentID_Call) RunAndReturn(run func(context.Context, string, string) (string, error)) *DeunaOrderRepository_GetTokenByOrderAndPaymentID_Call {
	_c.Call.Return(run)
	return _c
}

// NewDeunaOrderRepository creates a new instance of DeunaOrderRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDeunaOrderRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *DeunaOrderRepository {
	mock := &DeunaOrderRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
