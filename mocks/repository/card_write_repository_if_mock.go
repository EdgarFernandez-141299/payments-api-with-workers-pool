// Code generated by mockery. DO NOT EDIT.

package repository

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"
)

// CardWriteRepositoryIF is an autogenerated mock type for the CardWriteRepositoryIF type
type CardWriteRepositoryIF struct {
	mock.Mock
}

type CardWriteRepositoryIF_Expecter struct {
	mock *mock.Mock
}

func (_m *CardWriteRepositoryIF) EXPECT() *CardWriteRepositoryIF_Expecter {
	return &CardWriteRepositoryIF_Expecter{mock: &_m.Mock}
}

// CreateCard provides a mock function with given fields: ctx, entity
func (_m *CardWriteRepositoryIF) CreateCard(ctx context.Context, entity *entities.CardEntity) error {
	ret := _m.Called(ctx, entity)

	if len(ret) == 0 {
		panic("no return value specified for CreateCard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entities.CardEntity) error); ok {
		r0 = rf(ctx, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CardWriteRepositoryIF_CreateCard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateCard'
type CardWriteRepositoryIF_CreateCard_Call struct {
	*mock.Call
}

// CreateCard is a helper method to define mock.On call
//   - ctx context.Context
//   - entity *entities.CardEntity
func (_e *CardWriteRepositoryIF_Expecter) CreateCard(ctx interface{}, entity interface{}) *CardWriteRepositoryIF_CreateCard_Call {
	return &CardWriteRepositoryIF_CreateCard_Call{Call: _e.mock.On("CreateCard", ctx, entity)}
}

func (_c *CardWriteRepositoryIF_CreateCard_Call) Run(run func(ctx context.Context, entity *entities.CardEntity)) *CardWriteRepositoryIF_CreateCard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entities.CardEntity))
	})
	return _c
}

func (_c *CardWriteRepositoryIF_CreateCard_Call) Return(_a0 error) *CardWriteRepositoryIF_CreateCard_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CardWriteRepositoryIF_CreateCard_Call) RunAndReturn(run func(context.Context, *entities.CardEntity) error) *CardWriteRepositoryIF_CreateCard_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteCard provides a mock function with given fields: ctx, cardId
func (_m *CardWriteRepositoryIF) DeleteCard(ctx context.Context, cardId string) error {
	ret := _m.Called(ctx, cardId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, cardId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CardWriteRepositoryIF_DeleteCard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteCard'
type CardWriteRepositoryIF_DeleteCard_Call struct {
	*mock.Call
}

// DeleteCard is a helper method to define mock.On call
//   - ctx context.Context
//   - cardId string
func (_e *CardWriteRepositoryIF_Expecter) DeleteCard(ctx interface{}, cardId interface{}) *CardWriteRepositoryIF_DeleteCard_Call {
	return &CardWriteRepositoryIF_DeleteCard_Call{Call: _e.mock.On("DeleteCard", ctx, cardId)}
}

func (_c *CardWriteRepositoryIF_DeleteCard_Call) Run(run func(ctx context.Context, cardId string)) *CardWriteRepositoryIF_DeleteCard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CardWriteRepositoryIF_DeleteCard_Call) Return(_a0 error) *CardWriteRepositoryIF_DeleteCard_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CardWriteRepositoryIF_DeleteCard_Call) RunAndReturn(run func(context.Context, string) error) *CardWriteRepositoryIF_DeleteCard_Call {
	_c.Call.Return(run)
	return _c
}

// NewCardWriteRepositoryIF creates a new instance of CardWriteRepositoryIF. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCardWriteRepositoryIF(t interface {
	mock.TestingT
	Cleanup(func())
}) *CardWriteRepositoryIF {
	mock := &CardWriteRepositoryIF{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
