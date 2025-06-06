// Code generated by mockery. DO NOT EDIT.

package repository

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"

	projections "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/projections"

	time "time"
)

// CardReadRepositoryIF is an autogenerated mock type for the CardReadRepositoryIF type
type CardReadRepositoryIF struct {
	mock.Mock
}

type CardReadRepositoryIF_Expecter struct {
	mock *mock.Mock
}

func (_m *CardReadRepositoryIF) EXPECT() *CardReadRepositoryIF_Expecter {
	return &CardReadRepositoryIF_Expecter{mock: &_m.Mock}
}

// CheckCardExistence provides a mock function with given fields: ctx, userID, enterpriseID, lastFour
func (_m *CardReadRepositoryIF) CheckCardExistence(ctx context.Context, userID string, enterpriseID string, lastFour string) (bool, error) {
	ret := _m.Called(ctx, userID, enterpriseID, lastFour)

	if len(ret) == 0 {
		panic("no return value specified for CheckCardExistence")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) (bool, error)); ok {
		return rf(ctx, userID, enterpriseID, lastFour)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) bool); ok {
		r0 = rf(ctx, userID, enterpriseID, lastFour)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, userID, enterpriseID, lastFour)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CardReadRepositoryIF_CheckCardExistence_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckCardExistence'
type CardReadRepositoryIF_CheckCardExistence_Call struct {
	*mock.Call
}

// CheckCardExistence is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - enterpriseID string
//   - lastFour string
func (_e *CardReadRepositoryIF_Expecter) CheckCardExistence(ctx interface{}, userID interface{}, enterpriseID interface{}, lastFour interface{}) *CardReadRepositoryIF_CheckCardExistence_Call {
	return &CardReadRepositoryIF_CheckCardExistence_Call{Call: _e.mock.On("CheckCardExistence", ctx, userID, enterpriseID, lastFour)}
}

func (_c *CardReadRepositoryIF_CheckCardExistence_Call) Run(run func(ctx context.Context, userID string, enterpriseID string, lastFour string)) *CardReadRepositoryIF_CheckCardExistence_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *CardReadRepositoryIF_CheckCardExistence_Call) Return(_a0 bool, _a1 error) *CardReadRepositoryIF_CheckCardExistence_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CardReadRepositoryIF_CheckCardExistence_Call) RunAndReturn(run func(context.Context, string, string, string) (bool, error)) *CardReadRepositoryIF_CheckCardExistence_Call {
	_c.Call.Return(run)
	return _c
}

