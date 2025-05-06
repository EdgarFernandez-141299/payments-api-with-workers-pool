package http

import (
	"context"
	"encoding/json"
	"gitlab.com/clubhub.ai1/go-libraries/observability/mocks/apm"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources/dto/response"
	config2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/mockresponse"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources/dto/request"
)

func createMockServer(responseCode int, responseBody string, verifyHeaders func(r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifyHeaders(r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseCode)
		w.Write([]byte(responseBody))
	}))
}

func createApiClient(mockServerURL, apiKey string, t *testing.T) resources.DeunaCardResourceIF {
	config := &config2.DeUnaApiConfig{
		URL:    mockServerURL,
		ApiKey: apiKey,
	}

	tracer := apm.NewTracer(t)
	tracer.On("GetTracer").Return(nil)

	return NewDeUnaCardHTTPClient(config, tracer)
}

func TestDeUnaCardHTTPClient_CreateCard(t *testing.T) {
	ctx := context.Background()
	apiKey := "2jkmdj4-sd4j3k4-3j4k3j4-3j4k3j4"

	t.Run("successful card creation", func(t *testing.T) {
		userID := "user123"
		expectedResponse := response.CardResponseDTO{Data: response.CardResponseDataDTO{
			ID:                      "12345",
			UserID:                  "67890",
			CardHolder:              "John Doe",
			CardHolderDni:           "987654321",
			Company:                 "Visa",
			LastFour:                "1234",
			FirstSix:                "123456",
			ExpirationDate:          "12/24",
			IsValid:                 true,
			IsExpired:               false,
			VerifiedBy:              "System",
			VerifiedWithTransaction: "abc123",
			VerifiedAt:              "2023-10-12T15:00:00Z",
			LastUsed:                "2023-10-15T09:30:00Z",
			CreatedAt:               "2023-01-01T08:00:00Z",
			UpdatedAt:               "2023-10-01T10:00:00Z",
			DeletedAt:               "",
			BankName:                "Bank of Somewhere",
			CountryIso:              "US",
			CardType:                "Credit",
			Source:                  "Online",
			ZipCode:                 "12345",
			Vault:                   "VaultService",
		}}

		mockServer := createMockServer(http.StatusCreated, mockresponse.SuccessCardResponseMock, func(r *http.Request) {
			assert.Equal(t, "/users/user123/cards", r.URL.Path)
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		})
		defer mockServer.Close()

		client := createApiClient(mockServer.URL, apiKey, t)

		body := request.CreateCardRequestDTO{
			CardHolder: "John Doe",
		}

		gotResponse, err := client.CreateCard(ctx, body, userID, "token123")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.Data, gotResponse, "Response does not match expected")
	})

	t.Run("Card already exists", func(t *testing.T) {
		userID := "user123"
		mockServer := createMockServer(http.StatusConflict, `{"error": "Card already exists"}`, func(r *http.Request) {
			assert.Equal(t, "/users/user123/cards", r.URL.Path)
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		})
		defer mockServer.Close()

		client := createApiClient(mockServer.URL, apiKey, t)

		body := request.CreateCardRequestDTO{
			CardHolder: "John Doe",
		}

		_, err := client.CreateCard(ctx, body, userID, "token123")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "the card already exist with status")
	})

	t.Run("Operation denied by anti fraud rules", func(t *testing.T) {
		userID := "user123"
		mockServer := createMockServer(http.StatusUnprocessableEntity, `{"error": "Operation denied by anti fraud rules"}`, func(r *http.Request) {
			assert.Equal(t, "/users/user123/cards", r.URL.Path)
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		})
		defer mockServer.Close()

		client := createApiClient(mockServer.URL, apiKey, t)

		body := request.CreateCardRequestDTO{
			CardHolder: "John Doe",
		}

		_, err := client.CreateCard(ctx, body, userID, "token123")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "operation denied by anti fraud rules")
	})

	t.Run("should fail when unsupported scheme", func(t *testing.T) {
		userID := "user123"

		client := createApiClient("api.co", apiKey, t)

		body := request.CreateCardRequestDTO{
			CardHolder: "John Doe",
		}

		_, err := client.CreateCard(ctx, body, userID, "token123")

		assert.Error(t, err)
		assert.Containsf(t, err.Error(), "unsupported protocol scheme", "Error message does not match expected")
	})

	t.Run("should fail when empty data is timeout", func(t *testing.T) {
		userID := "user123"
		mockServer := createMockServer(http.StatusCreated, mockresponse.EmptyDataMock, func(r *http.Request) {
			assert.Equal(t, "/users/user123/cards", r.URL.Path)
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		})

		defer mockServer.Close()

		client := createApiClient(mockServer.URL, apiKey, t)

		body := request.CreateCardRequestDTO{
			CardHolder: "John Doe",
		}

		_, err := client.CreateCard(ctx, body, userID, "token123")

		assert.Error(t, err, DeunaEmptyBody.Error())
	})

	t.Run("fail card creation when is not status code", func(t *testing.T) {
		userID := "user123"

		mockServer := createMockServer(http.StatusOK, mockresponse.SuccessCardResponseMock, func(r *http.Request) {
			assert.Equal(t, "/users/user123/cards", r.URL.Path)
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		})
		defer mockServer.Close()

		client := createApiClient(mockServer.URL, apiKey, t)

		body := request.CreateCardRequestDTO{
			CardHolder: "John Doe",
		}

		_, err := client.CreateCard(ctx, body, userID, "token123")

		assert.Error(t, err)
	})

	t.Run("fail card creation when is not status StatusBadRequest", func(t *testing.T) {
		userID := "user123"

		mockServer := createMockServer(http.StatusBadRequest, mockresponse.SuccessCardResponseMock, func(r *http.Request) {
			assert.Equal(t, "/users/user123/cards", r.URL.Path)
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		})
		defer mockServer.Close()

		client := createApiClient(mockServer.URL, apiKey, t)

		body := request.CreateCardRequestDTO{
			CardHolder: "John Doe",
		}

		_, err := client.CreateCard(ctx, body, userID, "token123")

		assert.Error(t, err)
	})

	t.Run("failed card creation", func(t *testing.T) {
		userID := "user123"

		mockServer := createMockServer(http.StatusNotFound, mockresponse.SuccessCardResponseMock, func(r *http.Request) {
			assert.Equal(t, "/users/user123/cards", r.URL.Path)
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		})
		defer mockServer.Close()

		client := createApiClient(mockServer.URL, apiKey, t)

		body := request.CreateCardRequestDTO{
			CardHolder: "John Doe",
		}

		_, err := client.CreateCard(ctx, body, userID, "token123")

		assert.Error(t, err)
	})
}

