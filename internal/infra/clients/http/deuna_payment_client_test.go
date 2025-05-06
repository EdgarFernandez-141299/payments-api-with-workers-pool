package http

import (
	"context"
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/mocks/apm"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/request"
	config2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/mockrequest"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/mockresponse"
)

func createMockPaymentServer(responseCode int, responseBody string, verifyRequest func(r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifyRequest(r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseCode)
		w.Write([]byte(responseBody))
	}))
}

func createPaymentApiClient(mockServerURL, apiKey string, t *testing.T) resources.DeunaPaymentResourceIF {
	config := &config2.DeUnaApiConfig{
		URL:    mockServerURL,
		ApiKey: apiKey,
	}

	tracer := apm.NewTracer(t)
	tracer.On("GetTracer").Return(nil)

	return NewDeUnaPaymentHTTPClient(config, tracer)
}

const paymentApiKey = "2jkmdj4-sd4j3k4-3j4k3j4-3j4k3j4"
const paymentOrderToken = "c5c5abc9-67a7-4a44-aa83-ee7e81d7052c"

var (
	CardHolder    = "John Doe"
	CardHolderDNI = "123456789"
	Country       = "US"
	State         = "NY"
	City          = "New York"
	Zip           = "10001"
	Address1      = "123 Main St"
	Phone         = "+1234567890"
	ExpiryMonth   = "12"
	ExpiryYear    = "2025"
	CardCVV       = "123"
)

// Test functions
func TestDeUnaPaymentHTTPClient_MakeOrderPayment(t *testing.T) {
	ctx := context.Background()

	// Create a mock request for MakeOrderPayment
	mockPaymentRequest := request.DeunaOrderPaymentRequest{
		OrderToken: paymentOrderToken,
		MethodType: "card",
		Email:      "test@example.com",
		StoreCode:  "store-123",
		CreditCard: &request.CreditCardInfo{
			CardHolder:    &CardHolder,
			CardHolderDNI: &CardHolderDNI,
			Country:       &Country,
			State:         &State,
			City:          &City,
			Zip:           &Zip,
			Address1:      &Address1,
			Phone:         &Phone,
			ExpiryMonth:   &ExpiryMonth,
			ExpiryYear:    &ExpiryYear,
			CardCVV:       &CardCVV,
		},
		SaveUserInfo: true,
	}

	t.Run("successful payment operation", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusOK, mockresponse.SuccessOrderPaymentResponseMock, func(r *http.Request) {
			assert.Equal(t, "/merchants/transactions/purchase", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, paymentApiKey, r.Header.Get("x-api-key"))
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		err := client.MakeOrderPayment(ctx, mockPaymentRequest, "")

		assert.NoError(t, err)
	})

	t.Run("failure when server returns 500 error", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusInternalServerError, `{"error": "Internal server error"}`, func(r *http.Request) {
			assert.Equal(t, "/merchants/transactions/purchase", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		err := client.MakeOrderPayment(ctx, mockPaymentRequest, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "make order payment failed with status: 500")
	})

	t.Run("failure when resource not found (404)", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusNotFound, `{"error": "Resource not found"}`, func(r *http.Request) {
			assert.Equal(t, "/merchants/transactions/purchase", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		err := client.MakeOrderPayment(ctx, mockPaymentRequest, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "make order payment failed with status: 404")
	})
}

func TestDeUnaPaymentHTTPClient_MakeTotalRefund(t *testing.T) {
	ctx := context.Background()

	mockRefundRequest := utils.DeunaTotalRefundRequest{
		Reason: "Customer request",
	}

	t.Run("successful refund operation", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusOK, mockresponse.SuccessRefundPaymentResponseMock, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/v2/merchants/orders/%s/refund", paymentOrderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, paymentApiKey, r.Header.Get("x-api-key"))
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		gotResponse, err := client.MakeTotalRefund(ctx, mockRefundRequest, paymentOrderToken)

		assert.NoError(t, err)
		assert.Equal(t, "refund-123", gotResponse.Data.RefundID, "PartialRefund ID does not match expected")
		assert.Equal(t, "approved", gotResponse.Data.Status, "PartialRefund status does not match expected")
		assert.Equal(t, "125", gotResponse.Data.RefundAmount.Amount, "PartialRefund amount does not match expected")
	})

	t.Run("failure when server returns 500 error", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusInternalServerError, `{"error": "Internal server error"}`, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/v2/merchants/orders/%s/refund", paymentOrderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		_, err := client.MakeTotalRefund(ctx, mockRefundRequest, paymentOrderToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "refund payment failed with status: 500")
	})

	t.Run("failure when order not found (404)", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusNotFound, `{"error": "Order not found"}`, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/v2/merchants/orders/%s/refund", paymentOrderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		_, err := client.MakeTotalRefund(ctx, mockRefundRequest, paymentOrderToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "refund payment failed with status: 404")
	})
}

