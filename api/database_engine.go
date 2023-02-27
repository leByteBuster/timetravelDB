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
			log.Println("checkpoint1")
			queryResult, err = getShallow(queryInfo, queryInfo.WhereClause)

			if err != nil {
				return nil, fmt.Errorf("error executing shallow query with no property lookups: %v", err)
			}
		} else {

			isValid, isWhere, isReturn := getPropertyLookupParentClause(queryInfo.PropertyClauseInsights)
			if !isValid {
				return nil, errors.New("invalid query, property lookup only allowed in WHERE or RETURN clause")
			} else if isWhere && isReturn {
				log.Println("checkpoint2")

				queryResult, err = propertyLookupWhereReturnShallow(queryInfo)
				if err != nil {
					if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing shallow query with lookup in RETURN & WHERE: %v", err)); !ok {
						return nil, err
					}
				}

			} else if isWhere {
				log.Println("checkpoint3")
				queryResult, err = propertyLookupWhereShallow(queryInfo)
				if err != nil {
					if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing shallow query with lookup in WHERE: %v", err)); !ok {
						return nil, err
					}
				}
			} else if isReturn {
				log.Println("checkpoint4")
				queryResult, err = propertyLookupReturnShallow(queryInfo)
				if err != nil {
					if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing shallow query with lookup in RETURN: %v", err)); !ok {
						return nil, err
					}
				}
			} else {
				fmt.Printf("\nReturn: %v, Where: %v, Valid: %v\n", isReturn, isWhere, isValid)
				return nil, errors.New("this option should not be possible")
			}
		}
	} else {
		if !queryInfo.ContainsPropertyLookup || queryInfo.ContainsPropertyLookup && queryInfo.ContainsOnlyNullPredicate {

			log.Println("checkpoint5")

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
				log.Println("checkpoint6")
				queryResult, err = propertyLookupWhereReturn(queryInfo)
				if err != nil {
					if ok, err := handleErrorOnResult(queryResult, fmt.Errorf("error executing non-shallow query with lookup in RETURN & WHERE: %v", err)); !ok {
						return nil, err
					}
				}
			} else if isWhere {
				log.Println("checkpoint7")
				propertyLookupWhere(queryInfo)
			} else if isReturn {
				log.Println("checkpoint8")
				propertyLookupReturn(queryInfo)
			} else {
				return nil, errors.New("this option should not be possible")
			}
		}
	}

	fmt.Printf("\n\n\n                      QUERY RESULT                         \n%+v\n\n\n", queryResult)
	if len(queryInfo.ReturnProjections) > 0 {
		fmt.Printf("\n\n\n                      Printed ordered                         \n\n\n\n")
		fmt.Println(utils.JsonStringFromMapOrdered(queryResult, queryInfo.ReturnProjections))
	} else {
		fmt.Printf("\n\n\n                      Printed unordered                         \n\n\n\n")
		fmt.Println(utils.JsonStringFromMap(queryResult))
	}

	// return errors.New("no option choosen, this should not occour")
	return queryResult, nil
}

func noPropertyLookup(queryInfo parser.ParseResult) (map[string][]any, error) {
	queryResult, err := getShallow(queryInfo, queryInfo.WhereClause)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving graph data", err)
	}

	fmt.Printf("\n\n  Query Result:\n  %+v", queryResult)

	returnVariables := queryInfo.GraphElements.ReturnGraphElements

	// ############################
	// ############### TODO: delete if this error does not occur
	// ############################
	// if !reflect.DeepEqual(returnVariables, queryInfo.GraphElements.ReturnGraphElementsNoLookup) {
	// 	fmt.Println(returnVariables)
	// 	fmt.Println(queryInfo.GraphElements.ReturnGraphElementsNoLookup)
	// 	return nil, errors.New("in this case the return variables should be the same as the ones without lookups")
	// }

	// this can be changed to "if len(queryInfo.ReturnProjections) == 0" in case we merge this funciton with propertyLookupWhereReturn
	// here we can use GraphElements.ReturnGraphElements because in this case no lookups are occuring
	if len(returnVariables) == 0 {
		returnVariables = queryInfo.GraphElements.MatchGraphElements
	}

	queryResult, err = getAllProperties(queryInfo, returnVariables, queryResult)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving time-series data", err)
	}

	return queryResult, nil
}

