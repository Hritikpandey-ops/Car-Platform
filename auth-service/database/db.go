package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	var db *sql.DB
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				DB = db
				fmt.Println("Connected to database ✅")
				return nil
			}
		}
		fmt.Printf("Retrying DB connection... (%d/10)\n", i+1)
		time.Sleep(3 * time.Second)
	}

	return fmt.Errorf("could not connect to DB: %v", err)
}
