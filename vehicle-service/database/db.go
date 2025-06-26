package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	for i := 1; i <= 10; i++ {
		DB, err = sql.Open("postgres", dsn)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				log.Println("🚗 Vehicle DB connected")
				return
			}
		}
		log.Printf("Retrying Vehicle DB connection... (%d/10)\n", i)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf(" Vehicle DB unreachable after retries: %v", err)
}
