package config

import "os"

type Config struct {
	Port        string
	ServiceName string
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8003"),
		ServiceName: getEnv("SERVICE_NAME", "cart-service"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
