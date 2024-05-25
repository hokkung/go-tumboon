package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hokkung/go-tumboon/config"
	"github.com/stretchr/testify/assert"
)

const (
	testDonationFileAddr = "test_data/test.csv"
	testPublicKey        = "test_public_key"
	testPrivateKey       = "test_private_key"
	testMaxConcurrent    = 10
)

func TestNewConfiguration(t *testing.T) {
	os.Setenv("APP_DONATION_FILE_ADDR", testDonationFileAddr)
	os.Setenv("APP_OMISE_PUBLIC_KEY", testPublicKey)
	os.Setenv("APP_OMISE_PRIVATE_KEY", testPrivateKey)
	os.Setenv("APP_MAX_CONCURRENT", fmt.Sprintf("%d", testMaxConcurrent))

	cfg, err := config.NewConfiguration()

	assert.Nil(t, err)
	assert.Equal(t, testDonationFileAddr, cfg.DonationFileConfiguration.DonationFileAddr)
	assert.Equal(t, testMaxConcurrent, cfg.DonationFileConfiguration.MaxConcurrent)
	assert.Equal(t, testPublicKey, cfg.OmiseConfiguration.PublicKey)
	assert.Equal(t, testPrivateKey, cfg.OmiseConfiguration.PrivateKey)
}

func TestDefaultConfiguration(t *testing.T) {
	os.Unsetenv("APP_DONATION_FILE_ADDR")
	os.Unsetenv("APP_MAX_CONCURRENT")
	os.Unsetenv("APP_OMISE_PUBLIC_KEY")
	os.Unsetenv("APP_OMISE_PRIVATE_KEY")

	cfg, err := config.NewConfiguration()

	assert.Nil(t, err)
	assert.Equal(t, "internal/data/fng.1000.csv.rot128", cfg.DonationFileConfiguration.DonationFileAddr)
	assert.Equal(t, 2, cfg.DonationFileConfiguration.MaxConcurrent)
	assert.Equal(t, "pkey_test_no1t4tnemucod0e51mo", cfg.OmiseConfiguration.PublicKey)
	assert.Equal(t, "skey_test_no1t4tnemucod0e51mo", cfg.OmiseConfiguration.PrivateKey)
}
