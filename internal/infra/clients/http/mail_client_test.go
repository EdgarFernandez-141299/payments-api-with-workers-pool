package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"
	commonAdapters "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common/adapters"
)

const (
	MailServiceTimeout = 10 * time.Second
)

func createMockMailServer(responseCode int, responseBody string, verifyRequest func(r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifyRequest(r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseCode)
		w.Write([]byte(responseBody))
	}))
}

func createMailApiClient(baseURL string) commonAdapters.MailAdapterIF {
	return &MailHTTPClientImpl{
		client: instrument.NewInstrumentedClient(
			instrument.WithBaseUrl(baseURL),
			instrument.WithRequestTimeout(MailServiceTimeout),
		),
	}
}

func TestMailHTTPClientImpl_Send_Success(t *testing.T) {
	mailRequest := commonAdapters.MailRequest{
		Recipient: "user1234",
		Content:   "lorem ipsum dolor sit amet",
		Title:     "Hello, world!",
		Metadata:  commonAdapters.Metadata{"key": "value"},
	}

	mockServer := createMockMailServer(http.StatusCreated, `{"success": true}`, func(r *http.Request) {
		assert.Equal(t, "/api/v1/notifications/broadcast", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		var requestBody commonAdapters.MailRequest
		err = json.Unmarshal(body, &requestBody)
		assert.NoError(t, err)

		assert.Equal(t, mailRequest.Recipient, requestBody.Recipient)
		assert.Equal(t, mailRequest.Content, requestBody.Content)
		assert.Equal(t, mailRequest.Title, requestBody.Title)
		assert.Equal(t, mailRequest.Metadata, requestBody.Metadata)
	})
	defer mockServer.Close()

	client := createMailApiClient(mockServer.URL)

	err := client.Send(context.Background(), mailRequest)

	assert.NoError(t, err)
}

func TestMailHTTPClientImpl_Send_Error(t *testing.T) {
	mailRequest := commonAdapters.MailRequest{
		Recipient: "user1234",
		Content:   "lorem ipsum dolor sit amet",
		Title:     "Hello, world!",
		Metadata:  commonAdapters.Metadata{"key": "value"},
	}

	mockServer := createMockMailServer(http.StatusInternalServerError, `{"error": "Internal server error"}`, func(r *http.Request) {
	})
	defer mockServer.Close()

	client := createMailApiClient(mockServer.URL)

	err := client.Send(context.Background(), mailRequest)

	assert.Error(t, err)
}

func TestMailHTTPClientImpl_Send_InvalidStatusCode(t *testing.T) {
	mailRequest := commonAdapters.MailRequest{
		Recipient: "user1234",
		Content:   "lorem ipsum dolor sit amet",
		Title:     "Hello, world!",
		Metadata:  commonAdapters.Metadata{"key": "value"},
	}

	mockServer := createMockMailServer(http.StatusOK, `{"success": true}`, func(r *http.Request) {
	})
	defer mockServer.Close()

	client := createMailApiClient(mockServer.URL)

	err := client.Send(context.Background(), mailRequest)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code from mail service: 200")
}

func TestMailHTTPClientImpl_Send_ResponseNotSuccess(t *testing.T) {
	mailRequest := commonAdapters.MailRequest{
		Recipient: "user1234",
		Content:   "lorem ipsum dolor sit amet",
		Title:     "Hello, world!",
		Metadata:  commonAdapters.Metadata{"key": "value"},
	}

	mockServer := createMockMailServer(http.StatusCreated, `{"success": false, "total_sent": 0}`, func(r *http.Request) {
	})
	defer mockServer.Close()

	client := createMailApiClient(mockServer.URL)

	err := client.Send(context.Background(), mailRequest)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mail service reported failure: success=false, total_sent=0")
}

func TestMailHTTPClientImpl_Send_RequestError(t *testing.T) {
	mailRequest := commonAdapters.MailRequest{
		Recipient: "user1234",
		Content:   "lorem ipsum dolor sit amet",
		Title:     "Hello, world!",
		Metadata:  commonAdapters.Metadata{"key": "value"},
	}

	// Use an invalid URL to force a request error
	client := createMailApiClient("http://invalid-server-that-does-not-exist")

	err := client.Send(context.Background(), mailRequest)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send mail notification")
}
