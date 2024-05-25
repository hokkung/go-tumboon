package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	service "github.com/hokkung/go-tumboon/internal/service/payment"
	mock_omisecli "github.com/hokkung/go-tumboon/pkg/omise/mock"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"github.com/stretchr/testify/suite"
)

type PaymentServiceTestSuite struct {
	suite.Suite

	mockOmiseClient *mock_omisecli.MockOmiseClient
	underTest       *service.OmisePaymentService
}

func (suite *PaymentServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	mockOmiseClient := mock_omisecli.NewMockOmiseClient(ctrl)
	suite.mockOmiseClient = mockOmiseClient

	suite.underTest = service.NewOmisePaymentService(mockOmiseClient)
}

func TestPaymentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}

func (suite *PaymentServiceTestSuite) TestDo() {
	mockReq := service.PaymentRequest{
		Type:           service.Card,
		Name:           "John Doe",
		CCNumber:       "4242424242424242",
		ExpMonth:       12,
		ExpYear:        2025,
		CVV:            "123",
		AmountSubunits: 1000,
	}

	suite.mockOmiseClient.EXPECT().Do(gomock.Any(), &operations.CreateToken{
		Name:            mockReq.Name,
		Number:          mockReq.CCNumber,
		ExpirationMonth: mockReq.ExpMonth,
		ExpirationYear:  mockReq.ExpYear,
		SecurityCode:    mockReq.CVV,
	}).DoAndReturn(func(result interface{}, _ interface{}) error {
		token := result.(*omise.Card)
		token.ID = "card_test_5fz2lvcrbnao9mkwz80"
		return nil
	})

	suite.mockOmiseClient.EXPECT().Do(gomock.Any(), &operations.CreateCharge{
		Amount:   mockReq.AmountSubunits,
		Currency: "thb",
		Card:     "card_test_5fz2lvcrbnao9mkwz80",
	}).DoAndReturn(func(result interface{}, _ interface{}) error {
		charge := result.(*omise.Charge)
		charge.ID = "charge_test_5fz2lvcrbnao9mkwz80"
		return nil
	})

	err := suite.underTest.Do(mockReq)

	suite.NoError(err)
}

func (suite *PaymentServiceTestSuite) TestDoUnSupportedPaymentMethod() {
	err := suite.underTest.Do(service.PaymentRequest{
		Type:           service.Cash,
		Name:           "John Doe",
		CCNumber:       "4242424242424242",
		ExpMonth:       12,
		ExpYear:        2025,
		CVV:            "123",
		AmountSubunits: 1000,
	})

	suite.Error(err)
	suite.Equal(service.ErrUnSupportedPaymentMethod, err)
}
