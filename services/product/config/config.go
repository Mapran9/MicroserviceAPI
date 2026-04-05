package config

import "os"

type Config struct {
	Port        string
	ServiceName string
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8002"),
		ServiceName: getEnv("SERVICE_NAME", "product-service"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
