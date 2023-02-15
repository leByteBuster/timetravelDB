package parser

import (
	"errors"
	"fmt"
	"strings"

	li "github.com/LexaTRex/timetravelDB/parser/listeners"
	tti "github.com/LexaTRex/timetravelDB/parser/ttql_interface"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type ParseResult struct {
	IsShallow                 bool   // is the query shallow
	ContainsPropertyLookup    bool   // does it contain any property lookup
	ContainsOnlyNullPredicate bool   // if it contains property lookups - do all of them have a NullPredicates suffix ?
	From                      string // start time
	To                        string // end time
	MatchClause               string
	WhereClause               string
	ReturnClause              string
	GraphElements             li.GraphElements                                                   // all element variables occouring in the query
	LookupsWhere              map[string][]string                                                // all relevant lookups in Where (lookups that are relevant for binary querying) - mapped onto their variable: n: {property1,property2} s: {property1,property4}..
	LookupsReturn             map[string][]string                                                // all relevant lookups in Return (lookups that are relevant for binary querying)
	PropertyClauseInsights    map[*tti.OC_ComparisonExpressionContext][]li.PropertyClauseInsight // insights of Comparison expressions / Property Clauses
}

func ParseQuery(query string) (ParseResult, error) {

	qS := antlr.NewInputStream(query)
	fmt.Println("Lexer Tokens: ")

	lexer := tti.NewTTQLLexer(qS)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := tti.NewTTQLParser(tokens)
	parser.RemoveErrorListeners()
	errorListener := li.NewErrorListener()
	parser.AddErrorListener(errorListener)

	// interrupt parsing as soon as an error occours
	//parser.SetErrorHandler(antlr.NewBailErIrorStrategy())

	treectx := parser.TtQL()

	if errorListener.Errors != nil {
		var sb strings.Builder
		for _, error := range errorListener.Errors {
			sb.WriteString(error)
			sb.WriteString("\n")
		}
		return ParseResult{}, errors.New(sb.String())
	}

	// initialize listeners
	timeShallowListener := li.NewTimeShallowListener()
	propertyListener := li.NewPropertyOrLabelsExpressionListener()
	elementListener := li.NewElementListener()

	// HOW TO GET THE ERRORS ?????
	//		parser.GetErrorHandler()
	//		parser.GetErrorListenerDispatch()
	//		parser._SyntaxErrors
	//
	fmt.Println()
	fmt.Println("............................................")
	fmt.Println()
	fmt.Printf("Tree: %v", treectx.GetText())
	fmt.Println("Query is valid.")
	fmt.Println()
	fmt.Println("............................................")
	fmt.Println()
	antlr.ParseTreeWalkerDefault.Walk(timeShallowListener, treectx)
	fmt.Printf("\nTimeClauseInsights: \n from: %v\nto: %v\nisShallow: %v", timeShallowListener.TimePeriod.From,
		timeShallowListener.TimePeriod.To, timeShallowListener.IsShallow)
	fmt.Println()
	fmt.Println("............................................")
	fmt.Println()
	antlr.ParseTreeWalkerDefault.Walk(elementListener, treectx)
	fmt.Printf("\nMatchClause variables: %v\nWhereClause variables: %v\nReturnClause variables: %v",
		elementListener.GraphElements.MatchGraphElements, elementListener.GraphElements.WhereGraphElements, elementListener.GraphElements.ReturnGraphElements)
	fmt.Println()
	fmt.Println("............................................")
	fmt.Println()

	antlr.ParseTreeWalkerDefault.Walk(propertyListener, treectx)
	propertyClauseInsights := propertyListener.Insights

	fmt.Println("............................................")
	containsPropertyLookup := false
	containsOnlyNullPredicate := true

	lookupsWhere := map[string][]string{}
	lookupsReturn := map[string][]string{}
	for comparisonCtx, listOfInsights := range propertyClauseInsights {
		for _, insight := range listOfInsights {
			insightClause := comparisonCtx.GetText() // this should be the part of the string to be cut out
			insightClauseLookup := insight.PropertyLookupClause
			field := insight.Element
			propKeys := insight.PropertyKey
			labels := insight.Labels
			compareOp := insight.CompareOperator

			isWhere := insight.IsWhere
			isReturn := insight.IsReturn
			isComparison := insight.IsComparison
			isPartialComparison := insight.IsPartialComparison
			isPropertyLookup := insight.IsPropertyLookup
			isAppendixOfNullPredicate := insight.IsAppendixOfNullPredicate

			if insight.IsPropertyLookup {
				containsPropertyLookup = true

				// collect property lookups that are relevant for binary databae fetching (neo4j, timescaleDB)
				if !insight.IsAppendixOfNullPredicate {
					containsOnlyNullPredicate = false
					if isWhere {
						lookupsWhere[insight.Element] = append(lookupsWhere[insight.Element], insight.PropertyKey)
					} else if isReturn {
						lookupsReturn[insight.Element] = append(lookupsReturn[insight.Element], insight.PropertyKey)
					}
				}
			}

			isValid := insight.IsValid

			fmt.Printf("\nComparisonWithPropertyLookupQuery: %v\nPropertyLookupinsight: %v \ncomparisonCtx: %v \nfield: %v \npropKeys: %v \nlabels: %v \ncompareOp: %v", insightClause,
				insightClauseLookup, comparisonCtx, field, propKeys, labels, compareOp)

			// print all of the insight insights
			fmt.Printf("\nIsWhere: %v	\nIsReturn: %v	\nIsComparison: %v	\nIsPartialComparison: %v	\nIsPropertyLookup: %v \nIsAppendixOfNullPredicate: %v	\nIsValid: %v",
				isWhere, isReturn, isComparison, isPartialComparison, isPropertyLookup, isAppendixOfNullPredicate, isValid)

			fmt.Println("")
			fmt.Println("............................................")
			fmt.Println("............................................")
		}
	}

	// String cypherQuery2 = "MATCH (n) WHERE n.ping > 22.33" + " RETURN n.ping, n ";
	// String cypherQuery3 = "MATCH (a)-[x]->(b) " + " RETURN a.ping, b "; // should parse
	// String cypherQuery4 = "MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b "; // should parse

	return ParseResult{
		IsShallow:                 timeShallowListener.IsShallow,
		ContainsPropertyLookup:    containsPropertyLookup,
		ContainsOnlyNullPredicate: containsOnlyNullPredicate,
		From:                      timeShallowListener.TimePeriod.From,
		To:                        timeShallowListener.TimePeriod.To,
		MatchClause:               elementListener.MatchClause,
		WhereClause:               elementListener.WhereClause,
		ReturnClause:              elementListener.ReturnClause,
		GraphElements:             elementListener.GraphElements,
		LookupsWhere:              lookupsWhere,
		LookupsReturn:             lookupsReturn,
		PropertyClauseInsights:    propertyClauseInsights,
	}, nil
}
