package listeners

import (
	"fmt"

	tti "github.com/LexaTRex/timetravelDB/parser/ttql_interface"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// TreeShapeListener is a listener that enters every node in the parser tree and prints the tree shape of the parse tree

type TreeShapeListener struct {
	*tti.BaseTTQLListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (treeShapeListener *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(ctx.GetText())
}

// PropertyLookupListener is a listener that collects insights about property lookup clauses

type PropertyClauseInsights struct {
	PropertyLookupClause string
	ComparisonContext    *tti.OC_ComparisonExpressionContext
	Field                string
	PropertyKeys         []string
	Labels               []string
	CompareOperator      string
	//take care. this can be a property lookup as well. how to handle ?
	// maybe set a bool. if true then use the next propertyClauseInsights in line as comparison value
	// or maybe just a pointer ?
	IsPropertyLookup    bool
	IsComparison        bool
	IsPartialComparison bool
	IsWhere             bool
	IsReturn            bool
	IsValid             bool
}

// type PropertyLookupListener struct {
// 	*tti.BaseTTQLListener
// 	Insights []PropertyClauseInsights
// }
//
// func NewPropertyLookupListener() *PropertyLookupListener {
// 	return new(PropertyLookupListener)
// }
//
// func (pL *PropertyLookupListener) EnterOC_PropertyLookup(ctx *tti.OC_PropertyLookupContext) {
// 	fmt.Println(ctx.GetText())
// 	parent := ctx.GetParent()
// 	field := parent.GetChild(0)
// 	propertyLookupClause := ctx.GetText()
// 	isWhere := false
// 	isReturn := false
// 	for parent != nil && !(isWhere || isReturn) {
// 		if _, ok := parent.(*tti.OC_WhereContext); ok {
// 			isWhere = true
// 		} else if _, ok := parent.(*tti.OC_ReturnContext); ok {
// 			isReturn = true
// 		}
// 		parent = parent.GetParent()
// 	}
// 	if pL.Insights == nil {
// 		pL.Insights = []PropertyClauseInsights{}
// 	}
//
// 	pL.Insights = append(pL.Insights, PropertyClauseInsights{
// 		PropertyLookupClause: propertyLookupClause,
// 		IsWhere:              isWhere,
// 		IsReturn:             isReturn,
// 		IsValid:              isWhere || isReturn,
// 	})
//
// }

// Listener For PropertyOrLabelsExpression (to get all property lookups inside the clause, with all concatenated propertylookups
// and all concatenated labels)

type PropertyOrLabelsExpressionListener struct {
	*tti.BaseTTQLListener
	Insights []PropertyClauseInsights
}

func NewPropertyOrLabelsExpressionListener() *PropertyOrLabelsExpressionListener {
	return new(PropertyOrLabelsExpressionListener)
}

func (listener *PropertyOrLabelsExpressionListener) EnterOC_PropertyOrLabelsExpression(pOLE *tti.OC_PropertyOrLabelsExpressionContext) {
	fmt.Println(pOLE.GetText())
	propertyLookupClause := pOLE.GetText()

	var isWhere = false
	var isReturn = false
	var isPropertyLookup = false
	var isComparison = false
	var isPartialComparison = false

	var comparison *tti.OC_ComparisonExpressionContext
	var propertyLookup *tti.OC_PropertyLookupContext
	var partialComparison *tti.OC_PartialComparisonExpressionContext

	var field string = pOLE.GetChild(0).(*tti.OC_AtomContext).GetText()
	var propertyKey string
	var compareOperator string // if comparisonExpression get compare operator and compare value

	// get the property lookup if existing, looping is necessary for eventual white space
	// (until now only one. but this is easy to extend for multiple property lookups)
	for _, child := range pOLE.GetChildren() {
		if t, ok := child.(*tti.OC_PropertyLookupContext); ok {
			propertyLookup = t
			isPropertyLookup = true
			break
		}
	}

	if isPropertyLookup {
		// get the property key of the lookup, looping is necessary for eventual white space
		for _, child := range propertyLookup.GetChildren() {
			if t, ok := child.(*tti.OC_PropertyKeyNameContext); ok {
				propertyKey = t.GetText()
				break
			}
		}
	}

	parent := pOLE.GetParent()

	// get the comparison expression & check if WHERE or RETURN clause. If not, then it is invalid

	// 1: property Lookup is true: run until WHERE or RETURN is found. If not found, then it is invalid
	// 2: property Lookup is false: run and see if CompareExpression is found. If not, then it is NOT invalid
	// Note: when Searching for WHERE or RETURN, the ComparissonExpression would be passed on the way anywas so no extra check needed
	for parent != nil && (!isPropertyLookup && !isComparison || isPropertyLookup && (isReturn || isWhere)) {
		if e, ok := parent.(*tti.OC_PartialComparisonExpressionContext); ok {
			partialComparison = e
			isPartialComparison = true
		}
		if e, ok := parent.(*tti.OC_ComparisonExpressionContext); ok {
			comparison = e
			isComparison = true
		} else if _, ok := parent.(*tti.OC_WhereContext); ok {
			isWhere = true
			break // no need to check further
		} else if _, ok := parent.(*tti.OC_ReturnContext); ok {
			isReturn = true
			break // no need to check further
		}
		parent = parent.GetParent()
	}

	if isPartialComparison {
		// the first child of the PartialComparisonExpression is always a compare token
		compareOperator = partialComparison.GetChild(0).GetPayload().(antlr.Token).GetText()
		fmt.Printf("PartialComparisonExpression: %v", compareOperator)
	}

	listener.Insights = append(listener.Insights, PropertyClauseInsights{
		PropertyLookupClause: propertyLookupClause,
		ComparisonContext:    comparison,
		Field:                field,
		PropertyKeys:         []string{propertyKey},
		Labels:               []string{},
		CompareOperator:      compareOperator,
		IsComparison:         isComparison,
		IsPartialComparison:  isPartialComparison,
		IsWhere:              isWhere,
		IsReturn:             isReturn,
		IsValid:              isWhere || isReturn,
	})

}

// TimeClauseListener is a listener that collects insights about property lookup clauses

type TimePeriod struct {
	From string
	To   string
}

type TimeClauseListener struct {
	*tti.BaseTTQLListener
	TimePeriod
}

func NewTimeClauseListener() *TimeClauseListener {
	return new(TimeClauseListener)
}

func (listener *TimeClauseListener) EnterTtQL_TimeClause(tC *tti.TtQL_TimeClauseContext) {
	fmt.Println("\nTimeclause: ", tC.GetText())

	var from antlr.Token
	var to antlr.Token
	for _, child := range tC.GetChildren() {
		leaf := child.GetPayload().(antlr.Token)
		if leaf.GetTokenType() == tti.TTQLParserDATETIME {
			if from == nil {
				from = leaf
			} else {
				to = leaf
			}
		}
	}

	listener.TimePeriod = TimePeriod{
		From: from.GetText(),
		To:   to.GetText(),
	}
}
