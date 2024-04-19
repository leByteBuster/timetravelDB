package dataadapter

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	databaseapi "github.com/LexaTRex/timetravelDB/database-api"
	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/google/uuid"
)

// TODO: make all structs private

// Define a new type for timeSeriesMap
type TimeSeriesMap struct {
	data  map[uuid.UUID][]map[string]interface{}
	mutex sync.Mutex
}

func (tsm *TimeSeriesMap) Set(key uuid.UUID, value []map[string]interface{}) {
	tsm.mutex.Lock()
	tsm.data[key] = value
	tsm.mutex.Unlock()
}

func NewTimeSeriesMap() *TimeSeriesMap {
	return &TimeSeriesMap{data: map[uuid.UUID][]map[string]interface{}{}}
}

type NodesResult struct {
	nodeData []RelevantNodeData
	err      error
}

type EdgesResult struct {
	edgeData []RelevantEdgeData
	err      error
}

type RelevantEdgeData struct {
	RelevantNodeData
	From string // from node
	To   string // to node
}

type RelevantNodeData struct {
	Start          string   // start time - start of validity interval
	End            string   // end time - end of validity interval
	Label          string   // label TODO: allow multiple labels for nodes
	FlatProperties []string //
}

// extracts all relevant data from nodes
func retrieveDataFromNodes(graph_nodes []map[string]interface{}, timeSeriesMap *TimeSeriesMap, nodesChan chan<- NodesResult) {

	var allNodesData []RelevantNodeData
	allNodesData = []RelevantNodeData{}

	// Loop through the data array
	for _, node := range graph_nodes {

		// TODO: make this a goroutine

		nodeData, err := retrieveDataFromNode(node)

		if err != nil {
			nodesChan <- NodesResult{nodeData: nil, err: err}
			return
		}

		// TODO: is this really necessary ? why do we not just pass egde.properties to generateFlatproperties ?
		//       test and debug this.
		//       I believe there was a reason but it might have been a logical error.
		//       answer: The reason is probably that nodeid / relation id are stored as neo4j properties. So
		//       we remove all data we already retrieved here and the rest is stored as properties.
		//       We probably can change this to simply read out the node id here, remove it from the properties
		//       so properties are only actual data, and add the id to the query string later. this way we could just
		//       pass node.properties to generateFlatProperties instead of passing the node

		// delete the retrieved fields to ease the recursive traversing of the node
		delete(node, "start")
		delete(node, "end")
		delete(node, "labels")

		// generate flat property strings 'property: value'
		flatPropertyStrings, err := generateFlatProperties(node, timeSeriesMap)

		nodeData.FlatProperties = flatPropertyStrings

		allNodesData = append(allNodesData, nodeData)

	}
	nodesChan <- NodesResult{nodeData: allNodesData, err: nil}
}

// extracts all relevant data from edges
func retrieveDataFromEdges(graph_edges []map[string]interface{}, timeSeriesMap *TimeSeriesMap, edgesChan chan<- EdgesResult) {

	var allEdgesData []RelevantEdgeData
	allEdgesData = []RelevantEdgeData{}

	// Loop through the data array
	for _, edge := range graph_edges {

		// TODO: make this a goroutine

		edgeData, err := retrieveDataFromEdge(edge)

		if err != nil {
			edgesChan <- EdgesResult{edgeData: nil, err: err}
			return
		}

		// delete the retrieved fields to ease the recursive traversing of the edge
		// TODO: is this really necessary ? why do we not just pass egde.properties to generateFlatproperties ?
		//       test and debug this.
		//       I believe there was a reason but it might have been a logical error.
		//       answer: The reason is probably that nodeid / relation id are stored as neo4j properties. So
		//       we remove all data we already retrieved here and the rest is stored as properties.
		//       We probably can change this to simply read out the relation id here, remove it from the properties
		//       so properties are only actual data, and add the id to the query string later. this way we could just
		//       pass edge.properties to generateFlatProperties instead of passing the edge
		delete(edge, "start")
		delete(edge, "end")
		delete(edge, "from")
		delete(edge, "to")
		delete(edge, "labels")

		// why do we not delete the relationid field here ? maybe so properties.relationid is set ?

		// generate flat property strings 'property: value'
		flatPropertyStrings, err := generateFlatProperties(edge, timeSeriesMap)

		edgeData.FlatProperties = flatPropertyStrings

		allEdgesData = append(allEdgesData, edgeData)

	}
	edgesChan <- EdgesResult{edgeData: allEdgesData, err: nil}
}

