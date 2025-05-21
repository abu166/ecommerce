package main

import (
	"ecommerce/internal/apigateway"
	"ecommerce/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := apigateway.Run(cfg); err != nil {
		log.Fatalf("Consumer service failed: %v", err)
	}
}
