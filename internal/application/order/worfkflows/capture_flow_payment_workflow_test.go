package worfkflows

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/go-libraries/saga/workflow"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows/activities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

func TestCaptureFlowWorkflow(t *testing.T) {
	referenceOrderID := "order-123"
	paymentID := "payment-456"

	t.Run("debe ejecutar correctamente la captura cuando se recibe la señal de captura", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		env := ts.Env

		defer ts.AfterTest()

		capturePaymentActivity := new(activities.CapturePaymentActivity)
		releasePaymentActivity := new(activities.ReleasePaymentActivity)

		env.RegisterActivity(capturePaymentActivity)
		env.RegisterActivity(releasePaymentActivity)

		env.OnActivity(
			capturePaymentActivity.CapturePayment,
			mock.Anything,
			referenceOrderID,
			paymentID,
		).Return(nil).Once()

		env.RegisterDelayedCallback(func() {
			env.SignalWorkflow(CaptureFlowSignalName, CompleteCaptureFlowSignal{
				OrderReferecenId: referenceOrderID,
				PaymentID:        paymentID,
				Action:           enums.CapturePayment,
			})
		}, 10*time.Millisecond)

		env.ExecuteWorkflow(CaptureFlowWorkflow, referenceOrderID, paymentID)

		ts.NoError(env.GetWorkflowError())
	})

	t.Run("debe ejecutar correctamente la liberación cuando se recibe la señal de liberación", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		env := ts.Env

		defer ts.AfterTest()

		capturePaymentActivity := new(activities.CapturePaymentActivity)
		releasePaymentActivity := new(activities.ReleasePaymentActivity)

		env.RegisterActivity(capturePaymentActivity)
		env.RegisterActivity(releasePaymentActivity)

		env.OnActivity(
			releasePaymentActivity.ReleasePayment,
			mock.Anything,
			referenceOrderID,
			paymentID,
			"razón de liberación",
		).Return(nil).Once()

		env.RegisterDelayedCallback(func() {
			env.SignalWorkflow(CaptureFlowSignalName, CompleteCaptureFlowSignal{
				OrderReferecenId: referenceOrderID,
				PaymentID:        paymentID,
				Action:           enums.ReleasePayment,
				Reason:           "razón de liberación",
			})
		}, 10*time.Millisecond)

		env.ExecuteWorkflow(CaptureFlowWorkflow, referenceOrderID, paymentID)

		ts.NoError(env.GetWorkflowError())
	})

	t.Run("debe ejecutar la liberación automáticamente cuando se alcanza el timeout", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		env := ts.Env

		defer ts.AfterTest()

		capturePaymentActivity := new(activities.CapturePaymentActivity)
		releasePaymentActivity := new(activities.ReleasePaymentActivity)

		env.RegisterActivity(capturePaymentActivity)
		env.RegisterActivity(releasePaymentActivity)

		env.OnActivity(
			releasePaymentActivity.ReleasePayment,
			mock.Anything,
			referenceOrderID,
			paymentID,
			TimeoutReason,
		).Return(nil).Once()

		env.ExecuteWorkflow(CaptureFlowWorkflow, referenceOrderID, paymentID)

		ts.NoError(env.GetWorkflowError())
	})

	t.Run("debe manejar el error cuando la actividad de captura falla", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		env := ts.Env

		defer ts.AfterTest()

		capturePaymentActivity := new(activities.CapturePaymentActivity)
		releasePaymentActivity := new(activities.ReleasePaymentActivity)

		env.RegisterActivity(capturePaymentActivity)
		env.RegisterActivity(releasePaymentActivity)

		expectedError := assert.AnError

		env.OnActivity(
			capturePaymentActivity.CapturePayment,
			mock.Anything,
			referenceOrderID,
			paymentID,
		).Return(expectedError).Times(2)

		env.RegisterDelayedCallback(func() {
			env.SignalWorkflow(CaptureFlowSignalName, CompleteCaptureFlowSignal{
				OrderReferecenId: referenceOrderID,
				PaymentID:        paymentID,
				Action:           enums.CapturePayment,
			})
		}, 10*time.Millisecond)

		env.ExecuteWorkflow(CaptureFlowWorkflow, referenceOrderID, paymentID)

		gotErr := env.GetWorkflowError()
		ts.Error(gotErr)
		assert.Contains(t, gotErr.Error(), expectedError.Error())
	})

	t.Run("debe manejar el error cuando la actividad de liberación falla", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		env := ts.Env

		defer ts.AfterTest()

		capturePaymentActivity := new(activities.CapturePaymentActivity)
		releasePaymentActivity := new(activities.ReleasePaymentActivity)

		env.RegisterActivity(capturePaymentActivity)
		env.RegisterActivity(releasePaymentActivity)

		expectedError := assert.AnError

		env.OnActivity(
			releasePaymentActivity.ReleasePayment,
			mock.Anything,
			referenceOrderID,
			paymentID,
			"razón de liberación",
		).Return(expectedError).Times(2)

		env.RegisterDelayedCallback(func() {
			env.SignalWorkflow(CaptureFlowSignalName, CompleteCaptureFlowSignal{
				OrderReferecenId: referenceOrderID,
				PaymentID:        paymentID,
				Action:           enums.ReleasePayment,
				Reason:           "razón de liberación",
			})
		}, 10*time.Millisecond)

		env.ExecuteWorkflow(CaptureFlowWorkflow, referenceOrderID, paymentID)

		gotErr := env.GetWorkflowError()
		ts.Error(gotErr)
		assert.Contains(t, gotErr.Error(), expectedError.Error())
	})
}
