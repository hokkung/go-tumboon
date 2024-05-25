package service

import "time"

type PaymentType int

const (
	Cash PaymentType = 1 + iota
	Card
)

type PaymentRequest struct {
	Name           string
	AmountSubunits int64
	CCNumber       string
	CVV            string
	ExpMonth       time.Month
	ExpYear        int
	Type           PaymentType
}
