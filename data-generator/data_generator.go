// TODO
// -[] die adjacency list stimmt noch nicht ganz {anzahl der edges/multiedges stimmt nicht}
// 	-[] should be fixed but double check
// -[] restrictios for time values einführen? ggf. erst später. Ist erst später wichtig

package datagenerator

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/LexaTRex/timetravelDB/utils"
	"gopkg.in/yaml.v3"
)

type TmpPropVal struct {
	Start time.Time
	End   time.Time
	Value interface{}
}

type TimePeriod struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type NodeTemplate struct {
	Labels           []string       `yaml:"labels"`
	Count            int            `yaml:"count"`
	PropertyTemplate map[string]any `yaml:"template"`
}

type EdgeTemplate struct {
	Label            string         `yaml:"label"`
	Count            int            `yaml:"count"`
	From             string         `yaml:"from"`
	To               string         `yaml:"to"`
	PropertyTemplate map[string]any `yaml:"template"`
}

type GraphTemplate struct {
	TimePeriod TimePeriod     `yaml:"timePeriod"`
	Nodes      []NodeTemplate `yaml:"nodes"`
	Edges      []EdgeTemplate `yaml:"edges"`
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var current_node_id = 0
var current_relation_id = 0

// the index x of the array represents the starting node of and edge
// every map at position x contains all the ending nodes of the edges
// of node x as key with the amount as values (because between two nodes
// there can be multiple edges because we work with a multi graph model)
// example:
// 0: {1: 3, 3: 1} 		// three edges of type (0,1) and one edge of type (0,3)
// 1: ...
var adjacency_list []map[int]int

var R = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateData(template string) {

	templatePath := ""

	if template == "" {
		templatePath = "data-generator/graph_template.yaml"
	} else {
		templatePath = "data-generator/" + template
	}

	graphFile, err := os.ReadFile(templatePath)
	if err != nil {
		log.Printf("error reading edge file: %v", err)
		return
	}

	var graphTemplate GraphTemplate
	err = yaml.Unmarshal(graphFile, &graphTemplate)
	if err != nil {
		log.Printf("error unmarshalling node data: %v", err)
		return
	}

	timePeriod := graphTemplate.TimePeriod
	nodeTemplates := graphTemplate.Nodes
	edgeTemplates := graphTemplate.Edges

	utils.Debugf("%#v\n", nodeTemplates)
	utils.Debugf("%#v\n", edgeTemplates)

	var begin time.Time
	var end time.Time

	utils.Debugf("Time Frame Generated Data: %+v", timePeriod)

	// TODO: parse correctly
	if graphTemplate.TimePeriod.From == "" || graphTemplate.TimePeriod.To == "" {
		begin, err = time.Parse("2006-01-02 15:04:05.0000000 -0700 MST", "2023-01-01 00:00:00.0000000 +0000 UTC")
		if err != nil {
			log.Printf("couldn't parse time: %v", err)
			return
		}
		end, err = time.Parse("2006-01-02 15:04:05.0000000 -0700 MST", "2023-01-02 00:00:00.0000000 +0000 UTC")
		if err != nil {
			log.Printf("couldn't parse time: %v", err)
			return
		}
	} else {
		begin, err = time.Parse("2006-01-02T15:04:05.99999999999999Z", timePeriod.From)
		if err != nil {
			log.Printf("couldn't parse time: %v", err)
			return
		}
		end, err = time.Parse("2006-01-02T15:04:05.99999999999999Z", timePeriod.To)
		if err != nil {
			log.Printf("couldn't parse time: %v", err)
			return
		}
	}

	var graphNodes [][]map[string]interface{}
	var graphNodesRes []map[string]any
	var graphEdgesRes []map[string]any

	for _, node := range nodeTemplates {
		nodes := generatePropertyNodes(node.Count, node.Labels, node.PropertyTemplate, begin, end)
		graphNodes = append(graphNodes, nodes)
	}
	for _, edgeGroup := range edgeTemplates {
		var edges []map[string]any
		if edgeGroup.From == edgeGroup.To {
			for _, nodeGroup := range graphNodes {

				// this is kind of unstable because we are only alloed single labels for the data generator
				labels := nodeGroup[0]["labels"].([]string)
				if labels[0] == edgeGroup.From {
					utils.Debugf("generate intra relations from %v to %v", edgeGroup.From, edgeGroup.To)
					edges = generateIntraRelations(edgeGroup.Count, edgeGroup.Label, nodeGroup, edgeGroup.PropertyTemplate)
				}
			}
		} else {
			fromGroups := make([][]map[string]interface{}, 0)
			toGroups := make([][]map[string]interface{}, 0)
			for _, nodeGroup := range graphNodes {

				// this is kind of unstable because we are only alloed single labels for the data generator
				labels := nodeGroup[0]["labels"].([]string)
				if labels[0] == edgeGroup.From {
					fromGroups = append(fromGroups, nodeGroup)
				}
				if labels[0] == edgeGroup.To {
					toGroups = append(toGroups, nodeGroup)
				}
			}

			for _, fromGroup := range fromGroups {
				for _, toGroup := range toGroups {
					utils.Debugf("generate inter relations from %v to %v", edgeGroup.From, edgeGroup.To)
					edges = generateInterRelations(edgeGroup.Count, edgeGroup.Label, fromGroup, toGroup, edgeGroup.PropertyTemplate)
				}
			}
		}
		graphEdgesRes = append(graphEdgesRes, edges...)
	}

	for _, nodeGroup := range graphNodes {
		graphNodesRes = append(graphNodesRes, nodeGroup...)
	}

	exportGraphAsJson(graphNodesRes, graphEdgesRes, "data-generator/generated-data/")
}

func exportGraphAsJson(graph_nodes []map[string]interface{}, graph_edges []map[string]interface{}, file_path string) {

	err := os.MkdirAll(file_path, 0755)
	if err != nil {
		log.Printf("couldn't generate directory: %v\n", err)
		return
	}

	edgeFile, err := os.OpenFile(file_path+"graph_edges.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		log.Printf("error open edge file: %v\n", err)
		edgeFile.Close()
		return
	}
	defer edgeFile.Close()
	encoderEdges := json.NewEncoder(edgeFile)
	encoderEdges.Encode(graph_edges)

	nodeFile, err := os.OpenFile(file_path+"graph_nodes.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		log.Printf("error open node file: %v\n", err)
		nodeFile.Close()
		return
	}
	defer nodeFile.Close()
	encoderNodes := json.NewEncoder(nodeFile)
	encoderNodes.Encode(graph_nodes)

}

func generatePropertyNodes(numberNodes int, nodelabels []string, property_fields map[string]interface{}, begin time.Time, end time.Time) []map[string]interface{} {
	nodes := []map[string]interface{}{}
	for i := 0; i < numberNodes; i++ {
		node := make(map[string]interface{})
		node["nodeid"] = current_node_id
		node["labels"] = nodelabels
		node["start"] = begin
		node["end"] = end
		properties := make(map[string]interface{})
		for key, value := range property_fields {
			properties = setProperty(key, value, properties, begin, end)
		}
		node["ts"] = properties
		nodes = append(nodes, node)

		// create adjacency entry for the newly generated node (if missing also for nodes before - should not happen)
		for len(adjacency_list) < current_node_id+1 {
			adjacency_list = append(adjacency_list, make(map[int]int))
		}
		current_node_id++
	}
	return nodes
}

// This function generates >numberRelations< random relations between the passed nodes. It sets reandom values for the passed property_fields. It returns
// the relations as a tuple of json-like Objects which store all the data as well as an array of AdjacenctPairs which just represent the edges via node ids
// in the form of (from, to)
func generateIntraRelations(numberRelations int, relation_label string, nodes []map[string]interface{}, property_fields map[string]interface{}) []map[string]interface{} {

	relations := []map[string]interface{}{}

	for i := 0; i < numberRelations; i++ {

		relation := make(map[string]interface{})

		relation["relationid"] = current_relation_id
		relation["label"] = relation_label

		random_from_index := R.Intn(len(nodes))
		from_node := nodes[random_from_index]

		random_to_index := R.Intn(len(nodes))
		to_node := nodes[random_to_index]

		node_id_from := from_node["nodeid"].(int)
		node_id_to := to_node["nodeid"].(int)

		begin, end := minTimeBoundaries(from_node, to_node)

		relation["start"] = begin
		relation["end"] = end

		relation["from"] = node_id_from
		relation["to"] = node_id_to

		// add edge to adjacency list
		adjacency_list[node_id_from][node_id_to] = adjacency_list[node_id_from][node_id_to] + 1

		properties := make(map[string]interface{})
		for key, value := range property_fields {
			properties = setProperty(key, value, properties, begin, end)
		}
		relation["properties"] = properties
		relations = append(relations, relation)
		current_relation_id++
	}
	return relations
}

// returns the biggest of the two start values and the lowest of the two end values
// so a relation can only exist when both nodes exist
func minTimeBoundaries(from_node, to_node map[string]interface{}) (time.Time, time.Time) {
	from_start := from_node["start"].(time.Time)
	to_start := to_node["start"].(time.Time)
	from_end := from_node["end"].(time.Time)
	to_end := to_node["end"].(time.Time)

	var start, end time.Time

	if from_start.After(to_start) {
		start = from_start
	} else {
		start = to_start
	}

	if from_end.Before(to_end) {
		end = from_end
	} else {
		end = to_end
	}

	return start, end
}

// generatInterRelations generates <numberRelations> random relations between the passed from_nodes and to_nodes. It sets reandom values for the passed property_fields.
// It returns the relations as a tuple of json-like Objects which store all the data as well as an array of AdjacenctPairs which just represent the edges via node ids
// in the form of (from, to)
func generateInterRelations(numberRelations int, relation_label string, from_nodes []map[string]interface{}, to_nodes []map[string]interface{}, property_fields map[string]interface{}) []map[string]interface{} {

	// an array of relations of type <relation_label>
	relations := []map[string]interface{}{}

	// generate <numberRelations> relations of type <relation_label>
	for i := 0; i < numberRelations; i++ {
		relation := make(map[string]interface{})

		// set all the obligatory fields
		relation["relationid"] = current_relation_id
		relation["label"] = relation_label

		random_from_index := R.Intn(len(from_nodes))
		from_node := from_nodes[random_from_index]
		random_to_index := R.Intn(len(to_nodes))
		to_node := to_nodes[random_to_index]
		node_id_from := from_node["nodeid"].(int)
		node_id_to := to_node["nodeid"].(int)

		begin, end := minTimeBoundaries(from_node, to_node)

		relation["start"] = begin
		relation["end"] = end
		relation["from"] = node_id_from
		relation["to"] = node_id_to

		// increment edge count between the two nodes in the adjacency list
		adjacency_list[node_id_from][node_id_to] = adjacency_list[node_id_from][node_id_to] + 1

		// Generate all the properties of <property_fields>
		var properties = make(map[string]interface{})
		for key, value := range property_fields {
			properties = setProperty(key, value, properties, begin, end)
		}
		relation["properties"] = properties
		relations = append(relations, relation)
		current_relation_id++
	}
	return relations
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[R.Intn(len(letters))]
	}
	return string(b)
}

