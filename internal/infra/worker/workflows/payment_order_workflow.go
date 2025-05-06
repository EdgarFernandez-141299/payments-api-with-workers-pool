package workflows

import (
	"context"

	"gitlab.com/clubhub.ai1/go-libraries/saga/workflow"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows/activities"
	"go.temporal.io/sdk/client"
	"go.uber.org/fx"
)

const (
	paymentOrderWorkflow = "paymentOrderWorkflow"
	paymentOrderQueue    = "paymentOrderQueue"
)

type PaymentWorkflowImpl struct {
	workflow workflow.Workflow
}

func NewPaymentWorkflowImpl(workflow workflow.Workflow) worfkflows.PaymentOrderWorkflow {
	return &PaymentWorkflowImpl{workflow: workflow}
}

func (p *PaymentWorkflowImpl) Call(
	ctx context.Context, referenceID string, cmd worfkflows.PaymentWorkflowInput,
) (client.WorkflowRun, error) {
	return p.workflow.Execute(ctx, referenceID, cmd)
}

func (p *PaymentWorkflowImpl) SendProcessedSignal(
	ctx context.Context, paymentOrderID string, cmd worfkflows.PaymentProcessedSignal,
) error {
	workflowID := worfkflows.GetPostProcessingWorkflowName(paymentOrderID)
	return p.workflow.SendSignal(ctx, workflowID, worfkflows.PaymentIntegrationProcessedSignalName, cmd)
}

func (p *PaymentWorkflowImpl) SendCaptureFlowSignal(
	ctx context.Context, paymentOrderID string, cmd worfkflows.CompleteCaptureFlowSignal,
) error {
	workflowID := worfkflows.GetCaptureFlowWorkflowName(paymentOrderID)
	return p.workflow.SendSignal(ctx, workflowID, worfkflows.CaptureFlowSignalName, cmd)
}

func NewPaymentOrderWorkflow(
	lc fx.Lifecycle,
	temporalClient client.Client,
	checkOrderActivity *activities.CheckOrderActivity,
	createPaymentOrderActivity *activities.CreatePaymentOrderActivity,
	postProcessingActivity *activities.PostProcessingPaymentOrderActivity,
	notifyOrderChangeActivity *activities.NotifyOrderChangeActivity,
	capturePaymentActivity *activities.CapturePaymentActivity,
	releasePaymentActivity *activities.ReleasePaymentActivity,
	generateReceiptActivity *activities.GeneratePaymentReceiptActivity,
) worfkflows.PaymentOrderWorkflow {
	wf := workflow.NewTemporalWorkflowBuilder().
		WithWorkflow(worfkflows.ProcessOrderWorkflow).
		WithChildWorkflows(
			worfkflows.ProcessOrderPaymentWorkflow,
			worfkflows.PostProcessingPaymentWorkflow,
			worfkflows.CaptureFlowWorkflow,
		).
		WithActivities(
			workflow.NewActivity(activities.CheckOrderExistenceActivityName, checkOrderActivity.CheckOrderExistence),
			workflow.NewActivity(activities.CreatePaymentOrderActivityName, createPaymentOrderActivity.CreatePaymentOrder),
			workflow.NewActivity(activities.PostProcessingPaymentOrderActivityName, postProcessingActivity.PostProcessingPaymentOrder),
			workflow.NewActivity(activities.NotifyOrderChangeActivityName, notifyOrderChangeActivity.NotifyOrderChange),
			workflow.NewActivity(activities.CapturePaymentActivityName, capturePaymentActivity.CapturePayment),
			workflow.NewActivity(activities.ReleasePaymentActivityName, releasePaymentActivity.ReleasePayment),
			workflow.NewActivity(activities.GenerateReceiptActivityName, generateReceiptActivity.GenerateReceipt),
		).
		WithName(paymentOrderWorkflow).
		WithQueue(paymentOrderQueue).
		WithClient(temporalClient).
		Build()

	paymentWorkflow := NewPaymentWorkflowImpl(wf)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := wf.Run(); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := wf.Stop(); err != nil {
				panic(err)
			}
			return nil
		},
	})

	return paymentWorkflow
}
