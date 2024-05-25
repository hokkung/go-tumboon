package di

import (
	"github.com/google/wire"
	"github.com/hokkung/go-tumboon/config"
	"github.com/hokkung/go-tumboon/internal/runner"
	ds "github.com/hokkung/go-tumboon/internal/service/donation"
	ps "github.com/hokkung/go-tumboon/internal/service/payment"
	omisecli "github.com/hokkung/go-tumboon/pkg/omise"
	"github.com/hokkung/go-tumboon/internal/validator"
)

var MakePermitRunnerSet = wire.NewSet(
	ConfigSet,
	ClientSet,
	ValidatorSet,
	ServiceSet,
	RunnerSet,
	wire.Struct(new(MakePermitApplication), "*"),
)

var ValidatorSet = wire.NewSet(
	validator.NewCustomValidator,
)

var ConfigSet = wire.NewSet(
	config.ProvideConfiguration,
)

var ClientSet = wire.NewSet(
	omisecli.ProvideOmiseClient,
)

var ServiceSet = wire.NewSet(
	ds.ProvideDonationService,
	ps.ProvideOmisePaymentService,
)

var RunnerSet = wire.NewSet(
	runner.NewDonationRunner,
)

type MakePermitApplication struct {
	Runner *runner.DonationRunner
}
