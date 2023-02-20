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
func ProcessQuery(query string) error {

	var queryResult neo4j.ResultWithContext
	utils.UNUSED(queryResult)

	queryInfo, err := parser.ParseQuery(query)
	if err != nil {
		return err
	}

	// maybe doulbe check if from, to is valid ISO8601

	// is shallow
	if queryInfo.IsShallow {
		// no or no relevant lookups
		if !queryInfo.ContainsPropertyLookup || queryInfo.ContainsPropertyLookup && queryInfo.ContainsOnlyNullPredicate {
			qRes, err := getShallow(queryInfo, queryInfo.WhereClause, false)

			if err != nil {
				return err
			}
			log.Println("checkpoint1")
			log.Printf("QUERIED DATA: \n")
			utils.PrettyPrintMapOfArrays(qRes)
			return nil
		} else {

			isValid, isWhere, isReturn := getPropertyLookupParentClause(queryInfo.PropertyClauseInsights)
			if !isValid {
				return errors.New("invalid query, property lookup only allowed in WHERE or RETURN clause")
			} else if isWhere && isReturn {
				log.Println("checkpoint2")

				// TODO: implement this
				// propertyLookupWhereReturnShallow(queryInfo)

				// NOTE: I just copied the code of checkpoint 3 to test it! change it!
				// ##### start:
				res, err := propertyLookupWhereShallow(queryInfo)
				if err != nil && res == nil {
					return fmt.Errorf("error executing shallow query with lookup in RETURN: %v", err)
				} else if err != nil {
					log.Printf("Not all elements contained the property: %v", err)
				}
				log.Printf("\nto return:\n ")
				utils.PrettyPrintMapOfArrays(res)
				// ##### end

			} else if isWhere {
				log.Println("checkpoint3")
				res, err := propertyLookupWhereShallow(queryInfo)
				if err != nil && res == nil {
					return fmt.Errorf("error executing shallow query with lookup in RETURN: %v", err)
				} else if err != nil {
					log.Printf("Not all elements contained the property: %v", err)
				}
				log.Printf("\nto return:\n ")
				utils.PrettyPrintMapOfArrays(res)
			} else if isReturn {
				log.Println("checkpoint4")
				res, err := propertyLookupReturnShallow(queryInfo)
				if err != nil && res == nil {
					return fmt.Errorf("error executing shallow query with lookup in RETURN: %v", err)
				} else if err != nil {
					log.Printf("Not all elements contained the property: %v", err)
				}
				log.Printf("\nto return:\n ")
				utils.PrettyPrintMapOfArrays(res)
				return nil
			} else {
				fmt.Printf("\nReturn: %v, Where: %v, Valid: %v\n", isReturn, isWhere, isValid)
				return errors.New("this option should not be possible")
			}
		}
	} else {
		if !queryInfo.ContainsPropertyLookup || queryInfo.ContainsPropertyLookup && queryInfo.ContainsOnlyNullPredicate {

			tmpWhere := addTempToWhereQuery(queryInfo.From, queryInfo.To, queryInfo.WhereClause, queryInfo.GraphElements.MatchGraphElements)
			query := buildFinalQuery(queryInfo, tmpWhere, false)
			queryRes, err := queryNeo4j(query)

			if err != nil {
				return err
			}
			// TODO: get all properties for variables in return clause:
			//getPropertyUUIDS of elements
			//queryTimeScale()
			log.Println("checkpoint5")
			log.Printf("to process further: %v", queryRes)
		} else {

			isValid, isWhere, isReturn := getPropertyLookupParentClause(queryInfo.PropertyClauseInsights)
			if !isValid {
				return errors.New("invalid query, property lookup only allowed in WHERE or RETURN clause")
			} else if isWhere && isReturn {
				log.Println("checkpoint6")
				propertyLookupWhereReturn(queryInfo)
			} else if isWhere {
				log.Println("checkpoint7")
				// TODO
				propertyLookupWhere(queryInfo)
			} else if isReturn {
				log.Println("checkpoint8")
				propertyLookupReturn(queryInfo)
			} else {
				return errors.New("this option should not be possible")
			}
		}
	}

	// return errors.New("no option choosen, this should not occour")
	return nil
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

func propertyLookupWhereReturn(res parser.ParseResult) {
	panic("unimplemented")
}

// functions for shallow queries with lookup that requires double database access

func propertyLookupWhereShallow(queryInfo parser.ParseResult) (map[string][]interface{}, error) {

	where, err := manipulateWhereClause(queryInfo, queryInfo.WhereClause)
	if err != nil {
		return nil, fmt.Errorf("%w; error manipulating WHERE query for neo4j", err)
	}
	graphData, err := getShallow(queryInfo, where, true)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving graph data", err)
	}
	// relevantLookups is an array of relevant lookups in the where clause. for every lookup we have a LookupInfo struct
	// which contains all necessarythe information about the lookup
	// [n.name, n.age, n.address, e.name, e.age, e.address]
	// up to now relevant lookups in the WHERE clause considers only considers lookups which are part of a comparison like (n.name = "Max")
	relevantLookups := queryInfo.LookupsWhereRelevant

	// query the data from timescaleDB according to the property lookups in the original WHERE clause
	res, err := getPropertiesforLookupsInWhere(queryInfo, graphData, relevantLookups)
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
	graphData, err := getShallow(queryInfo, queryInfo.WhereClause, true)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving graph data", err)
	}
	mergedRes, err := getPropertiesforLookupsInReturn(queryInfo, graphData)
	if err != nil {
		return mergedRes, fmt.Errorf("%w; Not all properties could be fetched", err)
	}

	return mergedRes, nil

}

