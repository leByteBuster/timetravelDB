// TODO
// -[] die adjacency list stimmt noch nicht ganz {anzahl der edges/multiedges stimmt nicht}
// 	-[] should be fixed but double check
// -[] restrictios for time values einführen? ggf. erst später. Ist erst später wichtig

package datagenerator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type TmpPropVal struct {
	Start time.Time
	End   time.Time
	Value interface{}
}

type PropFeatures struct {
	DataType string
	Quantity uint
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var current_node_id = 0
var current_relation_id = 0
var graph_nodes []map[string]interface{}
var graph_edges []map[string]interface{}

// the index x of the array represents the starting node of and edge
// every map at position x contains all the ending nodes of the edges
// of node x as key with the amount as values (because between two nodes
// there can be multiple edges because we work with a multi graph model)
// example:
// 0: {1: 3, 3: 1} 		// there are three edges of type (0,1) and one edge of type (0,3)
// 1: ...
var adjacency_list []map[int]int

func GenerateData() {
	rand.Seed(time.Now().UnixNano())
	var client_interface map[string]interface{} = map[string]interface{}{
		"IP": PropFeatures{"string", 5},
		// "IP":       PropFeatures{"string", 2},
		"firewall": PropFeatures{"boolean", 30},
		"root":     PropFeatures{"string", 2},
		"Risc":     PropFeatures{"int", 300},
		// "Risc": PropFeatures{"int", 3},
		"components": map[string]interface{}{
			"gpu": PropFeatures{"string", 1},
			"cpu": PropFeatures{"string", 1},
			"ram": PropFeatures{"string", 2},
		},
	}

	var server_interface map[string]interface{} = map[string]interface{}{
		"IP":       PropFeatures{"string", 1},
		"firewall": PropFeatures{"boolean", 3},
		"root":     PropFeatures{"string", 10},
		// "Risc": PropFeatures{"int", 800},
		"Risc": PropFeatures{"int", 20},
		"components": map[string]interface{}{
			"cpu": PropFeatures{"string", 1},
			"ram": PropFeatures{"string", 1},
		},
	}

	var printer_interface map[string]interface{} = map[string]interface{}{
		"IP":   PropFeatures{"string", 1},
		"root": PropFeatures{"string", 1},
		"Risc": PropFeatures{"int", 100},
		//"Risc": PropFeatures{"int", 2},
		"components": map[string]interface{}{
			"wifi": PropFeatures{"string", 1},
		},
	}

	var generic_traffic_interface map[string]interface{} = map[string]interface{}{
		"TCPUDP": PropFeatures{"string", 200},
		//"TCPUDP": PropFeatures{"string", 2},
		"IPv4IPv6": PropFeatures{"string", 200},
		//"IPv4IPv6": PropFeatures{"string", 2},
		"Risc": PropFeatures{"int", 1000},
		//"Risc": PropFeatures{"int", 1},
		"Count": PropFeatures{"int", 10000},
		//"Count": PropFeatures{"int", 1},
	}

	// TODO: parse correctly
	begin, err := time.Parse("2006-01-02 15:04:05.0000000 -0700 MST", "2022-12-22 15:33:13.0000005 +0000 UTC")
	if err != nil {
		fmt.Printf("couldn't parse time: %v", err)
	}
	end, err := time.Parse("2006-01-02 15:04:05.0000000 -0700 MST", "2023-01-12 15:33:13.0000005 +0000 UTC")
	if err != nil {
		fmt.Printf("couldn't parse time: %v", err)
	}

	server_nodes := generatePropertyNodes(5, "Server", server_interface, begin, end)
	printer_nodes := generatePropertyNodes(10, "Server", printer_interface, begin, end)
	client_nodes := generatePropertyNodes(50, "Server", client_interface, begin, end)

	graph_nodes = append(server_nodes, printer_nodes...)
	graph_nodes = append(graph_nodes, client_nodes...)

	// fmt.Printf("Servers: %v\n, Printers: %v\n, Clients: %v\n", server_nodes, printer_nodes, client_nodes)

	// generate 10 "Traffic" relations between server nodes. Generate properties from traffic_property_struct for each
	server_relation_objects := generateIntraRelations(10, "Traffic", server_nodes, generic_traffic_interface)

	// generate 10 "Traffic" relations between server nodes and client_nodes. Generate properties from traffic_property_struct for each
	server_client_relation_objects := generateInterRelations(10, "Traffic", server_nodes, client_nodes, generic_traffic_interface)

	// generate 10 "Traffic" relations between server nodes and printer_nodes. Generate properties from traffic_property_struct for each
	server_printer_relation_objects := generateInterRelations(10, "Traffic", server_nodes, printer_nodes, generic_traffic_interface)

	graph_edges = append(server_relation_objects, server_client_relation_objects...)
	graph_edges = append(graph_edges, server_printer_relation_objects...)

	//fmt.Printf("Server-Server relations: %v\n", server_relation_objects)
	//fmt.Printf("Server-Client relations: %v\n", server_client_relation_objects)
	//fmt.Printf("Server-Printer relations: %v\n", server_printer_relation_objects)
	//fmt.Printf("Adjacency List: %v\n", adjacency_list)

	//exportGraphAsJson(graph_nodes, graph_edges, "")
	exportGraphAsJson(graph_nodes, graph_edges, "data-adapter/")
	exportGraphAsJson(graph_nodes, graph_edges, "data-adapter-neo4j-only/")
}

func exportGraphAsJson(graph_nodes []map[string]interface{}, graph_edges []map[string]interface{}, file_path string) {

	edgeFile, err := os.OpenFile(file_path+"graph_edges.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		fmt.Printf("error open edge file: %v\n", err)
	}
	defer edgeFile.Close()
	encoderEdges := json.NewEncoder(edgeFile)
	encoderEdges.Encode(graph_edges)

	nodeFile, err := os.OpenFile(file_path+"graph_nodes.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		fmt.Printf("error open node file: %v\n", err)
	}
	defer nodeFile.Close()
	encoderNodes := json.NewEncoder(nodeFile)
	encoderNodes.Encode(graph_nodes)

}

func generatePropertyNodes(numberNodes int, nodelabel string, property_fields map[string]interface{}, begin time.Time, end time.Time) []map[string]interface{} {
	nodes := []map[string]interface{}{}
	for i := 0; i < numberNodes; i++ {
		node := make(map[string]interface{})
		node["nodeid"] = current_node_id
		node["label"] = nodelabel
		node["start"] = begin
		node["end"] = end
		properties := make(map[string]interface{})
		for key, value := range property_fields {
			properties = setProperty(key, value, properties, begin, end)
		}
		node["ts"] = properties
		nodes = append(nodes, node)

		// create adjacency entry for the newly generated node (if missing also for nodes before - should not happen)
		for len(adjacency_list) < current_node_id {
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
		// TODO: time frame condition for timestamps (between the time of the nodes)
		random_from_index := rand.Intn(len(nodes))
		from_node := nodes[random_from_index]
		random_to_index := rand.Intn(len(nodes))
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

// This function generates >numberRelations< random relations between the passed from_nodes and to_nodes. It sets reandom values for the passed property_fields.
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

		random_from_index := rand.Intn(len(from_nodes))
		from_node := from_nodes[random_from_index]
		random_to_index := rand.Intn(len(to_nodes))
		to_node := to_nodes[random_to_index]
		node_id_from := from_node["nodeid"].(int)
		node_id_to := to_node["nodeid"].(int)

		begin, end := minTimeBoundaries(from_node, to_node)

		relation["start"] = begin
		relation["end"] = end
		relation["from"] = node_id_from
		relation["to"] = node_id_to

		// add edge to adjacency list
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

// func randomTimestamp() time.Time {
// 	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000
//
// 	randomNow := time.Unix(randomTime, 500)
//
// 	return randomNow
// }

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// assign a random value to a property field depending on the type. Type may be a struct which consists of
// a struct PropFeatures which contains the DataType field which describes the data type in the form of a string
// : "boolean", "string", "integer" and the field Quantity which holds an integer which defines the number of values.
// Or it contains an object containing properties itself. In this case, every field of the object is assigned a
// value recursively.
// TODO:
// if a TmpPropVal is returned it has to be added to the array of values. If a [string]map{interface} is returned it has
// to be the value itself. Either differentiate by the returned value everywhere where setField is called or find another
// solution
func setField(fieldtype interface{}, parentBegin, parentEnd time.Time) (interface{}, []interface{}) {
	if features, ok := fieldtype.(PropFeatures); ok {
		if features.DataType == "string" {
			arr := make([]interface{}, 0)
			var lastTimestampEnd time.Time = parentBegin
			for i := uint(0); i < features.Quantity; i++ {
				begin, end := generateTimeseriesTimestamp(lastTimestampEnd, parentBegin, parentEnd, features.Quantity)
				arr = append(arr, TmpPropVal{
					Start: begin,
					End:   end,
					Value: randSeq(rand.Intn(10-3) + 3),
				})
				lastTimestampEnd = end
			}
			return nil, arr
		} else if features.DataType == "int" {
			arr := make([]interface{}, 0)
			var lastTimestampEnd time.Time = parentBegin
			for i := uint(0); i < features.Quantity; i++ {
				begin, end := generateTimeseriesTimestamp(lastTimestampEnd, parentBegin, parentEnd, features.Quantity)
				arr = append(arr, TmpPropVal{
					Start: begin,
					End:   end,
					Value: rand.Intn(100),
				})
				lastTimestampEnd = end
			}
			return nil, arr
		} else if features.DataType == "boolean" {
			arr := make([]interface{}, 0)
			var lastTimestampEnd time.Time = parentBegin
			for i := uint(0); i < features.Quantity; i++ {
				begin, end := generateTimeseriesTimestamp(lastTimestampEnd, parentBegin, parentEnd, features.Quantity)
				arr = append(arr, TmpPropVal{
					Start: begin,
					End:   end,
					Value: false,
				})
				lastTimestampEnd = end
			}
			return nil, arr
		}
	}

	// only left option for type of fieldtype should be map[string]interface{}
	// so let panic if not the case
	property_val := fieldtype.(map[string]interface{})

	// set nested properties
	var nested_properties = make(map[string]interface{})
	for key, val := range property_val {
		nested_properties = setProperty(key, val, nested_properties, parentBegin, parentEnd)
	}
	return nested_properties, nil
}

func generateTimeseriesTimestamp(lastTimestampEnd, parentBegin, parentEnd time.Time, numberTimestamps uint) (time.Time, time.Time) {
	var start, end time.Time
	if lastTimestampEnd.IsZero() {
		start = parentBegin
	} else {
		start = lastTimestampEnd
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	duration := (parentEnd.Sub(parentBegin)) / time.Duration(numberTimestamps)
	//fmt.Printf("Max duration: %v", duration)
	randDuration := time.Duration(r.Int63n(int64(duration)))
	if randDuration == time.Duration(0) {
		randDuration = time.Duration(int64(time.Millisecond))
	}
	//fmt.Printf("Rand duration: %v", duration)
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
			// if a value for the key exists already, add it to the array
		} else {
			switch x := nested_properties[key].(type) {
			case []interface{}:
				//fmt.Printf("\narr_val: %v\n", array_val)
				//fmt.Printf("\nnested_properties: %v\n", nested_properties)
				nested_properties[key] = append(x, array_val...)
			default:
				err := fmt.Errorf("unexpected type of value")
				panic(err)
			}
		}
	}
	return nested_properties
}
