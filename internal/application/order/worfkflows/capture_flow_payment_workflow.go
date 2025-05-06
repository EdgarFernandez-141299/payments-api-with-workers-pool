package worfkflows

import (
	"fmt"
	"time"

	"gitlab.com/clubhub.ai1/go-libraries/saga/activity"
	"gitlab.com/clubhub.ai1/go-libraries/saga/workflow"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows/activities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	enums2 "go.temporal.io/api/enums/v1"
)

const (
	CaptureFlowSignalName = "captureFlowSignal"
	CaptureFlowTimeout    = 2 * time.Minute
	TimeoutReason         = "Timeout: No se recibió señal de captura en el tiempo esperado"
)

type CompleteCaptureFlowSignal struct {
	OrderReferecenId string
	PaymentID        string
	Reason           string
	Action           enums.PaymentFlowActionEnum
}

func GetCaptureFlowSignalName(paymentOrderID string) string {
	return fmt.Sprintf("%s-%s", CaptureFlowSignalName, paymentOrderID)
}

func GetCaptureFlowWorkflowName(paymentOrderID string) string {
	return fmt.Sprintf("CaptureFlow-%s", paymentOrderID)
}

func CallCaptureFlowWorkflow(ctx workflow.Context, referenceOrderID string, paymentID string) workflow.ChildWorkflowFuture {
	deunaOrderID := utils.NewDeunaOrderID(referenceOrderID, paymentID)
	cwo := workflow.ChildWorkflowOptions{
		WorkflowID:        GetCaptureFlowWorkflowName(deunaOrderID.GetID()),
		ParentClosePolicy: enums2.PARENT_CLOSE_POLICY_ABANDON,
	}

	ctx = workflow.WithChildOptions(ctx, cwo)

	return workflow.ExecuteChildWorkflow(ctx, CaptureFlowWorkflow, referenceOrderID, paymentID)
}

func CaptureFlowWorkflow(ctx workflow.Context, referenceOrderID string, paymentID string) error {
	var processSignal CompleteCaptureFlowSignal

	ao := activity.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &activity.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    2,
		},
	}

	ctx = activity.WithActivityOptions(ctx, ao)

	var capturePaymentActivity *activities.CapturePaymentActivity
	var releasePaymentActivity *activities.ReleasePaymentActivity
	var paymentCaptureFlowErr error

	selector := workflow.NewSelector(ctx)

	signalChan := workflow.GetSignalChannel(ctx, CaptureFlowSignalName)

	timeoutChan := workflow.NewTimer(ctx, CaptureFlowTimeout)

	selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &processSignal)
	})

	selector.AddFuture(timeoutChan, func(f workflow.Future) {
		processSignal.Action = enums.ReleasePayment
		processSignal.OrderReferecenId = referenceOrderID
		processSignal.PaymentID = paymentID
		processSignal.Reason = TimeoutReason
	})

	selector.Select(ctx)

	if processSignal.Action.IsCapture() {
		paymentCaptureFlowErr = workflow.ExecuteActivity(ctx, capturePaymentActivity.CapturePayment,
			processSignal.OrderReferecenId, processSignal.PaymentID).Get(ctx, nil)
	} else if processSignal.Action.IsRelease() {
		paymentCaptureFlowErr = workflow.ExecuteActivity(ctx, releasePaymentActivity.ReleasePayment,
			processSignal.OrderReferecenId, processSignal.PaymentID, processSignal.Reason).Get(ctx, nil)
	}

	return paymentCaptureFlowErr
}
