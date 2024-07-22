package kademlia

import (
	"fmt"
	"go-dht/bson"
	"go-dht/bsonrpc"
	"math/big"
)

func (s Server) SendPing(other Server) error {
	if s.Id() == other.Id() {
		return nil
	}
	client, err := s.Contact(other)
	if err != nil {
		return err
	}

	args := bson.M{
		"q":    "Ping",
		"id":   s.Node.ID,
		"host": s.Node.Host,
		"port": int32(s.Node.Port),
	}

	reply := bson.M{}
	err = client.Call(args, reply)
	if err != nil {
		return err
	}

	node := NodeFromTuple(bson.A{other.Node.Host, other.Node.Port, other.Node.ID})
	s.updateRoutingTable(node)
	//fmt.Println("PONG", reply)
	return nil
}

func (s Server) SendFindNode(key *big.Int, other Node) ([]Node, error) {
	client, err := s.ContactNode(other)
	if err != nil {
		return nil, err
	}

	args := bson.M{
		"q":    "FindNode",
		"id":   s.Node.ID,
		"host": s.Node.Host,
		"port": s.Node.Port,
		"key":  key,
	}

	reply := bson.M{}
	err = client.Call(args, reply)
	if err != nil {
		return nil, err
	}

	tuples := reply["nodes"].(bson.M)
	var nodes []Node
	for _, tuple := range tuples {
		node, err := NodeFromMap(tuple.(bson.M))
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (s Server) SendStore(key string, val any, other Server) error {
	client, err := s.Contact(other)
	if err != nil {
		return err
	}

	args := bson.M{
		"q":   "Store",
		"id":  s.Node.ID,
		"key": key,
		"val": val,
	}

	reply := bson.M{}
	err = client.Call(args, reply)
	if err != nil {
		return err
	}

	return nil
}

func (s Server) Ping(callArgs bson.M, reply bson.M) error {
	id, host, port := callArgs["id"].(string), callArgs["host"].(string), callArgs["port"].(int32)
	ID, ok := new(big.Int).SetString(id, 16)
	if !ok {
		return fmt.Errorf("PONG: Invalid ID: %s", id)
	}
	node := NodeFromTuple(bson.A{host, port, ID})
	//fmt.Printf("PING %s\n", node)
	s.updateRoutingTable(node)
	reply["id"] = s.Node.ID
	return nil
}

func (s Server) FindNode(callArgs bson.M, reply bson.M) error {
	key := callArgs["key"].(string)
	intKey, ok := new(big.Int).SetString(key, 16)
	if !ok {
		return fmt.Errorf("invalid key string: %s", key)
	}
	nodes := s.routingTable.GetNearest(intKey)
	tuples := bson.A{}
	//fmt.Println(nodes)
	for _, node := range nodes {
		tuples = append(tuples, node.Tuple())
	}
	reply["nodes"] = tuples
	return nil
}

func (s Server) FindValue(callArgs bson.M, reply bson.M) error {
	key := callArgs["key"].(*big.Int)
	reply["nodes"] = s.routingTable.GetNearest(key)
	return nil
}

func (s Server) Store(callArgs bson.M, reply bson.M) error {
	s.dataStore[callArgs["key"].(string)] = callArgs["val"]
	return nil
}

func (s Server) Contact(other Server) (*bsonrpc.Client, error) {
	client, err := bsonrpc.Dial(other.Node.Host, int(other.Node.Port))
	if err != nil {
		return nil, fmt.Errorf("error contacting (UDP) node at %s", other.Node)
	}

	return client, nil
}

func (s Server) ContactNode(node Node) (*bsonrpc.Client, error) {
	client, err := bsonrpc.Dial(node.Host, int(node.Port))
	if err != nil {
		return nil, fmt.Errorf("error contacting (UDP) node at %s", node)
	}

	return client, nil
}
