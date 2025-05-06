package capture_flow

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"go.temporal.io/sdk/client"
)

type mockWorkflow struct {
	mock.Mock
}

func (m *mockWorkflow) Call(ctx context.Context, referenceID string, cmd worfkflows.PaymentWorkflowInput) (client.WorkflowRun, error) {
	args := m.Called(ctx, referenceID, cmd)
	return args.Get(0).(client.WorkflowRun), args.Error(1)
}

func (m *mockWorkflow) SendProcessedSignal(ctx context.Context, paymentOrderID string, cmd worfkflows.PaymentProcessedSignal) error {
	args := m.Called(ctx, paymentOrderID, cmd)
	return args.Error(0)
}

func (m *mockWorkflow) SendCaptureFlowSignal(ctx context.Context, paymentOrderID string, cmd worfkflows.CompleteCaptureFlowSignal) error {
	args := m.Called(ctx, paymentOrderID, cmd)
	return args.Error(0)
}

func setupTestContext(body string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Configurar parámetros de autenticación
	authParams := auth.AuthParams{
		UserID:   "test-user",
		UserType: "test-type",
	}
	c.Set(auth.AuthParamsKey, authParams)

	return c
}

func TestPaymentCapture(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		workflowError  error
	}{
		{
			name:           "Captura exitosa",
			requestBody:    `{"reference_order_id": "123", "payment_id": "456"}`,
			expectedStatus: http.StatusCreated,
			workflowError:  nil,
		},
		{
			name:           "Error en el workflow",
			requestBody:    `{"reference_order_id": "123", "payment_id": "456"}`,
			expectedStatus: http.StatusInternalServerError,
			workflowError:  errors.New("error en el workflow"),
		},
		{
			name:           "Request inválido",
			requestBody:    `{"invalid": "data"}`,
			expectedStatus: http.StatusBadRequest,
			workflowError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestContext(tt.requestBody)

			mockWorkflow := new(mockWorkflow)
			if tt.expectedStatus != http.StatusBadRequest {
				mockWorkflow.On("SendCaptureFlowSignal", mock.Anything, mock.Anything, mock.Anything).Return(tt.workflowError)
			}

			handler := NewCaptureFlowHandler(mockWorkflow)
			err := handler.PaymentCapture(c)

			if tt.expectedStatus == http.StatusCreated {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, c.Response().Status)
			} else {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, he.Code)
			}
		})
	}
}

func TestPaymentRelease(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		workflowError  error
	}{
		{
			name:           "Liberación exitosa",
			requestBody:    `{"reference_order_id": "123", "payment_id": "456", "reason": "test"}`,
			expectedStatus: http.StatusCreated,
			workflowError:  nil,
		},
		{
			name:           "Error en el workflow",
			requestBody:    `{"reference_order_id": "123", "payment_id": "456", "reason": "test"}`,
			expectedStatus: http.StatusInternalServerError,
			workflowError:  errors.New("error en el workflow"),
		},
		{
			name:           "Request inválido",
			requestBody:    `{"invalid": "data"}`,
			expectedStatus: http.StatusBadRequest,
			workflowError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestContext(tt.requestBody)

			mockWorkflow := new(mockWorkflow)
			if tt.expectedStatus != http.StatusBadRequest {
				mockWorkflow.On("SendCaptureFlowSignal", mock.Anything, mock.Anything, mock.Anything).Return(tt.workflowError)
			}

			handler := NewCaptureFlowHandler(mockWorkflow)
			err := handler.PaymentRelease(c)

			if tt.expectedStatus == http.StatusCreated {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, c.Response().Status)
			} else {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, he.Code)
			}
		})
	}
}

func TestSendFlowSignal(t *testing.T) {
	tests := []struct {
		name          string
		referenceID   string
		paymentID     string
		action        enums.PaymentFlowActionEnum
		reason        string
		expectedError error
	}{
		{
			name:          "Señal de captura exitosa",
			referenceID:   "123",
			paymentID:     "456",
			action:        enums.CapturePayment,
			reason:        "",
			expectedError: nil,
		},
		{
			name:          "Señal de liberación exitosa",
			referenceID:   "123",
			paymentID:     "456",
			action:        enums.ReleasePayment,
			reason:        "test",
			expectedError: nil,
		},
		{
			name:          "Error en el workflow",
			referenceID:   "123",
			paymentID:     "456",
			action:        enums.CapturePayment,
			reason:        "",
			expectedError: errors.New("error en el workflow"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestContext("")

			mockWorkflow := new(mockWorkflow)
			mockWorkflow.On("SendCaptureFlowSignal", mock.Anything, mock.Anything, mock.Anything).Return(tt.expectedError)

			handler := &CaptureFlowHandler{
				workflow: mockWorkflow,
			}

			err := handler.sendFlowSignal(c, tt.referenceID, tt.paymentID, tt.action, tt.reason)

			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError.Error())
			}
		})
	}
}
