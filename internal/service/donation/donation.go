package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"

	v "github.com/go-playground/validator/v10"
	"github.com/hokkung/go-tumboon/config"
	"github.com/hokkung/go-tumboon/internal/model"
	service "github.com/hokkung/go-tumboon/internal/service/payment"
	"github.com/hokkung/go-tumboon/pkg/cipher"
	"github.com/leekchan/accounting"
)

//go:generate mockgen -source=donation.go -destination=./mock/mock_donation.go

// DonationService manages donation domain.
type DonationService interface {
	MakePermit() error
	Donate(donation model.Donation) (*service.PaymentResponse, error)
	Donates(donations []model.Donation) (*SummaryDetail, error)
}

type donationService struct {
	paymentService            service.PaymentService
	donationFileConfiguration config.DonationFileConfiguration
	validator                 *v.Validate
}

// NewDonationService creates donation service.
func NewDonationService(
	paymentService service.PaymentService,
	cfg config.Configuration,
	validator *v.Validate,
) *donationService {
	return &donationService{
		donationFileConfiguration: *cfg.DonationFileConfiguration,
		validator:                 validator,
		paymentService:            paymentService,
	}
}

// ProvideDonationService provides donation service for dependency injection.
func ProvideDonationService(
	paymentService service.PaymentService,
	cfg config.Configuration,
	validator *v.Validate,
) DonationService {
	return NewDonationService(
		paymentService,
		cfg,
		validator,
	)
}

// MakePermit performs donation process by reading all information from CSV file and reports the summary result.
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

	s.buildSummaryReport(summary)

	return nil
}

func (s donationService) buildSummaryReport(summaryDetail *SummaryDetail) {
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

// Donates performs donation process by given a list of donations and returns summary result.
// The process in this method will be performed concurrently using Go channel and wait group.
// The error arises throughout the donation process will not be addressed until the entire procedure is finished.
func (s donationService) Donates(donations []model.Donation) (*SummaryDetail, error) {
	if len(donations) <= 0 {
		return &SummaryDetail{}, nil
	}

	var faultyDonated int64
	var successfulDonated int64
	var totalReceived int64
	donorToTotalAmount := make(map[string]int64)

	var wg sync.WaitGroup
	ch := make(chan *service.PaymentResponse, s.donationFileConfiguration.MaxConcurrent)
	defer close(ch)

	go func() {
		for res := range ch {
			if res.IsSuccess {
				successfulDonated += res.Amount
			} else {
				faultyDonated += res.Amount
			}
			totalReceived += res.Amount
			donorToTotalAmount[res.Source.Name] += res.Amount
		}
	}()

	for _, donation := range donations {
		wg.Add(1)
		go func(donation model.Donation) {
			defer wg.Done()
			res, err := s.Donate(donation)
			if err != nil {
				fmt.Println(err, "donation has been failed")
			}
			ch <- res
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

	topDonors := aggregatedDonors
	if len(aggregatedDonors) > numberOfTopHighestDonation {
		topDonors = aggregatedDonors[:numberOfTopHighestDonation]
	}

	return topDonors
}

// Donate performs a single donation process.
func (s donationService) Donate(donation model.Donation) (*service.PaymentResponse, error) {
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
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		if s.isHeaderFile(row) {
			continue
		}

		donation, err := model.NewDonation(row)
		if err != nil {
			fmt.Println("read row from csv failed", err)
			continue
		}

		err = s.validator.Struct(donation)
		if err != nil {
			fmt.Println("validate donation struct failed", donation, err)
			continue
		}

		donations = append(donations, *donation)
	}

	return donations, nil
}

func (s donationService) isHeaderFile(row []string) bool {
	return row[0] == "Name" &&
		row[1] == "AmountSubunits" &&
		row[2] == "CCNumber" &&
		row[3] == "CVV" &&
		row[4] == "ExpMonth" &&
		row[5] == "ExpYear"
}
