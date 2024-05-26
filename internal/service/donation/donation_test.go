package service_test

import (
	"testing"
	"time"

	"github.com/hokkung/go-tumboon/config"
	ds "github.com/hokkung/go-tumboon/internal/service/donation"
	service "github.com/hokkung/go-tumboon/internal/service/donation"
	ps "github.com/hokkung/go-tumboon/internal/service/payment"
	mock_service "github.com/hokkung/go-tumboon/internal/service/payment/mock"
	"github.com/hokkung/go-tumboon/internal/validator"
	"github.com/hokkung/go-tumboon/testutils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

const (
	filePath                           = "../../../testutils/example_donation_file.csv.rot128"
	filePathNoData                     = "../../../testutils/example_donation_file_no_data.csv.rot128"
	filePathContainsWrongDataOneRow = "../../../testutils/example_donation_file_wrong_data_one_row.csv.rot128"
)

type DonationServiceTestSuite struct {
	suite.Suite
	mockService *mock_service.MockPaymentService
	underTest   service.DonationService
}

func (suite *DonationServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockService = mock_service.NewMockPaymentService(ctrl)
}

func TestDonationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DonationServiceTestSuite))
}

func (suite *DonationServiceTestSuite) SetupSuite() {
	testutils.GenerateMockDonationFile(filePath, "Name,AmountSubunits,CCNumber,CVV,ExpMonth,ExpYear\nJohn,100,1111111111111111,111,3,2026\nAmenda,20,2222222222222222,111,3,2026\nJohn,50,1111111111111111,111,3,2026\n")
	testutils.GenerateMockDonationFile(filePathNoData, "Name,AmountSubunits,CCNumber,CVV,ExpMonth,ExpYear\n")
	testutils.GenerateMockDonationFile(filePathContainsWrongDataOneRow, "Name,AmountSubunits,CCNumber,CVV,ExpMonth,ExpYear\nJohn,100,1111111111111111,111,3,2026\nAmenda,20,a,111,3,2026\nJohn,50,1111111111111111,111,3,2026\n")

}

func (suite *DonationServiceTestSuite) TearDownSuite() {
	testutils.RemoveMockFile(filePath)
	testutils.RemoveMockFile(filePathNoData)
	testutils.RemoveMockFile(filePathContainsWrongDataOneRow)
}

func (suite *DonationServiceTestSuite) TestMakePermit() {
	suite.underTest = ds.NewDonationService(
		suite.mockService,
		config.Configuration{
			DonationFileConfiguration: &config.DonationFileConfiguration{
				DonationFileAddr: filePath,
				MaxConcurrent:    8,
			},
		},
		validator.NewCustomValidator(),
	)

	p1 := ps.PaymentRequest{
		Name:           "John",
		AmountSubunits: 100,
		CCNumber:       "1111111111111111",
		CVV:            "111",
		ExpMonth:       time.March,
		ExpYear:        2026,
		Type:           ps.Card,
	}
	p2 := ps.PaymentRequest{
		Name:           "Amenda",
		AmountSubunits: 20,
		CCNumber:       "2222222222222222",
		CVV:            "111",
		ExpMonth:       time.March,
		ExpYear:        2026,
		Type:           ps.Card,
	}
	p3 := ps.PaymentRequest{
		Name:           "John",
		AmountSubunits: 50,
		CCNumber:       "1111111111111111",
		CVV:            "111",
		ExpMonth:       time.March,
		ExpYear:        2026,
		Type:           ps.Card,
	}
	r1 := &ps.PaymentResponse{
		IsSuccess: true,
		Amount:    p1.AmountSubunits,
		Source:    p1,
	}
	r2 := &ps.PaymentResponse{
		IsSuccess: true,
		Amount:    p2.AmountSubunits,
		Source:    p2,
	}
	r3 := &ps.PaymentResponse{
		IsSuccess: true,
		Amount:    p3.AmountSubunits,
		Source:    p3,
	}

	suite.mockService.EXPECT().Do(p1).Return(r1, nil)
	suite.mockService.EXPECT().Do(p2).Return(r2, nil)
	suite.mockService.EXPECT().Do(p3).Return(r3, nil)

	err := suite.underTest.MakePermit()

	suite.NoError(err)
}

func (suite *DonationServiceTestSuite) TestMakePermitNoData() {
	suite.underTest = ds.NewDonationService(
		suite.mockService,
		config.Configuration{
			DonationFileConfiguration: &config.DonationFileConfiguration{
				DonationFileAddr: filePathNoData,
				MaxConcurrent:    8,
			},
		},
		validator.NewCustomValidator(),
	)

	err := suite.underTest.MakePermit()

	suite.NoError(err)
}

func (suite *DonationServiceTestSuite) TestMakePermitContainsWrongDataOneRow() {
	suite.underTest = ds.NewDonationService(
		suite.mockService,
		config.Configuration{
			DonationFileConfiguration: &config.DonationFileConfiguration{
				DonationFileAddr: filePathContainsWrongDataOneRow,
				MaxConcurrent:    8,
			},
		},
		validator.NewCustomValidator(),
	)

	p1 := ps.PaymentRequest{
		Name:           "John",
		AmountSubunits: 100,
		CCNumber:       "1111111111111111",
		CVV:            "111",
		ExpMonth:       time.March,
		ExpYear:        2026,
		Type:           ps.Card,
	}
	p3 := ps.PaymentRequest{
		Name:           "John",
		AmountSubunits: 50,
		CCNumber:       "1111111111111111",
		CVV:            "111",
		ExpMonth:       time.March,
		ExpYear:        2026,
		Type:           ps.Card,
	}
	r1 := &ps.PaymentResponse{
		IsSuccess: true,
		Amount:    p1.AmountSubunits,
		Source:    p1,
	}
	r3 := &ps.PaymentResponse{
		IsSuccess: true,
		Amount:    p3.AmountSubunits,
		Source:    p3,
	}

	suite.mockService.EXPECT().Do(p1).Return(r1, nil)
	suite.mockService.EXPECT().Do(p3).Return(r3, nil)

	err := suite.underTest.MakePermit()

	suite.NoError(err)
}


