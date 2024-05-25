package main

import (
	"github.com/hokkung/go-tumboon/internal/di"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	app, cleanFn, err := di.InitializeMakePermitRunnerApplication()
	if err != nil {
		panic(err)
	}

	defer cleanFn()

	err = app.Runner.Run()
	if err != nil {
		panic(err)
	}
}