func getPropertiesforLookupsInReturn(queryInfo parser.ParseResult, graphData map[string][]interface{}) (map[string][]interface{}, error) {

	// this maps one CYPHER query variable onto all its lookups which occour in the query RETURN clause
	lookupsMap := queryInfo.LookupsReturn
	var err error

	// elVar represents the element variable of the RETURN clause the lookups are happening on (i.e. n, s, e)
	// lookups represents all the lookups which are happening on the element variable elVar (i.e. n.name, n.age, n.address)
	for elVar, lookups := range lookupsMap {
		fits := graphData[elVar]

		for _, lookup := range lookups {
			graphData, err = mergeTimeSeriesLookupInReturn(queryInfo.From, queryInfo.To, graphData, fits, lookup, elVar)
			// graphData = filterMatches(graphData, rowsToRemove, []string{})
		}
	}
	return graphData, err
}

func propertyLookupWhereReturnShallow(res parser.ParseResult) {
	panic("unimplemented")
}

// Replaces the RETURN clause of the query with "RETURN *" if returnAll == true, add temporal boundaries in the WHERE clause
// and receives the according data from neo4j

// TODO: implement test
func getShallow(queryInfo parser.ParseResult, whereManipulated string, returnAll bool) (map[string][]interface{}, error) {

	tmpWhere := addTempToWhereQuery(queryInfo.From, queryInfo.To, whereManipulated, queryInfo.GraphElements.MatchGraphElements)
	query := buildFinalQuery(queryInfo, tmpWhere, returnAll)
	fmt.Printf("\n\n    NEO4j QUERY\n\n     %v\n", query)
	res, err := queryNeo4j(query)

	if err != nil {
		return nil, err
	}

	if res.Err() != nil {
		return nil, err
	}

	return resultToMap(res)
}

func getPropertiesforLookupsInWhere(queryInfo parser.ParseResult, graphData map[string][]interface{}, relevantLookups []parser.LookupInfo) (map[string][]interface{}, error) {

	var err error
	var toRemove []int
	fmt.Printf("\nrelevantLookups: \n%+v\n", relevantLookups)

	filteredData := graphData
	// TODO: retrieve a TREE for AND, OR, XOR, NOT expressions and evaluate accordingly
	for _, lookupInfo := range relevantLookups {

		// elements are all elements which came back from a MATCH pattern for one variable i.e. n of "MATCH (n)-[e]->(s)"
		// lookup.ElementVariable is n in this case
		elements := graphData[lookupInfo.ElementVariable]

		filteredData, toRemove, err = mergeTimeSeriesLookupCmpWhere(queryInfo.From, queryInfo.To, graphData, elements, lookupInfo.Property, lookupInfo.ElementVariable, lookupInfo.CompareOperator, lookupInfo.CompareValue, lookupInfo.LookupLeft)
		filteredData = filterMatches(filteredData, toRemove, []string{})
		if err != nil {
			return nil, fmt.Errorf("%w; error - couldnt merge time series in property", err)
		}

		// andClauses := strings.Split(queryInfo.WhereClause, "AND")
		// orClauses := strings.Split(queryInfo.WhereClause, "OR")

	}
	return filteredData, nil
}

