package value_objects

const (
	processing         = "PROCESSING"
	failed             = "FAILED"
	refunded           = "REFUNDED"
	partiallyRefunded  = "PARTIALLY_REFUNDED"
	partiallyProcessed = "PARTIALLY_PROCESSED"
	processed          = "PROCESSED"
	authorized         = "AUTHORIZED"
	canceled           = "CANCELED"
)

type OrderStatus struct {
	status string
}

func OrderStatusFailed() OrderStatus {
	return OrderStatus{status: failed}
}

func OrderStatusProcessing() OrderStatus {
	return OrderStatus{status: processing}
}

func OrderPartialProcessed() OrderStatus {
	return OrderStatus{status: partiallyProcessed}
}

func OrderStatusProcessed() OrderStatus {
	return OrderStatus{status: processed}
}

func OrderStatusAuthorized() OrderStatus {
	return OrderStatus{status: authorized}
}

func OrderStatusCanceled() OrderStatus {
	return OrderStatus{status: canceled}
}

func OrderStatusRefunded() OrderStatus {
	return OrderStatus{status: refunded}
}

func OrderStatusPartiallyRefunded() OrderStatus {
	return OrderStatus{status: partiallyRefunded}
}

func (o OrderStatus) Get() string {
	return o.status
}