// setField assigns a random value to a property field depending on the type. Type may be a struct which consists of
// a struct TemplateField which contains the DataType field which describes the data type in the form of a string
// : "boolean", "string", "integer" and the field Quantity which holds an integer which defines the number of values.
// Or it contains an object containing properties itself. In this case, every field of the object is assigned a
// value recursively.
func setField(fieldtype interface{}, parentBegin, parentEnd time.Time) (interface{}, []interface{}) {

	var nested_properties = make(map[string]interface{})

	if property_val, ok := fieldtype.(map[string]interface{}); ok {

		if property_val["DataType"] != nil && property_val["Quantity"] != nil {
			var lastTimestampEnd time.Time = parentBegin
			var arr = make([]interface{}, 0)
			switch property_val["DataType"].(string) {
			case "string":
				// generate random time-series string
				for i := int(0); i < property_val["Quantity"].(int); i++ {
					begin, end := generateTimeseriesTimestamp(lastTimestampEnd, parentBegin, parentEnd, property_val["Quantity"].(int))
					arr = append(arr, TmpPropVal{
						Start: begin,
						End:   end,
						Value: randSeq(R.Intn(10-3) + 3),
					})
					lastTimestampEnd = end
				}
			case "int":
				// generate random time-series int
				for i := int(0); i < property_val["Quantity"].(int); i++ {
					begin, end := generateTimeseriesTimestamp(lastTimestampEnd, parentBegin, parentEnd, property_val["Quantity"].(int))
					arr = append(arr, TmpPropVal{
						Start: begin,
						End:   end,
						Value: R.Intn(100),
					})
					lastTimestampEnd = end
				}
			case "boolean":
				// generate random time-series bool
				for i := int(0); i < property_val["Quantity"].(int); i++ {
					begin, end := generateTimeseriesTimestamp(lastTimestampEnd, parentBegin, parentEnd, property_val["Quantity"].(int))
					arr = append(arr, TmpPropVal{
						Start: begin,
						End:   end,
						Value: R.Intn(2) == 1,
					})
					lastTimestampEnd = end
				}
			}
			return nil, arr
		}
		for key, val := range property_val {
			// set nested properties
			nested_properties = setProperty(fmt.Sprintf("%v", key), val, nested_properties, parentBegin, parentEnd)
		}
	} else {
		panic("Fieldtype is not of type map[string]interface{}")
	}

	return nested_properties, nil
}

