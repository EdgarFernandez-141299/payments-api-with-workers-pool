package card

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/constants"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/response"
	auth2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/infra/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/queries"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/usecases"
)

const (
	userID       = "test-user-id"
	enterpriseID = "test-enterprise-id"
	userType     = "member"
)

func TestCardHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		mockResponse response.CardResponse
		mockError    error
		expectedCode int
		expectedBody string
		httpError    *echo.HTTPError
	}{
		{
			name: "valid request",
			requestBody: request.CardRequest{
				UserID:         "member123",
				CardId:         "card123",
				Alias:          "test_card",
				Status:         "active",
				CardType:       "credit_card",
				CardBrand:      "visa",
				LastFour:       "2345",
				FirstSix:       "411111",
				ExpirationDate: "12/25",
				IsDefault:      false,
				IsRecurrent:    false,
			},
			mockResponse: response.CardResponse{
				ID:       "card123",
				Alias:    "test_card",
				LastFour: "2345",
				Brand:    "visa",
			},
			mockError:    nil,
			expectedCode: http.StatusCreated,
			expectedBody: `{"message":"Card created successfully","card_id":"card123","alias":"test_card","last_four":"2345","brand":"visa","status":"active","is_recurrent":false}`,
		},
		{
			name:         "invalid request body",
			requestBody:  []byte(`invalid`),
			mockResponse: response.CardResponse{},
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Unmarshal type error: expected=object"}`,
			httpError: echo.NewHTTPError(http.StatusBadRequest,
				"code=400, message=Syntax error: offset=1, "+
					"error=invalid character 'i' looking for beginning of value, "+
					"internal=invalid character 'i' looking for beginning of value"),
		},
		{
			name: "validation error",
			requestBody: request.CardRequest{
				UserID:         "member123",
				CardId:         "",
				Alias:          "test_card",
				Status:         "active",
				CardType:       "credit_card",
				CardBrand:      "visa",
				LastFour:       "2345",
				FirstSix:       "411111",
				ExpirationDate: "12/25",
				IsDefault:      false,
				IsRecurrent:    false,
			},
			mockResponse: response.CardResponse{},
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"CardId is required"}`,
			httpError: echo.NewHTTPError(http.StatusBadRequest,
				"Key: 'CardRequest.CardId' Error:Field validation for 'CardId' failed on the 'required' tag"),
		},
		{
			name: "usecase error",
			requestBody: request.CardRequest{
				UserID:         "member123",
				CardId:         "card123",
				Alias:          "test_card",
				Status:         "active",
				CardType:       "credit_card",
				CardBrand:      "visa",
				LastFour:       "2345",
				FirstSix:       "411111",
				ExpirationDate: "12/25",
				IsDefault:      false,
				IsRecurrent:    false,
			},
			mockResponse: response.CardResponse{},
			mockError:    errors.New("internal server error"),
			expectedCode: http.StatusInternalServerError,
			httpError:    echo.NewHTTPError(http.StatusInternalServerError, "internal server error"),
			expectedBody: `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoInstance := echo.New()
			mockUsecase := usecases.NewCardUsecaseIF(t)
			getCardUseCase := queries.NewGetCardUsecaseIF(t)
			handler := NewCardHandler(mockUsecase, getCardUseCase)

			var requestBody []byte
			if r, ok := tt.requestBody.([]byte); ok {
				requestBody = r
			} else {
				requestBody, _ = json.Marshal(tt.requestBody)
			}

			if tt.mockError == nil && tt.mockResponse != (response.CardResponse{}) {
				mockUsecase.On("CreateCard", mock.Anything, mock.Anything, enterpriseID, mock.Anything).Return(tt.mockResponse, nil)
			} else if tt.mockError != nil {
				mockUsecase.On("CreateCard", mock.Anything, mock.Anything, enterpriseID, mock.Anything).Return(tt.mockResponse, tt.mockError)
			}

			authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)

			ctx := context.Background()
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/cards", bytes.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			responseRecorder := httptest.NewRecorder()
			c := echoInstance.NewContext(req, responseRecorder)
			c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())

			err := handler.CreateCard(c)

			if tt.httpError != nil {
				assert.Equal(t, tt.httpError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, responseRecorder.Code)
			}
		})
	}
}

