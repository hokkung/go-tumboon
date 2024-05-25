package runner

import (
	service "github.com/hokkung/go-tumboon/internal/service/donation"
)

type DonationRunner struct {
	donationService service.DonationService
}

func (r DonationRunner) Run() error {
	return r.donationService.MakePermit()
}

func NewDonationRunner(
	donationService service.DonationService,
) *DonationRunner {
	return &DonationRunner{
		donationService: donationService,
	}
}
