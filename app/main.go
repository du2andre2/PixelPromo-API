package main

import (
	"fmt"
	"go.uber.org/fx"
	"pixelPromo/bootstrap"
)

func main() {
	fmt.Println("Microservice started")

	fx.New(
		bootstrap.Module,
	).Run()
}
