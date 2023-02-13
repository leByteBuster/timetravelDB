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
	IsShallow                 bool
	ContainsPropertyLookup    bool
	ContainsOnlyNullPredicate bool
	From                      string
	To                        string
	MatchClause               string
	WhereClause               string
	ReturnClause              string
	GraphElements             li.GraphElements
	PropertyClauseInsights    map[tti.OC_ComparisonExpressionContext][]li.PropertyClauseInsight
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
	for _, subquery := range propertyClauseInsights {
		subqueryClause := subquery.ComparisonContext.GetText() // this should be the part of the string to be cut out
		subqueryClauseLookup := subquery.PropertyLookupClause
		comparisonCtx := subquery.ComparisonContext
		field := subquery.Field
		propKeys := subquery.PropertyKeys
		labels := subquery.Labels
		compareOp := subquery.CompareOperator

		isWhere := subquery.IsWhere
		isReturn := subquery.IsReturn
		isComparison := subquery.IsComparison
		isPartialComparison := subquery.IsPartialComparison
		isPropertyLookup := subquery.IsPropertyLookup
		if subquery.IsPropertyLookup {
			containsPropertyLookup = true
		}
		if !subquery.IsAppendixOfNullPredicate {
			containsOnlyNullPredicate = false
		}
		isValid := subquery.IsValid

		fmt.Printf("\nComparisonWithPropertyLookupQuery: %v\nPropertyLookupSubquery: %v \ncomparisonCtx: %v \nfield: %v \npropKeys: %v \nlabels: %v \ncompareOp: %v", subqueryClause,
			subqueryClauseLookup, comparisonCtx, field, propKeys, labels, compareOp)

		// print all of the subquery insights
		fmt.Printf("\nIsWhere: %v	\nIsReturn: %v	\nIsComparison: %v	\nIsPartialComparison: %v	\nIsPropertyLookup: %v	\nIsValid: %v",
			isWhere, isReturn, isComparison, isPartialComparison, isPropertyLookup, isValid)

		fmt.Println("")
		fmt.Println("............................................")
		fmt.Println("............................................")
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
		PropertyClauseInsights:    propertyClauseInsights,
	}, nil
}
