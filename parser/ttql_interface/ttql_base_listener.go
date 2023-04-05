// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package ttql_interface // TTQL
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// BaseTTQLListener is a complete listener for a parse tree produced by TTQLParser.
type BaseTTQLListener struct{}

var _ TTQLListener = &BaseTTQLListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseTTQLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseTTQLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseTTQLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseTTQLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterTtQL is called when production ttQL is entered.
func (s *BaseTTQLListener) EnterTtQL(ctx *TtQLContext) {}

// ExitTtQL is called when production ttQL is exited.
func (s *BaseTTQLListener) ExitTtQL(ctx *TtQLContext) {}

// EnterTtQL_Statement is called when production ttQL_Statement is entered.
func (s *BaseTTQLListener) EnterTtQL_Statement(ctx *TtQL_StatementContext) {}

// ExitTtQL_Statement is called when production ttQL_Statement is exited.
func (s *BaseTTQLListener) ExitTtQL_Statement(ctx *TtQL_StatementContext) {}

// EnterTtQL_Query is called when production ttQL_Query is entered.
func (s *BaseTTQLListener) EnterTtQL_Query(ctx *TtQL_QueryContext) {}

// ExitTtQL_Query is called when production ttQL_Query is exited.
func (s *BaseTTQLListener) ExitTtQL_Query(ctx *TtQL_QueryContext) {}

// EnterTtQL_TimeClause is called when production ttQL_TimeClause is entered.
func (s *BaseTTQLListener) EnterTtQL_TimeClause(ctx *TtQL_TimeClauseContext) {}

// ExitTtQL_TimeClause is called when production ttQL_TimeClause is exited.
func (s *BaseTTQLListener) ExitTtQL_TimeClause(ctx *TtQL_TimeClauseContext) {}

// EnterOC_Cypher is called when production oC_Cypher is entered.
func (s *BaseTTQLListener) EnterOC_Cypher(ctx *OC_CypherContext) {}

// ExitOC_Cypher is called when production oC_Cypher is exited.
func (s *BaseTTQLListener) ExitOC_Cypher(ctx *OC_CypherContext) {}

// EnterOC_Statement is called when production oC_Statement is entered.
func (s *BaseTTQLListener) EnterOC_Statement(ctx *OC_StatementContext) {}

// ExitOC_Statement is called when production oC_Statement is exited.
func (s *BaseTTQLListener) ExitOC_Statement(ctx *OC_StatementContext) {}

// EnterOC_Query is called when production oC_Query is entered.
func (s *BaseTTQLListener) EnterOC_Query(ctx *OC_QueryContext) {}

// ExitOC_Query is called when production oC_Query is exited.
func (s *BaseTTQLListener) ExitOC_Query(ctx *OC_QueryContext) {}

// EnterOC_RegularQuery is called when production oC_RegularQuery is entered.
func (s *BaseTTQLListener) EnterOC_RegularQuery(ctx *OC_RegularQueryContext) {}

// ExitOC_RegularQuery is called when production oC_RegularQuery is exited.
func (s *BaseTTQLListener) ExitOC_RegularQuery(ctx *OC_RegularQueryContext) {}

// EnterOC_Union is called when production oC_Union is entered.
func (s *BaseTTQLListener) EnterOC_Union(ctx *OC_UnionContext) {}

// ExitOC_Union is called when production oC_Union is exited.
func (s *BaseTTQLListener) ExitOC_Union(ctx *OC_UnionContext) {}

// EnterOC_SingleQuery is called when production oC_SingleQuery is entered.
func (s *BaseTTQLListener) EnterOC_SingleQuery(ctx *OC_SingleQueryContext) {}

// ExitOC_SingleQuery is called when production oC_SingleQuery is exited.
func (s *BaseTTQLListener) ExitOC_SingleQuery(ctx *OC_SingleQueryContext) {}

// EnterOC_SinglePartQuery is called when production oC_SinglePartQuery is entered.
func (s *BaseTTQLListener) EnterOC_SinglePartQuery(ctx *OC_SinglePartQueryContext) {}

// ExitOC_SinglePartQuery is called when production oC_SinglePartQuery is exited.
func (s *BaseTTQLListener) ExitOC_SinglePartQuery(ctx *OC_SinglePartQueryContext) {}

// EnterOC_MultiPartQuery is called when production oC_MultiPartQuery is entered.
func (s *BaseTTQLListener) EnterOC_MultiPartQuery(ctx *OC_MultiPartQueryContext) {}

// ExitOC_MultiPartQuery is called when production oC_MultiPartQuery is exited.
func (s *BaseTTQLListener) ExitOC_MultiPartQuery(ctx *OC_MultiPartQueryContext) {}

// EnterOC_UpdatingClause is called when production oC_UpdatingClause is entered.
func (s *BaseTTQLListener) EnterOC_UpdatingClause(ctx *OC_UpdatingClauseContext) {}

// ExitOC_UpdatingClause is called when production oC_UpdatingClause is exited.
func (s *BaseTTQLListener) ExitOC_UpdatingClause(ctx *OC_UpdatingClauseContext) {}

// EnterOC_ReadingClause is called when production oC_ReadingClause is entered.
func (s *BaseTTQLListener) EnterOC_ReadingClause(ctx *OC_ReadingClauseContext) {}

