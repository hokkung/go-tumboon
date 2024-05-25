package service

//go:generate mockgen -source=payment.go -destination=./mock/mock_payment.go
type PaymentService interface {
	Do(payment PaymentRequest) error
}
