package api

import (
	"errors"
	"fmt"
	"log"

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
		fmt.Printf("%+v\n", utils.JsonStringFromMapOrdered(queryResult, queryInfo.ReturnProjections))
	} else {
		fmt.Printf("\n\n\n                      Printed unordered                         \n\n\n\n")
		fmt.Printf("%+v\n", utils.JsonStringFromMap(queryResult))
	}

	// return errors.New("no option choosen, this should not occour")
	return queryResult, nil
}
