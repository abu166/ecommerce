package main

import (
	"ecommerce/internal/apigateway"
	"ecommerce/internal/config"
	"ecommerce/internal/inventory"
	"ecommerce/internal/order"
	"ecommerce/internal/user"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	go func() {
		if err := inventory.Run(cfg); err != nil { // Call transport.Run for inventory
			log.Fatalf("Inventory service failed: %v", err)
		}
	}()

	go func() {
		if err := order.Run(cfg); err != nil { // Call transport.Run for order
			log.Fatalf("Order service failed: %v", err)
		}
	}()

	go func() {
		if err := user.Run(cfg); err != nil { // Call transport.Run for user
			log.Fatalf("User service failed: %v", err)
		}
	}()

	if err := apigateway.Run(cfg); err != nil {
		log.Fatalf("API Gateway failed: %v", err)
	}
}