// ExitOC_ReadingClause is called when production oC_ReadingClause is exited.
func (s *BaseTTQLListener) ExitOC_ReadingClause(ctx *OC_ReadingClauseContext) {}

// EnterOC_Match is called when production oC_Match is entered.
func (s *BaseTTQLListener) EnterOC_Match(ctx *OC_MatchContext) {}

// ExitOC_Match is called when production oC_Match is exited.
func (s *BaseTTQLListener) ExitOC_Match(ctx *OC_MatchContext) {}

// EnterOC_Unwind is called when production oC_Unwind is entered.
func (s *BaseTTQLListener) EnterOC_Unwind(ctx *OC_UnwindContext) {}

// ExitOC_Unwind is called when production oC_Unwind is exited.
func (s *BaseTTQLListener) ExitOC_Unwind(ctx *OC_UnwindContext) {}

// EnterOC_Merge is called when production oC_Merge is entered.
func (s *BaseTTQLListener) EnterOC_Merge(ctx *OC_MergeContext) {}

// ExitOC_Merge is called when production oC_Merge is exited.
func (s *BaseTTQLListener) ExitOC_Merge(ctx *OC_MergeContext) {}

// EnterOC_MergeAction is called when production oC_MergeAction is entered.
func (s *BaseTTQLListener) EnterOC_MergeAction(ctx *OC_MergeActionContext) {}

// ExitOC_MergeAction is called when production oC_MergeAction is exited.
func (s *BaseTTQLListener) ExitOC_MergeAction(ctx *OC_MergeActionContext) {}

// EnterOC_Create is called when production oC_Create is entered.
func (s *BaseTTQLListener) EnterOC_Create(ctx *OC_CreateContext) {}

// ExitOC_Create is called when production oC_Create is exited.
func (s *BaseTTQLListener) ExitOC_Create(ctx *OC_CreateContext) {}

// EnterOC_Set is called when production oC_Set is entered.
func (s *BaseTTQLListener) EnterOC_Set(ctx *OC_SetContext) {}

// ExitOC_Set is called when production oC_Set is exited.
func (s *BaseTTQLListener) ExitOC_Set(ctx *OC_SetContext) {}

// EnterOC_SetItem is called when production oC_SetItem is entered.
func (s *BaseTTQLListener) EnterOC_SetItem(ctx *OC_SetItemContext) {}

// ExitOC_SetItem is called when production oC_SetItem is exited.
func (s *BaseTTQLListener) ExitOC_SetItem(ctx *OC_SetItemContext) {}

// EnterOC_Delete is called when production oC_Delete is entered.
func (s *BaseTTQLListener) EnterOC_Delete(ctx *OC_DeleteContext) {}

// ExitOC_Delete is called when production oC_Delete is exited.
func (s *BaseTTQLListener) ExitOC_Delete(ctx *OC_DeleteContext) {}

// EnterOC_Remove is called when production oC_Remove is entered.
func (s *BaseTTQLListener) EnterOC_Remove(ctx *OC_RemoveContext) {}

// ExitOC_Remove is called when production oC_Remove is exited.
func (s *BaseTTQLListener) ExitOC_Remove(ctx *OC_RemoveContext) {}

// EnterOC_RemoveItem is called when production oC_RemoveItem is entered.
func (s *BaseTTQLListener) EnterOC_RemoveItem(ctx *OC_RemoveItemContext) {}

// ExitOC_RemoveItem is called when production oC_RemoveItem is exited.
func (s *BaseTTQLListener) ExitOC_RemoveItem(ctx *OC_RemoveItemContext) {}

// EnterOC_InQueryCall is called when production oC_InQueryCall is entered.
func (s *BaseTTQLListener) EnterOC_InQueryCall(ctx *OC_InQueryCallContext) {}

// ExitOC_InQueryCall is called when production oC_InQueryCall is exited.
func (s *BaseTTQLListener) ExitOC_InQueryCall(ctx *OC_InQueryCallContext) {}

// EnterOC_StandaloneCall is called when production oC_StandaloneCall is entered.
func (s *BaseTTQLListener) EnterOC_StandaloneCall(ctx *OC_StandaloneCallContext) {}

// ExitOC_StandaloneCall is called when production oC_StandaloneCall is exited.
func (s *BaseTTQLListener) ExitOC_StandaloneCall(ctx *OC_StandaloneCallContext) {}

// EnterOC_YieldItems is called when production oC_YieldItems is entered.
func (s *BaseTTQLListener) EnterOC_YieldItems(ctx *OC_YieldItemsContext) {}

// ExitOC_YieldItems is called when production oC_YieldItems is exited.
func (s *BaseTTQLListener) ExitOC_YieldItems(ctx *OC_YieldItemsContext) {}

// EnterOC_YieldItem is called when production oC_YieldItem is entered.
func (s *BaseTTQLListener) EnterOC_YieldItem(ctx *OC_YieldItemContext) {}

// ExitOC_YieldItem is called when production oC_YieldItem is exited.
func (s *BaseTTQLListener) ExitOC_YieldItem(ctx *OC_YieldItemContext) {}

// EnterOC_With is called when production oC_With is entered.
func (s *BaseTTQLListener) EnterOC_With(ctx *OC_WithContext) {}

