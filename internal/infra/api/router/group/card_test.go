package group

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCardHandler is a mock implementation of the CardHandlerIF interface
type MockCardHandler struct {
	mock.Mock
}

func (m *MockCardHandler) CreateCard(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCardHandler) GetCardsByUserID(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCardHandler) TriggerExpiringSoonNotifications(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func TestCardRoutes_GetCardsByUserID(t *testing.T) {
	e := echo.New()
	group := e.Group("/api")
	mockHandler := new(MockCardHandler)

	NewCardRoutes(group, mockHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/card/by-user/test-user-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/card/by-user/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues("test-user-id")

	mockHandler.On("GetCardsByUserID", c).Return(nil)

	if assert.NoError(t, mockHandler.GetCardsByUserID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	mockHandler.AssertExpectations(t)
}
