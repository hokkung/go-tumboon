package service


// SummaryDetail collects summary report information.
type SummaryDetail struct {
	TotalReceived     int64
	SuccessfulDonated int64
	FaultyDonated     int64
	AveragePerPerson  float64
	TopDonors         []Donor
}

// Donor collects donor information.
type Donor struct {
	Name   string
	Amount int64
}
