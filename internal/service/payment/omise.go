package service

import (
	omisecli "github.com/hokkung/go-tumboon/pkg/omise"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

// OmisePaymentService manages Omise payment.
// The struct implements `PaymentService“ interface for handling Omise transactions.
type OmisePaymentService struct {
	omiseClient omisecli.OmiseClient
}

// NewOmisePaymentService creates Omise payment service.
func NewOmisePaymentService(omiseClient omisecli.OmiseClient) *OmisePaymentService {
	return &OmisePaymentService{omiseClient: omiseClient}
}

// ProvideOmisePaymentService provides Omise payment service for dependency injection.
func ProvideOmisePaymentService(omiseClient omisecli.OmiseClient) PaymentService {
	return NewOmisePaymentService(omiseClient)
}

// Do performs a payment transaction.
func (s *OmisePaymentService) Do(payment PaymentRequest) (*PaymentResponse, error) {
	switch payment.Type {
	case Card:
		return s.payByCard(payment)
	default:
		return s.getFailedResponse(payment), ErrUnSupportedPaymentMethod
	}
}

func (s *OmisePaymentService) payByCard(payment PaymentRequest) (*PaymentResponse, error) {
	var token omise.Token
	err := s.omiseClient.CreateToken(&token, &operations.CreateToken{
		Name:            payment.Name,
		Number:          payment.CCNumber,
		ExpirationMonth: payment.ExpMonth,
		ExpirationYear:  payment.ExpYear,
		SecurityCode:    payment.CVV,
	})
	if err != nil {
		return s.getFailedResponse(payment), err
	}

	var result omise.Charge
	err = s.omiseClient.CreateCharge(&result, &operations.CreateCharge{
		Amount:   payment.AmountSubunits,
		Currency: "thb",
		Card:     token.ID,
	})
	if err != nil {
		return s.getFailedResponse(payment), err
	}

	return &PaymentResponse{
		Amount:    result.Amount,
		IsSuccess: true,
		Source:    payment,
	}, nil
}

func (s *OmisePaymentService) getFailedResponse(
	payment PaymentRequest,
) *PaymentResponse {
	return &PaymentResponse{
		IsSuccess: false,
		Source:    payment,
	}
}
