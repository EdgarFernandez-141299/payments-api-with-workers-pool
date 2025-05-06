package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources"
	config2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/mockrequest"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/mockresponse"
)

func createMockOrderServer(responseCode int, responseBody string, verifyHeaders func(r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifyHeaders(r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseCode)
		w.Write([]byte(responseBody))
	}))
}

func createOrderApiClient(mockServerURL, apiKey string) resources.DeunaOrderResourceIF {
	config := &config2.DeUnaApiConfig{
		URL:    mockServerURL,
		ApiKey: apiKey,
	}

	return NewDeUnaOrderHTTPClient(config)
}

const apiKey = "2jkmdj4-sd4j3k4-3j4k3j4-3j4k3j4"
const orderToken = "c5c5abc9-67a7-4a44-aa83-ee7e81d7052c"

func TestDeUnaOrderHTTPClient_CreateOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("successful order creation", func(t *testing.T) {
		expectedResponse := mockresponse.ExpectedSuccessResponse

		bodyRequest := mockrequest.CreateOrderRequest

		mockServer := createMockOrderServer(http.StatusOK, mockresponse.SuccessOrderResponseMock, func(r *http.Request) {
			assert.Equal(t, "/merchants/orders", r.URL.Path)
		})
		defer mockServer.Close()

		client := createOrderApiClient(mockServer.URL, apiKey)

		gotResponse, err := client.CreateOrder(ctx, bodyRequest)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, gotResponse, "Response does not match expected")
	})

	t.Run("fail card creation when is not status code 200", func(t *testing.T) {
		bodyRequest := mockrequest.CreateOrderRequest

		mockServer := createMockOrderServer(http.StatusInternalServerError, mockresponse.FailOrderResponseMock, func(r *http.Request) {
			assert.Equal(t, "/merchants/orders", r.URL.Path)
		})

		defer mockServer.Close()

		client := createOrderApiClient(mockServer.URL, apiKey)

		_, err := client.CreateOrder(ctx, bodyRequest)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create order request failed with status: 500")
	})

	t.Run("fail card creation when payment method is invalid", func(t *testing.T) {
		bodyRequest := mockrequest.CreateOrderRequest
		expectedResponse := mockresponse.FailOrderResponseMockWithInvalidPaymentMethod

		mockServer := createMockOrderServer(http.StatusConflict, expectedResponse, func(r *http.Request) {
			assert.Equal(t, "/merchants/orders", r.URL.Path)
		})
		defer mockServer.Close()

		client := createOrderApiClient(mockServer.URL, apiKey)

		_, err := client.CreateOrder(ctx, bodyRequest)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid payment method with status: 409")
	})
}

func TestDeUnaOrderHTTPClient_GetOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("successful order retrieval", func(t *testing.T) {
		expectedResponse := mockresponse.ExpectedSuccessResponse

		mockServer := createMockOrderServer(http.StatusOK, mockresponse.SuccessOrderResponseMock, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/merchants/orders/%s", orderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
		})
		defer mockServer.Close()

		client := createOrderApiClient(mockServer.URL, apiKey)

		gotResponse, err := client.GetOrder(ctx, orderToken)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, gotResponse, "Response does not match expected")
	})

	t.Run("failure when server returns non-200 status code", func(t *testing.T) {
		mockServer := createMockOrderServer(http.StatusInternalServerError, mockresponse.FailOrderResponseMock, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/merchants/orders/%s", orderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
		})
		defer mockServer.Close()

		client := createOrderApiClient(mockServer.URL, apiKey)

		_, err := client.GetOrder(ctx, orderToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get order failed with status: 500")
	})

	t.Run("failure when order does not exist", func(t *testing.T) {
		mockServer := createMockOrderServer(http.StatusInternalServerError, mockresponse.OrderNotFoundResponseMock, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/merchants/orders/%s", orderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
		})
		defer mockServer.Close()

		client := createOrderApiClient(mockServer.URL, apiKey)

		_, err := client.GetOrder(ctx, orderToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get order failed with status: 500")
	})
}

func TestDeUnaOrderHTTPClient_ExpireOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("successful order expiration", func(t *testing.T) {
		mockServer := createMockOrderServer(http.StatusOK, `{"success": true}`, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/merchants/orders/%s/expire", orderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)
		})
		defer mockServer.Close()

		client := createOrderApiClient(mockServer.URL, apiKey)

		err := client.ExpireOrder(ctx, orderToken)

		assert.NoError(t, err)
	})

	t.Run("failure when server returns non-200 status code", func(t *testing.T) {
		mockServer := createMockOrderServer(http.StatusInternalServerError, `{"error": "Internal server error"}`, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/merchants/orders/%s/expire", orderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)
		})
		defer mockServer.Close()

		client := createOrderApiClient(mockServer.URL, apiKey)

		err := client.ExpireOrder(ctx, orderToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expire order request failed with status: 500")
	})

	t.Run("failure when order does not exist", func(t *testing.T) {
		mockServer := createMockOrderServer(http.StatusNotFound, `{"error": "Order not found"}`, func(r *http.Request) {
			expectedPath := fmt.Sprintf("/merchants/orders/%s/expire", orderToken)
			assert.Equal(t, expectedPath, r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)
		})
		defer mockServer.Close()

		client := createOrderApiClient(mockServer.URL, apiKey)

		err := client.ExpireOrder(ctx, orderToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expire order request failed with status: 404")
	})
}
