package main

import (
	"ecommerce/internal/config"
	"ecommerce/internal/consumer"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := consumer.Run(cfg); err != nil {
		log.Fatalf("Consumer service failed: %v", err)
	}
}
