package service

//go:generate mockgen -source=payment.go -destination=./mock/mock_payment.go

// PaymentService manages payment transaction.
type PaymentService interface {
	Do(payment PaymentRequest) (*PaymentResponse, error)
}
