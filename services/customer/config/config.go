package config

import "os"

type Config struct {
	ServiceName string
	Port        string
	DBDSN       string
}

func Load() Config {
	return Config{
		ServiceName: getenv("SERVICE_NAME", "service"),
		Port:        getenv("PORT", "8000"),
		DBDSN:       getenv("DB_DSN", ""),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
