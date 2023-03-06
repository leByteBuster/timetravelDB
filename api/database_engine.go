package api

import (
	"errors"
	"fmt"

	"github.com/LexaTRex/timetravelDB/parser"
	"github.com/LexaTRex/timetravelDB/utils"
)

// process a TTQL Query
func ProcessQuery(query string) (map[string][]any, error) {

	var queryResult map[string][]any

	queryInfo, err := parser.ParseQuery(query)
	if err != nil {
		return nil, err
	}

	// maybe doulbe check if from, to is valid ISO8601

	// is shallow
	if queryInfo.IsShallow {
		// no or no relevant lookups
		if !queryInfo.ContainsPropertyLookup || queryInfo.ContainsPropertyLookup && queryInfo.ContainsOnlyNullPredicate {
			utils.Debug("SHALLOW NO PROPERTY LOOKUPS")
			queryResult, err = getShallow(queryInfo, queryInfo.WhereClause)

			if err != nil {
				return nil, fmt.Errorf("error executing shallow query with no property lookups: %v", err)
			}
		} else {

			isValid, isWhere, isReturn := getPropertyLookupParentClause(queryInfo.PropertyClauseInsights)
			if !isValid {
				return nil, errors.New("invalid query, property lookup only allowed in WHERE or RETURN clause")
			} else if isWhere && isReturn {
				utils.Debug("SHALLOW WHERE AND RETURN PROPERTY LOOKUPS")

				queryResult, err = propertyLookupWhereReturnShallow(queryInfo)
				if err != nil {
					if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing shallow query with lookup in RETURN & WHERE: %v", err)); !ok {
						return nil, err
					}
				}

			} else if isWhere {
				utils.Debug("SHALLOW WHERE PROPERTY LOOUKPS")
				queryResult, err = propertyLookupWhereShallow(queryInfo)
				if err != nil {
					if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing shallow query with lookup in WHERE: %v", err)); !ok {
						return nil, err
					}
				}
			} else if isReturn {
				utils.Debug("SHALLOW RETURN PROPERTY LOOKUPS")
				queryResult, err = propertyLookupReturnShallow(queryInfo)
				if err != nil {
					if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing shallow query with lookup in RETURN: %v", err)); !ok {
						return nil, err
					}
				}
			} else {
				return nil, errors.New("this option should not be possible")
			}
		}
	} else {
		if !queryInfo.ContainsPropertyLookup || queryInfo.ContainsPropertyLookup && queryInfo.ContainsOnlyNullPredicate {

			utils.Debug("NON-SHALLOW NO PROPERTY LOOKUPS")

			queryResult, err = noPropertyLookup(queryInfo)
			if err != nil {
				if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing non-shallow query with no property lookups: %v", err)); !ok {
					return nil, err
				}
			}
		} else {

			isValid, isWhere, isReturn := getPropertyLookupParentClause(queryInfo.PropertyClauseInsights)
			if !isValid {
				return nil, errors.New("invalid query, property lookup only allowed in WHERE or RETURN clause")
			} else if isWhere && isReturn {
				utils.Debug("NON-SHALLOW WHERE AND RETURN PROPERTY LOOKUPS")
				queryResult, err = propertyLookupWhereReturn(queryInfo)
				if err != nil {
					if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing non-shallow query with lookup in RETURN & WHERE: %v", err)); !ok {
						return nil, err
					}
				}
			} else if isWhere {
				utils.Debug("NON-SHALLOW WHERE PROPERTY LOOKUPS")
				queryResult, err = propertyLookupWhere(queryInfo)
				if err != nil {
					if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing non-shallow query with lookup in RETURN & WHERE: %v", err)); !ok {
						return nil, err
					}
				}
			} else if isReturn {
				utils.Debug("NON-SHALLOW RETURN PROPERTY LOOKUPS")
				queryResult, err = propertyLookupReturn(queryInfo)
				if err != nil {
					if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing non-shallow query with lookup in RETURN & WHERE: %v", err)); !ok {
						return nil, err
					}
				}
			} else {
				return nil, errors.New("this option should not be possible")
			}
		}
	}

	utils.Debugf("\n\n\n                 		 QUERY RESULT\n						%+v\n\n\n", queryResult)
	if len(queryInfo.ReturnProjections) > 0 {
		utils.Debug("\n\n\n                      Printed ordered                         \n\n\n\n")
		utils.Debugf("%+v\n", utils.JsonStringFromMapOrdered(queryResult, queryInfo.ReturnProjections))
	} else {
		utils.Debug("\n\n\n                      Printed unordered                         \n\n\n\n")
		utils.Debugf("%+v\n", utils.JsonStringFromMap(queryResult))
	}

	// return errors.New("no option choosen, this should not occour")
	return queryResult, nil
}

