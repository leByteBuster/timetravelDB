// TODO
// -[] die adjacency list stimmt noch nicht ganz {anzahl der edges/multiedges stimmt nicht}
// 	-[] should be fixed but double check
// -[] restrictios for time values einführen? ggf. erst später. Ist erst später wichtig

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type TmpPropVal[T any] struct {
	Start string
	End   string
	Value T
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

func main() {
	rand.Seed(time.Now().UnixNano())
	var client_interface map[string]interface{} = map[string]interface{}{
		"IP": "string",
		//"firewall": "boolean",
		//"root":     "string",
		"Risc": "int",
		"components": map[string]interface{}{
			"gpu": "string",
			//"cpu": "string",
			"ram": "string",
		},
	}

	var server_interface map[string]interface{} = map[string]interface{}{
		"IP": "string",
		//"firewall": "boolean",
		//"root":     "string",
		"Risc": "int",
		"components": map[string]interface{}{
			"cpu": "string",
			"ram": "string",
		},
	}

	var printer_interface map[string]interface{} = map[string]interface{}{
		"IP":   "string",
		"root": "string",
		"Risc": "int",
		"components": map[string]interface{}{
			"wifi": "string",
		},
	}

	var generic_traffic_interface map[string]interface{} = map[string]interface{}{
		//"TCP/UDP":   "string",
		//"IPv4/IPv6": "string",
		"Risc":  "int",
		"Count": "int",
	}

	server_nodes := generatePropertyNodes(1, "Server", server_interface)
	printer_nodes := generatePropertyNodes(3, "Server", printer_interface)
	client_nodes := generatePropertyNodes(3, "Server", client_interface)

	graph_nodes = append(server_nodes, printer_nodes...)
	graph_nodes = append(graph_nodes, client_nodes...)

	fmt.Printf("Servers: %v\n, Printers: %v\n, Clients: %v\n", server_nodes, printer_nodes, client_nodes)

	// generate 10 "Traffic" relations between server nodes. Generate properties from traffic_property_struct for each
	server_relation_objects := generateIntraRelations(2, "Traffic", server_nodes, generic_traffic_interface)
	// generate 10 "Traffic" relations between server nodes and client_nodes. Generate properties from traffic_property_struct for each
	server_client_relation_objects := generateInterRelations(3, "Traffic", server_nodes, client_nodes, generic_traffic_interface)
	// generate 10 "Traffic" relations between server nodes and printer_nodes. Generate properties from traffic_property_struct for each
	server_printer_relation_objects := generateInterRelations(3, "Traffic", server_nodes, printer_nodes, generic_traffic_interface)

	graph_edges = append(server_relation_objects, server_client_relation_objects...)
	graph_edges = append(graph_edges, server_printer_relation_objects...)

	//fmt.Printf("Server-Server relations: %v\n", server_relation_objects)
	//fmt.Printf("Server-Client relations: %v\n", server_client_relation_objects)
	//fmt.Printf("Server-Printer relations: %v\n", server_printer_relation_objects)
	//fmt.Printf("Adjacency List: %v\n", adjacency_list)

	exportGraphAsJson(graph_nodes, graph_edges, "")
}

func exportGraphAsJson(graph_nodes []map[string]interface{}, graph_edges []map[string]interface{}, file_path string) {

	//// ### I'm using the method below because it's faster ###

	// node_bytes, err := json.Marshal(graph_nodes)
	// fmt.Printf("Error: %v\n", err)
	// fmt.Printf("Marshalled nodes: %v\n", string(node_bytes))

	// edge_bytes, err := json.Marshal(graph_edges)
	// fmt.Printf("Error: %v\n", err)
	// fmt.Printf("Marshalled edges: %v\n", string(edge_bytes))

	// err = ioutil.WriteFile(file_path+"graph_edges", edge_bytes, 0644)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// }

	// err = ioutil.WriteFile(file_path+"graph_nodes", node_bytes, 0644)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// }

	edgeFile, err := os.OpenFile(file_path+"graph_edges.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	defer edgeFile.Close()
	encoderEdges := json.NewEncoder(edgeFile)
	encoderEdges.Encode(graph_edges)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	nodeFile, err := os.OpenFile(file_path+"graph_nodes.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	defer nodeFile.Close()
	encoderNodes := json.NewEncoder(nodeFile)
	encoderNodes.Encode(graph_nodes)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func generatePropertyNodes(numberNodes int, nodelabel string, property_fields map[string]interface{}) []map[string]interface{} {
	nodes := []map[string]interface{}{}
	for i := 0; i < numberNodes; i++ {
		node := make(map[string]interface{})
		node["nodeid"] = current_node_id
		node["label"] = nodelabel
		node["start"] = randomTimestamp().String()
		node["end"] = randomTimestamp().String() // (in the range 1 - 1000)
		properties := make(map[string]interface{})
		for key, value := range property_fields {
			properties = setProperty(key, value, properties)
		}
		node["properties"] = properties
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
		relation["start"] = randomTimestamp().String()
		relation["end"] = randomTimestamp().String() // (in the range 1 - 1000)
		random_from_index := rand.Intn(len(nodes))
		from_node := nodes[random_from_index]
		random_to_index := rand.Intn(len(nodes))
		to_node := nodes[random_to_index]
		node_id_from := from_node["nodeid"].(int)
		node_id_to := to_node["nodeid"].(int)
		relation["from"] = node_id_from
		relation["to"] = node_id_to

		// add edge to adjacency list
		adjacency_list[node_id_from][node_id_to] = adjacency_list[node_id_from][node_id_to] + 1

		properties := make(map[string]interface{})
		for key, value := range property_fields {
			properties = setProperty(key, value, properties)
		}
		relation["properties"] = properties
		relations = append(relations, relation)
		current_relation_id++
	}
	return relations
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
		// TODO: time frame condition for timestamps (between the time of the nodes)
		relation["start"] = randomTimestamp().String()
		relation["end"] = randomTimestamp().String() // (in the range 1 - 1000)
		random_from_index := rand.Intn(len(from_nodes))
		from_node := from_nodes[random_from_index]
		random_to_index := rand.Intn(len(to_nodes))
		to_node := to_nodes[random_to_index]
		node_id_from := from_node["nodeid"].(int)
		node_id_to := to_node["nodeid"].(int)
		relation["from"] = node_id_from
		relation["to"] = node_id_to

		// add edge to adjacency list
		adjacency_list[node_id_from][node_id_to] = adjacency_list[node_id_from][node_id_to] + 1

		// Generate all the properties of <property_fields>
		var properties = make(map[string]interface{})
		for key, value := range property_fields {
			properties = setProperty(key, value, properties)
		}
		relation["properties"] = properties
		relations = append(relations, relation)
		current_relation_id++
	}
	return relations
}

func randomTimestamp() time.Time {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000

	randomNow := time.Unix(randomTime, 0)

	return randomNow
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// assign a random value to a property field depending on the type. Type may be string integer or an object
// containing properties itself. In this case, every field of the object is assigned a value recursively.
// TODO:
// if a TmpPropVal is returned it hast to be added to the array of values. If a [string]map{interface} is returned it has
// to be the value itself. Either differentiate by the returned value everywhere where setField is called or find another
// solution
func setField(fieldtype any) (any, any) {
	if str, ok := fieldtype.(string); ok {
		if str == "string" {
			return nil, TmpPropVal[string]{
				Start: randomTimestamp().String(),
				End:   randomTimestamp().String(),
				Value: randSeq(rand.Intn(10-3) + 3),
			}
		}
		if str == "int" {
			return nil, TmpPropVal[int]{
				Start: randomTimestamp().String(),
				End:   randomTimestamp().String(),
				Value: rand.Intn(100),
			}
		}
		if str == "boolean" {
			return nil, TmpPropVal[bool]{
				Start: randomTimestamp().String(),
				End:   randomTimestamp().String(),
				Value: false,
			}
		}
	}

	// only left option for type of fieldtype should be map[string]interface{}
	// so let panic if not the case
	property_val := fieldtype.(map[string]interface{})

	// set nested properties
	var nested_properties = make(map[string]interface{})
	for key, val := range property_val {
		nested_properties = setProperty(key, val, nested_properties)
	}
	return nested_properties, nil
}

func setProperty(key string, val any, nested_properties map[string]interface{}) map[string]interface{} {
	nested_val, array_val := setField(val)
	if nested_val != nil {
		nested_properties[key] = nested_val
	} else {
		if nested_properties[key] == nil {
			nested_properties[key] = [...]any{array_val}
		} else {
			switch x := nested_properties[key].(type) {
			case []interface{}:
				nested_properties[key] = append(x, array_val)
			}
		}
	}
	return nested_properties
}

func exportToJson() {

}

//		// fill property_fields struct with random values and time values
//		// time values need to be smaller than parent time vaues
//		// merge the structure in the form of
//		{
//		  nodeid:
//			start:
//			end:
//			nodelabel:
//			properties:
//				{
//					property_name_1: {
// 						values: [(startval,endval,value),
//              (startval,endval,{
//                 property_name_1_1: {
// 										values: [(startval,endval,value), .., (startval,endval,value)]
//                 }
// 							})
//							(startval,endval,value), (startval,endval,value),(startval,endval,value)]
// 				  }
//					property_name_2: {
// 						values: [(startval,endval,value), (startval,endval,value), (startval,endval,value),(startval,endval,value)]
// 					 }
//					property_name_3: {
// 						values: [(startval,endval,value), (startval,endval,value), (startval,endval,value),(startval,endval,value)]
//				}
//  	}
//