func buildEdgeQueries(allEdgesData []RelevantEdgeData) []string {

	var queries []string

	for _, edgeData := range allEdgesData {

		queryPrefix := `MATCH (a),(b) WHERE a.nodeid = ` + edgeData.From + ` AND b.nodeid = ` + edgeData.To + ` CREATE (a)-[r:` + edgeData.Label + ` {`
		queryTemporalProperties := `start: datetime('` + edgeData.Start + `'), end: datetime('` + edgeData.End + `'),`
		querySuffix := `}]->(b)`

		queryProperties := ""

		for _, fragment := range edgeData.FlatProperties {
			queryProperties += fragment
		}

		// clip last ',' of property list
		lastIndex := strings.LastIndex(queryProperties, ",")
		queryProperties = queryProperties[:lastIndex]

		query := queryPrefix + queryTemporalProperties + queryProperties + querySuffix

		queries = append(queries, query)

		utils.Debugf("Edge Query: %v", query)

	}

	return queries
}

func buildNodeQueries(allNodesData []RelevantNodeData) []string {

	var queries []string

	for _, nodeData := range allNodesData {

		// TODO: node["label"] should be able to contain a list of labels
		queryPrefix := `CREATE (n:` + nodeData.Label + ` { `
		queryTemporalProperties := `start: datetime('` + nodeData.Start + `'), end: datetime('` + nodeData.End + `'),`
		querySuffix := `})`

		queryProperties := ""

		for _, fragment := range nodeData.FlatProperties {
			queryProperties += fragment
		}

		// clip last ',' of property list
		lastIndex := strings.LastIndex(queryProperties, ",")
		queryProperties = queryProperties[:lastIndex]

		query := queryPrefix + queryTemporalProperties + queryProperties + querySuffix

		queries = append(queries, query)

		utils.Debugf("Node Query: %v", query)
	}

	return queries
}

// read out all the relevant fields of an edge beside the properties
func retrieveDataFromEdge(edge map[string]interface{}) (RelevantEdgeData, error) {

	var relevantData RelevantEdgeData

	// temporality
	startStr, ok := edge["start"].(string)
	if !ok {
		return RelevantEdgeData{}, errors.New("Data corruption, start time of edge must be string")
	}

	endStr, ok := edge["end"].(string)
	if !ok {
		return RelevantEdgeData{}, errors.New("Data corruption, end time of edge must be string")
	}

	// relation vector
	from, ok := edge["from"].(float64)
	if !ok {
		return RelevantEdgeData{}, errors.New("Data corruption, 'from' node of edge must be float")
	}
	fromStr := strconv.FormatFloat(from, 'f', -1, 64)

	to, ok := edge["to"].(float64)
	if !ok {
		return RelevantEdgeData{}, errors.New("Data corruption, 'to' node of edge must be float")
	}
	toStr := strconv.FormatFloat(to, 'f', -1, 64)

	// NOTE: in Neo4j only one label is possible per edge. we could store more within properties but then
	//       we would have rewrite functionalities in either ttql (cypher) or neo4j (fork)
	labelStr, ok := edge["label"].(string)
	if !ok {
		return RelevantEdgeData{}, errors.New("Data corruption, labelStr must be string")
	}

	relevantData.Start = startStr
	relevantData.End = endStr
	relevantData.From = fromStr
	relevantData.To = toStr
	relevantData.Label = labelStr

	return relevantData, nil
}

// read out all the relevant fields of an edge beside the properties
func retrieveDataFromNode(node map[string]interface{}) (RelevantNodeData, error) {

	fmt.Println("Complete node: ", node)

	var relevantData RelevantNodeData

	// temporality
	startStr, ok := node["start"].(string)
	if !ok {
		return RelevantNodeData{}, errors.New("Data corruption, start time of node must be string")
	}
	relevantData.Start = startStr

	endStr, ok := node["end"].(string)
	if !ok {
		return RelevantNodeData{}, errors.New("Data corruption, end time of node must be string")
	}
	relevantData.End = endStr

	labels, ok := node["labels"].([]interface{})
	if !ok {
		return RelevantNodeData{}, errors.New("Data corruption, labels of node must be array")
	}
	// NOTE: right now this is an array and we only read the first one
	// TODO: multiple labels should be possible
	labelStr, ok := labels[0].(string)
	if !ok {
		return RelevantNodeData{}, errors.New("Data corruption, labelStr must be a string")
	}

	relevantData.Start = startStr
	relevantData.End = endStr
	relevantData.Label = labelStr

	utils.Debugf("\n Single Node Data inner: %v", relevantData)

	return relevantData, nil
}

