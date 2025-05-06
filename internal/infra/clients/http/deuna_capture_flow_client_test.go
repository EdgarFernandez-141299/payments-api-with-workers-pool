package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	config2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
)

func createMockCaptureServer(responseCode int, responseBody string, verifyRequest func(r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifyRequest(r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseCode)
		w.Write([]byte(responseBody))
	}))
}

func createCaptureApiClient(mockServerURL, apiKey string) *DeUnaCaptureFlowHTTPClient {
	config := &config2.DeUnaApiConfig{
		URL:    mockServerURL,
		ApiKey: apiKey,
	}

	client := NewDeUnaCaptureFlowHTTPClient(config)
	return client.(*DeUnaCaptureFlowHTTPClient)
}

func TestDeUnaCaptureFlowHTTPClient_Release(t *testing.T) {
	apiKey := "test-api-key"

	tests := []struct {
		name              string
		responseCode      int
		responseBody      string
		expectedError     bool
		expectedErrorMsg  string
		expectedResult    bool
		verifyRequestFunc func(r *http.Request)
		orderToken        string
		reason            string
	}{
		{
			name:         "successful release",
			responseCode: http.StatusNoContent,
			responseBody: "",
			verifyRequestFunc: func(r *http.Request) {
				assert.Equal(t, "/merchants/orders/test-token/void", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, apiKey, r.Header.Get(config2.DeunaApiKeyHeader))
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				var requestBody map[string]string
				err := json.NewDecoder(r.Body).Decode(&requestBody)
				assert.NoError(t, err)
				assert.Equal(t, "test reason", requestBody["reason"])
			},
			orderToken:     "test-token",
			reason:         "test reason",
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:         "conflict response",
			responseCode: http.StatusConflict,
			responseBody: `{"error": {"code": "EMA-6004", "description": "the provided payment method is not configured for merchant and store"}}`,
			verifyRequestFunc: func(r *http.Request) {
				assert.Equal(t, "/merchants/orders/test-token/void", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
			},
			orderToken:       "test-token",
			reason:           "test reason",
			expectedResult:   false,
			expectedError:    true,
			expectedErrorMsg: "Release request failed with status: 409",
		},
		{
			name:         "error response",
			responseCode: http.StatusInternalServerError,
			responseBody: `{"error": {"code": "EMA-3002", "description": "failed order request"}}`,
			verifyRequestFunc: func(r *http.Request) {
				assert.Equal(t, "/merchants/orders/test-token/void", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
			},
			orderToken:       "test-token",
			reason:           "test reason",
			expectedResult:   false,
			expectedError:    true,
			expectedErrorMsg: "Release request failed with status: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := createMockCaptureServer(tt.responseCode, tt.responseBody, tt.verifyRequestFunc)
			defer mockServer.Close()

			client := createCaptureApiClient(mockServer.URL, apiKey)

			result, err := client.Release(context.Background(), tt.orderToken, tt.reason)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestDeUnaCaptureFlowHTTPClient_Capture(t *testing.T) {
	apiKey := "test-api-key"

	tests := []struct {
		name              string
		responseCode      int
		responseBody      string
		expectedError     bool
		expectedErrorMsg  string
		expectedResult    bool
		verifyRequestFunc func(r *http.Request)
		orderToken        string
		amount            int64
	}{
		{
			name:         "successful capture",
			responseCode: http.StatusNoContent,
			responseBody: "",
			verifyRequestFunc: func(r *http.Request) {
				assert.Equal(t, "/merchants/orders/test-token/capture", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, apiKey, r.Header.Get(config2.DeunaApiKeyHeader))
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				var requestBody map[string]int64
				err := json.NewDecoder(r.Body).Decode(&requestBody)
				assert.NoError(t, err)
				assert.Equal(t, int64(1000), requestBody["amount"])
			},
			orderToken:     "test-token",
			amount:         1000,
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:         "conflict response",
			responseCode: http.StatusConflict,
			responseBody: `{"error": {"code": "EMA-6004", "description": "the provided payment method is not configured for merchant and store"}}`,
			verifyRequestFunc: func(r *http.Request) {
				assert.Equal(t, "/merchants/orders/test-token/capture", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
			},
			orderToken:       "test-token",
			amount:           1000,
			expectedResult:   false,
			expectedError:    true,
			expectedErrorMsg: "Capture request failed with status: 409",
		},
		{
			name:         "error response",
			responseCode: http.StatusInternalServerError,
			responseBody: `{"error": {"code": "EMA-3002", "description": "failed order request"}}`,
			verifyRequestFunc: func(r *http.Request) {
				assert.Equal(t, "/merchants/orders/test-token/capture", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
			},
			orderToken:       "test-token",
			amount:           1000,
			expectedResult:   false,
			expectedError:    true,
			expectedErrorMsg: "Capture request failed with status: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := createMockCaptureServer(tt.responseCode, tt.responseBody, tt.verifyRequestFunc)
			defer mockServer.Close()

			client := createCaptureApiClient(mockServer.URL, apiKey)

			result, err := client.Capture(context.Background(), tt.orderToken, tt.amount)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
