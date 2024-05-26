package model_test

import (
	"testing"
	"time"

	"github.com/hokkung/go-tumboon/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewDonation(t *testing.T) {
	name := "John"
	amountSubunits := "100"
	CCNumber := "1234"
	CVV := "000"
	expMonth := "1"
	expYear := "2025"

	recs := []string{name, amountSubunits, CCNumber, CVV, expMonth, expYear}

	donation, err := model.NewDonation(recs)

	assert.Nil(t, err)
	assert.Equal(t, name, donation.Name)
	assert.Equal(t, int64(100), donation.AmountSubunits)
	assert.Equal(t, CCNumber, donation.CCNumber)
	assert.Equal(t, CVV, donation.CVV)
	assert.Equal(t, time.January, donation.ExpMonth)
	assert.Equal(t, 2025, donation.ExpYear)
}

func TestNewDonationInvalidRawData(t *testing.T) {
	name := "John"
	amountSubunits := "100"
	CCNumber := "1234"
	CVV := "000"
	expMonth := "1"

	recs := []string{name, amountSubunits, CCNumber, CVV, expMonth}

	donation, err := model.NewDonation(recs)

	assert.Error(t, err)
	assert.EqualError(t, err, model.ErrInvalidRawDataSize.Error())
	assert.Nil(t, donation)
}

func TestNewDonationWrongAmountSubunits(t *testing.T) {
	name := "John"
	amountSubunits := "a"
	CCNumber := "1234"
	CVV := "000"
	expMonth := "1"
	expYear := "2025"

	recs := []string{name, amountSubunits, CCNumber, CVV, expMonth, expYear}

	donation, err := model.NewDonation(recs)

	assert.Error(t, err)
	assert.Nil(t, donation)
}

func TestNewDonationWrongMonth(t *testing.T) {
	name := "John"
	amountSubunits := "100"
	CCNumber := "1234"
	CVV := "000"
	expMonth := "a"
	expYear := "2025"

	recs := []string{name, amountSubunits, CCNumber, CVV, expMonth, expYear}

	donation, err := model.NewDonation(recs)

	assert.Error(t, err)
	assert.Nil(t, donation)
}

func TestNewDonationWrongYear(t *testing.T) {
	name := "John"
	amountSubunits := "100"
	CCNumber := "1234"
	CVV := "000"
	expMonth := "1"
	expYear := "a"

	recs := []string{name, amountSubunits, CCNumber, CVV, expMonth, expYear}

	donation, err := model.NewDonation(recs)

	assert.Error(t, err)
	assert.Nil(t, donation)
}
