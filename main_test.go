package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestActiveConnections(t *testing.T) {
	// Set up the PostgreSQL database connection string.
	connStr := os.Getenv("CONNECTION_STRING")
	if connStr == "" {
		connStr = "user=test dbname=test sslmode=disable"
	}

	connection_number := os.Getenv("CONNECTION_NUMBER")
	numConnections, _ := strconv.Atoi(connection_number)
	if numConnections == 0 {
		numConnections = 10
		log.Println("connection number config not set, using default to: ", numConnections)
	}
	log.Println("configured connection number ", numConnections)

	// Create a connection pool.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error opening the database: %v", err)
	}

	conns := make([]*sql.Tx, 0)
	for i := 0; i < numConnections; i++ {
		conn, err := db.Begin()
		if err != nil {
			log.Fatalf("Error opening connection: %v", err)
		}
		conns = append(conns, conn)
	}

	numActiveConnections := db.Stats().OpenConnections
	fmt.Printf("active connections: %d\n", numActiveConnections)
	if numActiveConnections != numConnections {
		t.Fatalf("Incorrect number of active connections. Expected: %d, Got: %d", numConnections, numActiveConnections)
	} else {
		log.Printf("number of active connections satisfied. Expected: %d, Got: %d", numConnections, numActiveConnections)
	}
	time.Sleep(1000 * time.Millisecond)

	for _, c := range conns {
		c.Rollback()
	}

	db.Close()
	numActiveConnections = db.Stats().OpenConnections
	fmt.Printf("active connections at the end: %d\n", numActiveConnections)
}
