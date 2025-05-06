package worfkflows

import (
	"context"

	"go.temporal.io/sdk/client"
)

type PaymentOrderWorkflow interface {
	Call(ctx context.Context, referenceID string, cmd PaymentWorkflowInput) (client.WorkflowRun, error)
	SendProcessedSignal(
		ctx context.Context, paymentOrderID string, cmd PaymentProcessedSignal,
	) error
	SendCaptureFlowSignal(
		ctx context.Context, paymentOrderID string, cmd CompleteCaptureFlowSignal,
	) error
}