// in comparison to mergeTimeSeriesLookupCmpWhere this function does not filter out elements at all
func mergeTimeSeriesLookupInReturn(from string, to string, graphData map[string][]interface{}, elements []interface{}, property string, elementVar string) (map[string][]interface{}, error) {
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
				_, properties, err := getPropertyFromTable(from, to, "", tablename)
				if err != nil {
					return nil, fmt.Errorf("%w; error - couldnt fetch  properties for %v of element", err, property)
				}
				propertyMapOfElement[property] = properties
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
func mergeTimeSeriesLookupCmpWhere(from string, to string, graphData map[string][]interface{}, elements []interface{}, property string, elementVar string, compareOp string, compareVal any, lookupLeft bool) (map[string][]interface{}, []int, error) {

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
					return nil, nil, fmt.Errorf("%w; error - check if value with condidtion exists for time-series %v of element %v", err, property, e.GetElementId())
				} else if exists {
					propertyMapOfElement := graphData[elementVar][i].(neo4j.Entity).GetProperties()
					fmt.Println("reached 2")
					_, properties, err := getPropertyFromTable(from, to, "", tablename)
					if err != nil {
						return nil, nil, fmt.Errorf("%w; error - couldnt fetch  properties for %v of element", err, property)
					}
					fmt.Println("reached 3 - all properties should be fetched")
					propertyMapOfElement[property] = properties
					//	// if no matched properties delete matched structure from result set
					//	if filteredProperties == nil {
					//		for elVar := range graphData {
					//			elements := graphData[elVar]
					//			fmt.Println()
					//			fmt.Printf("\nPrint elements before deletion for element %v: %+v\n", elVar, elements)
					//			fmt.Println("i: ", i)
					//			graphData[elVar] = append(elements[:i], elements[i+1:]...)
					//			fmt.Printf("\nPrint elements before deletion: %+v\n", graphData[elVar])
					//			fmt.Println()
					//		}
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
				return nil, nil, errors.New("error - uuid is not a string - this should not happen")
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
	return graphData, rowsToRemove, nil
}

// TODO: handle exceptions (not in the sense of errors but for example if some matches should explicitely not be removed)
// expects a valid list of indices in ascending order to remove elements from graphData arrays
func filterMatches(graphData map[string][]interface{}, rowsToRemove []int, exceptions []string) map[string][]interface{} {
	// remove elements which are filtered from the match
	// note: the indices in rowsToRemove are sorted in ascending order. Iterate over the indices in reverse order so removing an element does not change the indices of the remaining elements
	log.Printf("\nTo remove: %+v\n", rowsToRemove)
	for elVar, elements := range graphData {
		for i := len(rowsToRemove) - 1; i >= 0; i-- {
			log.Printf("\nElement to be removed %v: %+v\n", elVar, elements[i])
			log.Printf("\nList before removal %+v\n", elements)
			elements = utils.RemoveIdxFromSlice(elements, rowsToRemove[i])
			graphData[elVar] = elements
			log.Printf("\nList after removal %+v\n", graphData[elVar])
		}
	}
	return graphData
}

// i have to determine the direction of the comparison before
func filterPropertiesByCompareOp(properties []TimeSeriesRow, compareOp string, compareVal any) ([]TimeSeriesRow, error) {
	if strings.TrimSpace(compareOp) == "" || compareVal == nil {
		// this is not an error case! it just means that the property is not part of a comparison and nothing has to be filtered
		return properties, nil
	}
	var filtered []TimeSeriesRow
	for _, prop := range properties {
		if matched, err := utils.CompareValues(prop.Value, compareVal, compareOp); err != nil {
			return nil, err
		} else if matched {
			fmt.Println("matched: ", prop.Value)
			filtered = append(filtered, prop)
		}
	}

	return filtered, nil
}

// shouldnt I be able to get a list of these in the listener??
