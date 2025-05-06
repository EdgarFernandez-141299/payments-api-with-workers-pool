package worfkflows

import (
	"gitlab.com/clubhub.ai1/go-libraries/saga/errors"
	errors2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	"time"

	"gitlab.com/clubhub.ai1/go-libraries/saga/activity"
	"gitlab.com/clubhub.ai1/go-libraries/saga/workflow"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows/activities"
	"go.temporal.io/api/enums/v1"
)

func ProcessOrderWorkflow(ctx workflow.Context, workflowPaymentInput PaymentWorkflowInput) (*PaymentOrderWorkflowOut, error) {
	logger := workflow.GetLogger(ctx)
	orderID := workflowPaymentInput.OrderID

	paymentWorkflowOut := NewPaymentOrderWorkflowOut(orderID)

	for _, cmd := range workflowPaymentInput.PaymentCommands {
		paymentWorkflowOut.InitStatus(cmd.Payment.ID)
	}

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

	var createOrderActivity *activities.CheckOrderActivity
	var existsOrder bool

	err := workflow.ExecuteActivity(ctx, createOrderActivity.CheckOrderExistence, orderID).
		Get(ctx, &existsOrder)

	if err != nil {
		logger.Error("GetTokenByOrderAndPaymentID greeting failed.", "Error", err)
		return paymentWorkflowOut, nil
	}

	if !existsOrder {
		return nil, errors.WrapActivityError(errors2.NewOrderNotFoundError(orderID))
	}

	for _, paymentCommand := range workflowPaymentInput.PaymentCommands {
		paymentOrderID := workflowPaymentInput.OrderID + "_" + paymentCommand.Payment.ID
		cwo := workflow.ChildWorkflowOptions{
			WorkflowID:        GetProcessOrderPaymentWorkflowName(paymentOrderID),
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
		}

		ctx = workflow.WithChildOptions(ctx, cwo)

		logger.Info("Starting child workflow", "WorkflowID", cwo.WorkflowID)

		callErr := CallProcessOrderPaymentChildWorkflow(ctx, paymentCommand).Get(ctx, nil)

		if callErr != nil {
			paymentWorkflowOut.ChangePaymentStatus(NewFailedPaymentStatus(paymentCommand.Payment.ID, ""))
			return paymentWorkflowOut, nil
		} else {
			paymentWorkflowOut.ChangePaymentStatus(NewProcessingPaymentStatus(paymentCommand.Payment.ID))
		}
	}

	return paymentWorkflowOut, nil
}