// ExitOC_With is called when production oC_With is exited.
func (s *BaseTTQLListener) ExitOC_With(ctx *OC_WithContext) {}

// EnterOC_Return is called when production oC_Return is entered.
func (s *BaseTTQLListener) EnterOC_Return(ctx *OC_ReturnContext) {}

// ExitOC_Return is called when production oC_Return is exited.
func (s *BaseTTQLListener) ExitOC_Return(ctx *OC_ReturnContext) {}

// EnterOC_ProjectionBody is called when production oC_ProjectionBody is entered.
func (s *BaseTTQLListener) EnterOC_ProjectionBody(ctx *OC_ProjectionBodyContext) {}

// ExitOC_ProjectionBody is called when production oC_ProjectionBody is exited.
func (s *BaseTTQLListener) ExitOC_ProjectionBody(ctx *OC_ProjectionBodyContext) {}

// EnterOC_ProjectionItems is called when production oC_ProjectionItems is entered.
func (s *BaseTTQLListener) EnterOC_ProjectionItems(ctx *OC_ProjectionItemsContext) {}

// ExitOC_ProjectionItems is called when production oC_ProjectionItems is exited.
func (s *BaseTTQLListener) ExitOC_ProjectionItems(ctx *OC_ProjectionItemsContext) {}

// EnterOC_ProjectionItem is called when production oC_ProjectionItem is entered.
func (s *BaseTTQLListener) EnterOC_ProjectionItem(ctx *OC_ProjectionItemContext) {}

// ExitOC_ProjectionItem is called when production oC_ProjectionItem is exited.
func (s *BaseTTQLListener) ExitOC_ProjectionItem(ctx *OC_ProjectionItemContext) {}

// EnterOC_Order is called when production oC_Order is entered.
func (s *BaseTTQLListener) EnterOC_Order(ctx *OC_OrderContext) {}

// ExitOC_Order is called when production oC_Order is exited.
func (s *BaseTTQLListener) ExitOC_Order(ctx *OC_OrderContext) {}

// EnterOC_Skip is called when production oC_Skip is entered.
func (s *BaseTTQLListener) EnterOC_Skip(ctx *OC_SkipContext) {}

// ExitOC_Skip is called when production oC_Skip is exited.
func (s *BaseTTQLListener) ExitOC_Skip(ctx *OC_SkipContext) {}

// EnterOC_Limit is called when production oC_Limit is entered.
func (s *BaseTTQLListener) EnterOC_Limit(ctx *OC_LimitContext) {}

// ExitOC_Limit is called when production oC_Limit is exited.
func (s *BaseTTQLListener) ExitOC_Limit(ctx *OC_LimitContext) {}

// EnterOC_SortItem is called when production oC_SortItem is entered.
func (s *BaseTTQLListener) EnterOC_SortItem(ctx *OC_SortItemContext) {}

// ExitOC_SortItem is called when production oC_SortItem is exited.
func (s *BaseTTQLListener) ExitOC_SortItem(ctx *OC_SortItemContext) {}

// EnterOC_Where is called when production oC_Where is entered.
func (s *BaseTTQLListener) EnterOC_Where(ctx *OC_WhereContext) {}

// ExitOC_Where is called when production oC_Where is exited.
func (s *BaseTTQLListener) ExitOC_Where(ctx *OC_WhereContext) {}

// EnterOC_Pattern is called when production oC_Pattern is entered.
func (s *BaseTTQLListener) EnterOC_Pattern(ctx *OC_PatternContext) {}

// ExitOC_Pattern is called when production oC_Pattern is exited.
func (s *BaseTTQLListener) ExitOC_Pattern(ctx *OC_PatternContext) {}

// EnterOC_PatternPart is called when production oC_PatternPart is entered.
func (s *BaseTTQLListener) EnterOC_PatternPart(ctx *OC_PatternPartContext) {}

// ExitOC_PatternPart is called when production oC_PatternPart is exited.
func (s *BaseTTQLListener) ExitOC_PatternPart(ctx *OC_PatternPartContext) {}

// EnterOC_AnonymousPatternPart is called when production oC_AnonymousPatternPart is entered.
func (s *BaseTTQLListener) EnterOC_AnonymousPatternPart(ctx *OC_AnonymousPatternPartContext) {}

// ExitOC_AnonymousPatternPart is called when production oC_AnonymousPatternPart is exited.
func (s *BaseTTQLListener) ExitOC_AnonymousPatternPart(ctx *OC_AnonymousPatternPartContext) {}

// EnterOC_PatternElement is called when production oC_PatternElement is entered.
func (s *BaseTTQLListener) EnterOC_PatternElement(ctx *OC_PatternElementContext) {}

// ExitOC_PatternElement is called when production oC_PatternElement is exited.
func (s *BaseTTQLListener) ExitOC_PatternElement(ctx *OC_PatternElementContext) {}

// EnterOC_RelationshipsPattern is called when production oC_RelationshipsPattern is entered.
func (s *BaseTTQLListener) EnterOC_RelationshipsPattern(ctx *OC_RelationshipsPatternContext) {}

// ExitOC_RelationshipsPattern is called when production oC_RelationshipsPattern is exited.
func (s *BaseTTQLListener) ExitOC_RelationshipsPattern(ctx *OC_RelationshipsPatternContext) {}

