package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// TODO: change this so it isnt hard coded anymore. Should be read from config or so
// AUTH timescaleDB
var UserTS = "postgres"
var PassTS = "password"
var PortTS = "5432"
var DBnameTS = "postgres"

// AUTH neo4j
var UriNeo = "neo4j://localhost:7687"
var UserNeo = "neo4j"
var PassNeo = "rhebo"
var driverNeo neo4j.DriverWithContext

type TimeSeriesRow struct {
	timestamp   time.Time
	isTimestamp bool
	value       interface{}
}

// send any read query and return the results as a key value map
func queryReadNeo4j(ctx context.Context, cypherQueryString string) (any, error) {
	// Connect to the Neo4j database

	session := driverNeo.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	res, errReq := session.Run(ctx, cypherQueryString, map[string]interface{}{})
	// Loop through the records and get the values

	if errReq != nil {
		log.Printf("Query failed")
		return nil, errReq
	}

	// try res.Collect(ctx) instead of res.Next(ctx)
	if res.Next(ctx) {
		record := res.Record()
		if len(record.Keys) > 1 {
			panic("Query returned more than one node (or other value)")
		}
		fmt.Println("as expected")
		return record.Values[0], nil
	}

	return []any{}, nil
}

func queryWriteNeo4j(ctx context.Context, uri, username, password, cypherQueryString string) {

	session := driverNeo.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, errReq := session.Run(ctx, cypherQueryString, map[string]interface{}{})
	if errReq != nil {
		fmt.Println(errReq)
		return
	}
}

func queryWriteMultipleNeo4j(ctx context.Context, uri, username, password string, cypherQueryStrings []string) {
	// Connect to the Neo4j database

	session := driverNeo.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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
func queryTimeScale(query, username, password, port, dbname string) string {
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())
	res, err := conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Println(err)
	}
	return res.String()
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

// multiple row READ queries
func readRowsTimescale(query string, parameters [][]interface{}, username, password, port, dbname string) ([]TimeSeriesRow, error) {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Println("Query failed:", err)
		return nil, err
	}
	defer rows.Close()

	res := make([]TimeSeriesRow, 0)

	for rows.Next() {

		var timestamp time.Time
		var isTimestamp bool
		var value interface{}
		if err := rows.Scan(&timestamp, &isTimestamp, &value); err != nil {
			log.Println("Scan failed:", err)
			return nil, err
		}
		res = append(res, TimeSeriesRow{timestamp, isTimestamp, value})
	}

	if err := rows.Err(); err != nil {
		log.Println("Error during iteration:", err)
		return nil, err
	}

	return res, nil
}

// single row READ queries
func readRowTimescale(query string, parameters []interface{}, username, password, port, dbname string) TimeSeriesRow {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())

	var timestamp time.Time
	var isTimestamp bool
	var value interface{}

	err := conn.QueryRow(context.Background(), query, parameters...).Scan(&timestamp, &isTimestamp, &value)
	if err != nil {
		fmt.Println(err)
	}

	return TimeSeriesRow{timestamp, isTimestamp, value}
}

// single row only return time READ queries
func readRowTimestampTimescale(query string, parameters []interface{}, username, password, port, dbname string) (interface{}, error) {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())

	var timestamp time.Time

	err := conn.QueryRow(context.Background(), query, parameters...).Scan(&timestamp)
	if err != nil {
		return nil, err
	}

	return timestamp, nil
}

// single row only return time READ queries
func readRowValueTimescale(query string, parameters []interface{}, username, password, port, dbname string) (interface{}, error) {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())

	var value interface{}

	err := conn.QueryRow(context.Background(), query, parameters...).Scan(&value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// single row only return time READ queries
func readRowIsTimestampTimescale(query string, parameters []interface{}, username, password, port, dbname string) (interface{}, error) {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())

	var isTimestamp bool

	err := conn.QueryRow(context.Background(), query, parameters...).Scan(&isTimestamp)
	if err != nil {
		return nil, err
	}

	return isTimestamp, nil
}
