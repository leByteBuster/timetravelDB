package qpengine

import (
	"errors"
	"fmt"

	"github.com/LexaTRex/timetravelDB/query-processor/parser"
	"github.com/LexaTRex/timetravelDB/utils"
)

// ProcessQuery uses the collected parsing information in the ParseResult object to obtain the requested data by
// performing binary querying (neo4j & timescaledb) and merging the result if necessary
func ProcessQuery(queryInfo parser.ParseResult) (map[string][]any, error) {

	var queryResult map[string][]any

	isValid := isValidPropertyLookups(queryInfo.PropertyClauseInsights)
	if !isValid {
		return nil, errors.New("invalid query, property lookup only allowed in WHERE or RETURN clause")
	}

	shallowResult, err := getShallow(queryInfo)

	if err != nil {
		return nil, fmt.Errorf("error fetching shallow graph data: %v", err)
	}

	if queryInfo.IsShallow {
		queryResult, err = applyBinaryQueryShallow(queryInfo, shallowResult)
		if err != nil {
			return nil, fmt.Errorf("%w; error applying shallow query", err)
		}
	} else {
		queryResult, err = applyBinaryQueryDeep(queryInfo, shallowResult)
		if err != nil {
			return nil, fmt.Errorf("%w; error applying deep query", err)
		}
	}

	return queryResult, nil
}

func applyBinaryQueryShallow(queryInfo parser.ParseResult, shallowResult map[string][]interface{}) (map[string][]interface{}, error) {
	var queryResult map[string][]interface{}
	// no or no relevant lookups
	// NOTE: in the case of prop.a IS NOT NULL > 20 this doesnt work right yet
	// even though this does not make much sense in practice in the first place we should handle it correctly
	if !queryInfo.ContainsPropertyLookup || queryInfo.ContainsPropertyLookup && queryInfo.ContainsOnlyNullPredicate {
		utils.Debug("SHALLOW QUERY CONTAINING LOOKUPS")

		// only case where no binary querying is not neccessary
		return shallowResult, nil
	} else {

		filteredResult, err := filterForCondLookups(queryInfo.From, queryInfo.To, queryInfo.LookupsWhereRelevant, shallowResult)
		if err != nil {
			return nil, fmt.Errorf("%w; error filtering query result on WHERE conditions", err)
		}
		queryResult, err = propertyLookupShallow(queryInfo, filteredResult)
		if err != nil {
			if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing shallow query with lookup: %v", err)); !ok {
				return nil, err
			}
		}
	}
	return queryResult, nil
}

func applyBinaryQueryDeep(queryInfo parser.ParseResult, shallowResult map[string][]interface{}) (map[string][]interface{}, error) {

	var queryResult map[string][]interface{}

	utils.Debug("DEEP QUERY CONTAINING LOOKUPS")

	filteredResult, err := filterForCondLookups(queryInfo.From, queryInfo.To, queryInfo.LookupsWhereRelevant, shallowResult)
	if err != nil {
		return nil, fmt.Errorf("%w; error filtering query result on WHERE conditions", err)
	}
	queryResult, err = propertyLookupDeep(queryInfo, filteredResult)

	if err != nil {
		if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing deep query with lookup: %v", err)); !ok {
			return nil, err
		}
	}
	return queryResult, nil
}

func propertyLookupShallow(queryInfo parser.ParseResult, queryResult map[string][]interface{}) (map[string][]interface{}, error) {

	returnVariables := queryInfo.ReturnProjections
	lookupsReturn := queryInfo.LookupsReturn

	queryResult, err := getSelectedTimeSeries(queryInfo, lookupsReturn, returnVariables, queryResult)
	if err != nil {
		return nil, fmt.Errorf("%w; error filtering and merging queried data", err)
	}

	return queryResult, nil
}

func propertyLookupDeep(queryInfo parser.ParseResult, queryResult map[string][]interface{}) (map[string][]any, error) {

	returnVariables := queryInfo.QueryVariables.ReturnQueryVariables
	if len(returnVariables) == 0 {
		returnVariables = queryInfo.QueryVariables.MatchQueryVariables
	}
	queryResult, err := getAllTimeseries(queryInfo, queryInfo.LookupsReturn, returnVariables, queryResult)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving time-series data", err)
	}
	return queryResult, nil
}

// func noPropertyLookupDeep(queryInfo parser.ParseResult, shallowResult map[string][]any) (map[string][]any, error) {
//
// 	returnVariables := queryInfo.QueryVariables.ReturnQueryVariables
//
// 	if len(returnVariables) == 0 {
// 		returnVariables = queryInfo.QueryVariables.MatchQueryVariables
// 	}
//
// 	queryResult, err := getAllTimeSeries(queryInfo, returnVariables, shallowResult)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w; error retrieving time-series data", err)
// 	}
//
// 	return queryResult, nil
// }
