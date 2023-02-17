package listeners

import (
	"fmt"
	"strings"

	tti "github.com/LexaTRex/timetravelDB/parser/ttql_interface"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

type TimePeriod struct {
	From string
	To   string
}

type GraphElements struct {
	MatchGraphElements  []string
	WhereGraphElements  []string
	ReturnGraphElements []string
}

// TreeShapeListener is a listener that enters every node in the parser tree and prints the tree shape of the parse tree

type TtqlTreeListener struct {
	*tti.BaseTTQLListener
	// every comparison expression is mapped onto a list of its containing property clause insights
	// a property clause insight is not always a property lookup
	Insights map[*tti.OC_ComparisonExpressionContext][]PropertyClauseInsight
	TimePeriod
	IsShallow bool
	GraphElements
	MatchClause  string
	WhereClause  string
	ReturnClause string
}

func NewTtqlTreeListener() *TtqlTreeListener {
	tl := new(TtqlTreeListener)
	tl.Insights = map[*tti.OC_ComparisonExpressionContext][]PropertyClauseInsight{}
	return tl
}

type PropertyClauseInsight struct {
	PropertyLookupClause string // string of the lookup or literal value: "a.property" or "22"
	// ComparisonContext    *tti.OC_ComparisonExpressionContext // if isPropertyLookup is true, the comparison string can be retrieved from this to split the Clause
	Element         string   // element of the lookup (variable of node/edge) represented as string
	PropertyKey     string   // the property key (only one is necessary because of flattening n.property1.property2 should not be necessary/possible)
	Labels          []string // labels are not supporte yet here
	CompareOperator string   // compariosn operator - if the Comparison expression contains a CompareOperator
	//take care. this can be a property lookup as well. how to handle ?
	// maybe set a bool. if true then use the next propertyClauseInsights in line as comparison value
	// or maybe just a pointer ?
	IsPropertyLookup bool // contains this PropertyClause a property lookup ?
	IsComparison     bool // is it part of a comparison ? (i think we can delete this because property clauses are safed in a map with the comparison object
	// as the key for a list of Insight)
	IsPartialComparison       bool // is it part of a partial comparison ?
	IsWhere                   bool // part of a WHERE clause ?
	IsReturn                  bool // part of a RETURN clause ?
	IsValid                   bool // is Part of a RETURN or WHERE clause ? if not invalid
	CountPartialComparison    int  // how many PartialComparisons
	IsAppendixOfNullPredicate bool // does a NullPredicate exist as suffix ?
}

// Listener For PropertyOrLabelsExpression (to get all property lookups inside the clause, with all concatenated propertylookups
// and all concatenated labels)

func (listener *TtqlTreeListener) EnterOC_PropertyOrLabelsExpression(pOLE *tti.OC_PropertyOrLabelsExpressionContext) {
	propertyLookupClause := pOLE.GetText()

	var isWhere = false
	var isReturn = false
	var isPropertyLookup = false
	var isComparison = false
	var isPartialComparison = false
	var isAppendixOfNullPredicate = false

	var comparison *tti.OC_ComparisonExpressionContext
	var propertyLookup *tti.OC_PropertyLookupContext
	var partialComparison *tti.OC_PartialComparisonExpressionContext
	var countPartialComparison int = 0

	var element string = pOLE.GetChild(0).(*tti.OC_AtomContext).GetText()
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
	// Note: invalid if: property lookup true && !isReturn && !isWhere

	// Note: this did not work because also for literals (e.g. 22) I want to know if isReturn or isWhere
	//for parent != nil && (!isPropertyLookup && !isComparison || isPropertyLookup && !(isReturn || isWhere)) {
	for parent != nil && !(isReturn || isWhere) {
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

	if isComparison {
		for _, ctx := range comparison.GetChildren() {
			if _, ok := ctx.(*tti.OC_PartialComparisonExpressionContext); ok {
				countPartialComparison++
			}
			// NOTE: this might only work in some cases. Because PartialComparisons itself
			//			 contain a StringListNullPredicate and therefore can contain
			//			 a NullPredicate.
			if sCtx, ok := ctx.(*tti.OC_StringListNullPredicateExpressionContext); ok {
				for _, ctx2 := range sCtx.GetChildren() {
					if _, ok := ctx2.(*tti.OC_NullPredicateExpressionContext); ok {
						isAppendixOfNullPredicate = true
					}
				}
			}
		}

	}

	if isPartialComparison {
		// the first child of the PartialComparisonExpression is always a compare token
		// this might be wrong since it maybe does not consider two character comparison operators
		compareOperator = partialComparison.GetChild(0).GetPayload().(antlr.Token).GetText()
	}

	listener.Insights[comparison] = append(listener.Insights[comparison], PropertyClauseInsight{
		PropertyLookupClause:      propertyLookupClause,
		Element:                   element,
		PropertyKey:               propertyKey,
		Labels:                    []string{},
		IsWhere:                   isWhere,
		IsReturn:                  isReturn,
		IsValid:                   (isWhere || isReturn) || !isPropertyLookup,
		IsComparison:              isComparison,
		CompareOperator:           compareOperator,
		IsPropertyLookup:          isPropertyLookup,
		IsPartialComparison:       isPartialComparison,
		CountPartialComparison:    countPartialComparison,
		IsAppendixOfNullPredicate: isAppendixOfNullPredicate,
	})

}

func (listener *TtqlTreeListener) EnterTtQL_Query(qC *tti.TtQL_QueryContext) {

	isShallow := false
	tC := qC.TtQL_TimeClause()

	lastChild := qC.GetChild(qC.GetChildCount() - 1)
	switch c := lastChild.GetPayload().(type) {
	case antlr.Token:
		if c.GetTokenType() == tti.TTQLParserSHALLOW {
			isShallow = true
		}
	}

	var from antlr.Token
	var to antlr.Token
	for _, child := range tC.GetChildren() {

		// all children are antlr Tokens
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
	listener.IsShallow = isShallow
}

func (listener *TtqlTreeListener) EnterOC_Variable(vC *tti.OC_VariableContext) {
	el := vC.GetText()
	parent := vC.GetParent()
	for parent != nil {
		if _, ok := parent.(*tti.OC_MatchContext); ok {
			listener.GraphElements.MatchGraphElements = append(listener.GraphElements.MatchGraphElements, el)
			break
		} else if _, ok := parent.(*tti.OC_WhereContext); ok {
			listener.GraphElements.WhereGraphElements = append(listener.GraphElements.WhereGraphElements, el)
			break // no need to check further
		} else if _, ok := parent.(*tti.OC_ReturnContext); ok {
			listener.GraphElements.ReturnGraphElements = append(listener.GraphElements.ReturnGraphElements, el)
			break // no need to check further
		}
		parent = parent.GetParent()
	}

	// after parsing the query, the tree walk is not conducted in case there have been any parse errors.
	// This should not happen and therefore panics if it somehow does.
	if parent == nil {
		panic("This should not happen. variable not in match, where or return clause")
	}

}

func (listener *TtqlTreeListener) EnterOC_Match(wC *tti.OC_MatchContext) {
	var pC *tti.OC_PatternContext
	var sb strings.Builder
	children := wC.GetChildren()

	for _, child := range children {
		if c, ok := child.(*tti.OC_PatternContext); ok {
			pC = c
			break
		}
	}
	sb.WriteString("Match ")
	sb.WriteString(pC.GetText())
	listener.MatchClause = sb.String()
}

func (listener *TtqlTreeListener) EnterOC_Where(wC *tti.OC_WhereContext) {
	listener.WhereClause = wC.GetText()
}

func (listener *TtqlTreeListener) EnterOC_Return(rC *tti.OC_ReturnContext) {
	listener.ReturnClause = rC.GetText()
}

type ErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []string
}

func (el *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	errorMessage := fmt.Sprintf("line %d:%d %s", line, column, msg)
	el.Errors = append(el.Errors, errorMessage)
}

func NewErrorListener() *ErrorListener {
	return new(ErrorListener)
}
