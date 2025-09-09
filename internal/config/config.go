package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL        string
	GRPCPort     string
	HTTPPort     string
	RedisAddr    string
	KafkaBrokers string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		DBURL:        os.Getenv("DB_URL"),        
		GRPCPort:     os.Getenv("GRPC_PORT"),     
		HTTPPort:     os.Getenv("HTTP_PORT"),     
		RedisAddr:    os.Getenv("REDIS_ADDR"),    
		KafkaBrokers: os.Getenv("KAFKA_BROKERS"), 
	}, nil
}
