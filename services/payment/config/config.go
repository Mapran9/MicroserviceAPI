package config

import "os"

type Config struct {
	Port        string
	ServiceName string
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8005"),
		ServiceName: getEnv("SERVICE_NAME", "payment-service"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
