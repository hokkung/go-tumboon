package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"sync/atomic"

	v "github.com/go-playground/validator/v10"
	"github.com/hokkung/go-tumboon/config"
	"github.com/hokkung/go-tumboon/internal/model"
	service "github.com/hokkung/go-tumboon/internal/service/payment"
	"github.com/hokkung/go-tumboon/pkg/cipher"
	"github.com/leekchan/accounting"
)

//go:generate mockgen -source=donation.go -destination=./mock/mock_donation.go
type DonationService interface {
	MakePermit() error
	Donate(donation model.Donation) error
	Donates(donations []model.Donation) (*SummaryDetail, error)
}

type donationService struct {
	paymentService            service.PaymentService
	donationFileConfiguration config.DonationFileConfiguration
	validator                 *v.Validate
}

func NewDonationService(
	paymentService service.PaymentService,
	cfg config.Configuration,
	v *v.Validate,
) *donationService {
	return &donationService{
		donationFileConfiguration: *cfg.DonationFileConfiguration,
		validator:                 v,
		paymentService:            paymentService,
	}
}

func ProvideDonationService(
	paymentService service.PaymentService,
	cfg config.Configuration,
	v *v.Validate,
) DonationService {
	return NewDonationService(paymentService, cfg, v)
}

func (s donationService) MakePermit() error {
	fmt.Println("performing donations...")

	donations, err := s.getDonationDetailFromFile()
	if err != nil {
		return err
	}

	summary, err := s.Donates(donations)
	if err != nil {
		return err
	}

	fmt.Println("done.")

	s.report(summary)

	return nil
}

func (s donationService) report(summaryDetail *SummaryDetail) {
	var topDonorsStr string
	for i, topDonors := range summaryDetail.TopDonors {
		if i != 0 {
			topDonorsStr += "\t\t"
		}
		topDonorsStr += topDonors.Name
		topDonorsStr += "\n"
	}

	ac := accounting.Accounting{Symbol: "", Precision: 2}
	fmt.Printf("total received:\t\tTHB  %s\nsuccessfully donated:\tTHB  %s\nfaulty donation:\tTHB  %s\naverage per person:\tTHB  %s\ntop donors:\t%s",
		ac.FormatMoney(summaryDetail.TotalReceived),
		ac.FormatMoney(summaryDetail.SuccessfulDonated),
		ac.FormatMoney(summaryDetail.FaultyDonated),
		ac.FormatMoney(summaryDetail.AveragePerPerson),
		topDonorsStr,
	)
}

func (s donationService) Donates(donations []model.Donation) (*SummaryDetail, error) {
	if len(donations) <= 0 {
		return &SummaryDetail{}, nil
	}

	var faultyDonated int64
	var successfulDonated int64
	var totalReceived int64
	donorToTotalAmount := make(map[string]int64)

	var wg sync.WaitGroup
	var mu sync.Mutex
	limiter := make(chan struct{}, s.donationFileConfiguration.MaxConcurrent)

	for _, donation := range donations {
		wg.Add(1)
		limiter <- struct{}{}

		go func(donation model.Donation) {
			defer wg.Done()
			defer func() { <-limiter }()

			err := s.Donate(donation)
			if err != nil {
				atomic.AddInt64(&faultyDonated, donation.AmountSubunits)
			} else {
				atomic.AddInt64(&successfulDonated, donation.AmountSubunits)
			}
			atomic.AddInt64(&totalReceived, donation.AmountSubunits)

			mu.Lock()
			donorToTotalAmount[donation.Name] += donation.AmountSubunits 
			mu.Unlock()
		}(donation)
	}

	wg.Wait()

	numberOfDonor := len(donorToTotalAmount)
	topThreeDonors := s.getTopDonors(
		donorToTotalAmount,
		numberOfDonor,
		3,
	)

	return &SummaryDetail{
		TotalReceived:     totalReceived,
		SuccessfulDonated: successfulDonated,
		FaultyDonated:     faultyDonated,
		AveragePerPerson:  float64(totalReceived) / float64(numberOfDonor),
		TopDonors:         topThreeDonors,
	}, nil
}

func (s donationService) getTopDonors(
	donorToTotalAmount map[string]int64,
	numberOfDonor int,
	numberOfTopHighestDonation int,
) []Donor {
	aggregatedDonors := make([]Donor, 0, numberOfDonor)
	for name, amount := range donorToTotalAmount {
		aggregatedDonors = append(aggregatedDonors, Donor{Name: name, Amount: amount})
	}

	sort.Slice(aggregatedDonors, func(i, j int) bool {
		return aggregatedDonors[i].Amount > aggregatedDonors[j].Amount
	})

	topThreeDonors := aggregatedDonors
	if len(aggregatedDonors) > numberOfTopHighestDonation {
		topThreeDonors = aggregatedDonors[:numberOfTopHighestDonation]
	}

	return topThreeDonors
}

func (s donationService) Donate(donation model.Donation) error {
	return s.paymentService.Do(service.PaymentRequest{
		Name:           donation.Name,
		AmountSubunits: donation.AmountSubunits,
		CCNumber:       donation.CCNumber,
		CVV:            donation.CVV,
		ExpMonth:       donation.ExpMonth,
		ExpYear:        donation.ExpYear,
		Type:           service.Card,
	})
}

func (s donationService) getDonationDetailFromFile() ([]model.Donation, error) {
	file, err := os.Open(s.donationFileConfiguration.DonationFileAddr)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader, err := cipher.NewRot128Reader(file)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	b := make([]byte, stat.Size())
	_, err = reader.Read(b)
	if err != nil {
		return nil, err
	}

	var donations []model.Donation
	r := csv.NewReader(bytes.NewReader(b))
	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		if s.isHeaderFile(record) {
			continue
		}

		donation, err := model.NewDonation(record)
		if err != nil {
			fmt.Println("read record failed", err)
			continue
		}

		err = s.validator.Struct(donation)
		if err != nil {
			fmt.Println("validate donation failed", donation, err)
			continue
		}

		donations = append(donations, *donation)
	}

	return donations, nil
}

func (s donationService) isHeaderFile(record []string) bool {
	return record[0] == "Name" &&
		record[1] == "AmountSubunits" &&
		record[2] == "CCNumber" &&
		record[3] == "CVV" &&
		record[4] == "ExpMonth" &&
		record[5] == "ExpYear"
}
