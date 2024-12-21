package main

import (
	"log"

	"github.com/idoyudha/eshop-order/config"
	"github.com/idoyudha/eshop-order/internal/app"
)

func main() {
	log.Println("Hello World! Eshop Order Service!")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
