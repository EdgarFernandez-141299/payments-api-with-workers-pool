package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/mockresponse"
)

func createMockMemberServer(responseCode int, responseBody string, verifyRequest func(r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifyRequest(r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseCode)
		w.Write([]byte(responseBody))
	}))
}

func createMemberApiClient(baseURL string) resources.MemberAPIResourceIF {
	return &MemberHTTPClientImpl{
		instrument.NewInstrumentedClient(
			instrument.WithBaseUrl(baseURL),
			instrument.WithRequestTimeout(MemberApiTimeout),
		),
	}
}

func TestMemberHTTPClientImpl_GetMemberById_Success(t *testing.T) {
	memberID := "123"
	enterpriseID := "456"
	successResponse := mockresponse.SuccessMemberResponseMock

	mockServer := createMockMemberServer(http.StatusOK, successResponse, func(r *http.Request) {
		assert.Equal(t, "/api/v1/members/"+memberID, r.URL.Path)
		assert.Equal(t, enterpriseID, r.Header.Get(enterpriseHeader))
		assert.Equal(t, memberID, r.Header.Get(userIDHeader))
		assert.Equal(t, memberID, r.Header.Get(usernameHeader))
	})
	defer mockServer.Close()

	client := createMemberApiClient(mockServer.URL)

	result, err := client.GetMemberByID(context.Background(), memberID, enterpriseID)

	assert.NoError(t, err)
	assert.Equal(t, memberID, result.ID)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.Equal(t, "1990-01-01", result.BirthDate)
	assert.Len(t, result.Phones, 1)
	assert.Len(t, result.Emails, 1)
	assert.Equal(t, "john.doe@example.com", result.Emails[0].Email)
}

func TestMemberHTTPClientImpl_GetMemberById_MemberNotFound(t *testing.T) {
	memberID := "123"
	enterpriseID := "456"
	notFoundResponse := mockresponse.MemberNotFoundResponseMock

	mockServer := createMockMemberServer(http.StatusOK, notFoundResponse, func(r *http.Request) {})
	defer mockServer.Close()

	client := createMemberApiClient(mockServer.URL)

	_, err := client.GetMemberByID(context.Background(), memberID, enterpriseID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "member not found")
}

func TestMemberHTTPClientImpl_GetUserProfileInfo_Success(t *testing.T) {
	userID := "123"
	enterpriseID := "456"
	successResponse := mockresponse.SuccessUserProfileInfoResponseMock

	mockServer := createMockMemberServer(http.StatusOK, successResponse, func(r *http.Request) {
		assert.Equal(t, fmt.Sprintf("/api/v1/members/profile-info"), r.URL.Path)
		assert.Equal(t, userID, r.Header.Get(userIDHeader))
		assert.Equal(t, userID, r.URL.Query().Get("userID"))
	})
	defer mockServer.Close()

	client := createMemberApiClient(mockServer.URL)

	result, err := client.GetUserProfileInfo(context.Background(), userID, enterpriseID)

	assert.NoError(t, err)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.Equal(t, "john.doe@example.com", result.Email)
	assert.Equal(t, "456", result.EnterpriseID)
}

func TestMemberHTTPClientImpl_GetUserProfileInfo_UserNotFound(t *testing.T) {
	userID := "123"
	enterpriseID := "456"
	notFoundResponse := mockresponse.UserProfileInfoNotFoundResponseMock

	mockServer := createMockMemberServer(http.StatusOK, notFoundResponse, func(r *http.Request) {})
	defer mockServer.Close()

	client := createMemberApiClient(mockServer.URL)

	_, err := client.GetUserProfileInfo(context.Background(), userID, enterpriseID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "member not found")
}
