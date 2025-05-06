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
	dto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/request"
	paymentmethodsDto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/request/payment_methods"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
	mocks "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/create"
	auth2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/infra/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/queries"
)

func TestHandlerCreatePaymentOrder(t *testing.T) {
	test := []struct {
		name         string
		requestBody  interface{}
		mockError    error
		expectedCode int
		expectedBody string
		httpError    *echo.HTTPError
		mockResponse response.PaymentOrderResponseDTO
	}{
		{
			name:         "error marshal body",
			requestBody:  "invalid",
			expectedCode: http.StatusBadRequest,
			httpError:    echo.NewHTTPError(http.StatusBadRequest, `code=400, message=Unmarshal type error: expected=dto.CreatePaymentOrderRequestDTO, got=string, field=, offset=9, internal=json: cannot unmarshal string into Go value of type dto.CreatePaymentOrderRequestDTO`),
			mockResponse: response.PaymentOrderResponseDTO{},
		},
		{
			name: "error validation fields failed",
			requestBody: dto.CreatePaymentOrderRequestDTO{
				UserID:   "",
				OrderID:  "123",
				UserType: "member",
				PaymentOrderRequestDTO: dto.PaymentOrderRequestDTO{
					AssociatedOrigin: enums.Downpayment.String(),
					Amount:           decimal.NewFromFloat(1000),
					PaymentMethod: dto.PaymentMethodDTO{
						Type: "CCData",
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "",
							CVV: "123",
						},
					},
				},
			},
			expectedCode: http.StatusBadRequest,
			httpError:    echo.NewHTTPError(http.StatusBadRequest, `Key: 'CreatePaymentOrderRequestDTO.UserID' Error:Field validation for 'UserID' failed on the 'required' tag`),
			mockResponse: response.PaymentOrderResponseDTO{},
		},
		{
			name: "error usecase create payment order",
			requestBody: dto.CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: dto.PaymentOrderRequestDTO{
					PaymentOrderID:   "123",
					AssociatedOrigin: enums.Downpayment.String(),
					PaymentMethod: dto.PaymentMethodDTO{
						Type: "CCData",
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "1",
							CVV: "123",
						},
					},
					Amount: decimal.NewFromFloat(1000),
				},
				UserID:       "111",
				OrderID:      "123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "MX",
			},
			mockError:    errors.New("error create payment order"),
			expectedCode: http.StatusBadRequest,
			httpError:    echo.NewHTTPError(http.StatusBadRequest, `error create payment order`),
			mockResponse: response.PaymentOrderResponseDTO{},
		},
		{
			name: "error validate command",
			requestBody: dto.CreatePaymentOrderRequestDTO{
				UserID:   "111",
				OrderID:  "123",
				UserType: "member",
				PaymentOrderRequestDTO: dto.PaymentOrderRequestDTO{
					Amount:           decimal.NewFromFloat(1000),
					AssociatedOrigin: enums.Downpayment.String(),
					PaymentMethod: dto.PaymentMethodDTO{
						Type: "CCData",
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "1",
							CVV: "123",
						},
					},
					PaymentOrderID: "123",
				},
				CurrencyCode: "PPP",
			},
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			httpError:    echo.NewHTTPError(http.StatusBadRequest, `error create payment order`),
			mockResponse: response.PaymentOrderResponseDTO{},
		},
		{
			name: "should create payment order successfully",
			requestBody: dto.CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: dto.PaymentOrderRequestDTO{
					PaymentOrderID:   "123",
					AssociatedOrigin: enums.Downpayment.String(),
					Amount:           decimal.NewFromFloat(1000),
					PaymentMethod: dto.PaymentMethodDTO{
						Type: "CCData",
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "1",
							CVV: "123",
						},
					},
				},
				UserID:       "111",
				OrderID:      "123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "MX",
			},
			mockError:    nil,
			expectedCode: http.StatusCreated,
			httpError:    nil,
			mockResponse: response.PaymentOrderResponseDTO{
				ReferenceOrderID: "123",
			},
		},
	}

	for _, tt := range test {
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

			if tt.mockError == nil && tt.mockResponse != (response.PaymentOrderResponseDTO{}) {
				paymentOrderUsecaseMock.On("CreatePaymentOrder", mock.Anything, mock.Anything).Return(response.PaymentOrderResponseDTO{}, nil)
			} else if tt.mockError != nil {
				paymentOrderUsecaseMock.On("CreatePaymentOrder", mock.Anything, mock.Anything).Return(response.PaymentOrderResponseDTO{}, tt.mockError)
			}
			authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)
			req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/payment-order", bytes.NewReader(requestBody))
			responseRecorder := httptest.NewRecorder()
			c := echoInstance.NewContext(req, responseRecorder)
			c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			err := handler.CreatePaymentOrder(c)

			if tt.httpError != nil {
				assert.Equal(t, tt.httpError.Code, err.(*echo.HTTPError).Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, responseRecorder.Code)
			}

		})
	}
}