// EnterOC_NodePattern is called when production oC_NodePattern is entered.
func (s *BaseTTQLListener) EnterOC_NodePattern(ctx *OC_NodePatternContext) {}

// ExitOC_NodePattern is called when production oC_NodePattern is exited.
func (s *BaseTTQLListener) ExitOC_NodePattern(ctx *OC_NodePatternContext) {}

// EnterOC_PatternElementChain is called when production oC_PatternElementChain is entered.
func (s *BaseTTQLListener) EnterOC_PatternElementChain(ctx *OC_PatternElementChainContext) {}

// ExitOC_PatternElementChain is called when production oC_PatternElementChain is exited.
func (s *BaseTTQLListener) ExitOC_PatternElementChain(ctx *OC_PatternElementChainContext) {}

// EnterOC_RelationshipPattern is called when production oC_RelationshipPattern is entered.
func (s *BaseTTQLListener) EnterOC_RelationshipPattern(ctx *OC_RelationshipPatternContext) {}

// ExitOC_RelationshipPattern is called when production oC_RelationshipPattern is exited.
func (s *BaseTTQLListener) ExitOC_RelationshipPattern(ctx *OC_RelationshipPatternContext) {}

// EnterOC_RelationshipDetail is called when production oC_RelationshipDetail is entered.
func (s *BaseTTQLListener) EnterOC_RelationshipDetail(ctx *OC_RelationshipDetailContext) {}

// ExitOC_RelationshipDetail is called when production oC_RelationshipDetail is exited.
func (s *BaseTTQLListener) ExitOC_RelationshipDetail(ctx *OC_RelationshipDetailContext) {}

// EnterOC_Properties is called when production oC_Properties is entered.
func (s *BaseTTQLListener) EnterOC_Properties(ctx *OC_PropertiesContext) {}

// ExitOC_Properties is called when production oC_Properties is exited.
func (s *BaseTTQLListener) ExitOC_Properties(ctx *OC_PropertiesContext) {}

// EnterOC_RelationshipTypes is called when production oC_RelationshipTypes is entered.
func (s *BaseTTQLListener) EnterOC_RelationshipTypes(ctx *OC_RelationshipTypesContext) {}

// ExitOC_RelationshipTypes is called when production oC_RelationshipTypes is exited.
func (s *BaseTTQLListener) ExitOC_RelationshipTypes(ctx *OC_RelationshipTypesContext) {}

// EnterOC_NodeLabels is called when production oC_NodeLabels is entered.
func (s *BaseTTQLListener) EnterOC_NodeLabels(ctx *OC_NodeLabelsContext) {}

// ExitOC_NodeLabels is called when production oC_NodeLabels is exited.
func (s *BaseTTQLListener) ExitOC_NodeLabels(ctx *OC_NodeLabelsContext) {}

// EnterOC_NodeLabel is called when production oC_NodeLabel is entered.
func (s *BaseTTQLListener) EnterOC_NodeLabel(ctx *OC_NodeLabelContext) {}

// ExitOC_NodeLabel is called when production oC_NodeLabel is exited.
func (s *BaseTTQLListener) ExitOC_NodeLabel(ctx *OC_NodeLabelContext) {}

// EnterOC_RangeLiteral is called when production oC_RangeLiteral is entered.
func (s *BaseTTQLListener) EnterOC_RangeLiteral(ctx *OC_RangeLiteralContext) {}

// ExitOC_RangeLiteral is called when production oC_RangeLiteral is exited.
func (s *BaseTTQLListener) ExitOC_RangeLiteral(ctx *OC_RangeLiteralContext) {}

// EnterOC_LabelName is called when production oC_LabelName is entered.
func (s *BaseTTQLListener) EnterOC_LabelName(ctx *OC_LabelNameContext) {}

// ExitOC_LabelName is called when production oC_LabelName is exited.
func (s *BaseTTQLListener) ExitOC_LabelName(ctx *OC_LabelNameContext) {}

// EnterOC_RelTypeName is called when production oC_RelTypeName is entered.
func (s *BaseTTQLListener) EnterOC_RelTypeName(ctx *OC_RelTypeNameContext) {}

// ExitOC_RelTypeName is called when production oC_RelTypeName is exited.
func (s *BaseTTQLListener) ExitOC_RelTypeName(ctx *OC_RelTypeNameContext) {}

// EnterOC_PropertyExpression is called when production oC_PropertyExpression is entered.
func (s *BaseTTQLListener) EnterOC_PropertyExpression(ctx *OC_PropertyExpressionContext) {}

// ExitOC_PropertyExpression is called when production oC_PropertyExpression is exited.
func (s *BaseTTQLListener) ExitOC_PropertyExpression(ctx *OC_PropertyExpressionContext) {}

// EnterOC_Expression is called when production oC_Expression is entered.
func (s *BaseTTQLListener) EnterOC_Expression(ctx *OC_ExpressionContext) {}

// ExitOC_Expression is called when production oC_Expression is exited.
func (s *BaseTTQLListener) ExitOC_Expression(ctx *OC_ExpressionContext) {}

// EnterOC_OrExpression is called when production oC_OrExpression is entered.
func (s *BaseTTQLListener) EnterOC_OrExpression(ctx *OC_OrExpressionContext) {}

