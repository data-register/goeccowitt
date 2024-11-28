# create_db.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	postgresHost     = os.Getenv("POSTGRES_HOST")
	postgresPort     = os.Getenv("POSTGRES_PORT")
	postgresDB       = os.Getenv("POSTGRES_DB")
	postgresUser     = os.Getenv("POSTGRES_USER")
	postgresPassword = os.Getenv("POSTGRES_PASSWORD")
)

func main() {
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", postgresHost, postgresPort, postgresDB, postgresUser, postgresPassword)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS weather_data (
		id SERIAL PRIMARY KEY,
		data JSONB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Unable to create table: %v", err)
	}

	log.Println("Table created successfully.")
}
