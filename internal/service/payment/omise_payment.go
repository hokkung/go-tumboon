package service

import (
	omisecli "github.com/hokkung/go-tumboon/pkg/omise"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type OmisePaymentService struct {
	omiseClient omisecli.OmiseClient
}

func NewOmisePaymentService(omiseClient omisecli.OmiseClient) *OmisePaymentService {
	return &OmisePaymentService{omiseClient: omiseClient}
}

func ProvideOmisePaymentService(omiseClient omisecli.OmiseClient) PaymentService {
	return NewOmisePaymentService(omiseClient)
}

func (s OmisePaymentService) Do(payment PaymentRequest) error {
	switch payment.Type {
	case Card:
		return s.payByCard(payment)
	default:
		return ErrUnSupportedPaymentMethod
	}
}

func (s OmisePaymentService) payByCard(payment PaymentRequest) error {
	var card omise.Card
	err := s.omiseClient.Do(&card, &operations.CreateToken{
		Name:            payment.Name,
		Number:          payment.CCNumber,
		ExpirationMonth: payment.ExpMonth,
		ExpirationYear:  payment.ExpYear,
		SecurityCode:    payment.CVV,
	})
	if err != nil {
		return err
	}

	var result omise.Charge
	err = s.omiseClient.Do(&result, &operations.CreateCharge{
		Amount:   payment.AmountSubunits,
		Currency: "thb",
		Card:     card.ID,
	})
	if err != nil {
		return err
	}
	return nil
}
