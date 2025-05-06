package order

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
	mocks "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/create"
	auth2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/infra/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/queries"
)

const (
	userID       = "test-user-id"
	enterpriseID = "test-enterprise-id"
	userType     = "member"
)

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		mockError    error
		expectedCode int
		expectedBody string
		httpError    *echo.HTTPError
		mockResponse response.OrderResponseDTO
	}{
		{
			name:         "invalid request body",
			requestBody:  "invalid",
			expectedCode: http.StatusBadRequest,
			httpError:    echo.NewHTTPError(http.StatusBadRequest, `code=400, message=Unmarshal type error: expected=dto.CreateOrderRequestDTO, got=string, field=, offset=9, internal=json: cannot unmarshal string into Go value of type dto.CreateOrderRequestDTO`),
			mockResponse: response.OrderResponseDTO{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			echoInstance := echo.New()

			orderUsecaseMock := mocks.NewCreateOrderUseCaseIF(t)
			paymentOrderUsecaseMock := mocks.NewCreatePaymentOrderUseCaseIF(t)
			queriesOrderUsecaseMock := queries.NewQueriesOrderUseCaseIF(t)
			handler := NewOrderHandler(orderUsecaseMock, paymentOrderUsecaseMock, queriesOrderUsecaseMock)

			var requestBody []byte
			if r, ok := tt.requestBody.([]byte); ok {
				requestBody = r
			} else {
				requestBody, _ = json.Marshal(tt.requestBody)
			}
			if tt.mockError == nil && tt.mockResponse != (response.OrderResponseDTO{}) {
				orderUsecaseMock.On("Create", mock.Anything, mock.Anything, enterpriseID).Return(tt.mockResponse, nil)
			} else if tt.mockError != nil {
				orderUsecaseMock.On("Create", mock.Anything, mock.Anything, enterpriseID).Return(response.OrderResponseDTO{}, tt.mockError)
			}

			authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)
			req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/order", bytes.NewReader(requestBody))
			responseRecorder := httptest.NewRecorder()
			c := echoInstance.NewContext(req, responseRecorder)
			c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			err := handler.Create(c)

			if tt.httpError != nil {
				assert.Equal(t, tt.httpError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, responseRecorder.Code)
			}
		})
	}
}

func TestGetOrderPayment(t *testing.T) {
	tests := []struct {
		name         string
		mockError    error
		expectedCode int
		httpError    *echo.HTTPError
		mockResponse *response.GetOrderPaymentResponseDTO
	}{
		{
			name:         "valid order payment retrieval",
			expectedCode: http.StatusOK,
			mockError:    nil,
			httpError:    nil,
			mockResponse: &response.GetOrderPaymentResponseDTO{
				ReferenceOrderID: "payment_order_reference_46",
				Status:           enums.PaymentProcessing,
				Total:            decimal.NewFromFloat(100.50),
				Currency:         "USD",
				CountryCode:      "US",
				Metadata:         map[string]string{},
				Payments:         []response.PaymentDTO{},
			},
		},
		{
			name:         "database error",
			expectedCode: http.StatusBadRequest,
			httpError:    echo.NewHTTPError(http.StatusBadRequest, "database error"),
			mockError:    errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			echoInstance := echo.New()

			queriesOrderUsecaseMock := queries.NewQueriesOrderUseCaseIF(t)
			handler := NewOrderHandler(nil, nil, queriesOrderUsecaseMock)

			authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)
			req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/order/status/payment_order_reference_46", nil)
			responseRecorder := httptest.NewRecorder()
			c := echoInstance.NewContext(req, responseRecorder)
			c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			queriesOrderUsecaseMock.On("GetOrderPayments", mock.Anything, mock.Anything, enterpriseID).Return(tt.mockResponse, tt.mockError)

			err := handler.GetOrderPayments(c)

			if tt.httpError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, responseRecorder.Code)
			}
		})
	}
}
