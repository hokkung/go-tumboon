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
	mockToken := &omise.Token{}
	suite.mockOmiseClient.EXPECT().CreateToken(mockToken, &operations.CreateToken{
		Name:            mockReq.Name,
		Number:          mockReq.CCNumber,
		ExpirationMonth: mockReq.ExpMonth,
		ExpirationYear:  mockReq.ExpYear,
		SecurityCode:    mockReq.CVV,
	}).DoAndReturn(func(result interface{}, _ interface{}) error {
		token := result.(*omise.Token)
		token.ID = "card_test_5fz2lvcrbnao9mkwz80"
		return nil
	})

	mockCharge := &omise.Charge{}
	suite.mockOmiseClient.EXPECT().CreateCharge(mockCharge, &operations.CreateCharge{
		Amount:   mockReq.AmountSubunits,
		Currency: "thb",
		Card:     "card_test_5fz2lvcrbnao9mkwz80",
	}).DoAndReturn(func(result interface{}, _ interface{}) error {
		charge := result.(*omise.Charge)
		charge.Amount = 1000
		charge.ID = "charge_test_5fz2lvcrbnao9mkwz80"
		return nil
	})

	res, err := suite.underTest.Do(mockReq)

	suite.NoError(err)
	suite.Equal(mockReq.AmountSubunits, res.Amount)
	suite.Equal(true, res.IsSuccess)
	suite.Equal(mockReq, res.Source)
}

func (suite *PaymentServiceTestSuite) TestDoUnSupportedPaymentMethod() {
	res, err := suite.underTest.Do(service.PaymentRequest{
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
	suite.Equal(false, res.IsSuccess)
	suite.Equal(int64(0), res.Amount)
}
