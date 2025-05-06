package resources

import (
	"context"
)

type DeunaCaptureFlowResourceIF interface {
	Release(
		ctx context.Context,
		orderToken string,
		reason string,
	) (bool, error)

	Capture(
		ctx context.Context,
		orderToken string,
		amount int64,
	) (bool, error)
}
