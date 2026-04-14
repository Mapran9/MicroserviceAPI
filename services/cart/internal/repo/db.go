package repo

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsn := os.Getenv("DB_DSN")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open db error: %v", err)
	}

	db.SetMaxOpenConns(getEnvInt("DB_MAX_OPEN_CONNS", 100))
	db.SetMaxIdleConns(getEnvInt("DB_MAX_IDLE_CONNS", 25))
	db.SetConnMaxLifetime(getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute))
	db.SetConnMaxIdleTime(getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), getEnvDuration("DB_PING_TIMEOUT", 3*time.Second))
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping db error: %v", err)
	}

	DB = db
	log.Println("cart-service connected to database")
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(v)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
