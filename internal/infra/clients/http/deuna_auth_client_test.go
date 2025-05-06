package http

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"gitlab.com/clubhub.ai1/go-libraries/observability/mocks/apm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/request"
	config2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
)

func createMockAuthServer(responseCode int, responseBody string, verifyHeaders func(r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifyHeaders(r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseCode)
		w.Write([]byte(responseBody))
	}))
}

func createAuthApiClient(mockServerURL, apiKey string, t *testing.T) resources.DeunaAuthResourceIF {
	config := &config2.DeUnaApiConfig{
		URL:    mockServerURL,
		ApiKey: apiKey,
	}

	tracer := apm.NewTracer(t)
	tracer.On("GetTracer").Return(nil)

	return NewDeUnaAuthHTTPClient(config, tracer)
}

func generatePublicKeyPEM() (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", err
	}
	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(pubPEM), nil
}

func TestDeUnaAuthHTTPClient_AuthUser(t *testing.T) {
	validPublicKey, err := generatePublicKeyPEM()
	require.NoError(t, err, "should generate a valid public key without error")

	invalidPublicKey := `-----BEGIN PUBLIC KEY-----
invalidkeydata
-----END PUBLIC KEY-----`

	const apiKey = "test-api-key"
	testEmail := "test@example.com"

	t.Run("Successful authentication", func(t *testing.T) {
		mockResponse := `{"refresh_token":"refresh123","token":"auth123"}`

		headerVerification := func(r *http.Request) {
			assert.Equal(t, apiKey, r.Header.Get("x-api-key"))

			authMessage := r.Header.Get("X-Auth-Message")
			assert.NotEmpty(t, authMessage)

			_, err := base64.StdEncoding.DecodeString(authMessage)
			assert.NoError(t, err, "Auth message should be valid base64")

			assert.Equal(t, "/users/external-authorize", r.URL.Path)
			assert.Equal(t, "true", r.URL.Query().Get("load_profile_data"))

			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		}

		mockServer := createMockAuthServer(http.StatusOK, mockResponse, headerVerification)
		defer mockServer.Close()

		client := createAuthApiClient(mockServer.URL, apiKey, t)

		client.(*DeUnaAuthHTTPClient).Config.PublicKeyPEM = base64.RawStdEncoding.EncodeToString([]byte(validPublicKey))

		requestDTO := request.DeunaAuthUserRequestDTO{
			Email: testEmail,
		}

		ctx := context.Background()
		response, err := client.AuthUser(ctx, requestDTO)

		require.NoError(t, err)
		assert.Equal(t, "refresh123", response.RefreshToken)
		assert.Equal(t, "auth123", response.AuthToken)
	})

	t.Run("Authentication when user is not found", func(t *testing.T) {
		mockResponse := `{"error":"user not found"}`

		headerVerification := func(r *http.Request) {
			assert.Equal(t, apiKey, r.Header.Get("x-api-key"))
		}

		mockServer := createMockAuthServer(http.StatusNotFound, mockResponse, headerVerification)
		defer mockServer.Close()

		client := createAuthApiClient(mockServer.URL, apiKey, t)

		client.(*DeUnaAuthHTTPClient).Config.PublicKeyPEM = base64.RawStdEncoding.EncodeToString([]byte(validPublicKey))

		requestDTO := request.DeunaAuthUserRequestDTO{
			Email: testEmail,
		}

		ctx := context.Background()
		_, err := client.AuthUser(ctx, requestDTO)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("Authentication with HTTP server error", func(t *testing.T) {
		mockResponse := `{"error":"internal server error"}`

		headerVerification := func(r *http.Request) {
			assert.Equal(t, apiKey, r.Header.Get("x-api-key"))
		}

		mockServer := createMockAuthServer(http.StatusInternalServerError, mockResponse, headerVerification)
		defer mockServer.Close()

		client := createAuthApiClient(mockServer.URL, apiKey, t)

		client.(*DeUnaAuthHTTPClient).Config.PublicKeyPEM = base64.RawStdEncoding.EncodeToString([]byte(validPublicKey))

		requestDTO := request.DeunaAuthUserRequestDTO{
			Email: testEmail,
		}

		ctx := context.Background()
		_, err := client.AuthUser(ctx, requestDTO)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "authenticate user failed")
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("Authentication with malformed public key", func(t *testing.T) {
		client := createAuthApiClient("http://example.com", apiKey, t)

		client.(*DeUnaAuthHTTPClient).Config.PublicKeyPEM = invalidPublicKey

		requestDTO := request.DeunaAuthUserRequestDTO{
			Email: testEmail,
		}

		ctx := context.Background()
		_, err := client.AuthUser(ctx, requestDTO)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "error generating auth message")
		assert.Contains(t, strings.ToLower(err.Error()), "failed to parse")
	})

	t.Run("Authentication with empty public key", func(t *testing.T) {
		client := createAuthApiClient("http://example.com", apiKey, t)

		client.(*DeUnaAuthHTTPClient).Config.PublicKeyPEM = ""

		requestDTO := request.DeunaAuthUserRequestDTO{
			Email: testEmail,
		}

		ctx := context.Background()
		_, err := client.AuthUser(ctx, requestDTO)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "error generating auth message")
		assert.Contains(t, strings.ToLower(err.Error()), "failed to parse pem block")
	})

	t.Run("Authentication with HTTP client connection error", func(t *testing.T) {
		client := createAuthApiClient("http://invalid-server-that-does-not-exist.example", apiKey, t)

		client.(*DeUnaAuthHTTPClient).Config.PublicKeyPEM = base64.RawStdEncoding.EncodeToString([]byte(validPublicKey))

		requestDTO := request.DeunaAuthUserRequestDTO{
			Email: testEmail,
		}

		ctx := context.Background()
		_, err := client.AuthUser(ctx, requestDTO)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to authenticate user")
	})

}

func TestGenerateAuthMessage(t *testing.T) {
	testEmail := "test@example.com"

	t.Run("With empty public key", func(t *testing.T) {
		_, err := generateAuthMessage(testEmail, "")
		require.Error(t, err, "should return an error for empty public key")
		assert.Contains(t, err.Error(), "failed to parse PEM block containing the public key")
	})

	t.Run("With non-RSA key type", func(t *testing.T) {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		require.NoError(t, err, "should generate a valid ECDSA private key without error")

		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		require.NoError(t, err, "should marshal ECDSA public key without error")

		pubPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		})

		_, err = generateAuthMessage(testEmail, string(pubPEM))

		require.Error(t, err, "should return an error for invalid key type")
		assert.Contains(t, err.Error(), "key type is not RSA")
	})

	t.Run("With valid RSA key", func(t *testing.T) {
		validPublicKey, err := generatePublicKeyPEM()
		require.NoError(t, err, "should generate a valid public key without error")

		message, err := generateAuthMessage(testEmail, validPublicKey)

		require.NoError(t, err, "should generate auth message without error")
		require.NotEmpty(t, message, "auth message should not be empty")

		_, err = base64.StdEncoding.DecodeString(message)
		assert.NoError(t, err, "message should be a valid base64-encoded string")
	})
}
