package runner

import (
	service "github.com/hokkung/go-tumboon/internal/service/donation"
)

// DonationRunner runs overall donation process.
type DonationRunner struct {
	donationService service.DonationService
}

// Run runs the process.
func (r DonationRunner) Run() error {
	return r.donationService.MakePermit()
}

// NewDonationRunner creates donation application runner.
func NewDonationRunner(
	donationService service.DonationService,
) *DonationRunner {
	return &DonationRunner{
		donationService: donationService,
	}
}
