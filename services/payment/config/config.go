package config

import "os"

type Config struct {
	Port        string
	ServiceName string
	InstanceID  string
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8005"),
		ServiceName: getEnv("SERVICE_NAME", "payment-service"),
		InstanceID:  getEnv("INSTANCE_ID", getEnv("HOSTNAME", "unknown")),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
