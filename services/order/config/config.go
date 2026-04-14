package config

import "os"

type Config struct {
	Port           string
	ServiceName    string
	DBDSN          string
	CartBaseURL    string
	PaymentBaseURL string
	InstanceID     string
}

func Load() Config {
	return Config{
		Port:           getEnv("PORT", "8004"),
		ServiceName:    getEnv("SERVICE_NAME", "order-service"),
		DBDSN:          getEnv("DB_DSN", ""),
		CartBaseURL:    getEnv("CART_BASE_URL", "http://localhost:8003"),
		PaymentBaseURL: getEnv("PAYMENT_BASE_URL", "http://localhost:8005"),
		InstanceID:     getEnv("INSTANCE_ID", getEnv("HOSTNAME", "unknown")),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
