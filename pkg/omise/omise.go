package omisecli

import (
	"github.com/hokkung/go-tumboon/config"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

//go:generate mockgen -source=omise.go -destination=./mock/mock_omise.go

// OmiseClient manages all Omise transactions.
type OmiseClient interface {
	CreateToken(result *omise.Card, createToken *operations.CreateToken) error
	CreateCharge(result *omise.Charge, createCharge *operations.CreateCharge) error
}

type omiseClient struct {
	*omise.Client
}

// NewOmiseClient creates the wrapper for Omise client.
func NewOmiseClient(cfg config.Configuration) (*omiseClient, error) {
	cli, err := omise.NewClient(
		cfg.OmiseConfiguration.PublicKey,
		cfg.OmiseConfiguration.PrivateKey,
	)
	return &omiseClient{cli}, err
}

// ProvideOmiseClient provides the wrapper for Omise client for dependency injection.
func ProvideOmiseClient(
	cfg config.Configuration,
) (OmiseClient, error) {
	return NewOmiseClient(cfg)
}

// CreateToken creates token.
func (c *omiseClient) CreateToken(
	result *omise.Card,
	createToken *operations.CreateToken,
) error {
	return c.Client.Do(result, createToken)
}

// CreateCharge creates a charge.
func (c *omiseClient) CreateCharge(
	result *omise.Charge,
	createCharge *operations.CreateCharge,
) error {
	return c.Client.Do(result, createCharge)
}