// checks if any of the property clauses of the query is invalid
// checks if property lookups occour only in WHERE, only in RETURN or in both. Only if
// the lookup is not an appendix of a null predicate it has to be taken care of
// note: just becaues it is a propertyClausInsight it does not need to be a property lookup
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

// functions for non shallow queries with lookup that requires double database access

func propertyLookupWhere(res parser.ParseResult) {
	panic("unimplemented")
}

func propertyLookupReturn(res parser.ParseResult) {
	panic("unimplemented")
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

	fmt.Printf("\n\n  Query Result:\n  %+v", queryResult)

	// relevantLookups is an array of relevant lookups in the where clause.
	relevantLookups := queryInfo.LookupsWhereRelevant

	// query the data from timescaleDB according to the property lookups in the original WHERE clause
	queryResult, err = filterForCondLookupsInWhere(queryInfo, queryResult, relevantLookups)

	if err != nil {
		return nil, fmt.Errorf("%w; error filtering query result on WHERE conditions", err)
	}

	fmt.Printf("\n\n  Query Result after filtering:\n  %+v", queryResult)

	returnVariables := queryInfo.GraphElements.ReturnGraphElements

	// if return clause is empty then return all variables from match clause
	if len(returnVariables) == 0 {
		returnVariables = queryInfo.GraphElements.MatchGraphElements
	}

	queryResult, err = getAllPropertiesLookupsReturn(queryInfo, queryInfo.LookupsReturn, returnVariables, queryResult)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving time-series data", err)
	}

	return queryResult, nil
}

// functions for shallow queries with lookup that requires double database access

func propertyLookupWhereShallow(queryInfo parser.ParseResult) (map[string][]interface{}, error) {

	where, err := manipulateWhereClause(queryInfo, queryInfo.WhereClause)
	if err != nil {
		return nil, fmt.Errorf("%w; error manipulating WHERE query for neo4j", err)
	}
	graphData, err := getShallow(queryInfo, where)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving graph data", err)
	}
	// relevantLookups is an array of relevant lookups in the where clause. for every lookup we have a LookupInfo struct
	// which contains all necessarythe information about the lookup
	// [n.name, n.age, n.address, e.name, e.age, e.address]
	// up to now relevant lookups in the WHERE clause considers only considers lookups which are part of a comparison like (n.name = "Max")
	relevantLookups := queryInfo.LookupsWhereRelevant

	// query the data from timescaleDB according to the property lookups in the original WHERE clause
	res, err := filterForCondLookupsInWhere(queryInfo, graphData, relevantLookups)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving properties for lookups in WHERE", err)
	}

	// filter out all elements from the neo4j answer which do not match the original WHERE clause with

	if err != nil {
		return nil, fmt.Errorf("%w; error filtering and merging queried data", err)
	}

	// TODO NEXT: process res
	// iterate over the keys of the records and apply the original RETURN Clause
	// with  record.Get("movieTitle") it should be possible to get the "columns"
	// PROBLEM: if movieTitle is an alias - how to get the real name which was used in MATCH so we can merge
	// the results of both queries ? But here it doesnt matter because we Return *
	// but for the saved RETRUN clause (the original) it matters. Think about it.
	return res, nil
}

