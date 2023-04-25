package qpengine

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/LexaTRex/timetravelDB/query-processor/parser"
	"github.com/LexaTRex/timetravelDB/query-processor/parser/listeners"
	tti "github.com/LexaTRex/timetravelDB/query-processor/parser/ttql_interface"
	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// isValidPropertyLookups checks if all property lookups occuring in valid clauses i.e. WHERE or RETURN clause.  We do not allow
// property lookups elsewhere. This should be part of parsing in the future.
func isValidPropertyLookups(insights map[*tti.OC_ComparisonExpressionContext][]listeners.PropertyClauseInsight) (isValid bool) {
	isValid = true
	for _, listOfInsights := range insights {
		for _, insight := range listOfInsights {
			if !insight.IsValid {
				isValid = false
			}
		}
	}
	return isValid
}

// getSelectedTimeSeries fetches all time series of specific property lookups which are time series properties
func getSelectedTimeSeries(queryInfo parser.ParseResult, lookupsMap map[string][]string, returnProjections []string, graphData map[string][]interface{}) (map[string][]interface{}, error) {
	var err error
	for elVar, lookups := range lookupsMap {
		elements := graphData[elVar]
		mergeVariables := false
		if utils.Contains(returnProjections, elVar) {
			mergeVariables = true
		}
		for _, prop := range lookups {

			// if the explicit named property is no time series property - adopt the value from the shallow graph
			if !strings.HasPrefix(prop, "ts_") && !strings.HasPrefix(prop, "properties_") {

				for _, el := range elements {
					switch e := el.(type) {
					case neo4j.Entity:
						propVal := e.GetProperties()[prop]
						lookup := getLookupString(elVar, prop)
						graphData[lookup] = append(graphData[prop], propVal)
					}
				}

				// else: fetch the according time series from TimescaleDB for all elements that fall under the according element variable
			} else {
				graphData, err = fetchTimeSeries(queryInfo.From, queryInfo.To, graphData, elements, prop, elVar, mergeVariables)
			}
		}
	}
	return graphData, err
}

func getAllTimeseries(queryInfo parser.ParseResult, lookupsMap map[string][]string, returnVariables []string, graphData map[string][]interface{}) (map[string][]interface{}, error) {
	var err error
	plainReturnVariables := []string{}
	projections := queryInfo.ReturnProjections

	// if RETURN *
	if len(projections) == 0 {
		plainReturnVariables = returnVariables
	} else {
		for _, projection := range projections {
			if !strings.Contains(projection, ".") {
				plainReturnVariables = append(plainReturnVariables, projection)
			}
		}
	}

	for _, n := range plainReturnVariables {
		elements := graphData[n]
		graphData, err = fetchTimeSeriesAll(queryInfo.From, queryInfo.To, graphData, elements, n)
	}

	for elVar, lookups := range lookupsMap {
		alreadyFetched := utils.Contains(plainReturnVariables, elVar)
		for _, property := range lookups {
			utils.Debugf("                        GRAPH DATA BEFORE SINGLE TIME SERIES FETCH: \n         %+v\n", graphData)
			if !strings.HasPrefix(property, "ts_") && !strings.HasPrefix(property, "properties_") {
				continue
			}
			graphData, err = getTimeSeriesSingleLookup(queryInfo, graphData, elVar, property, alreadyFetched)
			utils.Debugf("                        GRAPH DATA AFTER SINGLE TIME SERIES FETCH : \n         %+v\n", graphData)
		}
	}
	return graphData, err
}

func getTimeSeriesSingleLookup(queryInfo parser.ParseResult, graphData map[string][]interface{}, elementVar, property string, alreadyFetched bool) (map[string][]interface{}, error) {
	lookup := getLookupString(elementVar, property)
	elements := graphData[elementVar]
	graphData[lookup] = make([]any, len(elements))
	if alreadyFetched {
		for i, el := range elements {

			if e, ok := el.(neo4j.Entity); ok {
				properties := e.GetProperties()
				graphData[lookup][i] = properties[property]
				utils.Debugf("\n\n	    merged time-series: %v      \n\n", e.GetProperties()[property])
			} else {
				return nil, fmt.Errorf("error fetching time-series for lookup %v: element is not neo4j.Entity", lookup)
			}
		}
	} else {
		graphData, err := fetchSinglePropTimeSeries(queryInfo, graphData, elementVar, property)
		if err != nil {
			return graphData, fmt.Errorf("error fetching time-series for lookup %v: %v", lookup, err)
		}
	}
	return graphData, nil
}

