package omisecli

import (
	"github.com/hokkung/go-tumboon/config"
	"github.com/omise/omise-go"
)

//go:generate mockgen -source=omise.go -destination=./mock/mock_omise.go
type OmiseClient interface {
	Do(result interface{}, op interface{}) error
}

type omiseClient struct {
	*omise.Client
}

func (c omiseClient) Do(result interface{}, op interface{}) error {
	return c.Do(result, op)
}

func NewOmiseClient(cfg config.Configuration) (*omiseClient, func(), error) {
	cli, err := omise.NewClient(
		cfg.OmiseConfiguration.PublicKey,
		cfg.OmiseConfiguration.PrivateKey,
	)
	return &omiseClient{cli}, func() {}, err
}

func ProvideOmiseClient(cfg config.Configuration) (OmiseClient, func(), error) {
	return NewOmiseClient(cfg)
}