func TestCardHandler_GetCardsByUserID(t *testing.T) {
	tests := []struct {
		name         string
		memberID     string
		mockResponse []response.CardResponse
		mockError    error
		expectedCode int
		expectedBody string
		httpError    *echo.HTTPError
	}{
		{
			name:     "valid request",
			memberID: "valid-member-id",
			mockResponse: []response.CardResponse{
				{
					ID:             "card123",
					CardTokenID:    "cardtoken123",
					Brand:          "visa",
					Alias:          "alias",
					LastFour:       "1111",
					ExpirationDate: "12/25",
					IsDefault:      false,
					IsRecurrent:    true,
					CardType:       "credit_card",
				},
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: `[{"alias":"alias", "brand":"visa", "card_id":"card123", "card_token_id":"cardtoken123", "expiration_date":"12/25", "is_default":false, "is_recurrent":true, "last_four":"1111", "card_type":"credit_card"}]`,
		},
		{
			name:         "non-existent user",
			memberID:     "nonexistent-member-id",
			mockResponse: make([]response.CardResponse, 0),
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: `[]`,
		},
		{
			name:         "server error",
			memberID:     "valid-member-id",
			mockResponse: make([]response.CardResponse, 0),
			mockError:    errors.New("internal server error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"internal server error"}`,
			httpError:    echo.NewHTTPError(http.StatusInternalServerError, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoInstance := echo.New()
			mockGetCardUseCase := queries.NewGetCardUsecaseIF(t)
			mockUsecase := usecases.NewCardUsecaseIF(t)
			handler := NewCardHandler(mockUsecase, mockGetCardUseCase)

			if tt.mockError == nil {
				mockGetCardUseCase.On("GetCardsByUserID", mock.Anything, tt.memberID, enterpriseID).
					Return(tt.mockResponse, nil)
			} else {
				mockGetCardUseCase.On("GetCardsByUserID", mock.Anything, tt.memberID, enterpriseID).
					Return(tt.mockResponse, tt.mockError)
			}

			authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)

			ctx := context.Background()
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/cards/"+tt.memberID, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			responseRecorder := httptest.NewRecorder()
			c := echoInstance.NewContext(req, responseRecorder)
			c.SetParamNames("user_id")
			c.SetParamValues(tt.memberID)
			c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())

			err := handler.GetCardsByUserID(c)

			if tt.httpError != nil {
				assert.Equal(t, tt.httpError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, responseRecorder.Code)
				assert.JSONEq(t, tt.expectedBody, responseRecorder.Body.String())
			}
		})
	}
}

func TestCardHandler_Delete(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		mockResponse response.DeleteCardResponse
		mockError    error
		expectedCode int
		expectedBody string
		httpError    *echo.HTTPError
	}{
		{
			name: "valid request",
			requestBody: request.DeleteCardRequest{
				UserID: "member123",
				CardID: "card123",
			},
			mockResponse: response.DeleteCardResponse{
				Status: "success",
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"", "status":"success"}`,
		},
		{
			name: "non-existent card",
			requestBody: request.DeleteCardRequest{
				UserID: "member123",
				CardID: "card123",
			},
			mockResponse: response.DeleteCardResponse{},
			mockError:    errors.New("card not found"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"card not found"}`,
			httpError:    echo.NewHTTPError(http.StatusInternalServerError, "card not found"),
		},
		{
			name: "server error",
			requestBody: request.DeleteCardRequest{
				UserID: "member123",
				CardID: "card123",
			},
			mockResponse: response.DeleteCardResponse{},
			mockError:    errors.New("internal server error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: ``,
			httpError:    echo.NewHTTPError(http.StatusInternalServerError, "internal server error"),
		},
		{
			name: "invalid request body with empty card id",
			requestBody: request.DeleteCardRequest{
				UserID: "member123",
			},
			mockResponse: response.DeleteCardResponse{},
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"CardId is required"}`,
			httpError: echo.NewHTTPError(http.StatusBadRequest,
				"Key: 'DeleteCardRequest.CardID' Error:Field validation for 'CardID' failed on the 'required' tag"),
		},
		{
			name:         "invalid request body",
			requestBody:  []byte(`invalid`),
			mockResponse: response.DeleteCardResponse{},
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Unmarshal type error: expected=object"}`,
			httpError: echo.NewHTTPError(http.StatusBadRequest,
				"code=400, message=Syntax error: offset=1, "+
					"error=invalid character 'i' looking for beginning of value, "+
					"internal=invalid character 'i' looking for beginning of value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoInstance := echo.New()
			mockUsecase := usecases.NewDeleteCardUseCaseIF(t)

			handler := NewDeleteCardHandler(mockUsecase)

			var requestBody []byte
			if r, ok := tt.requestBody.([]byte); ok {
				requestBody = r
			} else {
				requestBody, _ = json.Marshal(tt.requestBody)
			}

			if tt.mockError == nil && tt.mockResponse != (response.DeleteCardResponse{}) {
				mockUsecase.On("DeleteCard", mock.Anything, mock.Anything, enterpriseID, mock.Anything).
					Return(tt.mockResponse, nil)
			} else if tt.mockError != nil {
				mockUsecase.On("DeleteCard", mock.Anything, mock.Anything, enterpriseID, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)

			ctx := context.Background()
			req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, "/cards/", bytes.NewReader(requestBody))

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			responseRecorder := httptest.NewRecorder()
			c := echoInstance.NewContext(req, responseRecorder)
			c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())

			err := handler.DeleteCard(c)

			if tt.httpError != nil {
				assert.Equal(t, tt.httpError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, responseRecorder.Code)
				assert.JSONEq(t, tt.expectedBody, responseRecorder.Body.String())
			}
		})
	}
}

