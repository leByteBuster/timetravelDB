package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
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
var DriverNeo neo4j.DriverWithContext
var SessionNeo neo4j.SessionWithContext

type TimeSeriesRow struct {
	Timestamp   time.Time
	IsTimestamp bool
	Value       interface{}
}

// send any read query and return the results as a key value map
func queryReadNeo4j(ctx context.Context, cypherQueryString string) (neo4j.ResultWithContext, error) {
	// Connect to the Neo4j database

	res, errReq := SessionNeo.Run(ctx, cypherQueryString, map[string]interface{}{})

	// fmt.Printf("\nDirect print: %v\n", res)
	// res.Next(ctx)
	// fmt.Printf("\nKeys: %v", res.Record().Keys)
	// data, _ := res.Record().Get("n")
	// fmt.Printf("\nn: %v", data)

	if errReq != nil {
		log.Printf("Query failed")
		return nil, errReq
	}

	return res, nil
}

//lint:ignore U1000 Ignore unused function temporarily for debugging

func queryWriteNeo4j(ctx context.Context, uri, username, password, cypherQueryString string) {

	session := DriverNeo.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, errReq := session.Run(ctx, cypherQueryString, map[string]interface{}{})
	if errReq != nil {
		fmt.Println(errReq)
		return
	}
}

//lint:ignore U1000 Ignore unused function temporarily for debugging

func queryWriteMultipleNeo4j(ctx context.Context, uri, username, password string, cypherQueryStrings []string) {
	// Connect to the Neo4j database

	session := DriverNeo.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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

// send any query string to the database
// this is just for executing, CommandTag only contains status
func TimeScale(query string) pgconn.CommandTag {
	conn := connectTimescale(UserTS, PassTS, PortTS, DBnameTS)
	defer conn.Close(context.Background())
	res, err := conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

// send a list of query strings to the database. not sure if CommandTag contains a result or just status though
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
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
		log.Println("error querying rows from a table in timescaledb:", err)
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
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func readRowTimescale(query string, parameters []interface{}, username, password, port, dbname string) (TimeSeriesRow, error) {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())

	var timestamp time.Time
	var isTimestamp bool
	var value interface{}

	err := conn.QueryRow(context.Background(), query, parameters...).Scan(&timestamp, &isTimestamp, &value)
	if err != nil {
		return TimeSeriesRow{}, fmt.Errorf("error querying a row from a table in timescaledb: %w", err)
	}

	return TimeSeriesRow{timestamp, isTimestamp, value}, nil
}

func readRowExistsTimescale(query string, username, password, port, dbname string) (bool, error) {
	// create the table according to  the data type
	conn := connectTimescale(username, password, port, dbname)
	defer conn.Close(context.Background())

	var exists bool
	err := conn.QueryRow(context.Background(), query).Scan(&exists)
	if err != nil {
		return false, nil
		// TODO: reintroduce this as soon as we can be sure that for every UUID in a property in neo4j a
		// table exists in timescaledb
		// return false, fmt.Errorf("error executing an existence check in timescaledb: %w", err)
	}

	return exists, nil
}

// single row only return time READ queries
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
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
//lint:ignore U1000 Ignore unused function temporarily for debugging

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
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
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