func TestDeUnaCardHTTPClient_DeleteCard(t *testing.T) {
	ctx := context.Background()

	t.Run("successful card deletion", func(t *testing.T) {
		cardID := "card123"
		userID := "user123"

		res := response.DeleteCardResponseDTO{
			Code:    "200",
			Message: "Card deleted successfully",
		}

		b, _ := json.Marshal(res)

		mockServer := createMockServer(http.StatusNoContent, string(b), func(r *http.Request) {
			assert.Equal(t, "/users/user123/cards/card123", r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		})
		defer mockServer.Close()

		client := createApiClient(mockServer.URL, "2jkmdj4-sd4j3k4-3j4k3j4-3j4k3j4", t)

		_, err := client.DeleteCard(ctx, request.DeleteCardRequestDTO{
			CardId: cardID,
			UserId: userID,
		}, "token123")

		assert.NoError(t, err)
	})

	t.Run("should fail when unsupported scheme", func(t *testing.T) {
		cardID := "card123"
		userID := "user123"

		client := createApiClient("api.co", "", t)

		_, err := client.DeleteCard(ctx, request.DeleteCardRequestDTO{
			CardId: cardID,
			UserId: userID,
		}, "token123")

		assert.Error(t, err)
		assert.Containsf(t, err.Error(), "unsupported protocol scheme", "Error message does not match expected")
	})

	t.Run("failed card deletion", func(t *testing.T) {
		cardID := "card123"
		userID := "user123"

		mockServer := createMockServer(http.StatusNotFound, `{"error": "Card not found"}`, func(r *http.Request) {
			assert.Equal(t, "/users/user123/cards/card123", r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		})
		defer mockServer.Close()

		client := createApiClient(mockServer.URL, "2jkmdj4-sd4j3k4-3j4k3j4-3j4k3j4", t)

		_, err := client.DeleteCard(ctx, request.DeleteCardRequestDTO{
			CardId: cardID,
			UserId: userID,
		}, "token123")

		assert.Error(t, err)
	})

	t.Run("fail when is not status code", func(t *testing.T) {
		cardID := "card123"
		userID := "user123"

		res := response.DeleteCardResponseDTO{
			Code:    "200",
			Message: "Card deleted successfully",
		}

		b, _ := json.Marshal(res)

		mockServer := createMockServer(http.StatusOK, string(b), func(r *http.Request) {
			assert.Equal(t, "/users/user123/cards/card123", r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		})
		defer mockServer.Close()

		client := createApiClient(mockServer.URL, "2jkmdj4-sd4j3k4-3j4k3j4-3j4k3j4", t)

		_, err := client.DeleteCard(ctx, request.DeleteCardRequestDTO{
			CardId: cardID,
			UserId: userID,
		}, "token123")

		assert.Error(t, err)
	})
}