// This functions loops through the fields of a generated json graph element and creates flattened property strings while
// extracting time series data and storing them within a map.  The flattened property strings are of the form "p1_p2_p3:
// val". px depeds on the depth of the nesting of properties. val is either a single value for indefinite properties or a
// uuid for time series properties.  in case of a time series the time series is stored within the timeSeriesMap with the
// uuid as key.
func generateFlatProperties(fields map[string]interface{}, timeSeriesMap *TimeSeriesMap) ([]string, error) {

	// TODO: think about making this a [][]string, just collecting all the single property labels from each stage in an
	// array and putting them together after. could be cleaner. BUT consider that also the uuid must be stored. just
	// storing it as last entry of each array would be bad style and error prone (see cases when property value is a
	// value and not another property)

	// this is a list of strings
	queryPropertyFragments := make([]string, 0)

	for key, value := range fields {
		switch propertyValue := value.(type) {

		// if value is a nested property itself - recursive call passing the nested property
		case map[string]interface{}:
			propertyFragments, err := generateFlatProperties(propertyValue, timeSeriesMap)

			if err != nil {
				return nil, err
			}

			// reversely add property key of this level to every flattened property string. fragment is the returned
			// flattened property string from the next depth.
			for _, fragment := range propertyFragments {
				queryPropertyFragments = append(queryPropertyFragments, key+`_`+fragment)
			}

		// TODO: in the future maybe a list is not neccessarily a time series
		// if value is a list its a time series..
		case []interface{}:
			timeSeries := utils.ConvertMaps(propertyValue)

			// use an uuid to store the time series within a map.
			// use the same uuid to assign to the flattened property
			id := uuid.New()
			valueFragmentUuid := key + `: "` + id.String() + `", `
			queryPropertyFragments = append(queryPropertyFragments, valueFragmentUuid)

			timeSeriesMap.Set(id, timeSeries)

		// if value is a string its a simple vaue of an indefinite property..
		case string:
			propertyFragment := key + `: "` + propertyValue + `", `
			queryPropertyFragments = append(queryPropertyFragments, propertyFragment)

		// in any other case its a simple non complex value of an indefinite property as well..
		case interface{}:
			// TODO: use more type assertion options instead of fmt.Sprint(propertyValue)
			// TODO: überlege auch welche Datentypen nicht in Neo4j gespeichert werden können. Diese könnten
			//       nach postgresql ausgelagert werden
			propertyFragment := key + `: ` + fmt.Sprint(propertyValue) + `,  `
			queryPropertyFragments = append(queryPropertyFragments, propertyFragment)
		default:
			// don't delete this. this is a safeguard
			return nil, errors.New("Unexpecte type of element field. This should not happen.")
		}
	}

	return queryPropertyFragments, nil
}

func loadDataNeo4j(graphNodes, graphEdges []map[string]interface{}) (map[uuid.UUID][]map[string]interface{}, error) {

	utils.Debugf("\n Number nodes: %v", len(graphNodes))
	utils.Debugf("\n Number edges: %v", len(graphEdges))

	timeSeriesMap := NewTimeSeriesMap()

	nodesChan := make(chan NodesResult)
	edgesChan := make(chan EdgesResult)

	go retrieveDataFromNodes(graphNodes, timeSeriesMap, nodesChan)
	go retrieveDataFromEdges(graphEdges, timeSeriesMap, edgesChan)

	nodesResult := <-nodesChan
	edgesResult := <-edgesChan

	if nodesResult.err != nil {
		return nil, fmt.Errorf("Error retrieving data from nodes: %v", nodesResult.err)
	}
	if edgesResult.err != nil {
		return nil, fmt.Errorf("Error retrieving data from edges: %v", edgesResult.err)
	}

	utils.Debugf("\n Relevant Node Data: %v", nodesResult.nodeData)
	utils.Debugf("\n Relevant Edge Data: %v", edgesResult.edgeData)

	utils.Debugf("\n time-series: %v", timeSeriesMap)
	utils.Debugf("\n Number time-series entries: %v", len(timeSeriesMap.data))

	nodeQueries := buildNodeQueries(nodesResult.nodeData)
	edgeQueries := buildEdgeQueries(edgesResult.edgeData)

	databaseapi.WriteQueryMultipleNeo4j(context.Background(), nodeQueries, map[string]interface{}{})
	databaseapi.WriteQueryMultipleNeo4j(context.Background(), edgeQueries, map[string]interface{}{})

	utils.Debugf("\n node Queries: %v", nodeQueries)
	utils.Debugf("\n edge Queries: %v", edgeQueries)

	return timeSeriesMap.data, nil
}
