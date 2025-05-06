package http

import (
	"context"
	"encoding/json"
	apm2 "gitlab.com/clubhub.ai1/go-libraries/observability/mocks/apm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
	config2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
)

func createMockUserServer(responseCode int, responseBody string, verifyRequest func(r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifyRequest(r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseCode)
		w.Write([]byte(responseBody))
	}))
}

func createUserApiClient(mockServerURL, apiKey string, t *testing.T) resources.DeUnaUserResourceIF {
	config := &config2.DeUnaApiConfig{
		URL:    mockServerURL,
		ApiKey: apiKey,
	}

	tracer := apm2.NewTracer(t)
	tracer.On("GetTracer").Return(nil)

	return NewDeUnaHTTPClient(config, tracer)
}

func TestDeUnaHTTPClient_CreateUser_Success(t *testing.T) {
	expectedResponse := response.CreatedUserResponse{
		Token:  "test-token",
		UserID: "user123",
	}
	responseBody, _ := json.Marshal(expectedResponse)

	var capturedRequest request.CreateUserRequestDTO
	mockServer := createMockUserServer(http.StatusCreated, string(responseBody), func(r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/users/register", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("x-api-key"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(body, &capturedRequest)
	},
	)
	defer mockServer.Close()

	client := createUserApiClient(mockServer.URL, "test-api-key", t)

	ctx := context.Background()
	userRequest := request.CreateUserRequestDTO{
		Email:            "test@example.com",
		FirstName:        "John",
		LastName:         "Doe",
		Phone:            "+51999999999",
		IdentityDocument: "12345678",
	}

	gotResponse, err := client.CreateUser(ctx, userRequest)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, gotResponse)
}

func TestDeUnaHTTPClient_CreateUser_ServerError(t *testing.T) {
	mockServer := createMockUserServer(
		http.StatusInternalServerError,
		`{"error": "Internal server error"}`,
		func(r *http.Request) {
			assert.Equal(t, "/users/register", r.URL.Path)
		},
	)
	defer mockServer.Close()

	client := createUserApiClient(mockServer.URL, "test-api-key", t)

	ctx := context.Background()
	userRequest := request.CreateUserRequestDTO{
		Email:            "test@example.com",
		FirstName:        "John",
		LastName:         "Doe",
		Phone:            "+51999999999",
		IdentityDocument: "12345678",
	}

	result, err := client.CreateUser(ctx, userRequest)

	assert.Error(t, err)
	assert.Equal(t, "error creating user", err.Error())
	assert.Empty(t, result.Token)
	assert.Empty(t, result.UserID)
}

func TestDeUnaHTTPClient_CreateUser_InvalidJSON(t *testing.T) {
	mockServer := createMockUserServer(
		http.StatusCreated,
		`{invalid json response}`,
		func(r *http.Request) {
			assert.Equal(t, "/users/register", r.URL.Path)
		},
	)
	defer mockServer.Close()

	client := createUserApiClient(mockServer.URL, "test-api-key", t)

	ctx := context.Background()
	userRequest := request.CreateUserRequestDTO{
		Email:            "test@example.com",
		FirstName:        "John",
		LastName:         "Doe",
		Phone:            "+51999999999",
		IdentityDocument: "12345678",
	}

	result, err := client.CreateUser(ctx, userRequest)

	assert.Error(t, err)
	assert.Empty(t, result.Token)
	assert.Empty(t, result.UserID)
}
