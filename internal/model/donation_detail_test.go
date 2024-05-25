package model_test

import (
	"testing"
	"time"

	"github.com/hokkung/go-tumboon/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewDonationDetail(t *testing.T) {
	name := "John"
	amountSubunits := "100"
	CCNumber := "1234"
	CVV := "000"
	expMonth := "1"
	expYear := "2025"

	recs := []string{name, amountSubunits, CCNumber, CVV, expMonth, expYear}

	detail, err := model.NewDonationDetail(recs)

	assert.Nil(t, err)
	assert.Equal(t, name, detail.Name)
	assert.Equal(t, int64(100), detail.AmountSubunits)
	assert.Equal(t, CCNumber, detail.CCNumber)
	assert.Equal(t, CVV, detail.CVV)
	assert.Equal(t, time.January, detail.ExpMonth)
	assert.Equal(t, 2025, detail.ExpYear)
}

func TestNewDonationDetailInvalidRawData(t *testing.T) {
	name := "John"
	amountSubunits := "100"
	CCNumber := "1234"
	CVV := "000"
	expMonth := "1"

	recs := []string{name, amountSubunits, CCNumber, CVV, expMonth}

	detail, err := model.NewDonationDetail(recs)

	assert.Error(t, err)
	assert.EqualError(t, err, model.ErrInvalidRawDataSize.Error())
	assert.Nil(t, detail)
}

func TestNewDonationDetailWrongAmountSubunits(t *testing.T) {
	name := "John"
	amountSubunits := "a"
	CCNumber := "1234"
	CVV := "000"
	expMonth := "1"
	expYear := "2025"

	recs := []string{name, amountSubunits, CCNumber, CVV, expMonth, expYear}

	detail, err := model.NewDonationDetail(recs)

	assert.Error(t, err)
	assert.Nil(t, detail)
}

func TestNewDonationDetailWrongMonth(t *testing.T) {
	name := "John"
	amountSubunits := "100"
	CCNumber := "1234"
	CVV := "000"
	expMonth := "a"
	expYear := "2025"

	recs := []string{name, amountSubunits, CCNumber, CVV, expMonth, expYear}

	detail, err := model.NewDonationDetail(recs)

	assert.Error(t, err)
	assert.Nil(t, detail)
}

func TestNewDonationDetailWrongYear(t *testing.T) {
	name := "John"
	amountSubunits := "100"
	CCNumber := "1234"
	CVV := "000"
	expMonth := "1"
	expYear := "a"

	recs := []string{name, amountSubunits, CCNumber, CVV, expMonth, expYear}

	detail, err := model.NewDonationDetail(recs)

	assert.Error(t, err)
	assert.Nil(t, detail)
}
