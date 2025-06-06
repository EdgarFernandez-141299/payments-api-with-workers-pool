// Code generated by mockery. DO NOT EDIT.

package adapters

import (
	context "context"

	adapters "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common/adapters"

	mock "github.com/stretchr/testify/mock"
)

// MailAdapterIF is an autogenerated mock type for the MailAdapterIF type
type MailAdapterIF struct {
	mock.Mock
}

type MailAdapterIF_Expecter struct {
	mock *mock.Mock
}

func (_m *MailAdapterIF) EXPECT() *MailAdapterIF_Expecter {
	return &MailAdapterIF_Expecter{mock: &_m.Mock}
}

// Send provides a mock function with given fields: ctx, request
func (_m *MailAdapterIF) Send(ctx context.Context, request adapters.MailRequest) error {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, adapters.MailRequest) error); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MailAdapterIF_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type MailAdapterIF_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - ctx context.Context
//   - request adapters.MailRequest
func (_e *MailAdapterIF_Expecter) Send(ctx interface{}, request interface{}) *MailAdapterIF_Send_Call {
	return &MailAdapterIF_Send_Call{Call: _e.mock.On("Send", ctx, request)}
}

func (_c *MailAdapterIF_Send_Call) Run(run func(ctx context.Context, request adapters.MailRequest)) *MailAdapterIF_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(adapters.MailRequest))
	})
	return _c
}

func (_c *MailAdapterIF_Send_Call) Return(_a0 error) *MailAdapterIF_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MailAdapterIF_Send_Call) RunAndReturn(run func(context.Context, adapters.MailRequest) error) *MailAdapterIF_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewMailAdapterIF creates a new instance of MailAdapterIF. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMailAdapterIF(t interface {
	mock.TestingT
	Cleanup(func())
}) *MailAdapterIF {
	mock := &MailAdapterIF{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
