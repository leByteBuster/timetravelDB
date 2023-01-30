package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type TimeSeriesRow struct {
	time      string
	timestamp bool
	value     interface{}
}

// urlExample := "postgres://username:password@localhost:5431/database_name"
// NOTE: dont remember what this was supposed to do
func getRowTimescale(username, password, port, dbname, sqlQueryString string) {
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())

	var name string
	var weight int64
	err := conn.QueryRow(context.Background(), sqlQueryString, 41).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(0)
	}

	fmt.Println(name, weight)
}

func queryNeo4j(ctx context.Context, uri, username, password, cypherQueryString string) {
	// Connect to the Neo4j database
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, errReq := session.Run(ctx, cypherQueryString, map[string]interface{}{})
	if errReq != nil {
		fmt.Println(err)
		return
	}
}

func queryMultipleNeo4j(ctx context.Context, uri, username, password string, cypherQueryStrings []string) {
	// Connect to the Neo4j database
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	for _, cypherQueryString := range cypherQueryStrings {
		_, err := session.Run(ctx, cypherQueryString, map[string]interface{}{})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func connectTimescale(username, password, port, dbname string) *pgx.Conn {
	connStr := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s", username, password, port, dbname)
	// conn, err := pgxpool.Connect(context.Background(), connStr)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Println("\nUnable to establish connection:", err)
	}
	return conn
}

// send any query string to the database. not sure if CommandTag contains a result or just status though
func queryTimeScale(query, username, password, port, dbname string) pgconn.CommandTag {
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())
	res, err := conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

// send a list of query strings to the database. not sure if CommandTag contains a result or just status though
func queryMultipleTimeScale(queries []string, parameters [][]interface{}, username, password, port, dbname string) []pgconn.CommandTag {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	results := make([]pgconn.CommandTag, 0)
	defer conn.Close(context.Background())
	for i, query := range queries {
		_, err := conn.Exec(context.Background(), query, parameters[i]...)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, pgconn.CommandTag{})
	}
	return results
}

func readRowsTimescale(query string, parameters [][]interface{}, username, password, port, dbname string) []TimeSeriesRow {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())

	var time string
	var timestamps bool
	var value interface{}
	rows := make([]TimeSeriesRow, 0)
	for i, _ := range parameters {
		err := conn.QueryRow(context.Background(), query, parameters[i]...).Scan(&time, &timestamps, &value)
		if err != nil {
			fmt.Println(err)
		}
		rows = append(rows, TimeSeriesRow{time, timestamps, value})
	}

	return rows
}

func readRowTimescale(query string, parameters []interface{}, username, password, port, dbname string) TimeSeriesRow {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())

	var time string
	var timestamps bool
	var value interface{}

	err := conn.QueryRow(context.Background(), query, parameters...).Scan(&time, &timestamps, &value)
	if err != nil {
		fmt.Println(err)
	}

	return TimeSeriesRow{time, timestamps, value}
}