// ExitOC_OrExpression is called when production oC_OrExpression is exited.
func (s *BaseTTQLListener) ExitOC_OrExpression(ctx *OC_OrExpressionContext) {}

// EnterOC_XorExpression is called when production oC_XorExpression is entered.
func (s *BaseTTQLListener) EnterOC_XorExpression(ctx *OC_XorExpressionContext) {}

// ExitOC_XorExpression is called when production oC_XorExpression is exited.
func (s *BaseTTQLListener) ExitOC_XorExpression(ctx *OC_XorExpressionContext) {}

// EnterOC_AndExpression is called when production oC_AndExpression is entered.
func (s *BaseTTQLListener) EnterOC_AndExpression(ctx *OC_AndExpressionContext) {}

// ExitOC_AndExpression is called when production oC_AndExpression is exited.
func (s *BaseTTQLListener) ExitOC_AndExpression(ctx *OC_AndExpressionContext) {}

// EnterOC_NotExpression is called when production oC_NotExpression is entered.
func (s *BaseTTQLListener) EnterOC_NotExpression(ctx *OC_NotExpressionContext) {}

// ExitOC_NotExpression is called when production oC_NotExpression is exited.
func (s *BaseTTQLListener) ExitOC_NotExpression(ctx *OC_NotExpressionContext) {}

// EnterOC_ComparisonExpression is called when production oC_ComparisonExpression is entered.
func (s *BaseTTQLListener) EnterOC_ComparisonExpression(ctx *OC_ComparisonExpressionContext) {}

// ExitOC_ComparisonExpression is called when production oC_ComparisonExpression is exited.
func (s *BaseTTQLListener) ExitOC_ComparisonExpression(ctx *OC_ComparisonExpressionContext) {}

// EnterOC_PartialComparisonExpression is called when production oC_PartialComparisonExpression is entered.
func (s *BaseTTQLListener) EnterOC_PartialComparisonExpression(ctx *OC_PartialComparisonExpressionContext) {
}

// ExitOC_PartialComparisonExpression is called when production oC_PartialComparisonExpression is exited.
func (s *BaseTTQLListener) ExitOC_PartialComparisonExpression(ctx *OC_PartialComparisonExpressionContext) {
}

// EnterOC_StringListNullPredicateExpression is called when production oC_StringListNullPredicateExpression is entered.
func (s *BaseTTQLListener) EnterOC_StringListNullPredicateExpression(ctx *OC_StringListNullPredicateExpressionContext) {
}

// ExitOC_StringListNullPredicateExpression is called when production oC_StringListNullPredicateExpression is exited.
func (s *BaseTTQLListener) ExitOC_StringListNullPredicateExpression(ctx *OC_StringListNullPredicateExpressionContext) {
}

// EnterOC_StringPredicateExpression is called when production oC_StringPredicateExpression is entered.
func (s *BaseTTQLListener) EnterOC_StringPredicateExpression(ctx *OC_StringPredicateExpressionContext) {
}

// ExitOC_StringPredicateExpression is called when production oC_StringPredicateExpression is exited.
func (s *BaseTTQLListener) ExitOC_StringPredicateExpression(ctx *OC_StringPredicateExpressionContext) {
}

// EnterOC_ListPredicateExpression is called when production oC_ListPredicateExpression is entered.
func (s *BaseTTQLListener) EnterOC_ListPredicateExpression(ctx *OC_ListPredicateExpressionContext) {}

// ExitOC_ListPredicateExpression is called when production oC_ListPredicateExpression is exited.
func (s *BaseTTQLListener) ExitOC_ListPredicateExpression(ctx *OC_ListPredicateExpressionContext) {}

// EnterOC_NullPredicateExpression is called when production oC_NullPredicateExpression is entered.
func (s *BaseTTQLListener) EnterOC_NullPredicateExpression(ctx *OC_NullPredicateExpressionContext) {}

// ExitOC_NullPredicateExpression is called when production oC_NullPredicateExpression is exited.
func (s *BaseTTQLListener) ExitOC_NullPredicateExpression(ctx *OC_NullPredicateExpressionContext) {}

// EnterOC_AddOrSubtractExpression is called when production oC_AddOrSubtractExpression is entered.
func (s *BaseTTQLListener) EnterOC_AddOrSubtractExpression(ctx *OC_AddOrSubtractExpressionContext) {}

// ExitOC_AddOrSubtractExpression is called when production oC_AddOrSubtractExpression is exited.
func (s *BaseTTQLListener) ExitOC_AddOrSubtractExpression(ctx *OC_AddOrSubtractExpressionContext) {}

// EnterOC_MultiplyDivideModuloExpression is called when production oC_MultiplyDivideModuloExpression is entered.
func (s *BaseTTQLListener) EnterOC_MultiplyDivideModuloExpression(ctx *OC_MultiplyDivideModuloExpressionContext) {
}

// ExitOC_MultiplyDivideModuloExpression is called when production oC_MultiplyDivideModuloExpression is exited.
func (s *BaseTTQLListener) ExitOC_MultiplyDivideModuloExpression(ctx *OC_MultiplyDivideModuloExpressionContext) {
}

// EnterOC_PowerOfExpression is called when production oC_PowerOfExpression is entered.
func (s *BaseTTQLListener) EnterOC_PowerOfExpression(ctx *OC_PowerOfExpressionContext) {}

