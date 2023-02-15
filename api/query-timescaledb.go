package api

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

// get single property without aggregation
func getProperty(from, to, table string) interface{} {

	// build query string
	//query := fmt.Sprintf("SELECT %s(value) FROM %s WHERE time >= '%s' AND time < '%s'", aggr, table, from, to)
	//fmt.Printf("Query: %v\n", query)

	// res := queryTimeScale(query, username, password, port, dbname)
	val, err := readRowValueTimescale(table, nil, UserTS, PassTS, PortTS, DBnameTS)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}

	return val
}

// get property / single time-series for the period (from,to], apply aggregation on it if not empty
func getPropertyAggr(from, to, aggr, table string) interface{} {

	// build query string
	query := fmt.Sprintf("SELECT %s(value) FROM %s WHERE time >= '%s' AND time < '%s'", aggr, table, from, to)
	fmt.Printf("Query: %v\n", query)

	// res := queryTimeScale(query, username, password, port, dbname)
	val, err := readRowValueTimescale(table, nil, UserTS, PassTS, PortTS, DBnameTS)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}

	return val
}

func getPropertyFromTable(from, to, aggr string, tablename string) (interface{}, []TimeSeriesRow, error) {
	return getPropertiesFromTables(from, to, "", []string{tablename})
}

// get properties / multiple time-series and apply aggrergation on it if not empty
func getPropertiesFromTables(from, to, aggr string, tables []string) (interface{}, []TimeSeriesRow, error) {

	var builder strings.Builder

	// build query string
	builder.WriteString("SELECT ")
	// builder.WriteString(aggr) cannot do this here because then i only have one return value
	// have to call a different subfunction aggrTimescale or something
	builder.WriteString("*")
	builder.WriteString(" FROM (")
	for i, tablename := range tables {
		if i > 0 {
			builder.WriteString(" UNION ALL ")
		}
		builder.WriteString("SELECT ")
		builder.WriteString("time, timestamps, value FROM ")
		builder.WriteString(tablename)
		builder.WriteString(" WHERE time >= ")
		builder.WriteString("'")
		builder.WriteString(from)
		builder.WriteString("'")
		builder.WriteString(" AND time < ")

		//TODO: CHANGE THIS TO DATETIME
		builder.WriteString("'")
		builder.WriteString(to)
		builder.WriteString("'")
	}
	builder.WriteString(") genericAliasName;")
	fmt.Println(builder.String())

	// Get the query string from the StringBuilder
	query := builder.String()
	// fmt.Println(query)

	rows, err := readRowsTimescale(query, nil, UserTS, PassTS, PortTS, DBnameTS)

	if err != nil {
		log.Println(err)
	}

	if aggr != "" && len(rows) == 1 {
		return rows[0].Value, nil, nil
	} else if aggr != "" && len(rows) > 1 {
		return nil, nil, errors.New("aggregation returned more than one row")
	}

	return nil, rows, nil
}

func uuidToTablename(uuid string) string {
	var builder strings.Builder
	builder.WriteString("ts_")
	builder.WriteString(strings.Replace(uuid, "-", "_", -1))
	return builder.String()
}
