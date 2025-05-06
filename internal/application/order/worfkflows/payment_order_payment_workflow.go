package worfkflows

import (
	"gitlab.com/clubhub.ai1/go-libraries/saga/activity"
	"gitlab.com/clubhub.ai1/go-libraries/saga/workflow"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows/activities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
	enums2 "go.temporal.io/api/enums/v1"

	"time"
)

func GetProcessOrderPaymentWorkflowName(paymentID string) string {
	return "ProcessOrderPaymentWorkflow_" + paymentID
}

func CallProcessOrderPaymentChildWorkflow(ctx workflow.Context, cmd command.CreatePaymentOrderCommand) workflow.ChildWorkflowFuture {
	cwo := workflow.ChildWorkflowOptions{
		WorkflowID:        GetProcessOrderPaymentWorkflowName(cmd.Payment.ID),
		ParentClosePolicy: enums2.PARENT_CLOSE_POLICY_ABANDON,
	}

	ctx = workflow.WithChildOptions(ctx, cwo)

	return workflow.ExecuteChildWorkflow(ctx, ProcessOrderPaymentWorkflow, cmd)
}

func ProcessOrderPaymentWorkflow(ctx workflow.Context, cmd command.CreatePaymentOrderCommand) error {
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

	var paymentOrderActivity *activities.CreatePaymentOrderActivity
	var paymentOrderResponse response.PaymentOrderResponseDTO

	err := workflow.ExecuteActivity(ctx, paymentOrderActivity.CreatePaymentOrder, cmd).Get(ctx, &paymentOrderResponse)

	if err != nil {
		return err
	}

	wfFuture := CallPostProcessingPaymentWorkflow(ctx, cmd.ReferenceOrderID, cmd.Payment.ID)

	// Wait for the child workflow to start.
	if getChildWorkflowErr := wfFuture.GetChildWorkflowExecution().Get(ctx, nil); getChildWorkflowErr != nil {
		return getChildWorkflowErr
	}

	return nil
}