// ExitOC_PowerOfExpression is called when production oC_PowerOfExpression is exited.
func (s *BaseTTQLListener) ExitOC_PowerOfExpression(ctx *OC_PowerOfExpressionContext) {}

// EnterOC_UnaryAddOrSubtractExpression is called when production oC_UnaryAddOrSubtractExpression is entered.
func (s *BaseTTQLListener) EnterOC_UnaryAddOrSubtractExpression(ctx *OC_UnaryAddOrSubtractExpressionContext) {
}

// ExitOC_UnaryAddOrSubtractExpression is called when production oC_UnaryAddOrSubtractExpression is exited.
func (s *BaseTTQLListener) ExitOC_UnaryAddOrSubtractExpression(ctx *OC_UnaryAddOrSubtractExpressionContext) {
}

// EnterOC_ListOperatorExpression is called when production oC_ListOperatorExpression is entered.
func (s *BaseTTQLListener) EnterOC_ListOperatorExpression(ctx *OC_ListOperatorExpressionContext) {}

// ExitOC_ListOperatorExpression is called when production oC_ListOperatorExpression is exited.
func (s *BaseTTQLListener) ExitOC_ListOperatorExpression(ctx *OC_ListOperatorExpressionContext) {}

// EnterOC_PropertyOrLabelsExpression is called when production oC_PropertyOrLabelsExpression is entered.
func (s *BaseTTQLListener) EnterOC_PropertyOrLabelsExpression(ctx *OC_PropertyOrLabelsExpressionContext) {
}

// ExitOC_PropertyOrLabelsExpression is called when production oC_PropertyOrLabelsExpression is exited.
func (s *BaseTTQLListener) ExitOC_PropertyOrLabelsExpression(ctx *OC_PropertyOrLabelsExpressionContext) {
}

// EnterOC_PropertyLookup is called when production oC_PropertyLookup is entered.
func (s *BaseTTQLListener) EnterOC_PropertyLookup(ctx *OC_PropertyLookupContext) {}

// ExitOC_PropertyLookup is called when production oC_PropertyLookup is exited.
func (s *BaseTTQLListener) ExitOC_PropertyLookup(ctx *OC_PropertyLookupContext) {}

// EnterOC_Atom is called when production oC_Atom is entered.
func (s *BaseTTQLListener) EnterOC_Atom(ctx *OC_AtomContext) {}

// ExitOC_Atom is called when production oC_Atom is exited.
func (s *BaseTTQLListener) ExitOC_Atom(ctx *OC_AtomContext) {}

// EnterOC_CaseExpression is called when production oC_CaseExpression is entered.
func (s *BaseTTQLListener) EnterOC_CaseExpression(ctx *OC_CaseExpressionContext) {}

// ExitOC_CaseExpression is called when production oC_CaseExpression is exited.
func (s *BaseTTQLListener) ExitOC_CaseExpression(ctx *OC_CaseExpressionContext) {}

// EnterOC_CaseAlternative is called when production oC_CaseAlternative is entered.
func (s *BaseTTQLListener) EnterOC_CaseAlternative(ctx *OC_CaseAlternativeContext) {}

// ExitOC_CaseAlternative is called when production oC_CaseAlternative is exited.
func (s *BaseTTQLListener) ExitOC_CaseAlternative(ctx *OC_CaseAlternativeContext) {}

// EnterOC_ListComprehension is called when production oC_ListComprehension is entered.
func (s *BaseTTQLListener) EnterOC_ListComprehension(ctx *OC_ListComprehensionContext) {}

// ExitOC_ListComprehension is called when production oC_ListComprehension is exited.
func (s *BaseTTQLListener) ExitOC_ListComprehension(ctx *OC_ListComprehensionContext) {}

// EnterOC_PatternComprehension is called when production oC_PatternComprehension is entered.
func (s *BaseTTQLListener) EnterOC_PatternComprehension(ctx *OC_PatternComprehensionContext) {}

// ExitOC_PatternComprehension is called when production oC_PatternComprehension is exited.
func (s *BaseTTQLListener) ExitOC_PatternComprehension(ctx *OC_PatternComprehensionContext) {}

// EnterOC_Quantifier is called when production oC_Quantifier is entered.
func (s *BaseTTQLListener) EnterOC_Quantifier(ctx *OC_QuantifierContext) {}

// ExitOC_Quantifier is called when production oC_Quantifier is exited.
func (s *BaseTTQLListener) ExitOC_Quantifier(ctx *OC_QuantifierContext) {}

// EnterOC_FilterExpression is called when production oC_FilterExpression is entered.
func (s *BaseTTQLListener) EnterOC_FilterExpression(ctx *OC_FilterExpressionContext) {}

// ExitOC_FilterExpression is called when production oC_FilterExpression is exited.
func (s *BaseTTQLListener) ExitOC_FilterExpression(ctx *OC_FilterExpressionContext) {}

// EnterOC_PatternPredicate is called when production oC_PatternPredicate is entered.
func (s *BaseTTQLListener) EnterOC_PatternPredicate(ctx *OC_PatternPredicateContext) {}

// ExitOC_PatternPredicate is called when production oC_PatternPredicate is exited.
func (s *BaseTTQLListener) ExitOC_PatternPredicate(ctx *OC_PatternPredicateContext) {}

