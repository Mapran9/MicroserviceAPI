package repo

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsn := os.Getenv("DB_DSN")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open db error: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db error: %v", err)
	}

	DB = db
	log.Println("cart-service connected to database")
}
