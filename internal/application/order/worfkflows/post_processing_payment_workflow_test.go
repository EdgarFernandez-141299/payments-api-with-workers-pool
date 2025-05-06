package worfkflows

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/go-libraries/saga/workflow"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows/activities"
	postpaymentUseCase "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/post_payment"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

func TestPostProcessingOrderPaymentWorkflow(t *testing.T) {
	t.Run("must success when payment order is sent and received a signal", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		env := ts.Env

		orderID := "32422-asdfer23-242asd2"
		paymentOrderID := "234dfc2-34523-dsd234"
		authorizationCode := "AuthCode1"
		orderStatusString := enums.PaymentProcessed.String()
		paymentReason := "reason-1"

		paymentReceipt := entities.PaymentReceipt{PaymentID: paymentOrderID, ReferenceOrderID: orderID}

		defer ts.AfterTest()

		postProcessingActivity := new(activities.PostProcessingPaymentOrderActivity)
		notifyOrderChangeActivity := new(activities.NotifyOrderChangeActivity)
		GeneratePaymentActivity := new(activities.GeneratePaymentReceiptActivity)

		env.RegisterActivity(postProcessingActivity)
		env.RegisterActivity(notifyOrderChangeActivity)
		env.RegisterActivity(GeneratePaymentActivity)

		env.RegisterWorkflow(PostProcessingPaymentWorkflow)
		env.RegisterWorkflow(CaptureFlowWorkflow)

		expectedCommand := postpaymentUseCase.PostProcessingPaymentOrderCommand{
			ReferenceOrderID:  orderID,
			PaymentID:         paymentOrderID,
			Status:            enums.PaymentProcessed,
			AuthorizationCode: authorizationCode,
			OrderStatusString: orderStatusString,
			PaymentReason:     paymentReason,
		}

		env.OnActivity(postProcessingActivity.PostProcessingPaymentOrder,
			mock.Anything,
			expectedCommand,
		).
			Return(enums.Autocapture, nil).
			Once()

		env.OnActivity(notifyOrderChangeActivity.NotifyOrderChange,
			mock.Anything,
			activities.NotifyOrderChangeParams{OrderID: orderID, PaymentID: paymentOrderID},
		).
			Return(nil)

		env.OnActivity(GeneratePaymentActivity.GenerateReceipt,
			mock.Anything,
			orderID,
			paymentOrderID,
		).
			Return(paymentReceipt, nil).
			Once()

		env.RegisterDelayedCallback(func() {
			env.SignalWorkflow(PaymentIntegrationProcessedSignalName, PaymentProcessedSignal{
				Status:              enums.PaymentProcessed,
				AuthorizationCode:   authorizationCode,
				OrderStatusString:   orderStatusString,
				OrderID:             orderID,
				PaymentStatusString: orderStatusString,
				PaymentReason:       paymentReason,
			})
		}, 10*time.Millisecond)

		env.ExecuteWorkflow(PostProcessingPaymentWorkflow, orderID, paymentOrderID)

		ts.NoError(env.GetWorkflowError())
	})

	t.Run("must success when payment order is sent and received a signal and failed is payment si status FAILED", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		env := ts.Env

		orderID := "32422-asdfer23-242asd2"
		paymentOrderID := "234dfc2-34523-dsd234"
		authorizationCode := "AuthCode1"
		orderStatusString := enums.PaymentFailed.String()
		paymentReason := "reason-1"

		defer ts.AfterTest()

		expectedCommand := postpaymentUseCase.PostProcessingPaymentOrderCommand{
			ReferenceOrderID:  orderID,
			PaymentID:         paymentOrderID,
			Status:            enums.PaymentFailed,
			AuthorizationCode: authorizationCode,
			OrderStatusString: orderStatusString,
			PaymentReason:     paymentReason,
		}

		postProcessingActivity := new(activities.PostProcessingPaymentOrderActivity)
		notifyOrderChangeActivity := new(activities.NotifyOrderChangeActivity)
		GeneratePaymentActivity := new(activities.GeneratePaymentReceiptActivity)

		env.RegisterWorkflow(PostProcessingPaymentWorkflow)
		env.RegisterWorkflow(CaptureFlowWorkflow)

		env.RegisterActivity(postProcessingActivity)
		env.RegisterActivity(notifyOrderChangeActivity)
		env.RegisterActivity(GeneratePaymentActivity)

		env.OnActivity(postProcessingActivity.PostProcessingPaymentOrder,
			mock.Anything,
			expectedCommand,
		).
			Return(enums.Autocapture, nil).
			Once()

		env.OnActivity(notifyOrderChangeActivity.NotifyOrderChange,
			mock.Anything,
			activities.NotifyOrderChangeParams{OrderID: orderID, PaymentID: paymentOrderID},
		).
			Return(nil)

		env.RegisterDelayedCallback(func() {
			env.SignalWorkflow(PaymentIntegrationProcessedSignalName, PaymentProcessedSignal{
				Status:              enums.PaymentFailed,
				AuthorizationCode:   authorizationCode,
				OrderStatusString:   orderStatusString,
				OrderID:             orderID,
				PaymentStatusString: orderStatusString,
				PaymentReason:       paymentReason,
			})
		}, 10*time.Millisecond)

		env.ExecuteWorkflow(PostProcessingPaymentWorkflow, orderID, paymentOrderID)

		ts.NoError(env.GetWorkflowError())
	})

	t.Run("must success when payment order is sent and received a signal and capture flow is triggered", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		env := ts.Env

		orderID := "32422-asdfer23-242asd2"
		paymentOrderID := "234dfc2-34523-dsd234"
		authorizationCode := "AuthCode1"
		orderStatusString := enums.PaymentProcessed.String()
		paymentReason := "reason-1"

		paymentReceipt := entities.PaymentReceipt{PaymentID: paymentOrderID, ReferenceOrderID: orderID}

		defer ts.AfterTest()

		postProcessingActivity := new(activities.PostProcessingPaymentOrderActivity)
		notifyOrderChangeActivity := new(activities.NotifyOrderChangeActivity)
		GeneratePaymentActivity := new(activities.GeneratePaymentReceiptActivity)

		env.RegisterActivity(postProcessingActivity)
		env.RegisterActivity(notifyOrderChangeActivity)
		env.RegisterActivity(GeneratePaymentActivity)

		env.RegisterWorkflow(PostProcessingPaymentWorkflow)
		env.RegisterWorkflow(CaptureFlowWorkflow)

		expectedCommand := postpaymentUseCase.PostProcessingPaymentOrderCommand{
			ReferenceOrderID:  orderID,
			PaymentID:         paymentOrderID,
			Status:            enums.PaymentProcessed,
			AuthorizationCode: authorizationCode,
			OrderStatusString: orderStatusString,
			PaymentReason:     paymentReason,
		}

		env.OnActivity(postProcessingActivity.PostProcessingPaymentOrder,
			mock.Anything,
			expectedCommand,
		).
			Return(enums.Capture, nil).
			Once()

		env.OnActivity(notifyOrderChangeActivity.NotifyOrderChange,
			mock.Anything,
			activities.NotifyOrderChangeParams{OrderID: orderID, PaymentID: paymentOrderID},
		).
			Return(nil)

		env.OnActivity(GeneratePaymentActivity.GenerateReceipt,
			mock.Anything,
			orderID,
			paymentOrderID,
		).
			Return(paymentReceipt, nil).
			Once()

		env.OnWorkflow(CaptureFlowWorkflow, mock.Anything, orderID, paymentOrderID).
			Return(nil).
			Once()

		env.RegisterDelayedCallback(func() {
			env.SignalWorkflow(PaymentIntegrationProcessedSignalName, PaymentProcessedSignal{
				Status:              enums.PaymentProcessed,
				AuthorizationCode:   authorizationCode,
				OrderStatusString:   orderStatusString,
				OrderID:             orderID,
				PaymentStatusString: orderStatusString,
				PaymentReason:       paymentReason,
			})
		}, 10*time.Millisecond)

		env.ExecuteWorkflow(PostProcessingPaymentWorkflow, orderID, paymentOrderID)

		ts.NoError(env.GetWorkflowError())
	})

	t.Run("must handle error when notify order change fails", func(t *testing.T) {
		ts := new(workflow.WorkflowUnitTestSuite)
		ts.SetupTest(t)
		env := ts.Env

		orderID := "32422-asdfer23-242asd2"
		paymentOrderID := "234dfc2-34523-dsd234"
		authorizationCode := "AuthCode1"
		orderStatusString := enums.PaymentProcessed.String()
		paymentReason := "reason-1"

		defer ts.AfterTest()

		postProcessingActivity := new(activities.PostProcessingPaymentOrderActivity)
		notifyOrderChangeActivity := new(activities.NotifyOrderChangeActivity)
		generatePaymentActivity := new(activities.GeneratePaymentReceiptActivity)

		env.RegisterActivity(postProcessingActivity)
		env.RegisterActivity(notifyOrderChangeActivity)
		env.RegisterActivity(generatePaymentActivity)

		env.RegisterWorkflow(PostProcessingPaymentWorkflow)
		env.RegisterWorkflow(CaptureFlowWorkflow)

		expectedCommand := postpaymentUseCase.PostProcessingPaymentOrderCommand{
			ReferenceOrderID:  orderID,
			PaymentID:         paymentOrderID,
			Status:            enums.PaymentProcessed,
			AuthorizationCode: authorizationCode,
			OrderStatusString: orderStatusString,
			PaymentReason:     paymentReason,
		}

		env.OnActivity(postProcessingActivity.PostProcessingPaymentOrder,
			mock.Anything,
			expectedCommand,
		).
			Return(enums.Autocapture, nil).
			Once()

		env.OnActivity(notifyOrderChangeActivity.NotifyOrderChange,
			mock.Anything,
			activities.NotifyOrderChangeParams{OrderID: orderID, PaymentID: paymentOrderID},
		).
			Return(fmt.Errorf("failed to notify order change")).
			Once()

		paymentReceipt := entities.PaymentReceipt{
			PaymentID:        paymentOrderID,
			ReferenceOrderID: orderID,
		}

		env.OnActivity(generatePaymentActivity.GenerateReceipt,
			mock.Anything,
			orderID,
			paymentOrderID,
		).
			Return(paymentReceipt, nil).
			Once()

		env.RegisterDelayedCallback(func() {
			env.SignalWorkflow(PaymentIntegrationProcessedSignalName, PaymentProcessedSignal{
				Status:              enums.PaymentProcessed,
				AuthorizationCode:   authorizationCode,
				OrderStatusString:   orderStatusString,
				OrderID:             orderID,
				PaymentStatusString: orderStatusString,
				PaymentReason:       paymentReason,
			})
		}, 10*time.Millisecond)

		env.ExecuteWorkflow(PostProcessingPaymentWorkflow, orderID, paymentOrderID)

		ts.NoError(env.GetWorkflowError())
	})
}
