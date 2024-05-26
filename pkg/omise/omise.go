package omisecli

import (
	"github.com/hokkung/go-tumboon/config"
	"github.com/omise/omise-go"
)

//go:generate mockgen -source=omise.go -destination=./mock/mock_omise.go

// OmiseClient manages all Omise transactions.
type OmiseClient interface {
	Do(result interface{}, op interface{}) error
}

type omiseClient struct {
	*omise.Client
}

// Do performs tracsaction using Omise client.
func (c omiseClient) Do(result interface{}, op interface{}) error {
	return c.Do(result, op)
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
func ProvideOmiseClient(cfg config.Configuration) (OmiseClient, error) {
	return NewOmiseClient(cfg)
}