// shallow functions:

func propertyLookupWhereReturnShallow(queryInfo parser.ParseResult) (map[string][]interface{}, error) {
	where, err := manipulateWhereClause(queryInfo, queryInfo.WhereClause)
	if err != nil {
		return nil, fmt.Errorf("%w; error manipulating WHERE query for neo4j", err)
	}
	graphData, err := getShallow(queryInfo, where)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving graph data", err)
	}
	relevantLookups := queryInfo.LookupsWhereRelevant
	graphData, err = filterForCondLookupsInWhere(queryInfo, graphData, relevantLookups)
	if err != nil {
		return nil, fmt.Errorf("%w; error filtering query result on WHERE conditions", err)
	}
	returnVariables := queryInfo.ReturnProjections
	lookupsReturn := queryInfo.LookupsReturn
	graphData, err = getPropertiesforLookupsInReturn(queryInfo, lookupsReturn, returnVariables, graphData)
	if err != nil {
		return nil, fmt.Errorf("%w; error filtering and merging queried data", err)
	}

	return graphData, nil
}

func propertyLookupWhereShallow(queryInfo parser.ParseResult) (map[string][]interface{}, error) {
	where, err := manipulateWhereClause(queryInfo, queryInfo.WhereClause)
	if err != nil {
		return nil, fmt.Errorf("%w; error manipulating WHERE query for neo4j", err)
	}
	graphData, err := getShallow(queryInfo, where)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving graph data", err)
	}
	relevantLookups := queryInfo.LookupsWhereRelevant
	res, err := filterForCondLookupsInWhere(queryInfo, graphData, relevantLookups)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving properties for lookups in WHERE", err)
	}
	if err != nil {
		return nil, fmt.Errorf("%w; error filtering and merging queried data", err)
	}
	return res, nil
}

func propertyLookupReturnShallow(queryInfo parser.ParseResult) (map[string][]interface{}, error) {
	graphData, err := getShallow(queryInfo, queryInfo.WhereClause)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving graph data", err)
	}
	mergedRes, err := getPropertiesforLookupsInReturn(queryInfo, queryInfo.LookupsReturn, queryInfo.ReturnProjections, graphData)
	if err != nil {
		return mergedRes, fmt.Errorf("%w; Not all properties could be fetched", err)
	}
	return mergedRes, nil
}

// non-shallow functions:

func noPropertyLookup(queryInfo parser.ParseResult) (map[string][]any, error) {
	queryResult, err := getShallow(queryInfo, queryInfo.WhereClause)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving graph data", err)
	}
	returnVariables := queryInfo.GraphElements.ReturnGraphElements

	if len(returnVariables) == 0 {
		returnVariables = queryInfo.GraphElements.MatchGraphElements
	}

	queryResult, err = getAllProperties(queryInfo, returnVariables, queryResult)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving time-series data", err)
	}

	return queryResult, nil
}

func propertyLookupWhereReturn(queryInfo parser.ParseResult) (map[string][]any, error) {
	where, err := manipulateWhereClause(queryInfo, queryInfo.WhereClause)
	if err != nil {
		return nil, fmt.Errorf("%w; error manipulating WHERE query for neo4j", err)
	}
	queryResult, err := getShallow(queryInfo, where)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving graph data", err)
	}
	relevantLookups := queryInfo.LookupsWhereRelevant
	queryResult, err = filterForCondLookupsInWhere(queryInfo, queryResult, relevantLookups)
	if err != nil {
		return nil, fmt.Errorf("%w; error filtering query result on WHERE conditions", err)
	}
	returnVariables := queryInfo.GraphElements.ReturnGraphElements
	if len(returnVariables) == 0 {
		returnVariables = queryInfo.GraphElements.MatchGraphElements
	}
	queryResult, err = getAllPropertiesLookupsReturn(queryInfo, queryInfo.LookupsReturn, returnVariables, queryResult)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving time-series data", err)
	}
	return queryResult, nil
}

func propertyLookupWhere(queryInfo parser.ParseResult) (map[string][]any, error) {
	return propertyLookupWhereReturn(queryInfo)
}

func propertyLookupReturn(queryInfo parser.ParseResult) (map[string][]any, error) {
	return propertyLookupWhereReturn(queryInfo)
}
