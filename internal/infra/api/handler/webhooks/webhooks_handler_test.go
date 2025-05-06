package webhooks

import (
	"bytes"
	"encoding/json"
	worfkflows2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/worfkflows"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	dto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/webhooks/dto/request"
)

func TestDeunaWebhookNotifyOrder(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    dto.WebhookOrderDTO
		setupMock      func(workflow *worfkflows.PaymentOrderWorkflow)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful webhook processing",
			requestBody: dto.WebhookOrderDTO{
				Order: dto.Order{
					OrderID: "test-order-123",
					Status:  "completed",
					Payment: dto.Payment{
						Data: dto.PaymentData{
							Status:            "PROCESSED",
							AuthorizationCode: "auth123",
							Reason:            "success",
							Amount: dto.MoneyAmount{
								Amount:   10000,
								Currency: "USD",
							},
						},
					},
					Metadata: map[string]interface{}{
						"key": "value",
					},
				},
			},
			setupMock: func(m *worfkflows.PaymentOrderWorkflow) {
				m.On("SendProcessedSignal", mock.Anything, "test-order-123", worfkflows2.PaymentProcessedSignal{
					AuthorizationCode:   "auth123",
					Status:              enums.PaymentProcessed,
					OrderStatusString:   "completed",
					OrderID:             "test-order",
					PaymentID:           "123",
					PaymentStatusString: "PROCESSED",
					PaymentReason:       "success",
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid payment status",
			requestBody: dto.WebhookOrderDTO{
				Order: dto.Order{
					OrderID: "test-order-123",
					Status:  "completed",
					Payment: dto.Payment{
						Data: dto.PaymentData{
							Status: "invalid_status",
						},
					},
				},
			},
			setupMock:      func(m *worfkflows.PaymentOrderWorkflow) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid payment status",
		},
		{
			name: "invalid request body",
			requestBody: dto.WebhookOrderDTO{
				Order: dto.Order{
					OrderID: "test-order-123",
					Status:  "completed",
					Payment: dto.Payment{
						Data: dto.PaymentData{
							Status: "PROCESSED",
						},
					},
				},
			},
			setupMock:      func(m *worfkflows.PaymentOrderWorkflow) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid character 'i' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			mockWorkflow := worfkflows.NewPaymentOrderWorkflow(t)
			tt.setupMock(mockWorkflow)

			handler := &WebhookHandler{
				workflow: mockWorkflow,
			}

			// Create request
			var body []byte
			if tt.name == "invalid request body" {
				body = []byte("invalid json")
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}
			req := httptest.NewRequest(http.MethodPost, "/webhooks/deuna", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute
			err := handler.DeunaWebhookNotifyOrder(c)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				httpError, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, httpError.Code)
				assert.Contains(t, httpError.Message, tt.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockWorkflow.AssertExpectations(t)
		})
	}
}
