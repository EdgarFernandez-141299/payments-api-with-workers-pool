package order

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	dto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
	auth2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/infra/auth"
	mocks "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/use_cases"
)

func TestPaymentRefundHandler_Refund(t *testing.T) {
	// Probar caso de error de binding (JSON inválido)
	t.Run("error_binding_json", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		echoInstance := echo.New()
		refundUsecaseMock := mocks.NewRefundTotalUseCaseIF(t)
		partialRefundUsecaseMock := mocks.NewPartialRefundUseCaseIF(t)
		handler := NewPaymentRefundHandler(refundUsecaseMock, partialRefundUsecaseMock)

		// Cuerpo completamente inválido (no es JSON)
		requestBody := []byte(`this is not json at all`)

		// Setup request
		authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/refund", bytes.NewReader(requestBody))
		responseRecorder := httptest.NewRecorder()
		c := echoInstance.NewContext(req, responseRecorder)
		c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act - la función puede retornar nil porque maneja el error con context.JSON()
		handler.Refund(c)

		// Assert - Verificamos el código de respuesta HTTP, no el error retornado
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.NotEmpty(t, responseRecorder.Body.String())
	})

	t.Run("validation_error", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		echoInstance := echo.New()
		refundUsecaseMock := mocks.NewRefundTotalUseCaseIF(t)
		partialRefundUsecaseMock := mocks.NewPartialRefundUseCaseIF(t)
		handler := NewPaymentRefundHandler(refundUsecaseMock, partialRefundUsecaseMock)

		invalidRefund := dto.RefundDTO{
			IsTotal: true,
		}
		requestBody, _ := json.Marshal(invalidRefund)

		authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/refund", bytes.NewReader(requestBody))
		responseRecorder := httptest.NewRecorder()
		c := echoInstance.NewContext(req, responseRecorder)
		c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		handler.Refund(c)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.NotEmpty(t, responseRecorder.Body.String())
	})
	t.Run("usecase_error", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		echoInstance := echo.New()
		refundUsecaseMock := mocks.NewRefundTotalUseCaseIF(t)
		partialRefundUsecaseMock := mocks.NewPartialRefundUseCaseIF(t)
		handler := NewPaymentRefundHandler(refundUsecaseMock, partialRefundUsecaseMock)

		validRefund := dto.RefundDTO{
			ReferenceOrderID: "ref-123",
			PaymentOrderID:   "payment-123",
			Reason:           "customer request",
			Amount:           decimal.NewFromFloat(12),
			IsTotal:          false,
		}
		requestBody, _ := json.Marshal(validRefund)

		mockError := errors.NewRefundCanNotAppliedDueToPaymentStatusError()

		partialRefundUsecaseMock.On("PartialRefund", mock.Anything, mock.Anything, enterpriseID).
			Return(response.RefundResponseDTO{}, mockError)

		// Setup request
		authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/refund", bytes.NewReader(requestBody))
		responseRecorder := httptest.NewRecorder()
		c := echoInstance.NewContext(req, responseRecorder)
		c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		err := handler.Refund(c)
		if err != nil {
			if httpErr, ok := err.(*echo.HTTPError); ok {
				responseRecorder.WriteHeader(httpErr.Code)
				jsonBody, _ := json.Marshal(httpErr.Message)
				responseRecorder.Write(jsonBody)
			}
		}

		// Assert
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.NotEmpty(t, responseRecorder.Body.String())
	})

	t.Run("successful_refund", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		echoInstance := echo.New()
		refundUsecaseMock := mocks.NewRefundTotalUseCaseIF(t)
		partialRefundUsecaseMock := mocks.NewPartialRefundUseCaseIF(t)
		handler := NewPaymentRefundHandler(refundUsecaseMock, partialRefundUsecaseMock)

		validRefund := dto.RefundDTO{
			ReferenceOrderID: "ref-123",
			PaymentOrderID:   "payment-123",
			Reason:           "customer request",
			IsTotal:          true,
		}
		requestBody, _ := json.Marshal(validRefund)

		expectedResponse := response.RefundResponseDTO{
			ReferenceOrderID: "ref-123",
			PaymentOrderID:   "payment-123",
		}

		refundUsecaseMock.On("Refund", mock.Anything, mock.Anything, enterpriseID).
			Return(expectedResponse, nil)

		// Setup request
		authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/refund", bytes.NewReader(requestBody))
		responseRecorder := httptest.NewRecorder()
		c := echoInstance.NewContext(req, responseRecorder)
		c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		handler.Refund(c)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)

		var actualResponse response.RefundResponseDTO
		err := json.Unmarshal(responseRecorder.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.ReferenceOrderID, actualResponse.ReferenceOrderID)
		assert.Equal(t, expectedResponse.PaymentOrderID, actualResponse.PaymentOrderID)
	})

	t.Run("non_total_refund", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		echoInstance := echo.New()
		refundUsecaseMock := mocks.NewRefundTotalUseCaseIF(t)
		partialRefundUsecaseMock := mocks.NewPartialRefundUseCaseIF(t)
		handler := NewPaymentRefundHandler(refundUsecaseMock, partialRefundUsecaseMock)

		// Datos de reembolso parcial
		nonTotalRefund := dto.RefundDTO{
			ReferenceOrderID: "ref-123",
			PaymentOrderID:   "payment-123",
			Reason:           "customer request",
			Amount:           decimal.NewFromFloat(50.00),
			IsTotal:          false,
		}
		requestBody, _ := json.Marshal(nonTotalRefund)

		expectedResponse := response.RefundResponseDTO{
			ReferenceOrderID: "ref-123",
			PaymentOrderID:   "payment-123",
			Amount:           decimal.NewFromFloat(50.00),
		}

		partialRefundUsecaseMock.On("PartialRefund",
			mock.Anything,
			mock.AnythingOfType("command.CreatePartialPaymentRefundCommand"),
			"test-enterprise-id",
		).Return(expectedResponse, nil)

		// Setup request
		authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, "test-enterprise-id") // Asegurar que usamos el mismo enterprise ID
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/refund", bytes.NewReader(requestBody))
		responseRecorder := httptest.NewRecorder()
		c := echoInstance.NewContext(req, responseRecorder)
		c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		handler.Refund(c)

		// Assert
		assert.Equal(t, http.StatusOK, responseRecorder.Code)

		var actualResponse response.RefundResponseDTO
		err := json.Unmarshal(responseRecorder.Body.Bytes(), &actualResponse)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.ReferenceOrderID, actualResponse.ReferenceOrderID)
		assert.Equal(t, expectedResponse.PaymentOrderID, actualResponse.PaymentOrderID)

		partialRefundUsecaseMock.AssertExpectations(t)
	})

	t.Run("usecase_error_total_refund", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		echoInstance := echo.New()
		refundUsecaseMock := mocks.NewRefundTotalUseCaseIF(t)
		partialRefundUsecaseMock := mocks.NewPartialRefundUseCaseIF(t)
		handler := NewPaymentRefundHandler(refundUsecaseMock, partialRefundUsecaseMock)

		validRefund := dto.RefundDTO{
			ReferenceOrderID: "ref-123",
			PaymentOrderID:   "payment-123",
			Reason:           "customer request",
			IsTotal:          true,
		}
		requestBody, _ := json.Marshal(validRefund)

		mockError := errors.NewRefundCanNotAppliedDueToPaymentStatusError()

		refundUsecaseMock.On("Refund", mock.Anything, mock.Anything, enterpriseID).
			Return(response.RefundResponseDTO{}, mockError)

		// Setup request
		authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/refund", bytes.NewReader(requestBody))
		responseRecorder := httptest.NewRecorder()
		c := echoInstance.NewContext(req, responseRecorder)
		c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		err := handler.Refund(c)
		if err != nil {
			if httpErr, ok := err.(*echo.HTTPError); ok {
				responseRecorder.WriteHeader(httpErr.Code)
				jsonBody, _ := json.Marshal(httpErr.Message)
				responseRecorder.Write(jsonBody)
			}
		}

		// Assert
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.NotEmpty(t, responseRecorder.Body.String())
	})

	t.Run("command_error", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		echoInstance := echo.New()
		refundUsecaseMock := mocks.NewRefundTotalUseCaseIF(t)
		partialRefundUsecaseMock := mocks.NewPartialRefundUseCaseIF(t)
		handler := NewPaymentRefundHandler(refundUsecaseMock, partialRefundUsecaseMock)

		invalidRefund := dto.RefundDTO{
			IsTotal: true,
			// Configura el DTO para que Command() falle
		}
		requestBody, _ := json.Marshal(invalidRefund)

		authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/refund", bytes.NewReader(requestBody))
		responseRecorder := httptest.NewRecorder()
		c := echoInstance.NewContext(req, responseRecorder)
		c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		handler.Refund(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.NotEmpty(t, responseRecorder.Body.String())
	})

	t.Run("command_partial_error", func(t *testing.T) {
		// Setup
		ctx := context.Background()
		echoInstance := echo.New()
		refundUsecaseMock := mocks.NewRefundTotalUseCaseIF(t)
		partialRefundUsecaseMock := mocks.NewPartialRefundUseCaseIF(t)
		handler := NewPaymentRefundHandler(refundUsecaseMock, partialRefundUsecaseMock)

		invalidRefund := dto.RefundDTO{
			IsTotal: false,
			// Configura el DTO para que CommandPartial() falle
		}
		requestBody, _ := json.Marshal(invalidRefund)

		authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/refund", bytes.NewReader(requestBody))
		responseRecorder := httptest.NewRecorder()
		c := echoInstance.NewContext(req, responseRecorder)
		c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Act
		handler.Refund(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.NotEmpty(t, responseRecorder.Body.String())
	})
}
