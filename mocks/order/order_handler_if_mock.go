// Code generated by mockery. DO NOT EDIT.

package order

import (
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"
)

// OrderHandlerIF is an autogenerated mock type for the OrderHandlerIF type
type OrderHandlerIF struct {
	mock.Mock
}

type OrderHandlerIF_Expecter struct {
	mock *mock.Mock
}

func (_m *OrderHandlerIF) EXPECT() *OrderHandlerIF_Expecter {
	return &OrderHandlerIF_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: context
func (_m *OrderHandlerIF) Create(context echo.Context) error {
	ret := _m.Called(context)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(context)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrderHandlerIF_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type OrderHandlerIF_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - context echo.Context
func (_e *OrderHandlerIF_Expecter) Create(context interface{}) *OrderHandlerIF_Create_Call {
	return &OrderHandlerIF_Create_Call{Call: _e.mock.On("Create", context)}
}

func (_c *OrderHandlerIF_Create_Call) Run(run func(context echo.Context)) *OrderHandlerIF_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(echo.Context))
	})
	return _c
}

func (_c *OrderHandlerIF_Create_Call) Return(_a0 error) *OrderHandlerIF_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrderHandlerIF_Create_Call) RunAndReturn(run func(echo.Context) error) *OrderHandlerIF_Create_Call {
	_c.Call.Return(run)
	return _c
}

// CreatePaymentOrder provides a mock function with given fields: context
func (_m *OrderHandlerIF) CreatePaymentOrder(context echo.Context) error {
	ret := _m.Called(context)

	if len(ret) == 0 {
		panic("no return value specified for CreatePaymentOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(context)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrderHandlerIF_CreatePaymentOrder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePaymentOrder'
type OrderHandlerIF_CreatePaymentOrder_Call struct {
	*mock.Call
}

// CreatePaymentOrder is a helper method to define mock.On call
//   - context echo.Context
func (_e *OrderHandlerIF_Expecter) CreatePaymentOrder(context interface{}) *OrderHandlerIF_CreatePaymentOrder_Call {
	return &OrderHandlerIF_CreatePaymentOrder_Call{Call: _e.mock.On("CreatePaymentOrder", context)}
}

func (_c *OrderHandlerIF_CreatePaymentOrder_Call) Run(run func(context echo.Context)) *OrderHandlerIF_CreatePaymentOrder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(echo.Context))
	})
	return _c
}

func (_c *OrderHandlerIF_CreatePaymentOrder_Call) Return(_a0 error) *OrderHandlerIF_CreatePaymentOrder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrderHandlerIF_CreatePaymentOrder_Call) RunAndReturn(run func(echo.Context) error) *OrderHandlerIF_CreatePaymentOrder_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrder provides a mock function with given fields: context
func (_m *OrderHandlerIF) GetOrder(context echo.Context) error {
	ret := _m.Called(context)

	if len(ret) == 0 {
		panic("no return value specified for GetOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(context)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrderHandlerIF_GetOrder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrder'
type OrderHandlerIF_GetOrder_Call struct {
	*mock.Call
}

// GetOrder is a helper method to define mock.On call
//   - context echo.Context
func (_e *OrderHandlerIF_Expecter) GetOrder(context interface{}) *OrderHandlerIF_GetOrder_Call {
	return &OrderHandlerIF_GetOrder_Call{Call: _e.mock.On("GetOrder", context)}
}

func (_c *OrderHandlerIF_GetOrder_Call) Run(run func(context echo.Context)) *OrderHandlerIF_GetOrder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(echo.Context))
	})
	return _c
}

func (_c *OrderHandlerIF_GetOrder_Call) Return(_a0 error) *OrderHandlerIF_GetOrder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrderHandlerIF_GetOrder_Call) RunAndReturn(run func(echo.Context) error) *OrderHandlerIF_GetOrder_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrderPayments provides a mock function with given fields: context
func (_m *OrderHandlerIF) GetOrderPayments(context echo.Context) error {
	ret := _m.Called(context)

	if len(ret) == 0 {
		panic("no return value specified for GetOrderPayments")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(context)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrderHandlerIF_GetOrderPayments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrderPayments'
type OrderHandlerIF_GetOrderPayments_Call struct {
	*mock.Call
}

// GetOrderPayments is a helper method to define mock.On call
//   - context echo.Context
func (_e *OrderHandlerIF_Expecter) GetOrderPayments(context interface{}) *OrderHandlerIF_GetOrderPayments_Call {
	return &OrderHandlerIF_GetOrderPayments_Call{Call: _e.mock.On("GetOrderPayments", context)}
}

func (_c *OrderHandlerIF_GetOrderPayments_Call) Run(run func(context echo.Context)) *OrderHandlerIF_GetOrderPayments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(echo.Context))
	})
	return _c
}

func (_c *OrderHandlerIF_GetOrderPayments_Call) Return(_a0 error) *OrderHandlerIF_GetOrderPayments_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrderHandlerIF_GetOrderPayments_Call) RunAndReturn(run func(echo.Context) error) *OrderHandlerIF_GetOrderPayments_Call {
	_c.Call.Return(run)
	return _c
}

// NewOrderHandlerIF creates a new instance of OrderHandlerIF. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderHandlerIF(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrderHandlerIF {
	mock := &OrderHandlerIF{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