func generateTimeseriesTimestamp(lastTimestampEnd, parentBegin, parentEnd time.Time, numberTimestamps int) (time.Time, time.Time) {
	var start, end time.Time
	if lastTimestampEnd.IsZero() {
		start = parentBegin
	} else {
		start = lastTimestampEnd
	}
	duration := (parentEnd.Sub(parentBegin)) / time.Duration(numberTimestamps)
	randDuration := time.Duration(R.Int63n(int64(duration)))
	if randDuration == time.Duration(0) {
		randDuration = time.Duration(int64(time.Millisecond))
	}
	end = start.Add(randDuration)
	if end.After(parentEnd) {
		end = parentEnd
	}
	return start, end
}

func setProperty(key string, val interface{}, nested_properties map[string]interface{}, begin, end time.Time) map[string]interface{} {
	// nested_val is defined if the value is nested. Array_val is defined if the val holds an array of actual values.
	// This means there is no further level of nesting but The time-series of values has been found.
	// In the case of nested_value, it is set on all nested levels at this point already becauese setField()
	// contains a recursive call on the setProperty() function.
	nested_val, array_val := setField(val, begin, end)
	if nested_val != nil {
		nested_properties[key] = nested_val
	} else {

		// if value for key doesnt exist yet create array and add value
		if nested_properties[key] == nil {
			nested_properties[key] = array_val

			// else if a value for the key exists already, add it to the array
		} else {
			switch x := nested_properties[key].(type) {
			case []interface{}:
				nested_properties[key] = append(x, array_val...)
			default:
				panic("unexpected type of value")
			}
		}
	}
	return nested_properties
}