func TestCardHandler_TriggerExpiringSoonNotifications(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  any
		mockResponse response.NotificationCardExpiringSoonResponseDTO
		mockError    error
		expectedCode int
		expectedBody string
		httpError    *echo.HTTPError
	}{
		{
			name: "valid request",
			requestBody: request.NotificationCardExpiringSoonRequestDTO{
				NotificationChannels: []constants.NotificationChannel{constants.EmailChannel},
			},
			mockResponse: response.NotificationCardExpiringSoonResponseDTO{
				Message: "success",
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"success"}`,
		},
		{
			name: "server error",
			requestBody: request.NotificationCardExpiringSoonRequestDTO{
				NotificationChannels: []constants.NotificationChannel{constants.EmailChannel},
			},
			mockResponse: response.NotificationCardExpiringSoonResponseDTO{},
			mockError:    errors.New("internal server error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"message":"internal server error"}`,
			httpError:    echo.NewHTTPError(http.StatusInternalServerError, "internal server error"),
		},
		{
			name:         "invalid request body with empty notification channels",
			requestBody:  request.NotificationCardExpiringSoonRequestDTO{},
			mockResponse: response.NotificationCardExpiringSoonResponseDTO{},
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"notification_channels is required"}`,
			httpError:    echo.NewHTTPError(http.StatusBadRequest, "notification_channels is required"),
		},
		{
			name:         "invalid request body",
			requestBody:  []byte(`invalid`),
			mockResponse: response.NotificationCardExpiringSoonResponseDTO{},
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Unmarshal type error: expected=object"}`,
			httpError: echo.NewHTTPError(http.StatusBadRequest,
				"code=400, message=Syntax error: offset=1, "+
					"error=invalid character 'i' looking for beginning of value, "+
					"internal=invalid character 'i' looking for beginning of value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			echoInstance := echo.New()
			mockUsecase := usecases.NewCardUsecaseIF(t)
			handler := NewCardHandler(mockUsecase, nil)

			var requestBody []byte
			if r, ok := tt.requestBody.([]byte); ok {
				requestBody = r
			} else {
				requestBody, _ = json.Marshal(tt.requestBody)
			}

			if tt.mockError == nil && tt.mockResponse != (response.NotificationCardExpiringSoonResponseDTO{}) {
				mockUsecase.On("TriggerCardExpiringSoonNotifications", mock.Anything, mock.Anything).
					Return(tt.mockResponse, nil)
			} else if tt.mockError != nil {
				mockUsecase.On("TriggerCardExpiringSoonNotifications", mock.Anything, mock.Anything).
					Return(tt.mockResponse, tt.mockError)
			}

			authParamsFixture := *auth2.NewAuthParamsFixture(userID, userType, enterpriseID)

			ctx := context.Background()
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/cards/trigger-expiring-soon-notifications", bytes.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			responseRecorder := httptest.NewRecorder()
			c := echoInstance.NewContext(req, responseRecorder)
			c.Set(auth.AuthParamsKey, authParamsFixture.GetParams())

			err := handler.TriggerExpiringSoonNotifications(c)

			if tt.httpError != nil {
				assert.Equal(t, tt.httpError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, responseRecorder.Code)
				assert.JSONEq(t, tt.expectedBody, responseRecorder.Body.String())
			}
		})
	}
}
