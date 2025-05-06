package worfkflows

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"

type PaymentStatus struct {
	Reference     string `json:"reference"`
	Status        string `json:"status"`
	FailureReason string `json:"failure_reason,omitempty"`
}

func NewNotProcessedStatus(referenceID string) PaymentStatus {
	return PaymentStatus{
		Reference: referenceID,
		Status:    enums.PaymentNotProcessed.String(),
	}
}

func NewFailedPaymentStatus(referenceID, failureReason string) PaymentStatus {
	return PaymentStatus{
		Reference:     referenceID,
		Status:        enums.PaymentFailed.String(),
		FailureReason: failureReason,
	}
}

func NewProcessingPaymentStatus(referenceID string) PaymentStatus {
	return PaymentStatus{
		Reference: referenceID,
		Status:    enums.PaymentProcessing.String(),
	}
}

type PaymentOrderWorkflowOut struct {
	ReferenceOrderID string          `json:"reference_order_id,omitempty"`
	Payments         []PaymentStatus `json:"payments,omitempty"`
}

func NewPaymentOrderWorkflowOut(referenceOrderID string) *PaymentOrderWorkflowOut {
	return &PaymentOrderWorkflowOut{ReferenceOrderID: referenceOrderID}
}

func (p *PaymentOrderWorkflowOut) InitStatus(paymentReferenceID string) {
	p.Payments = append(p.Payments, NewNotProcessedStatus(paymentReferenceID))
}

func (p *PaymentOrderWorkflowOut) ChangePaymentStatus(paymentStatus PaymentStatus) {
	for i, payment := range p.Payments {
		if payment.Reference == paymentStatus.Reference {
			p.Payments[i] = paymentStatus
			return
		}
	}
}
