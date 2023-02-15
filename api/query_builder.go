package api

import (
	"context"
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

	res, err := parser.ParseQuery(query)
	if err != nil {
		return err
	}

	// maybe doulbe check if from, to is valid ISO8601

	// is shallow
	if res.IsShallow {
		if !res.ContainsPropertyLookup || res.ContainsPropertyLookup && res.ContainsOnlyNullPredicate {
			queryResult, err := getReturnShallow(res, false)

			if err != nil {
				return err
			}
			log.Println("checkpoint1")
			log.Printf("QUERIED DATA: \n")
			prettyPrintMapOfArrays(queryResult)
			return nil
		} else {
			// only nullpredicate lookups or TODO: no lookups at all
			isValid, isWhere, isReturn := getPropertyLookupParentClause(res.PropertyClauseInsights)
			if !isValid {
				return errors.New("invalid query, property lookup only allowed in WHERE or RETURN clause")
			} else if isWhere && isReturn {
				log.Println("checkpoint2")
				propertyLookupWhereReturnShallow(res)
			} else if isWhere {
				log.Println("checkpoint3")
				// TODO
				propertyLookupWhereShallow(res)
			} else if isReturn {
				log.Println("checkpoint4")
				res, err := propertyLookupReturnShallow(res)
				if err != nil && res == nil {
					return fmt.Errorf("error executing shallow query with lookup in RETURN: %v", err)
				} else if err != nil {
					log.Printf("Not all elements contained the property: %v", err)
				}
				UNUSED(res)
				log.Printf("\nto return:\n ")
				prettyPrintMapOfArrays(res)
				return nil
			} else {
				fmt.Printf("\nReturn: %v, Where: %v, Valid: %v\n", isReturn, isWhere, isValid)
				return errors.New("this option should not be possible")
			}
		}
	} else {
		if !res.ContainsPropertyLookup || res.ContainsPropertyLookup && res.ContainsOnlyNullPredicate {

			queryResult, err = getTmpGraphData(res, false)

			if err != nil {
				return err
			}
			// TODO: get all properties for variables in return clause:
			//getPropertyUUIDS of elements
			//queryTimeScale()
			log.Println("checkpoint5")
			log.Printf("to process further: %v", queryResult)
		} else {

			isValid, isWhere, isReturn := getPropertyLookupParentClause(res.PropertyClauseInsights)
			if !isValid {
				return errors.New("invalid query, property lookup only allowed in WHERE or RETURN clause")
			} else if isWhere && isReturn {
				log.Println("checkpoint6")
				propertyLookupWhereReturn(res)
			} else if isWhere {
				log.Println("checkpoint7")
				// TODO
				propertyLookupWhere(res)
			} else if isReturn {
				log.Println("checkpoint8")
				propertyLookupReturn(res)
			} else {
				return errors.New("this option should not be possible")
			}
		}
	}

	return errors.New("no option choosen, this should not occour")
}

// applies temporal boundaries the other clauses are not changed
func getTmpGraphData(res parser.ParseResult, returnAll bool) (neo4j.ResultWithContext, error) {
	tmpWhere := addTempToWhereQuery(res.From, res.To, res.WhereClause, res.GraphElements.MatchGraphElements)
	var sb strings.Builder
	sb.WriteString(res.MatchClause)
	sb.WriteString(tmpWhere)
	if returnAll {
		sb.WriteString(" Return *")
	} else {
		sb.WriteString(res.ReturnClause)
	}
	tmpQuery := sb.String()
	fmt.Println(tmpQuery)
	return queryNeo4j(tmpQuery)
}

