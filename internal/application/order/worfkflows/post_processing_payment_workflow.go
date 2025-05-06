package worfkflows

import (
	"errors"
	"fmt"
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"

	"gitlab.com/clubhub.ai1/go-libraries/saga/activity"
	"gitlab.com/clubhub.ai1/go-libraries/saga/workflow"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows/activities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	postpaymentUseCase "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/post_payment"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	enums2 "go.temporal.io/api/enums/v1"
)

const PaymentIntegrationProcessedSignalName = "paymentIntegrationProcessedSignal"

type PaymentProcessedSignal struct {
	AuthorizationCode   string
	OrderStatusString   string
	Status              enums.PaymentStatus
	OrderID             string
	PaymentID           string
	PaymentStatusString string
	PaymentReason       string
	PaymentCard         CardData
}

type CardData struct {
	CardBrand string
	CardLast4 string
	CardType  string
}

func GetPostProcessingWorkflowName(paymentOrderID string) string {
	return fmt.Sprintf("PostProcessingOrderPaymentWorkflow-%s", paymentOrderID)
}

func CallPostProcessingPaymentWorkflow(ctx workflow.Context, referenceOrderID string, paymentID string) workflow.ChildWorkflowFuture {
	deunaOrderID := utils.NewDeunaOrderID(referenceOrderID, paymentID)
	cwo := workflow.ChildWorkflowOptions{
		WorkflowID:        GetPostProcessingWorkflowName(deunaOrderID.GetID()),
		ParentClosePolicy: enums2.PARENT_CLOSE_POLICY_ABANDON,
	}

	ctx = workflow.WithChildOptions(ctx, cwo)

	return workflow.ExecuteChildWorkflow(ctx, PostProcessingPaymentWorkflow, referenceOrderID, paymentID)
}

func PostProcessingPaymentWorkflow(ctx workflow.Context, referenceOrderID string, paymentID string) error {
	var processSignal PaymentProcessedSignal
	var paymentErr error

	var postProcessingActivity *activities.PostProcessingPaymentOrderActivity
	var notifyOrderChangeActivity *activities.NotifyOrderChangeActivity
	var generatePaymentActivity *activities.GeneratePaymentReceiptActivity
	var captureFlow enums.PaymentFlowEnum = enums.Autocapture

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

	workflow.GetSignalChannel(ctx, PaymentIntegrationProcessedSignalName).Receive(ctx, &processSignal)

	if processSignal.Status.IsFailure() {
		paymentErr = errors.New("payment failed")
	}

	postProcessingCmd := postpaymentUseCase.PostProcessingPaymentOrderCommand{
		ReferenceOrderID:  referenceOrderID,
		PaymentID:         paymentID,
		Status:            processSignal.Status,
		AuthorizationCode: processSignal.AuthorizationCode,
		OrderStatusString: processSignal.OrderStatusString,
		PaymentReason:     processSignal.PaymentReason,
		PaymentCard: postpaymentUseCase.CardData{
			CardBrand: processSignal.PaymentCard.CardBrand,
			CardLast4: processSignal.PaymentCard.CardLast4,
			CardType:  processSignal.PaymentCard.CardType,
		},
	}

	paymentErr = workflow.ExecuteActivity(ctx, postProcessingActivity.PostProcessingPaymentOrder, postProcessingCmd).
		Get(ctx, &captureFlow)

	if paymentErr != nil {
		return fmt.Errorf("error in post processing: %w", paymentErr)
	}

	notifyOrderChangeCmd := activities.NotifyOrderChangeParams{
		OrderID:   referenceOrderID,
		PaymentID: paymentID,
	}

	notifyOrderChangeErr := workflow.ExecuteActivity(
		ctx, notifyOrderChangeActivity.NotifyOrderChange, notifyOrderChangeCmd,
	).
		Get(ctx, nil)

	if notifyOrderChangeErr != nil {
		fmt.Errorf("error in notify order change: %w", notifyOrderChangeErr)
	}

	if captureFlow == enums.Capture {
		wfFuture := CallCaptureFlowWorkflow(ctx, referenceOrderID, paymentID)

		if getChildWorkflowErr := wfFuture.GetChildWorkflowExecution().Get(ctx, nil); getChildWorkflowErr != nil {
			return getChildWorkflowErr
		}
	}

	var paymentReceipt entities.PaymentReceipt

	if !processSignal.Status.IsFailure() {
		generateReceiptErr := workflow.ExecuteActivity(
			ctx, generatePaymentActivity.GenerateReceipt, referenceOrderID, paymentID,
		).
			Get(ctx, &paymentReceipt)

		if generateReceiptErr != nil {
			return fmt.Errorf("error in generate receipt: %w", generateReceiptErr)
		}
	}

	return nil
}
