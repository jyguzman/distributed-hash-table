package kademlia

import (
	"fmt"
	"go-dht/bson"
	"go-dht/bsonrpc"
	"math/big"
)

func (s Server) SendPing(other Server) error {
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

	node := FromTuple(bson.A{other.Node.Host, other.Node.Port, other.Node.ID})
	s.UpdateRoutingTable(node)
	fmt.Println("PONG", reply)
	return nil
}

func (s Server) SendFindNode(key string, other Server) ([]bson.A, error) {
	client, err := s.Contact(other)
	if err != nil {
		return nil, err
	}

	args := bson.M{
		"q":   "FindNode",
		"id":  s.Node.ID,
		"key": HashToBigInt(GetHash(key)),
	}

	reply := bson.M{}
	err = client.Call(args, reply)
	if err != nil {
		return nil, err
	}

	return reply["nodes"].([]bson.A), nil
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
	id, host, port := callArgs["id"].(string), callArgs["host"].(string), int(callArgs["port"].(int32))
	ID, ok := new(big.Int).SetString(id, 16)
	if !ok {
		return fmt.Errorf("PONG: Invalid ID: %s", id)
	}
	node := FromTuple(bson.A{ID, host, port})
	fmt.Printf("PING %s\n", node)
	s.UpdateRoutingTable(node)
	reply["id"] = s.Node.ID
	return nil
}

func (s Server) FindNode(callArgs bson.M, reply bson.M) error {
	key := callArgs["key"].(*big.Int)
	reply["nodes"] = s.routingTable.GetNearest(key)
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
	client, err := bsonrpc.Dial(other.Node.Host, other.Node.Port)
	if err != nil {
		return nil, fmt.Errorf("error contacting (UDP) node at %s", other.Node.Host)
	}

	return client, nil
}
