//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
)

func InitializeMakePermitRunnerApplication() (*MakePermitApplication, func(), error) {
	wire.Build(MakePermitRunnerSet)
	return &MakePermitApplication{}, func(){}, nil
}