package main

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "amakhyon"
	password = "87654321"
	dbname   = "aqaristudent"
)

func main() {
	// Connect to the PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Retrieve data from the Seat table
	rows, err := db.Query("SELECT id, student FROM Seat ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create a map to store the new seat IDs
	newSeatIDs := make(map[int]int)

	// Process each row and determine the new seat ID
	for rows.Next() {
		var id, newID int
		var student string
		err := rows.Scan(&id, &student)
		if err != nil {
			log.Fatal(err)
		}

		if id%2 == 1 && id < getMaxSeatID(db) {
			newID = id + 1
		} else if id%2 == 0 {
			newID = id - 1
		} else {
			newID = id
		}

		newSeatIDs[id] = newID
	}

	// Update the Seat table with the new seat IDs
	for oldID, newID := range newSeatIDs {
		_, err := db.Exec("UPDATE Seat SET id = $1 WHERE id = $2", newID, oldID)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Seat IDs swapped successfully.")
}

// getMaxSeatID returns the maximum seat ID from the Seat table
func getMaxSeatID(db *sql.DB) int {
	var maxID int
	err := db.QueryRow("SELECT MAX(id) FROM Seat").Scan(&maxID)
	if err != nil {
		log.Fatal(err)
	}
	return maxID
}
