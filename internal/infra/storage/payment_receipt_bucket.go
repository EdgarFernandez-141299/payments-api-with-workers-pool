package storage

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/config"

type PaymentReceiptBucket struct {
	name string
}

func NewPaymentReceiptBucket() *PaymentReceiptBucket {
	bucketName := config.Config().Aws.ReceiptPaymentBucket
	return &PaymentReceiptBucket{bucketName}
}
