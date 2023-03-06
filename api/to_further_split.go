package api

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/LexaTRex/timetravelDB/parser"
	"github.com/LexaTRex/timetravelDB/parser/listeners"
	tti "github.com/LexaTRex/timetravelDB/parser/ttql_interface"
	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func getPropertyLookupParentClause(insights map[*tti.OC_ComparisonExpressionContext][]listeners.PropertyClauseInsight) (isValid bool, isWhere bool, isReturn bool) {
	isValid = true
	for _, listOfInsights := range insights {
		for _, insight := range listOfInsights {
			if !insight.IsValid {
				isValid = false
			}
			if insight.IsWhere && insight.IsPropertyLookup && !insight.IsAppendixOfNullPredicate {
				isWhere = true
			}
			if insight.IsReturn && insight.IsPropertyLookup && !insight.IsAppendixOfNullPredicate {
				isReturn = true
			}
		}
	}
	return isValid, isWhere, isReturn
}

func getPropertiesforLookupsInReturn(queryInfo parser.ParseResult, lookupsMap map[string][]string, returnProjections []string, graphData map[string][]interface{}) (map[string][]interface{}, error) {
	var err error
	for elVar, lookups := range lookupsMap {
		elements := graphData[elVar]
		mergeVariables := false
		if utils.Contains(returnProjections, elVar) {
			mergeVariables = true
		}
		for _, lookup := range lookups {
			graphData, err = fetchTimeSeries(queryInfo.From, queryInfo.To, graphData, elements, lookup, elVar, mergeVariables)
		}
	}
	return graphData, err
}

func getAllProperties(queryInfo parser.ParseResult, returnVariables []string, graphData map[string][]interface{}) (map[string][]interface{}, error) {
	var err error
	for _, n := range returnVariables {
		elements := graphData[n]
		graphData, err = fetchTimeSeriesAll(queryInfo.From, queryInfo.To, graphData, elements, n)
	}
	return graphData, err
}

func getAllPropertiesLookupsReturn(queryInfo parser.ParseResult, lookupsMap map[string][]string, returnVariables []string, graphData map[string][]interface{}) (map[string][]interface{}, error) {
	var err error
	plainReturnVariables := []string{}
	projections := queryInfo.ReturnProjections
	for _, projection := range projections {
		if !strings.Contains(projection, ".") {
			plainReturnVariables = append(plainReturnVariables, projection)
		}
	}

	for _, n := range plainReturnVariables {
		elements := graphData[n]
		graphData, err = fetchTimeSeriesAll(queryInfo.From, queryInfo.To, graphData, elements, n)
	}

	for elVar, lookups := range lookupsMap {
		alreadyFetched := utils.Contains(plainReturnVariables, elVar)
		for _, lookup := range lookups {
			utils.Debugf("                        GRAPH DATA BEFORE SINGLE TIME SERIES FETCH: \n         %+v\n", graphData)
			graphData, err = getTimeSeriesSingleLookup(queryInfo, graphData, elVar, lookup, alreadyFetched)
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
				fmt.Printf("\n\n	    merged time-series: %v      \n\n", e.GetProperties()[property])
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

func fetchSinglePropTimeSeries(queryInfo parser.ParseResult, graphData map[string][]interface{}, elementVar, property string) (map[string][]any, error) {

	var sb strings.Builder
	sb.WriteString(elementVar)
	sb.WriteString(".")
	sb.WriteString(property)
	lookup := sb.String()
	elements := graphData[elementVar]
	graphData[lookup] = make([]any, len(elements))

	for i, el := range elements {
		if e, ok := el.(neo4j.Entity); ok {
			properties := e.GetProperties()
			uuid := properties[property]
			if uuid == nil {
				// property not available - do nothing
				log.Printf("\nproperty %v not available on element with id : %v\n", property, e.GetElementId())
			} else if s, ok := uuid.(string); ok {
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

func fetchTimeSeriesAll(from string, to string, graphData map[string][]interface{}, elements []interface{}, elementVar string) (map[string][]interface{}, error) {
	for _, el := range elements {
		switch e := el.(type) {
		case neo4j.Entity:
			properties := e.GetProperties()
			for prop, uuid := range properties {
				if !strings.HasPrefix(prop, "properties") {
					continue
				}
				if uuid == nil {
					log.Printf("\nproperty %v not available on element with id : %v\n", prop, e.GetElementId())
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

func getShallow(queryInfo parser.ParseResult, whereManipulated string) (map[string][]interface{}, error) {

	tmpWhere := buildTmpWhereClause(queryInfo.From, queryInfo.To, whereManipulated, queryInfo.GraphElements.MatchGraphElements)
	returnClause := buildReturnClause(queryInfo.LookupsWhereRelevant, queryInfo.GraphElements.ReturnGraphElements)
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

func filterForCondLookupsInWhere(queryInfo parser.ParseResult, graphData map[string][]interface{}, relevantLookups []parser.LookupInfo) (map[string][]interface{}, error) {
	var err error
	var toRemove []int
	filteredData := graphData
	for _, lookupInfo := range relevantLookups {
		elements := graphData[lookupInfo.ElementVariable]
		toRemove, err = checkIfValueForConditionExists(queryInfo.From, queryInfo.To, graphData, elements, lookupInfo.Property, lookupInfo.ElementVariable, lookupInfo.CompareOperator, lookupInfo.CompareValue, lookupInfo.LookupLeft)
		filteredData = filterMatches(filteredData, toRemove, []string{})
		if err != nil {
			return nil, fmt.Errorf("%w; error - couldnt merge time series in property", err)
		}
	}
	return filteredData, nil
}

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

func checkIfValueForConditionExists(from string, to string, graphData map[string][]interface{}, elements []interface{}, property string, elementVar string, compareOp string, compareVal any, lookupLeft bool) ([]int, error) {
	rowsToRemove := []int{}

	for i, el := range elements {
		switch e := el.(type) {
		case neo4j.Entity:
			uuid := e.GetProperties()[property]
			if uuid == nil {
			} else if s, ok := uuid.(string); ok {
				tablename := uuidToTablename(s)
				exists, err := checkIfValueWithConditionExists(from, to, "", compareOp, compareVal, lookupLeft, tablename)
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
		log.Printf("Not all elements contained the property: %v", err)
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
