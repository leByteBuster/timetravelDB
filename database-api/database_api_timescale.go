package databaseapi

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/jackc/pgx/v5"
)

var ConfigTS = TimescaleConfig{}

var SessionTS *pgx.Conn

type TimeSeriesRow struct {
	Timestamp   time.Time
	IsTimestamp bool
	Value       interface{}
}

func ConnectTimescale(username, password, port, dbname string) (*pgx.Conn, error) {

	var sb strings.Builder
	sb.WriteString("postgresql://")
	sb.WriteString(username)
	sb.WriteString(":")
	sb.WriteString(password)
	sb.WriteString("@localhost:")
	sb.WriteString(port)
	sb.WriteString("/")
	sb.WriteString(dbname)

	// conn, err := pgxpool.Connect(context.Background(), connStr) // use pgxpool for managing multiple connections
	conn, err := pgx.Connect(context.Background(), sb.String())
	return conn, err
}

// multiple row READ queries
func ReadRowsTimescale(query string, parameters [][]interface{}) ([]TimeSeriesRow, error) {
	// create the table according to  the data type

	rows, err := SessionTS.Query(context.Background(), query)
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
		res = append(res, TimeSeriesRow{Timestamp: timestamp, IsTimestamp: isTimestamp, Value: value})
	}

	if err := rows.Err(); err != nil {
		log.Println("Error during iteration:", err)
		return nil, err
	}

	return res, nil
}

// read single row from a table
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func ReadRowTimescale(query string, parameters []interface{}) (TimeSeriesRow, error) {

	var timestamp time.Time
	var isTimestamp bool
	var value interface{}

	err := SessionTS.QueryRow(context.Background(), query, parameters...).Scan(&timestamp, &isTimestamp, &value)
	if err != nil {
		return TimeSeriesRow{}, fmt.Errorf("error querying a row from a table in timescaledb: %w", err)
	}

	return TimeSeriesRow{Timestamp: timestamp, IsTimestamp: isTimestamp, Value: value}, nil
}

func ReadRowExistsTimescale(query string) (bool, error) {
	// create the table according to  the data type

	var exists bool
	err := SessionTS.QueryRow(context.Background(), query).Scan(&exists)
	if err != nil {
		utils.Debug("error executing an existence check in timescaledb:", err)
		return false, nil
		// TODO: reintroduce this as soon as we can be sure that for every UUID in a property in neo4j a
		// table exists in timescaledb
		// return false, fmt.Errorf("error executing an existence check in timescaledb: %w", err)
	}

	return exists, nil
}

// read timestamp of a single row
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func ReadRowTimestampTimescale(query string, parameters []interface{}) (interface{}, error) {
	// create the table according to  the data type

	var timestamp time.Time

	err := SessionTS.QueryRow(context.Background(), query, parameters...).Scan(&timestamp)
	if err != nil {
		return nil, err
	}

	return timestamp, nil
}

// read value of a single row
func ReadRowValueTimescale(query string, parameters []interface{}) (interface{}, error) {
	// create the table according to  the data type

	var value interface{}

	err := SessionTS.QueryRow(context.Background(), query, parameters...).Scan(&value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// the following functions are used by the data-adapter

func WriteQueryTimeScale(query string, parameters []interface{}) {
	c, err := SessionTS.Exec(context.Background(), query, parameters...)
	utils.Debug("write query timescale: ", query, "answer: ", c)
	if err != nil {
		log.Printf("%v: error executing timescaledb query: %v", err, query)
	}
}

// writes the same query multiple times with different parameters
func WriteSameQueryMultipleTimeScale(query string, parameters [][]interface{}) {
	for i := range parameters {
		WriteQueryTimeScale(query, parameters[i])
	}
}