// EnterOC_ParenthesizedExpression is called when production oC_ParenthesizedExpression is entered.
func (s *BaseTTQLListener) EnterOC_ParenthesizedExpression(ctx *OC_ParenthesizedExpressionContext) {}

// ExitOC_ParenthesizedExpression is called when production oC_ParenthesizedExpression is exited.
func (s *BaseTTQLListener) ExitOC_ParenthesizedExpression(ctx *OC_ParenthesizedExpressionContext) {}

// EnterOC_IdInColl is called when production oC_IdInColl is entered.
func (s *BaseTTQLListener) EnterOC_IdInColl(ctx *OC_IdInCollContext) {}

// ExitOC_IdInColl is called when production oC_IdInColl is exited.
func (s *BaseTTQLListener) ExitOC_IdInColl(ctx *OC_IdInCollContext) {}

// EnterOC_FunctionInvocation is called when production oC_FunctionInvocation is entered.
func (s *BaseTTQLListener) EnterOC_FunctionInvocation(ctx *OC_FunctionInvocationContext) {}

// ExitOC_FunctionInvocation is called when production oC_FunctionInvocation is exited.
func (s *BaseTTQLListener) ExitOC_FunctionInvocation(ctx *OC_FunctionInvocationContext) {}

// EnterOC_FunctionName is called when production oC_FunctionName is entered.
func (s *BaseTTQLListener) EnterOC_FunctionName(ctx *OC_FunctionNameContext) {}

// ExitOC_FunctionName is called when production oC_FunctionName is exited.
func (s *BaseTTQLListener) ExitOC_FunctionName(ctx *OC_FunctionNameContext) {}

// EnterOC_ExistentialSubquery is called when production oC_ExistentialSubquery is entered.
func (s *BaseTTQLListener) EnterOC_ExistentialSubquery(ctx *OC_ExistentialSubqueryContext) {}

// ExitOC_ExistentialSubquery is called when production oC_ExistentialSubquery is exited.
func (s *BaseTTQLListener) ExitOC_ExistentialSubquery(ctx *OC_ExistentialSubqueryContext) {}

// EnterOC_ExplicitProcedureInvocation is called when production oC_ExplicitProcedureInvocation is entered.
func (s *BaseTTQLListener) EnterOC_ExplicitProcedureInvocation(ctx *OC_ExplicitProcedureInvocationContext) {
}

// ExitOC_ExplicitProcedureInvocation is called when production oC_ExplicitProcedureInvocation is exited.
func (s *BaseTTQLListener) ExitOC_ExplicitProcedureInvocation(ctx *OC_ExplicitProcedureInvocationContext) {
}

// EnterOC_ImplicitProcedureInvocation is called when production oC_ImplicitProcedureInvocation is entered.
func (s *BaseTTQLListener) EnterOC_ImplicitProcedureInvocation(ctx *OC_ImplicitProcedureInvocationContext) {
}

// ExitOC_ImplicitProcedureInvocation is called when production oC_ImplicitProcedureInvocation is exited.
func (s *BaseTTQLListener) ExitOC_ImplicitProcedureInvocation(ctx *OC_ImplicitProcedureInvocationContext) {
}

// EnterOC_ProcedureResultField is called when production oC_ProcedureResultField is entered.
func (s *BaseTTQLListener) EnterOC_ProcedureResultField(ctx *OC_ProcedureResultFieldContext) {}

// ExitOC_ProcedureResultField is called when production oC_ProcedureResultField is exited.
func (s *BaseTTQLListener) ExitOC_ProcedureResultField(ctx *OC_ProcedureResultFieldContext) {}

// EnterOC_ProcedureName is called when production oC_ProcedureName is entered.
func (s *BaseTTQLListener) EnterOC_ProcedureName(ctx *OC_ProcedureNameContext) {}

// ExitOC_ProcedureName is called when production oC_ProcedureName is exited.
func (s *BaseTTQLListener) ExitOC_ProcedureName(ctx *OC_ProcedureNameContext) {}

// EnterOC_Namespace is called when production oC_Namespace is entered.
func (s *BaseTTQLListener) EnterOC_Namespace(ctx *OC_NamespaceContext) {}

// ExitOC_Namespace is called when production oC_Namespace is exited.
func (s *BaseTTQLListener) ExitOC_Namespace(ctx *OC_NamespaceContext) {}

// EnterOC_Variable is called when production oC_Variable is entered.
func (s *BaseTTQLListener) EnterOC_Variable(ctx *OC_VariableContext) {}

// ExitOC_Variable is called when production oC_Variable is exited.
func (s *BaseTTQLListener) ExitOC_Variable(ctx *OC_VariableContext) {}

// EnterOC_Literal is called when production oC_Literal is entered.
func (s *BaseTTQLListener) EnterOC_Literal(ctx *OC_LiteralContext) {}

// ExitOC_Literal is called when production oC_Literal is exited.
func (s *BaseTTQLListener) ExitOC_Literal(ctx *OC_LiteralContext) {}

// EnterOC_BooleanLiteral is called when production oC_BooleanLiteral is entered.
func (s *BaseTTQLListener) EnterOC_BooleanLiteral(ctx *OC_BooleanLiteralContext) {}