func addTempToWhereQuery(from, to, whereClause string, matchElements []string) string {

	var sb strings.Builder
	if strings.TrimSpace(whereClause) == "" {
		sb.WriteString(" WHERE")
	} else {
		sb.WriteString(" ")
		sb.WriteString(whereClause)
		sb.WriteString(" AND")
	}
	for i, el := range matchElements {
		sb.WriteString(" ")
		sb.WriteString(el)
		sb.WriteString(".")
		sb.WriteString("start >= '")
		sb.WriteString(from)
		sb.WriteString("' AND ")
		sb.WriteString(el)
		sb.WriteString(".")
		sb.WriteString("end < '")
		sb.WriteString(to)
		sb.WriteString("' ")
		if i < len(matchElements)-1 {
			sb.WriteString("AND")
		}
	}
	return sb.String()
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

func propertyLookupWhereShallow(queryInfo parser.ParseResult) {
	matchClause := queryInfo.MatchClause
	// originalReturnClause := queryInfo.ReturnClause
	// TODO: put this and further processing in a function which merges the two database reuquests
	// var originalReturn string = res.ReturnClause // TODO: use this later for merging the results
	var sb strings.Builder
	sb.WriteString(matchClause)
	// TODO: split the match clause in parts which are meant for Neo4j or TimeseriesDB
	// TODO: add FROM, TO to WHERE clause
	sb.WriteString(" RETURN *")
	res, err := queryNeo4j(sb.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
	// TODO NEXT: process res
	// iterate over the keys of the records and apply the original RETURN Clause
	// with  record.Get("movieTitle") it should be possible to get the "columns"
	// PROBLEM: if movieTitle is an alias - how to get the real name which was used in MATCH so we can merge
	// the results of both queries ? But here it doesnt matter because we Return *
	// but for the saved RETRUN clause (the original) it matters. Think about it.

}

// only RETURN clause contains property lookups which needs double database querying
// so send everything to neo4j but with RETURN * instead the original RETURN clause
// and then take care of the original RETURN clause
func propertyLookupReturnShallow(queryInfo parser.ParseResult) (map[string][]interface{}, error) {

	// graph data is a map where the RETURN variables of the CYPHER query are mapped
	// to their results ala: {n: [{id:node1, properties}, {id:node2, properties}], s: ..., e: ...} for "...RETURN n, s, e"
	// NOTE: graphData can contain the same relations multiple times if it is returned multiple times if the pattern matches multiple times
	graphData, err := getReturnShallow(queryInfo, true)
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
	var err = errors.New("")

	// elVar represents the element variable of the RETURN clause the lookups are happening on (i.e. n, s, e)
	// lookups represents all the lookups which are happening on the element variable elVar (i.e. n.name, n.age, n.address)
	for elVar, lookups := range lookupsMap {
		fits := graphData[elVar]

		for _, lookup := range lookups {

			// one element in fits represents one graph element (node or relationship) which was matched by a variable of the neo4j query
			for i, el := range fits {
				switch e := el.(type) {
				case neo4j.Node:
					uuid := e.Props[lookup]
					if uuid == nil {
						log.Printf("property %v not available on node", lookup)
					} else if s, ok := uuid.(string); ok {
						tablename := uuidToTablename(s)
						_, properties, err2 := getPropertyFromTable(queryInfo.From, queryInfo.To, "", tablename)
						if err2 != nil {
							err = fmt.Errorf("%w; error - couldnt fetch properties for %v of node", err2, lookup)
						} else {
							// merge properties into graphData
							// NOTE: to do this i might want to store pointers to arrays in graphData instead of arrays itself so i dont have
							// to get through all the access finding the right element in graphData for the uuid to be updated
							graphData[elVar][i].(neo4j.Node).Props[lookup] = properties
						}
					} else {
						return nil, errors.New("error - uuid is not a string - this should not happen")
					}
				case neo4j.Relationship:
					uuid := e.Props[lookup]
					if uuid == nil {
						log.Printf("property %v not available on edge", lookup)
					} else if s, ok := uuid.(string); ok {
						tablename := uuidToTablename(s)
						_, properties, err2 := getPropertyFromTable(queryInfo.From, queryInfo.To, "", tablename)
						if err2 != nil {
							err = fmt.Errorf("%w; error - couldnt fetch properties for %v of edge", err2, lookup)
						} else {
							// merge properties into graphData
							// NOTE: to do this i might want to store pointers to arrays in graphData instead of arrays itself so i dont have
							// to get through all the access finding the right element in graphData for the uuid to be updated
							graphData[elVar][i].(neo4j.Node).Props[lookup] = properties
						}
					} else {
						return nil, errors.New("error - uuid is not a string - this should not happen")
					}
				default:
					panic("error - type not supportet")
				}

			}
		}
	}
	return graphData, err
}

func propertyLookupWhereReturnShallow(res parser.ParseResult) {
	panic("unimplemented")
}

// Replaces the RETURN clause of the query with "RETURN *", add temporal boundaries in the WHERE clause
// and receives the according data from neo4j
func getReturnShallow(queryInfo parser.ParseResult, returnAll bool) (map[string][]interface{}, error) {

	res, err := getTmpGraphData(queryInfo, returnAll)

	if err != nil {
		return nil, err
	}

	if res.Err() != nil {
		return nil, err
	}

	return resultToMap(res)
}

// format the result of a neo4j query to a map of arrays
// every entry in the map represents a column in the result (a variable of the match clause)
func resultToMap(res neo4j.ResultWithContext) (map[string][]interface{}, error) {
	var formRes map[string][]any = map[string][]any{}
	for res.Next(context.Background()) {

		record := res.Record()

		for _, el := range record.Keys {

			// TODO: test if indexed map is faster (see helpers)
			elRec, ok := record.Get(el)
			if !ok {
				return nil, errors.New("Error getting value for column: " + el)
			}
			if formRes[el] == nil {
				formRes[el] = []interface{}{elRec}
				// elProperties[el] = make([]any, 0)
			} else {
				formRes[el] = append(formRes[el], elRec)
			}

			// MAYBE: use this later instead if types of
			// 				elements must be known (in edges and nodes)
			// if node, ok := elRec.(neo4j.Node); ok {
			// properties := node.Props
			// elProperties[el] = append(elProperties[el], properties)
			// } else if edge, ok := elRec.(neo4j.Relationship); ok {
			// properties := edge.Props
			// elProperties[el] = append(elProperties[el], properties)
			// } else {
			// return nil, errors.New("couldnt assert type of returned record element")
			// }
		}
	}

	if res.Err() != nil {
		log.Fatalf("\nNext error on query result: %v", res.Err())
	}

	// prettyPrintMapOfArrays(formRes)
	// prettyPrintMapOfArrays(elProperties)

	// PROBLEM: if movieTitle is an alias - how to get the real name which was used in MATCH so we can merge
	// the results of both queries ? But here it doesnt matter because we Return *
	// but for the saved RETRUN clause (the original) it matters. Think about it.
	// SOLUTION: Get the keys of the recod: record.Keys (should be the alias already if we just pass the alias to Neo4j (which we dont do yet))
	return formRes, nil

}
