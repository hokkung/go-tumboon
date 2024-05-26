package service

import "time"

type PaymentType int

const (
	Cash PaymentType = 1 + iota
	Card
)

// PaymentRequest collects payment request information.
type PaymentRequest struct {
	Name           string
	AmountSubunits int64
	CCNumber       string
	CVV            string
	ExpMonth       time.Month
	ExpYear        int
	Type           PaymentType
}

// PaymentResponse collects payment response information.
type PaymentResponse struct {
	Amount    int64
	IsSuccess bool
	Source    PaymentRequest
}
