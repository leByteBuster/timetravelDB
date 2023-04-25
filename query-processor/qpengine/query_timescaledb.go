package qpengine

import (
	"errors"
	"fmt"
	"log"

	databaseapi "github.com/LexaTRex/timetravelDB/database-api"
	"github.com/LexaTRex/timetravelDB/utils"
)

// get single property without aggregation
func getProperty(from, to, table string) interface{} {

	// build query string
	//query := fmt.Sprintf("SELECT %s(value) FROM %s WHERE time >= '%s' AND time < '%s'", aggr, table, from, to)
	//fmt.Printf("Query: %v\n", query)

	// res := queryTimeScale(query, username, password, port, dbname)
	val, err := databaseapi.ReadRowValueTimescale(table, nil)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	return val
}

// get property / single time-series for the period (from,to], apply aggregation on it if not empty
func getPropertyAggr(from, to, aggr, table string) interface{} {

	// build query string
	query := fmt.Sprintf("SELECT %s(value) FROM %s WHERE time >= '%s' AND time < '%s'", aggr, table, from, to)
	utils.Debugf("\nTimescale Query: %v\n", query)

	// res := queryTimeScale(query, username, password, port, dbname)
	val, err := databaseapi.ReadRowValueTimescale(table, nil)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	return val
}

// applies a comparison filter on the time-series entries and returns the remaining entries
// fucntionality for aggregation function not yet implemented
func getTimeSeriesCmp(from, to, aggrOp string, cmpOp string, cmpVal any, lookupLeft bool, tablename string) (interface{}, []databaseapi.TimeSeriesRow, error) {
	queryString, err := buildQueryString(from, to, aggrOp, cmpOp, cmpVal, lookupLeft, []string{tablename})
	utils.Debugf("\n TIMESCALEDB QUERY: %v\n", queryString)
	if err != nil {
		return nil, nil, fmt.Errorf("error building query string: %v", err)
	}
	return queryProperties(queryString, aggrOp)
}

func checkCondInTimeseries(from, to, aggrOp string, cmpOp string, cmpVal any, lookupLeft bool, tablename string) (bool, error) {
	queryString, err := buildQueryStringCmpExists(from, to, aggrOp, cmpOp, cmpVal, lookupLeft, tablename)
	utils.Debugf("\n TIMESCALEDB QUERY: %v\n", queryString)
	if err != nil {
		return false, fmt.Errorf("error building query string: %v", err)
	}
	return existenceQuery(queryString)
}

// fucntionality for aggregation function not yet implemented
func getTimeSeries(from, to, aggrOp, tablename string) (interface{}, []databaseapi.TimeSeriesRow, error) {
	return getTimeSeriesCmp(from, to, aggrOp, "", "", false, tablename)
}

func queryProperties(query, aggr string) (interface{}, []databaseapi.TimeSeriesRow, error) {
	rows, err := databaseapi.ReadRowsTimescale(query, nil)

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

func existenceQuery(query string) (bool, error) {
	exists, err := databaseapi.ReadRowExistsTimescale(query)
	if err != nil {
		return false, err
	}

	return exists, nil
}
