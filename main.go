# main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

const apiUrl = "https://api.ecowitt.net/api/v3/device/real_time?application_key=61759DF4094EBA1A61E070A285A2DAF7&api_key=fe33f769-d487-433c-a31d-8e504df4076f&mac=48:E7:29:5F:72:44&call_back=all"

func main() {
	http.HandleFunc("/fetch_store", fetchAndStoreHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchAndStoreHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(apiUrl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data from Ecowitt API: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding API response: %v", err), http.StatusInternalServerError)
		return
	}

	if err := storeDataToPostgres(data); err != nil {
		http.Error(w, fmt.Sprintf("Error storing data to PostgreSQL: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Data fetched and stored successfully.")
}

func storeDataToPostgres(data map[string]interface{}) error {
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", postgresHost, postgresPort, postgresDB, postgresUser, postgresPassword)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer db.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v", err)
	}

	query := `INSERT INTO weather_data (data) VALUES ($1)`
	_, err = db.Exec(query, jsonData)
	if err != nil {
		return fmt.Errorf("unable to execute insert query: %v", err)
	}

	return nil
}