// ExitOC_BooleanLiteral is called when production oC_BooleanLiteral is exited.
func (s *BaseTTQLListener) ExitOC_BooleanLiteral(ctx *OC_BooleanLiteralContext) {}

// EnterOC_NumberLiteral is called when production oC_NumberLiteral is entered.
func (s *BaseTTQLListener) EnterOC_NumberLiteral(ctx *OC_NumberLiteralContext) {}

// ExitOC_NumberLiteral is called when production oC_NumberLiteral is exited.
func (s *BaseTTQLListener) ExitOC_NumberLiteral(ctx *OC_NumberLiteralContext) {}

// EnterOC_IntegerLiteral is called when production oC_IntegerLiteral is entered.
func (s *BaseTTQLListener) EnterOC_IntegerLiteral(ctx *OC_IntegerLiteralContext) {}

// ExitOC_IntegerLiteral is called when production oC_IntegerLiteral is exited.
func (s *BaseTTQLListener) ExitOC_IntegerLiteral(ctx *OC_IntegerLiteralContext) {}

// EnterOC_DoubleLiteral is called when production oC_DoubleLiteral is entered.
func (s *BaseTTQLListener) EnterOC_DoubleLiteral(ctx *OC_DoubleLiteralContext) {}

// ExitOC_DoubleLiteral is called when production oC_DoubleLiteral is exited.
func (s *BaseTTQLListener) ExitOC_DoubleLiteral(ctx *OC_DoubleLiteralContext) {}

// EnterOC_ListLiteral is called when production oC_ListLiteral is entered.
func (s *BaseTTQLListener) EnterOC_ListLiteral(ctx *OC_ListLiteralContext) {}

// ExitOC_ListLiteral is called when production oC_ListLiteral is exited.
func (s *BaseTTQLListener) ExitOC_ListLiteral(ctx *OC_ListLiteralContext) {}

// EnterOC_MapLiteral is called when production oC_MapLiteral is entered.
func (s *BaseTTQLListener) EnterOC_MapLiteral(ctx *OC_MapLiteralContext) {}

// ExitOC_MapLiteral is called when production oC_MapLiteral is exited.
func (s *BaseTTQLListener) ExitOC_MapLiteral(ctx *OC_MapLiteralContext) {}

// EnterOC_PropertyKeyName is called when production oC_PropertyKeyName is entered.
func (s *BaseTTQLListener) EnterOC_PropertyKeyName(ctx *OC_PropertyKeyNameContext) {}

// ExitOC_PropertyKeyName is called when production oC_PropertyKeyName is exited.
func (s *BaseTTQLListener) ExitOC_PropertyKeyName(ctx *OC_PropertyKeyNameContext) {}

// EnterOC_Parameter is called when production oC_Parameter is entered.
func (s *BaseTTQLListener) EnterOC_Parameter(ctx *OC_ParameterContext) {}

// ExitOC_Parameter is called when production oC_Parameter is exited.
func (s *BaseTTQLListener) ExitOC_Parameter(ctx *OC_ParameterContext) {}

// EnterOC_SchemaName is called when production oC_SchemaName is entered.
func (s *BaseTTQLListener) EnterOC_SchemaName(ctx *OC_SchemaNameContext) {}

// ExitOC_SchemaName is called when production oC_SchemaName is exited.
func (s *BaseTTQLListener) ExitOC_SchemaName(ctx *OC_SchemaNameContext) {}

// EnterOC_ReservedWord is called when production oC_ReservedWord is entered.
func (s *BaseTTQLListener) EnterOC_ReservedWord(ctx *OC_ReservedWordContext) {}

// ExitOC_ReservedWord is called when production oC_ReservedWord is exited.
func (s *BaseTTQLListener) ExitOC_ReservedWord(ctx *OC_ReservedWordContext) {}

// EnterOC_SymbolicName is called when production oC_SymbolicName is entered.
func (s *BaseTTQLListener) EnterOC_SymbolicName(ctx *OC_SymbolicNameContext) {}

// ExitOC_SymbolicName is called when production oC_SymbolicName is exited.
func (s *BaseTTQLListener) ExitOC_SymbolicName(ctx *OC_SymbolicNameContext) {}

// EnterOC_LeftArrowHead is called when production oC_LeftArrowHead is entered.
func (s *BaseTTQLListener) EnterOC_LeftArrowHead(ctx *OC_LeftArrowHeadContext) {}

// ExitOC_LeftArrowHead is called when production oC_LeftArrowHead is exited.
func (s *BaseTTQLListener) ExitOC_LeftArrowHead(ctx *OC_LeftArrowHeadContext) {}

// EnterOC_RightArrowHead is called when production oC_RightArrowHead is entered.
func (s *BaseTTQLListener) EnterOC_RightArrowHead(ctx *OC_RightArrowHeadContext) {}

// ExitOC_RightArrowHead is called when production oC_RightArrowHead is exited.
func (s *BaseTTQLListener) ExitOC_RightArrowHead(ctx *OC_RightArrowHeadContext) {}

// EnterOC_Dash is called when production oC_Dash is entered.
func (s *BaseTTQLListener) EnterOC_Dash(ctx *OC_DashContext) {}

// ExitOC_Dash is called when production oC_Dash is exited.
func (s *BaseTTQLListener) ExitOC_Dash(ctx *OC_DashContext) {}
