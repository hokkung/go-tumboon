package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/hokkung/go-tumboon/config"
	"github.com/hokkung/go-tumboon/internal/model"
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
	filePath                        = "../../../testutils/example_donation_file.csv.rot128"
	filePathNoData                  = "../../../testutils/example_donation_file_no_data.csv.rot128"
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

func (suite *DonationServiceTestSuite) TestDonates() {
	d1 := model.Donation{
		Name:           "a",
		AmountSubunits: 1,
		CCNumber:       "1111111111111111",
		CVV:            "111",
		ExpMonth:       time.August,
		ExpYear:        2026,
	}
	d2 := model.Donation{
		Name:           "b",
		AmountSubunits: 1,
		CCNumber:       "1111111111111111",
		CVV:            "111",
		ExpMonth:       time.August,
		ExpYear:        2026,
	}
	d3 := model.Donation{
		Name:           "c",
		AmountSubunits: 1,
		CCNumber:       "1111111111111111",
		CVV:            "111",
		ExpMonth:       time.August,
		ExpYear:        2026,
	}
	d4 := model.Donation{
		Name:           "d",
		AmountSubunits: 4,
		CCNumber:       "1111111111111111",
		CVV:            "111",
		ExpMonth:       time.August,
		ExpYear:        2026,
	}

	r1 := ps.PaymentRequest{
		Name:           d1.Name,
		AmountSubunits: d1.AmountSubunits,
		CCNumber:       d1.CCNumber,
		CVV:            d1.CVV,
		ExpMonth:       d1.ExpMonth,
		ExpYear:        d1.ExpYear,
		Type:           ps.Card,
	}
	r2 := ps.PaymentRequest{
		Name:           d2.Name,
		AmountSubunits: d2.AmountSubunits,
		CCNumber:       d2.CCNumber,
		CVV:            d2.CVV,
		ExpMonth:       d2.ExpMonth,
		ExpYear:        d2.ExpYear,
		Type:           ps.Card,
	}
	r3 := ps.PaymentRequest{
		Name:           d3.Name,
		AmountSubunits: d3.AmountSubunits,
		CCNumber:       d3.CCNumber,
		CVV:            d3.CVV,
		ExpMonth:       d3.ExpMonth,
		ExpYear:        d3.ExpYear,
		Type:           ps.Card,
	}
	r4 := ps.PaymentRequest{
		Name:           d4.Name,
		AmountSubunits: d4.AmountSubunits,
		CCNumber:       d4.CCNumber,
		CVV:            d4.CVV,
		ExpMonth:       d4.ExpMonth,
		ExpYear:        d4.ExpYear,
		Type:           ps.Card,
	}
	res1 := &ps.PaymentResponse{
		Amount:    r1.AmountSubunits,
		IsSuccess: true,
		Source:    r1,
	}
	res2 := &ps.PaymentResponse{
		Amount:    r2.AmountSubunits,
		IsSuccess: true,
		Source:    r2,
	}
	res3 := &ps.PaymentResponse{
		Amount:    r3.AmountSubunits,
		IsSuccess: false,
		Source:    r3,
	}
	res4 := &ps.PaymentResponse{
		Amount:    r4.AmountSubunits,
		IsSuccess: false,
		Source:    r4,
	}
	mockErr := errors.New("mock error")

	suite.mockService.EXPECT().Do(r1).Return(res1, nil)
	suite.mockService.EXPECT().Do(r2).Return(res2, nil)
	suite.mockService.EXPECT().Do(r3).Return(res3, nil)
	suite.mockService.EXPECT().Do(r4).Return(res4, mockErr)

	res, err := suite.underTest.Donates([]model.Donation{d1, d2, d3, d4})

	suite.NoError(err)
	suite.Equal(int64(7), res.TotalReceived)
	suite.Equal(int64(5), res.FaultyDonated)
	suite.Equal(int64(2), res.SuccessfulDonated)
	suite.Equal(float64(1.75), res.AveragePerPerson)
	suite.Equal("d", res.TopDonors[0].Name)
	suite.Equal("a", res.TopDonors[1].Name)
	suite.Equal("b", res.TopDonors[2].Name)
}
