package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// urlExample := "postgres://username:password@localhost:5432/database_name"
func getRowTSDB(url string, sqlQueryString string) {
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
