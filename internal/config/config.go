package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	APIGatewayAddr string
	InventoryAddr  string
	OrderAddr      string
	UserAddr       string
	ProducerAddr   string
	NATSAddr       string
	RedisAddr      string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	return &Config{
		APIGatewayAddr: getEnv("API_GATEWAY_ADDR", ":8080"),
		InventoryAddr:  getEnv("INVENTORY_ADDR", ":50051"),
		OrderAddr:      getEnv("ORDER_ADDR", ":50052"),
		UserAddr:       getEnv("USER_ADDR", ":50053"),
		ProducerAddr:   getEnv("PRODUCER_ADDR", ":50054"),
		NATSAddr:       getEnv("NATS_ADDR", "nats://localhost:4222"),
		RedisAddr:      getEnv("REDIS_ADDR", "redis:6379"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "admin"),
		DBName:         getEnv("DB_NAME", "ecommerce"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
