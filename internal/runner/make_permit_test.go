package runner_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hokkung/go-tumboon/internal/runner"
	mock_service "github.com/hokkung/go-tumboon/internal/service/donation/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DonationRunnerTestSuite struct {
	suite.Suite
	mockService *mock_service.MockDonationService
	underTest   *runner.DonationRunner
}

func (suite *DonationRunnerTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockService = mock_service.NewMockDonationService(ctrl)

	suite.underTest = runner.NewDonationRunner(suite.mockService)

}

func TestDonationRunnerTestSuite(t *testing.T) {
	suite.Run(t, new(DonationRunnerTestSuite))
}

func (suite *DonationRunnerTestSuite) TestRun() {
	suite.mockService.EXPECT().MakePermit().Return(nil).Times(1)

	err := suite.underTest.Run()
	suite.NoError(err)
}

func (suite *DonationRunnerTestSuite) TestRunError() {
	expectedError := assert.AnError
	suite.mockService.EXPECT().MakePermit().Return(expectedError)

	err := suite.underTest.Run()
	suite.Error(err)
	suite.Equal(expectedError, err)
}
