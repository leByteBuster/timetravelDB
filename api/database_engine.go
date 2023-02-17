package api

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/LexaTRex/timetravelDB/parser"
	"github.com/LexaTRex/timetravelDB/parser/listeners"
	tti "github.com/LexaTRex/timetravelDB/parser/ttql_interface"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// process a TTQL Query
func ProcessQuery(query string) error {

	var queryResult neo4j.ResultWithContext
	UNUSED(queryResult)

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
			prettyPrintMapOfArrays(qRes)
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
				res, err := propertyLookupWhereShallow(queryInfo)
				if err != nil && res == nil {
					return fmt.Errorf("error executing shallow query with lookup in RETURN: %v", err)
				} else if err != nil {
					log.Printf("Not all elements contained the property: %v", err)
				}
				log.Printf("\nto return:\n ")
				prettyPrintMapOfArrays(res)

				// TODO: implement this
				// propertyLookupWhereReturnShallow(queryInfo)
			} else if isWhere {
				log.Println("checkpoint3")
				res, err := propertyLookupWhereShallow(queryInfo)
				if err != nil && res == nil {
					return fmt.Errorf("error executing shallow query with lookup in RETURN: %v", err)
				} else if err != nil {
					log.Printf("Not all elements contained the property: %v", err)
				}
				log.Printf("\nto return:\n ")
				prettyPrintMapOfArrays(res)
			} else if isReturn {
				log.Println("checkpoint4")
				res, err := propertyLookupReturnShallow(queryInfo)
				if err != nil && res == nil {
					return fmt.Errorf("error executing shallow query with lookup in RETURN: %v", err)
				} else if err != nil {
					log.Printf("Not all elements contained the property: %v", err)
				}
				log.Printf("\nto return:\n ")
				prettyPrintMapOfArrays(res)
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
	res, err := filterResultPropertyCalc(queryInfo, graphData)

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
			graphData, err = mergeTimeSeriesInProperty(queryInfo.From, queryInfo.To, graphData, fits, lookup, elVar)
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

func filterResultPropertyCalc(queryInfo parser.ParseResult, graphData map[string][]any) (map[string][]any, error) {

	// relevantLookups is an array of relevant lookups in the where clause. for every lookup we have a LookupInfo struct
	// which contains all necessarythe information about the lookup
	// [n.name, n.age, n.address, e.name, e.age, e.address]
	// up to now relevant lookups in the WHERE clause considers only considers lookups which are part of a comparison like (n.name = "Max")
	relevantLookups, err := getRelevantLookupInfoWhere(queryInfo)
	if err != nil {
		return nil, fmt.Errorf("%w; error - failed to retrieve relevant lookups", err)
	}

	// query the data from timescaleDB according to the property lookups in the original WHERE clause
	res, err := getPropertiesforLookupsInWhere(queryInfo, graphData, relevantLookups)
	if err != nil {
		return nil, fmt.Errorf("%w; error retrieving properties for lookups in WHERE", err)
	}

	// filter out all elements from the neo4j answer which do not match the original WHERE clause with

	return res, nil
}

func getPropertiesforLookupsInWhere(queryInfo parser.ParseResult, graphData map[string][]interface{}, relevantLookups []LookupInfo) (map[string][]interface{}, error) {

	var err error
	fmt.Printf("\nrelevantLookups: \n%+v\n", relevantLookups)

	// TODO: retrieve a TREE for AND, OR, XOR, NOT expressions and evaluate accordingly
	for _, lookupInfo := range relevantLookups {

		// elements are all elements which came back from a MATCH pattern for one variable i.e. n of "MATCH (n)-[e]->(s)"
		// lookup.ElementVariable is n in this case
		elements := graphData[lookupInfo.ElementVariable]

		graphData, err = mergeTimeSeriesInPropertyCmp(queryInfo.From, queryInfo.To, graphData, elements, lookupInfo.Property, lookupInfo.ElementVariable, lookupInfo.CompareOperator, lookupInfo.CompareValue, lookupInfo.LookupLeft)
		if err != nil {
			return nil, fmt.Errorf("%w; error - couldnt merge time series in property", err)
		}

		// andClauses := strings.Split(queryInfo.WhereClause, "AND")
		// orClauses := strings.Split(queryInfo.WhereClause, "OR")

	}
	return graphData, nil
}

func mergeTimeSeriesInProperty(from string, to string, graphData map[string][]interface{}, elements []interface{}, property string, elementVar string) (map[string][]interface{}, error) {
	return mergeTimeSeriesInPropertyCmp(from, to, graphData, elements, property, elementVar, "", "", false)
}

func mergeTimeSeriesInPropertyCmp(from string, to string, graphData map[string][]interface{}, elements []interface{}, property string, elementVar string, compareOp string, compareVal any, lookupLeft bool) (map[string][]interface{}, error) {

	// These represent the matches from the clause where the MATCH claus did not return any elements so they are removed
	// from the map. Our implementation for comparisons only put matches in the result set if the time-series/property
	// it is compared over still contains entries. For example "FROM x TO y Match (n) where n.name = "Max" return n" will only return
	// n if
	rowsToRemove := []int{}

	for i, el := range elements {
		switch e := el.(type) {
		case neo4j.Entity:
			uuid := e.GetProperties()[property]
			if uuid == nil {
				log.Printf("property %v not available on element", property)
			} else if s, ok := uuid.(string); ok {
				tablename := uuidToTablename(s)
				_, filteredProperties, err := getPropertyFromTableCmp(from, to, "", compareOp, compareVal, lookupLeft, tablename)
				if err != nil {
					return nil, fmt.Errorf("%w; error - couldnt fetch filtered properties for %v of element", err, property)
				} else if len(filteredProperties) > 0 {
					fmt.Println("reached 1: ", filteredProperties)
					propertyMapOfElement := graphData[elementVar][i].(neo4j.Entity).GetProperties()
					fmt.Println("reached 2")
					_, properties, err := getPropertyFromTable(from, to, "", tablename)
					if err != nil {
						return nil, fmt.Errorf("%w; error - couldnt fetch  properties for %v of element", err, property)
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
					rowsToRemove = append(rowsToRemove, i)
				}
			} else {
				return nil, errors.New("error - uuid is not a string - this should not happen")
			}
		default:
			panic("error - type not supportet")
		}
	}

	graphData = removeMatchesFromGraphData(graphData, rowsToRemove, []string{})

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
	return graphData, nil
}

// TODO: handle exceptions (not in the sense of errors but for example if some matches should explicitely not be removed)
// expects a valid list of indices in ascending order to remove elements from graphData arrays
func removeMatchesFromGraphData(graphData map[string][]interface{}, rowsToRemove []int, exceptions []string) map[string][]interface{} {
	// remove elements which are filtered from the match
	// note: the indices in rowsToRemove are sorted in ascending order. Iterate over the indices in reverse order so removing an element does not change the indices of the remaining elements
	log.Printf("\nTo remove: %+v\n", rowsToRemove)
	for elVar, elements := range graphData {
		for i := len(rowsToRemove) - 1; i >= 0; i-- {
			log.Printf("\nElement to be removed %v: %+v\n", elVar, elements[i])
			log.Printf("\nList before removal %+v\n", elements)
			elements = removeIdxFromSlice(elements, rowsToRemove[i])
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
		if matched, err := compareValues(prop.Value, compareVal, compareOp); err != nil {
			return nil, err
		} else if matched {
			fmt.Println("matched: ", prop.Value)
			filtered = append(filtered, prop)
		}
	}

	return filtered, nil
}

// shouldnt I be able to get a list of these in the listener??
type LookupInfo struct {
	ElementVariable string
	Property        string
	CompareOperator string
	CompareValue    any
	LookupLeft      bool // a.prop > 5 -> true, 5 > a.prop -> false
}

func getRelevantLookupInfoWhere(queryInfo parser.ParseResult) ([]LookupInfo, error) {

	var lookupInfos []LookupInfo

	for compCtx, insights := range queryInfo.PropertyClauseInsights {
		log.Printf("\nInsights for compareContext: %v \ninsights: %+v\n", compCtx.GetText(), insights)

		var elVar string
		var property string
		var compareOperator string // check if this is retrieved the right way in listener. Test if two symbol operators like <= are recognized correctly
		var compareValueStr string
		var compareValue any
		var lookupLeft bool

		switch len(insights) {
		case 0:
			return nil, errors.New("no insights found for comparison. should be impossible if comparison is in list")
		case 1:
			if !insights[0].IsAppendixOfNullPredicate && insights[0].IsWhere {
				return nil, errors.New("single lookups withouth appendix of null predicate (IS NULL / IS NOT NULL) only alloed in return")
			}
			continue
		// in this case it should be a comparison like "a.prop > 3"
		case 2:
			fmt.Println("Insight Left: ", insights[0])
			fmt.Println("Insight Right: ", insights[1])
			insightLeft := insights[0]
			insightRight := insights[1]
			if !insightLeft.IsWhere || !insightRight.IsWhere {
				return nil, errors.New("comparison not in WHERE clause")
			}
			if insightLeft.IsPartialComparison {
				compareOperator = insightLeft.CompareOperator
			} else if insightRight.IsPartialComparison {
				compareOperator = insightRight.CompareOperator
			} else {
				return nil, errors.New("comparison expression with two propertylabel expressions that include no partial comparison")
			}
			if insightLeft.IsPropertyLookup {
				lookupLeft = true
				elVar = insightLeft.Element
				property = insightLeft.PropertyKey
				compareValueStr = insightRight.Element
			} else if insightRight.IsPropertyLookup {
				elVar = insightRight.Element
				property = insightRight.PropertyKey
				compareValueStr = insightLeft.Element // if insight represents literal then Element is the CompareValue
				lookupLeft = false
			} else {
				continue
			}
		default:
			return nil, errors.New("chained comparisons are not allowed")
		}

		compareValue = convertString(compareValueStr)

		// should only end up here if there is a comparison with a property lookup
		lookupInfos = append(lookupInfos, LookupInfo{ElementVariable: elVar, Property: property, CompareOperator: compareOperator, CompareValue: compareValue, LookupLeft: lookupLeft})
	}

	return lookupInfos, nil
}
