package service

type SummaryDetail struct {
	TotalReceived     int64
	SuccessfulDonated int64
	FaultyDonated     int64
	AveragePerPerson  float64
	TopDonors         []Donor
}

type Donor struct {
	Name   string
	Amount int64
}
