package model

import (
	"errors"
	"strconv"
	"time"
)

type Donation struct {
	Name           string     `validate:"required"`
	AmountSubunits int64      `validate:"required"`
	CCNumber       string     `validate:"required,len=16"`
	CVV            string     `validate:"required,len=3"`
	ExpMonth       time.Month `validate:"required"`
	ExpYear        int        `validate:"required"`
}

var (
	ErrInvalidRawDataSize = errors.New("invalid raw data size")
)

func NewDonation(raw []string) (*Donation, error) {
	if len(raw) < 6 {
		return nil, ErrInvalidRawDataSize
	}

	amont, err := strconv.ParseInt(raw[1], 10, 64)
	if err != nil {
		return nil, err
	}

	month, err := strconv.Atoi(raw[4])
	if err != nil {
		return nil, err
	}

	year, err := strconv.Atoi(raw[5])
	if err != nil {
		return nil, err
	}

	return &Donation{
		Name:           raw[0],
		AmountSubunits: amont,
		CCNumber:       raw[2],
		CVV:            raw[3],
		ExpMonth:       time.Month(month),
		ExpYear:        year,
	}, nil
}