// fetchSinglePropTimeSeries fetches a single time series of a graph element property
func fetchSinglePropTimeSeries(queryInfo parser.ParseResult, graphData map[string][]interface{}, elementVar, property string) (map[string][]any, error) {

	lookup := getLookupString(elementVar, property)
	elements := graphData[elementVar]
	graphData[lookup] = make([]any, len(elements))

	for i, el := range elements {
		if e, ok := el.(neo4j.Entity); ok {

			properties := e.GetProperties()
			uuid := properties[property]

			if uuid == nil {
				// property not available - do nothing
				utils.Debugf("\nproperty %v not available on element with id : %v\n", property, e.GetElementId())
			} else if s, ok := uuid.(string); ok {

				// fetch time series by uuid
				tablename := uuidToTablename(s)
				_, timeseries, err := getTimeSeries(queryInfo.From, queryInfo.To, "", tablename)

				graphData[lookup][i] = timeseries

				if err != nil {
					return graphData, fmt.Errorf("%w; error - couldnt fetch  properties for %v of element", err, property)
				}
			} else {
				return graphData, errors.New("error - uuid is not a string - this should not happen")
			}
		} else {
			return graphData, fmt.Errorf("unknown type of object %v", el)
		}
	}
	return graphData, nil
}

// fetchTimeSeriesAll fetches all time series of all properties of all passed elements
func fetchTimeSeriesAll(from string, to string, graphData map[string][]interface{}, elements []interface{}, elementVar string) (map[string][]interface{}, error) {
	for _, el := range elements {
		switch e := el.(type) {
		case neo4j.Entity:
			properties := e.GetProperties()
			for prop, uuid := range properties {
				if !strings.HasPrefix(prop, "ts_") && !strings.HasPrefix(prop, "properties_") {
					continue
				}
				if uuid == nil {
					utils.Debugf("\nproperty %v not available on element with id : %v\n", prop, e.GetElementId())
				} else if s, ok := uuid.(string); ok {

					tablename := uuidToTablename(s)

					_, timeseries, err := getTimeSeries(from, to, "", tablename)

					if err != nil {
						return graphData, fmt.Errorf("%w; error - couldnt fetch  properties for %v of element", err, prop)
					}

					properties[prop] = timeseries
				} else {
					return graphData, errors.New("error - uuid is not a string - this should not happen")
				}
			}
		default:
			panic("error - type not supportet")
		}
	}
	return graphData, nil
}

