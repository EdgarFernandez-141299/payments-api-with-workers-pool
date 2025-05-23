// Code generated by mockery. DO NOT EDIT.

package group

import (
	echo "github.com/labstack/echo/v4"

	mock "github.com/stretchr/testify/mock"
)

// DeleteCardRoutesIF is an autogenerated mock type for the DeleteCardRoutesIF type
type DeleteCardRoutesIF struct {
	mock.Mock
}

type DeleteCardRoutesIF_Expecter struct {
	mock *mock.Mock
}

func (_m *DeleteCardRoutesIF) EXPECT() *DeleteCardRoutesIF_Expecter {
	return &DeleteCardRoutesIF_Expecter{mock: &_m.Mock}
}

// Resource provides a mock function with given fields: c
func (_m *DeleteCardRoutesIF) Resource(c *echo.Group) {
	_m.Called(c)
}

// DeleteCardRoutesIF_Resource_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Resource'
type DeleteCardRoutesIF_Resource_Call struct {
	*mock.Call
}

// Resource is a helper method to define mock.On call
//   - c *echo.Group
func (_e *DeleteCardRoutesIF_Expecter) Resource(c interface{}) *DeleteCardRoutesIF_Resource_Call {
	return &DeleteCardRoutesIF_Resource_Call{Call: _e.mock.On("Resource", c)}
}

func (_c *DeleteCardRoutesIF_Resource_Call) Run(run func(c *echo.Group)) *DeleteCardRoutesIF_Resource_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*echo.Group))
	})
	return _c
}

func (_c *DeleteCardRoutesIF_Resource_Call) Return() *DeleteCardRoutesIF_Resource_Call {
	_c.Call.Return()
	return _c
}

func (_c *DeleteCardRoutesIF_Resource_Call) RunAndReturn(run func(*echo.Group)) *DeleteCardRoutesIF_Resource_Call {
	_c.Run(run)
	return _c
}

// NewDeleteCardRoutesIF creates a new instance of DeleteCardRoutesIF. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDeleteCardRoutesIF(t interface {
	mock.TestingT
	Cleanup(func())
}) *DeleteCardRoutesIF {
	mock := &DeleteCardRoutesIF{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