func TestDeUnaPaymentHTTPClient_MakePartialRefund(t *testing.T) {
	ctx := context.Background()
	amount := utils.DeunaAmountToAmount(12500)

	mockRefundRequest := utils.DeunaPartialRefundRequest{
		Amount: amount.IntPart(),
		Reason: "Customer request",
	}

	t.Run("successful refund operation", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusOK, mockresponse.SuccessRefundPaymentResponseMock, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/v2/merchants/orders/%s/refund", paymentOrderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, paymentApiKey, r.Header.Get("x-api-key"))
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		gotResponse, err := client.MakePartialRefund(ctx, mockRefundRequest, paymentOrderToken)

		assert.NoError(t, err)
		assert.Equal(t, "refund-123", gotResponse.Data.RefundID, "PartialRefund ID does not match expected")
		assert.Equal(t, "approved", gotResponse.Data.Status, "PartialRefund status does not match expected")
		assert.Equal(t, "125", gotResponse.Data.RefundAmount.Amount, "PartialRefund amount does not match expected")
	})

	t.Run("failure when server returns 500 error", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusInternalServerError, `{"error": "Internal server error"}`, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/v2/merchants/orders/%s/refund", paymentOrderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		_, err := client.MakePartialRefund(ctx, mockRefundRequest, paymentOrderToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "refund payment failed with status: 500")
	})

	t.Run("failure when order not found (404)", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusNotFound, `{"error": "Order not found"}`, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/v2/merchants/orders/%s/refund", paymentOrderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		_, err := client.MakePartialRefund(ctx, mockRefundRequest, paymentOrderToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "refund payment failed with status: 404")
	})
}

func TestDeUnaPaymentHTTPClient_MakeOrderPaymentV2(t *testing.T) {
	ctx := context.Background()

	mockPaymentRequestV2 := mockrequest.MakeOrderPaymentV2RequestMock

	t.Run("successful payment operation", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusOK, mockresponse.SuccessOrderPaymentV2ResponseMock, func(r *http.Request) {
			assert.Equal(t, "/merchants/orders/purchase", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, paymentApiKey, r.Header.Get("x-api-key"))
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		gotResponse, err := client.MakeOrderPaymentV2(ctx, mockPaymentRequestV2)

		assert.NoError(t, err)
		assert.Equal(t, paymentOrderToken, gotResponse.OrderToken, "Response token does not match expected")
		assert.Equal(t, "completed", gotResponse.Order.Status, "Order status does not match expected")
		assert.Equal(t, "approved", gotResponse.Order.Payment.Data.Status, "Payment status does not match expected")
	})

	t.Run("failure when server returns 500 error", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusInternalServerError, `{"error": "Internal server error"}`, func(r *http.Request) {
			assert.Equal(t, "/merchants/orders/purchase", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		_, err := client.MakeOrderPaymentV2(ctx, mockPaymentRequestV2)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "make order payment failed with status: 500")
	})

	t.Run("failure when resource not found (404)", func(t *testing.T) {
		mockServer := createMockPaymentServer(http.StatusNotFound, `{"error": "Resource not found"}`, func(r *http.Request) {
			assert.Equal(t, "/merchants/orders/purchase", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
		})
		defer mockServer.Close()

		client := createPaymentApiClient(mockServer.URL, paymentApiKey, t)

		_, err := client.MakeOrderPaymentV2(ctx, mockPaymentRequestV2)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "make order payment failed with status: 404")
	})
}