// GetCardAndUserEmailByUserID provides a mock function with given fields: ctx, userID, cardID, enterpriseID
func (_m *CardReadRepositoryIF) GetCardAndUserEmailByUserID(ctx context.Context, userID string, cardID string, enterpriseID string) (*projections.CardUserEmailProjection, error) {
	ret := _m.Called(ctx, userID, cardID, enterpriseID)

	if len(ret) == 0 {
		panic("no return value specified for GetCardAndUserEmailByUserID")
	}

	var r0 *projections.CardUserEmailProjection
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) (*projections.CardUserEmailProjection, error)); ok {
		return rf(ctx, userID, cardID, enterpriseID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *projections.CardUserEmailProjection); ok {
		r0 = rf(ctx, userID, cardID, enterpriseID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*projections.CardUserEmailProjection)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, userID, cardID, enterpriseID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CardReadRepositoryIF_GetCardAndUserEmailByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCardAndUserEmailByUserID'
type CardReadRepositoryIF_GetCardAndUserEmailByUserID_Call struct {
	*mock.Call
}

// GetCardAndUserEmailByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - cardID string
//   - enterpriseID string
func (_e *CardReadRepositoryIF_Expecter) GetCardAndUserEmailByUserID(ctx interface{}, userID interface{}, cardID interface{}, enterpriseID interface{}) *CardReadRepositoryIF_GetCardAndUserEmailByUserID_Call {
	return &CardReadRepositoryIF_GetCardAndUserEmailByUserID_Call{Call: _e.mock.On("GetCardAndUserEmailByUserID", ctx, userID, cardID, enterpriseID)}
}

func (_c *CardReadRepositoryIF_GetCardAndUserEmailByUserID_Call) Run(run func(ctx context.Context, userID string, cardID string, enterpriseID string)) *CardReadRepositoryIF_GetCardAndUserEmailByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *CardReadRepositoryIF_GetCardAndUserEmailByUserID_Call) Return(_a0 *projections.CardUserEmailProjection, _a1 error) *CardReadRepositoryIF_GetCardAndUserEmailByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CardReadRepositoryIF_GetCardAndUserEmailByUserID_Call) RunAndReturn(run func(context.Context, string, string, string) (*projections.CardUserEmailProjection, error)) *CardReadRepositoryIF_GetCardAndUserEmailByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// GetCardByUserID provides a mock function with given fields: ctx, userID, cardID, enterpriseID
func (_m *CardReadRepositoryIF) GetCardByUserID(ctx context.Context, userID string, cardID string, enterpriseID string) (entities.CardEntity, error) {
	ret := _m.Called(ctx, userID, cardID, enterpriseID)

	if len(ret) == 0 {
		panic("no return value specified for GetCardByUserID")
	}

	var r0 entities.CardEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) (entities.CardEntity, error)); ok {
		return rf(ctx, userID, cardID, enterpriseID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) entities.CardEntity); ok {
		r0 = rf(ctx, userID, cardID, enterpriseID)
	} else {
		r0 = ret.Get(0).(entities.CardEntity)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, userID, cardID, enterpriseID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CardReadRepositoryIF_GetCardByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCardByUserID'
type CardReadRepositoryIF_GetCardByUserID_Call struct {
	*mock.Call
}

// GetCardByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - cardID string
//   - enterpriseID string
func (_e *CardReadRepositoryIF_Expecter) GetCardByUserID(ctx interface{}, userID interface{}, cardID interface{}, enterpriseID interface{}) *CardReadRepositoryIF_GetCardByUserID_Call {
	return &CardReadRepositoryIF_GetCardByUserID_Call{Call: _e.mock.On("GetCardByUserID", ctx, userID, cardID, enterpriseID)}
}

func (_c *CardReadRepositoryIF_GetCardByUserID_Call) Run(run func(ctx context.Context, userID string, cardID string, enterpriseID string)) *CardReadRepositoryIF_GetCardByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *CardReadRepositoryIF_GetCardByUserID_Call) Return(_a0 entities.CardEntity, _a1 error) *CardReadRepositoryIF_GetCardByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CardReadRepositoryIF_GetCardByUserID_Call) RunAndReturn(run func(context.Context, string, string, string) (entities.CardEntity, error)) *CardReadRepositoryIF_GetCardByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// GetCardsByUserID provides a mock function with given fields: ctx, userID, enterpriseID
func (_m *CardReadRepositoryIF) GetCardsByUserID(ctx context.Context, userID string, enterpriseID string) (entities.CardEntities, error) {
	ret := _m.Called(ctx, userID, enterpriseID)

	if len(ret) == 0 {
		panic("no return value specified for GetCardsByUserID")
	}

	var r0 entities.CardEntities
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (entities.CardEntities, error)); ok {
		return rf(ctx, userID, enterpriseID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) entities.CardEntities); ok {
		r0 = rf(ctx, userID, enterpriseID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(entities.CardEntities)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userID, enterpriseID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CardReadRepositoryIF_GetCardsByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCardsByUserID'
type CardReadRepositoryIF_GetCardsByUserID_Call struct {
	*mock.Call
}

// GetCardsByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - enterpriseID string
func (_e *CardReadRepositoryIF_Expecter) GetCardsByUserID(ctx interface{}, userID interface{}, enterpriseID interface{}) *CardReadRepositoryIF_GetCardsByUserID_Call {
	return &CardReadRepositoryIF_GetCardsByUserID_Call{Call: _e.mock.On("GetCardsByUserID", ctx, userID, enterpriseID)}
}

func (_c *CardReadRepositoryIF_GetCardsByUserID_Call) Run(run func(ctx context.Context, userID string, enterpriseID string)) *CardReadRepositoryIF_GetCardsByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *CardReadRepositoryIF_GetCardsByUserID_Call) Return(_a0 entities.CardEntities, _a1 error) *CardReadRepositoryIF_GetCardsByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CardReadRepositoryIF_GetCardsByUserID_Call) RunAndReturn(run func(context.Context, string, string) (entities.CardEntities, error)) *CardReadRepositoryIF_GetCardsByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// GetCardsExpiringSoon provides a mock function with given fields: ctx, expirationMonth, expirationYear
func (_m *CardReadRepositoryIF) GetCardsExpiringSoon(ctx context.Context, expirationMonth time.Month, expirationYear int) ([]projections.NotificationCardExpiringSoonProjection, error) {
	ret := _m.Called(ctx, expirationMonth, expirationYear)

	if len(ret) == 0 {
		panic("no return value specified for GetCardsExpiringSoon")
	}

	var r0 []projections.NotificationCardExpiringSoonProjection
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Month, int) ([]projections.NotificationCardExpiringSoonProjection, error)); ok {
		return rf(ctx, expirationMonth, expirationYear)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Month, int) []projections.NotificationCardExpiringSoonProjection); ok {
		r0 = rf(ctx, expirationMonth, expirationYear)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]projections.NotificationCardExpiringSoonProjection)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Month, int) error); ok {
		r1 = rf(ctx, expirationMonth, expirationYear)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CardReadRepositoryIF_GetCardsExpiringSoon_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCardsExpiringSoon'
type CardReadRepositoryIF_GetCardsExpiringSoon_Call struct {
	*mock.Call
}

// GetCardsExpiringSoon is a helper method to define mock.On call
//   - ctx context.Context
//   - expirationMonth time.Month
//   - expirationYear int
func (_e *CardReadRepositoryIF_Expecter) GetCardsExpiringSoon(ctx interface{}, expirationMonth interface{}, expirationYear interface{}) *CardReadRepositoryIF_GetCardsExpiringSoon_Call {
	return &CardReadRepositoryIF_GetCardsExpiringSoon_Call{Call: _e.mock.On("GetCardsExpiringSoon", ctx, expirationMonth, expirationYear)}
}

func (_c *CardReadRepositoryIF_GetCardsExpiringSoon_Call) Run(run func(ctx context.Context, expirationMonth time.Month, expirationYear int)) *CardReadRepositoryIF_GetCardsExpiringSoon_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(time.Month), args[2].(int))
	})
	return _c
}

func (_c *CardReadRepositoryIF_GetCardsExpiringSoon_Call) Return(_a0 []projections.NotificationCardExpiringSoonProjection, _a1 error) *CardReadRepositoryIF_GetCardsExpiringSoon_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CardReadRepositoryIF_GetCardsExpiringSoon_Call) RunAndReturn(run func(context.Context, time.Month, int) ([]projections.NotificationCardExpiringSoonProjection, error)) *CardReadRepositoryIF_GetCardsExpiringSoon_Call {
	_c.Call.Return(run)
	return _c
}

// NewCardReadRepositoryIF creates a new instance of CardReadRepositoryIF. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCardReadRepositoryIF(t interface {
	mock.TestingT
	Cleanup(func())
}) *CardReadRepositoryIF {
	mock := &CardReadRepositoryIF{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
