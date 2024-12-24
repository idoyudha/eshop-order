package main

import (
	"log"

	"github.com/idoyudha/eshop-order/config"
	"github.com/idoyudha/eshop-order/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