// only RETURN clause contains property lookups which needs double database querying
// so send everything to neo4j but with RETURN * instead the original RETURN clause
// and then take care of the original RETURN clause
func propertyLookupReturnShallow(queryInfo parser.ParseResult) (map[string][]interface{}, error) {

	// graph data is a map where the RETURN variables of the CYPHER query are mapped
	// to their results ala: {n: [{id:node1, properties}, {id:node2, properties}], s: ..., e: ...} for "...RETURN n, s, e"
	// NOTE: graphData can contain the same relations multiple times if it is returned multiple times if the pattern matches multiple times
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

// lookupsMap: map of element variables (n, s, e) to the lookups which are happening on them in RETURN clause (n: [name, age, address] for n.name, n.age, n.address)
// returnVariables: array of all element variables which are returned plain (no lookup) in the RETURN clause: RETURN  n, s, e, e.name -> [n, s, e]
func getPropertiesforLookupsInReturn(queryInfo parser.ParseResult, lookupsMap map[string][]string, returnProjections []string, graphData map[string][]interface{}) (map[string][]interface{}, error) {

	var err error

	// elVar represents the element variable of the RETURN clause the lookups are happening on (i.e. n, s, e)
	// lookups represents all the lookups which are happening on the element variable elVar (i.e. n.name, n.age, n.address)
	for elVar, lookups := range lookupsMap {
		elements := graphData[elVar]

		// if the element variable the lookup is happening on also occours plain then we merge the properties into the element

		mergeVariables := false

		// note: i cannot just check if graphData[elVar] exists because if there is a property lookup the element is fetched form neo4j
		// 	 	   even if it is not returned plain
		// solution: check if the variable of the lookup is occuring in returnProjections
		// 	         returnProjections contains all RETURN projections: RETURN n, r, s.prop -> [n, r, s, s.prop]
		if utils.Contains(returnProjections, elVar) {
			mergeVariables = true
		}

		for _, lookup := range lookups {
			graphData, err = fetchTimeSeries(queryInfo.From, queryInfo.To, graphData, elements, lookup, elVar, mergeVariables)
			// graphData = filterMatches(graphData, rowsToRemove, []string{})
		}
	}
	return graphData, err
}

func getAllProperties(queryInfo parser.ParseResult, returnVariables []string, graphData map[string][]interface{}) (map[string][]interface{}, error) {

	var err error
	// lookups represents all the lookups which are happening on the element variable elVar (i.e. n.name, n.age, n.address)
	for _, n := range returnVariables {
		elements := graphData[n]

		graphData, err = fetchTimeSeriesAll(queryInfo.From, queryInfo.To, graphData, elements, n)
		// graphData = filterMatches(graphData, rowsToRemove, []string{})
	}
	return graphData, err
}

// lookupsMap: map of element variables (n, s, e) to the lookups which are happening on them in RETURN clause (n: [name, age, address] for n.name, n.age, n.address)
// returnProjections: all the projections in the RETURN clause (i.e. "RETURN n, s, e, e.name" -> [n, s, e, e.name])
// note: we can retrieve plain elements through filtering returnProjections with lookupsMap:
//
//			   to retrieve [n,s,e] 	from [n,s,e,e.name]:
//	       if  _, ok := lookupsMap["foo"]; !ok { remove element from list }
//			THIS ONLY WORKS IF I ADD THE ELEMENTS TO THE LOOKUPS MAP IF THEY ARE RETURNED PLAIN WITH EMPTY LIST
//			- which it shold now
//
// returnVariables: contains all variables which occur in WHERE, does not say anything about if in a lookup context or plain
func getAllPropertiesLookupsReturn(queryInfo parser.ParseResult, lookupsMap map[string][]string, returnVariables []string, graphData map[string][]interface{}) (map[string][]interface{}, error) {

	var err error

	fmt.Printf("\n    return variables: %+v", returnVariables)

	// TODO NEXT: find a way to check which lookups are happening in the RETURN clause and fetch those properties
	// advanded: if it's parent variable occours in WHERE: get the time-series from the parent node/edge to avoid
	// 					 fetching twice

	// possibility 1:
	//   - fetch all the elements first: a,e,n
	//   - then for every lookup: check if parent has been fetched. retrieve the result from result map

	plainReturnVariables := []string{}
	projections := queryInfo.ReturnProjections
	for _, projection := range projections {
		if !strings.Contains(projection, ".") {
			plainReturnVariables = append(plainReturnVariables, projection)
		}
	}

	fmt.Printf("\n\n    lookups map: %+v      \n\n", lookupsMap)
	fmt.Printf("\n\n    plain return variables: %+v      \n\n", lookupsMap)
	for _, el := range returnVariables {
		fmt.Printf("\n\n    return variable: %+v      \n\n", el)

		// if the returnVariable is not a key of the lookupMap then it is a plain return variable
		if lookups, ok := lookupsMap[el]; !ok {
			if len(lookups) == 0 {
				plainReturnVariables = append(plainReturnVariables, el)
			}
		}
	}

	fmt.Printf("\n\n    plain return variables: %+v      \n\n", plainReturnVariables)

	// for every plain return variable: for every graph element related to it merge all time-series into its properties
	for _, n := range plainReturnVariables {
		elements := graphData[n]

		graphData, err = fetchTimeSeriesAll(queryInfo.From, queryInfo.To, graphData, elements, n)
		// graphData = filterMatches(graphData, rowsToRemove, []string{})
	}

	for elVar, lookups := range lookupsMap {

		alreadyFetched := utils.Contains(plainReturnVariables, elVar)

		fmt.Printf("\n\n	    already fetched?: %v      \n\n", alreadyFetched)

		for _, lookup := range lookups {
			graphData, err = getTimeSeriesSingleLookup(queryInfo, graphData, elVar, lookup, alreadyFetched)
		}
	}

	return graphData, err
}

func getTimeSeriesSingleLookup(queryInfo parser.ParseResult, graphData map[string][]interface{}, elementVar, property string, alreadyFetched bool) (map[string][]interface{}, error) {

	var sb strings.Builder
	sb.WriteString(elementVar)
	sb.WriteString(".")
	sb.WriteString(property)
	lookup := sb.String()
	elements := graphData[elementVar]
	graphData[lookup] = make([]any, len(elements))

	if alreadyFetched {
		// get time-series from parent
		for i, el := range elements {

			fmt.Printf("\n\n	    already fetched:       \n\n")

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

				fmt.Printf("\n\n	    fetched time-series: %v      \n\n", timeseries)

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
		fmt.Printf("\n\n	    el: %+v      \n\n", el)
		switch e := el.(type) {
		case neo4j.Entity:
			properties := e.GetProperties()
			for prop, uuid := range properties {
				// skip if not a timeseries property
				fmt.Printf("\n\n	    prop: %+v      \n\n", prop)
				fmt.Printf("\n\n	    uuid: %+v      \n\n", uuid)
				fmt.Printf("\n\n	    e: %+v      \n\n", e)
				if !strings.HasPrefix(prop, "properties") {
					continue
				}
				if uuid == nil {
					// property not available - do nothing
					log.Printf("\nproperty %v not available on element with id : %v\n", prop, e.GetElementId())
				} else if s, ok := uuid.(string); ok {

					fmt.Printf("\n\n	    ERREISCHT:      \n\n")
					tablename := uuidToTablename(s)

					// IF properties[prop] = timeseries later is same behaviour delete this !
					// propertyMapOfElement := graphData[elementVar][i].(neo4j.Entity).GetProperties()

					_, timeseries, err := getTimeSeries(from, to, "", tablename)

					fmt.Printf("\n\n	    fetched time-series: %v      \n\n", timeseries)
					if err != nil {
						return graphData, fmt.Errorf("%w; error - couldnt fetch  properties for %v of element", err, prop)
					}
					// new element for the add property/time-series list
					// var sb strings.Builder
					// sb.WriteString(elementVar)
					// sb.WriteString(".")
					// sb.WriteString(prop)
					// lookup := sb.String()
					// graphData[lookup] = append(graphData[lookup], properties)

					// merge it into the element if it's part of the RETURN clause

					// if mergeVariables {
					properties[prop] = timeseries
					fmt.Printf("\n\n	    merged time-series: %v      \n\n", properties[prop])
					// }
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

func propertyLookupWhereReturnShallow(queryInfo parser.ParseResult) (map[string][]interface{}, error) {
	where, err := manipulateWhereClause(queryInfo, queryInfo.WhereClause)
	if err != nil {
		return nil, fmt.Errorf("%w; error manipulating WHERE query for neo4j", err)
	}
	graphData, err := getShallow(queryInfo, where)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving graph data", err)
	}

	// relevantLookups is an array of relevant lookups in the where clause.
	relevantLookups := queryInfo.LookupsWhereRelevant

	// query the data from timescaleDB according to the property lookups in the original WHERE clause
	graphData, err = filterForCondLookupsInWhere(queryInfo, graphData, relevantLookups)

	if err != nil {
		return nil, fmt.Errorf("%w; error filtering query result on WHERE conditions", err)
	}

	returnVariables := queryInfo.ReturnProjections
	lookupsReturn := queryInfo.LookupsReturn

	graphData, err = getPropertiesforLookupsInReturn(queryInfo, lookupsReturn, returnVariables, graphData)

	// filter out all elements from the neo4j answer which do not match the original WHERE clause with

	if err != nil {
		return nil, fmt.Errorf("%w; error filtering and merging queried data", err)
	}

	return graphData, nil
}

// Replaces the RETURN clause of the query with "RETURN *" if returnAll == true, add temporal boundaries in the WHERE clause
// and receives the according data from neo4j

// TODO: implement test
func getShallow(queryInfo parser.ParseResult, whereManipulated string) (map[string][]interface{}, error) {

	tmpWhere := buildTmpWhereClause(queryInfo.From, queryInfo.To, whereManipulated, queryInfo.GraphElements.MatchGraphElements)
	returnClause := buildReturnClause(queryInfo.LookupsWhereRelevant, queryInfo.GraphElements.ReturnGraphElements)
	query := buildFinalQuery(queryInfo.MatchClause, tmpWhere, returnClause)
	fmt.Printf("\n\n    RETURN CLAUSE    %v\n", returnClause)
	fmt.Printf("    NEO4j QUERY\n    %v\n", query)

	res, err := queryNeo4j(query)

	if err != nil {
		return nil, err
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return resultToMap(res)
}

func filterForCondLookupsInWhere(queryInfo parser.ParseResult, graphData map[string][]interface{}, relevantLookups []parser.LookupInfo) (map[string][]interface{}, error) {

	var err error
	var toRemove []int
	fmt.Printf("\nrelevantLookups: \n%+v\n", relevantLookups)

	filteredData := graphData
	// TODO: retrieve a TREE for AND, OR, XOR, NOT expressions and evaluate accordingly
	for _, lookupInfo := range relevantLookups {

		// elements are all elements which came back from a MATCH pattern for one variable i.e. n of "MATCH (n)-[e]->(s)"
		// lookup.ElementVariable is n in this case
		elements := graphData[lookupInfo.ElementVariable]

		toRemove, err = checkIfValueForConditionExists(queryInfo.From, queryInfo.To, graphData, elements, lookupInfo.Property, lookupInfo.ElementVariable, lookupInfo.CompareOperator, lookupInfo.CompareValue, lookupInfo.LookupLeft)
		filteredData = filterMatches(filteredData, toRemove, []string{})
		if err != nil {
			return nil, fmt.Errorf("%w; error - couldnt merge time series in property", err)
		}

		// andClauses := strings.Split(queryInfo.WhereClause, "AND")
		// orClauses := strings.Split(queryInfo.WhereClause, "OR")

	}
	return filteredData, nil
}

func fetchTimeSeries(from string, to string, graphData map[string][]interface{}, elements []interface{}, property string, elementVar string, mergeVariables bool) (map[string][]interface{}, error) {
	for i, el := range elements {
		switch e := el.(type) {
		case neo4j.Entity:
			uuid := e.GetProperties()[property]
			if uuid == nil {
				// property not available - do nothing
				log.Printf("\nproperty %v not available on element with id : %v\n", property, e.GetElementId())
			} else if s, ok := uuid.(string); ok {
				tablename := uuidToTablename(s)
				propertyMapOfElement := graphData[elementVar][i].(neo4j.Entity).GetProperties()
				_, properties, err := getTimeSeries(from, to, "", tablename)
				if err != nil {
					return nil, fmt.Errorf("%w; error - couldnt fetch  properties for %v of element", err, property)
				}
				// new element for the add property/time-series list
				var sb strings.Builder
				sb.WriteString(elementVar)
				sb.WriteString(".")
				sb.WriteString(property)
				lookup := sb.String()

				// todo: assign pre sized array, add via index instead of append
				graphData[lookup] = append(graphData[lookup], properties)

				// merge it into the element if it's part of the RETURN clause
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

// this function takes the results of all elements of one variable from the MATCH pattern. For example, if the MATCH pattern is "(n)-[e]->(s)",
// it would take all elements of n if a lookup is happening on n. In addition, it takes a lookup property p
// It then iterates over all elements, gets the uuid of the elements property p to fetch the time series from timescaleDB. If the property doesn't
// exist for this element nothing happens. If the property exists, the time series is merged in the result set.
// So if the property does not exist it is not automatically removed from the result set. This is only the case, if
//

// THE RETURNED RESULT SET CONTAINS STILL ALL ELEMENTS FROM THE MATCH PATTERN RETAINED BY NEO4j. IF IT IS PRE-FILTERED AFTER EXISTING PROERPTIES DEPENDS
// ON IF THERE IS A COMPARISON IN THE WHERE CLAUSE.
func checkIfValueForConditionExists(from string, to string, graphData map[string][]interface{}, elements []interface{}, property string, elementVar string, compareOp string, compareVal any, lookupLeft bool) ([]int, error) {

	// These represent the matches from the clause where the MATCH claus did not return any elements so they are removed
	// from the map. Our implementation for comparisons only put matches in the result set if the time-series/property
	// it is compared over still contains entries. For example "FROM x TO y Match (n) where n.name = "Max" return n" will only return
	// n if ...
	// Q: what happens for the case when we have a comparison and filter only the elements which fulfill the comparison?
	// A:	- in this case we request the time-series (filteredProperties) based on the comparison. If it is emtpy we add the
	//		 index of the MATCH position to the remove list to later remove it from the result set. If it it not empty we merge
	//		 it into the result set (into the single element of the MATCH it belongs to)
	// Q: what happens for the case when dont have a comparioson. Like a lookup inside the RETURN clause.
	// 	 Do we still call the getPropertyFromTableCmp function? like does it really send a query to timescaleDB? this would be immense overhead
	//   In this case we should split this up
	// A:
	rowsToRemove := []int{}

	for i, el := range elements {
		switch e := el.(type) {
		case neo4j.Entity:
			uuid := e.GetProperties()[property]
			if uuid == nil {
				// Q: shouldnt i add the element e to remove in the case of a lookup comparison clause in WHERE ?
				// A: not necessary because they are filtered already by the WHERE clause because I ask x.prop IS NOT NULL. So this cannot happen (should not)
				//    NOTE: BETTER DOUBLE CHECK THIS SOMEHOW. THEREFORE I HAVE TO SOLIT THIS..
				// in the case of a lookup in the RETURN clause, it is not necessary to remove the element
				log.Printf("\nproperty %v not available on element with id : %v\n", property, e.GetElementId())
			} else if s, ok := uuid.(string); ok {
				tablename := uuidToTablename(s)
				exists, err := checkIfValueWithConditionExists(from, to, "", compareOp, compareVal, lookupLeft, tablename)
				if err != nil {
					return nil, fmt.Errorf("%w; error - check if value with condidtion exists for time-series %v of element %v", err, property, e.GetElementId())
				} else if exists {
					//propertyMapOfElement := graphData[elementVar][i].(neo4j.Entity).GetProperties()
					fmt.Println("reached 2")
					// _, properties, err := getPropertyFromTable(from, to, "", tablename)
					if err != nil {
						return nil, fmt.Errorf("%w; error - couldnt fetch  properties for %v of element", err, property)
					}
					fmt.Println("reached 3 - all properties should be fetched")
					//propertyMapOfElement[property] = properties
				} else {
					fmt.Printf("\nfiltered properties is nil for %v on element %v\n", property, elementVar)
					// filtered properties is nil so we have to remove the match from the result set (be aware that the match can include multiple elements)

					// NOTE: DO THIS NEXT
					// we only want to do this if a property check has occoured in the WHERE clause
					// if it has occoured in the RETURN CLAUSE we want ro also return the elements which have no elements in the time series
					// if they do not have this property do not return this property
					rowsToRemove = append(rowsToRemove, i)
				}
			} else {
				return nil, errors.New("error - uuid is not a string - this should not happen")
			}
		default:
			panic("error - type not supportet")
		}
	}

	// // remove elements which are filtered from the match
	// // Iterate over the indices in reverse order so removing an element does not change the indices of the remaining elements
	// for i := len(rowsToRemove) - 1; i >= 0; i-- {
	// 	// Remove the i-th element from each slice in graphData
	// 	for _, elements := range graphData {
	// 		elements = append(elements[:rowsToRemove[i]], elements[rowsToRemove[i]+1:]...)
	// 	}
	// }

	// returning graphData is unnecessary because maps are always passed by reference
	// leave it like that. more readable - just a reference to the same map anyways
	return rowsToRemove, nil
}

// TODO: handle exceptions (not in the sense of errors but for example if some matches should explicitely not be removed)
// expects a valid list of indices in ascending order to remove elements from graphData arrays
func filterMatches(graphData map[string][]interface{}, rowsToRemove []int, exceptions []string) map[string][]interface{} {
	// remove elements which are filtered from the match
	// note: the indices in rowsToRemove are sorted in ascending order. Iterate over the indices in reverse order so removing an element does not change the indices of the remaining elements
	for elVar, elements := range graphData {
		for i := len(rowsToRemove) - 1; i >= 0; i-- {
			elements = utils.RemoveIdxFromSlice(elements, rowsToRemove[i])
			graphData[elVar] = elements
		}
	}
	return graphData
}

// // i have to determine the direction of the comparison before
// func filterPropertiesByCompareOp(properties []TimeSeriesRow, compareOp string, compareVal any) ([]TimeSeriesRow, error) {
// 	if strings.TrimSpace(compareOp) == "" || compareVal == nil {
// 		// this is not an error case! it just means that the property is not part of a comparison and nothing has to be filtered
// 		return properties, nil
// 	}
// 	var filtered []TimeSeriesRow
// 	for _, prop := range properties {
// 		if matched, err := utils.CompareValues(prop.Value, compareVal, compareOp); err != nil {
// 			return nil, err
// 		} else if matched {
// 			fmt.Println("matched: ", prop.Value)
// 			filtered = append(filtered, prop)
// 		}
// 	}
//
// 	return filtered, nil
// }

// shouldnt I be able to get a list of these in the listener??

func handleErrorOnResult(res map[string][]any, err error) (bool, error) {
	if err != nil && res == nil {
		return false, err
	} else if err != nil {
		// this still right ?
		log.Printf("Not all elements contained the property: %v", err)
	}
	return true, nil
}
