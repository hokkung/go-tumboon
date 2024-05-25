package config

import "github.com/kelseyhightower/envconfig"

const APP_PREFIX string = "APP"

type Configuration struct {
	DonationFileConfiguration *DonationFileConfiguration
	OmiseConfiguration        *OmiseConfiguration
}

type DonationFileConfiguration struct {
	DonationFileAddr string `envconfig:"DONATION_FILE_ADDR" default:"internal/data/fng.1000.csv.rot128"`
	MaxConcurrent    int    `envconfig:"MAX_CONCURRENT" default:"2"`
}

type OmiseConfiguration struct {
	PublicKey  string `envconfig:"OMISE_PUBLIC_KEY" default:"pkey_test_no1t4tnemucod0e51mo"`
	PrivateKey string `envconfig:"OMISE_PRIVATE_KEY" default:"skey_test_no1t4tnemucod0e51mo"`
}

func NewConfiguration() (*Configuration, error) {
	var donationFileConfiguration DonationFileConfiguration
	err := envconfig.Process(APP_PREFIX, &donationFileConfiguration)
	if err != nil {
		return nil, err
	}

	var omiseConfiguration OmiseConfiguration
	err = envconfig.Process(APP_PREFIX, &omiseConfiguration)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		DonationFileConfiguration: &donationFileConfiguration,
		OmiseConfiguration:        &omiseConfiguration,
	}, nil
}

func ProvideConfiguration() (Configuration, error) {
	cfg, err := NewConfiguration()
	return *cfg, err
}
