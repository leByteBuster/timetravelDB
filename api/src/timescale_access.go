package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// send query and return result
func sendQuery(url string, sqlQuery string) (pgx.Rows, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v", err)
	}

	// NOTE: maybe i need to remove this from here
	defer conn.Close(context.Background())

	return conn.Query(context.Background(), sqlQuery)
}

// urlExample := "postgres://username:password@localhost:5432/database_name"
func getRow(url string, sqlQueryString string) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var weight int64
	err = conn.QueryRow(context.Background(), sqlQueryString, 42).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, weight)
}