// getShallow returns the shallow graph without time series content
func getShallow(queryInfo parser.ParseResult) (map[string][]interface{}, error) {
	cndWhere, err := buildCondWhereClause(queryInfo.LookupsWhereRelevant, queryInfo.WhereClause)
	if err != nil {
		return nil, fmt.Errorf("%w; error manipulating WHERE query for neo4j", err)
	}
	tmpWhere := buildTmpWhereClause(queryInfo.From, queryInfo.To, cndWhere, queryInfo.QueryVariables.MatchQueryVariables)
	returnClause := buildReturnClause(queryInfo.LookupsWhereRelevant, queryInfo.QueryVariables.ReturnQueryVariables)
	query := buildFinalQuery(queryInfo.MatchClause, tmpWhere, returnClause)
	res, err := queryNeo4j(query)
	if err != nil {
		return nil, err
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	resMap, err := resultToMap(res)

	utils.Debugf("NEO4J RESULT:\n      %+v", resMap)

	return resMap, err
}

// filterForCondLookups filters the result set of graph elements after a conditional lookup in a time series. Right now the only possible
// filter is ANY() which is treated equivalent to not using an operator for lookup of a time series. For example if a.prop > 20 in the
// WHERE clause it will be checked if there is a value > 20 in the time series of a.prop (in the boundaries of the specified time frame
// of the query). If there is an element > 20 the element will be in the result set which is available for RETURN. If not it is removed.
// Right now saying a.prop > 20 is equivalent to ANY(a.prop) > 20. This is supposed to be changed when more operators are implemented.
func filterForCondLookups(from string, to string, relevantLookups []parser.LookupInfo, graphData map[string][]interface{}) (map[string][]interface{}, error) {
	var err error
	var toRemove []int
	filteredData := graphData
	for _, lookupInfo := range relevantLookups {
		elements := graphData[lookupInfo.ElementVariable]

		toRemove, err = checkAnyCondition(from, to, graphData, elements, lookupInfo.Property, lookupInfo.ElementVariable, lookupInfo.CompareOperator, lookupInfo.CompareValue, lookupInfo.LookupLeft)

		// TO IMPLEMENT:
		// check if AVG() condition is fullfilled.
		// toRemove, err = checkAVGCondition()

		// filter result set
		filteredData = filterMatches(filteredData, toRemove, []string{})

		if err != nil {
			return nil, fmt.Errorf("%w; error - couldnt merge time series in property", err)
		}
	}
	return filteredData, nil
}

// fetchTimeSeries fetches the time series of a specific property for all matching elements
// for example with "MATCH a RETURN a.prop" all time series of prop are fetched for all elements a
func fetchTimeSeries(from string, to string, graphData map[string][]interface{}, elements []interface{}, property string, elementVar string, mergeVariables bool) (map[string][]interface{}, error) {
	for i, el := range elements {
		switch e := el.(type) {
		case neo4j.Entity:
			uuid := e.GetProperties()[property]
			if uuid == nil {
			} else if s, ok := uuid.(string); ok {
				tablename := uuidToTablename(s)
				propertyMapOfElement := graphData[elementVar][i].(neo4j.Entity).GetProperties()
				_, properties, err := getTimeSeries(from, to, "", tablename)
				if err != nil {
					return nil, fmt.Errorf("%w; error - couldnt fetch  properties for %v of element", err, property)
				}
				lookup := getLookupString(elementVar, property)
				graphData[lookup] = append(graphData[lookup], properties)
				if mergeVariables {
					propertyMapOfElement[property] = properties
				}
			} else {
				return nil, errors.New("error - uuid is not a string - this should not happen")
			}
		default:
			panic("error - type not supportet")
		}
	}
	return graphData, nil
}

// checkAnyCondition conditionally checks time series of a property for the ANY() operator
// ANY() is treated equivalent to not using an operator for lookup of a time series. For example if a.prop > 20 in the
// WHERE clause it will be checked if there is a value > 20 in the time series of a.prop (in the boundaries of the specified time frame
// of the query). If there is an element > 20 the element will be in the result set which is available for RETURN. If not it is removed.
// Right now saying a.prop > 20 is equivalent to ANY(a.prop) > 20. This is supposed to be changed when more operators are implemented.
func checkAnyCondition(from string, to string, graphData map[string][]interface{}, elements []interface{}, property string, elementVar string, compareOp string, compareVal any, lookupLeft bool) ([]int, error) {
	rowsToRemove := []int{}

	for i, el := range elements {
		switch e := el.(type) {
		case neo4j.Entity:
			uuid := e.GetProperties()[property]
			if uuid == nil {
			} else if s, ok := uuid.(string); ok {
				tablename := uuidToTablename(s)
				exists, err := checkCondInTimeseries(from, to, "", compareOp, compareVal, lookupLeft, tablename)
				if err != nil {
					return nil, fmt.Errorf("%w; error - check if value with condidtion exists for time-series %v of element %v", err, property, e.GetElementId())
				} else if exists {
					if err != nil {
						return nil, fmt.Errorf("%w; error - couldnt fetch  properties for %v of element", err, property)
					}
				} else {
					rowsToRemove = append(rowsToRemove, i)
				}
			} else {
				return nil, errors.New("error - uuid is not a string - this should not happen")
			}
		default:
			panic("error - type not supportet")
		}
	}
	return rowsToRemove, nil
}

func filterMatches(graphData map[string][]interface{}, rowsToRemove []int, exceptions []string) map[string][]interface{} {
	for elVar, elements := range graphData {
		for i := len(rowsToRemove) - 1; i >= 0; i-- {
			elements = utils.RemoveIdxFromSlice(elements, rowsToRemove[i])
			graphData[elVar] = elements
		}
	}
	return graphData
}

// shouldnt I be able to get a list of these in the listener??

func handleErrorOnResult(res map[string][]any, err error) (bool, error) {
	if err != nil && res == nil {
		return false, err
	} else if err != nil {
		log.Fatalf("Not all elements contained the property: %v", err)
	}
	return true, nil
}

func getLookupString(elementVar, property string) string {
	var sb strings.Builder
	sb.WriteString(elementVar)
	sb.WriteString(".")
	sb.WriteString(property)
	return sb.String()
}
