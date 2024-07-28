package kademlia

import (
	"fmt"
	"go-dht/bsonrpc"
	"math/big"
)

type Args struct {
	Sender Node
	Key    string
	Data   any
}

type Response struct {
	Message string
	Code    uint8
	Nodes   []Node
}

func (s Server) SendPing(other Server) error {
	if s.Id() == other.Id() {
		return nil
	}

	client, err := s.Contact(other)
	if err != nil {
		return err
	}

	args := Args{Sender: s.Node}

	var resp Response
	err = client.Call("Server.Ping", args, &resp)
	if err != nil {
		return err
	}

	s.updateRoutingTable(other.Node)
	fmt.Println("PONG", resp)
	return nil
}

func (s Server) Ping(callArgs Args, response *Response) error {
	sender := callArgs.Sender
	fmt.Printf("PING %s\n", sender)
	s.updateRoutingTable(sender)
	response.Message = s.Node.Id.Text(16)
	response.Code = 1
	return nil
}

func (s Server) SendFindNode(key string, other Node) ([]Node, error) {
	client, err := s.ContactNode(other)
	if err != nil {
		return nil, err
	}

	args := Args{
		Sender: s.Node,
		Key:    key,
	}

	var resp Response
	err = client.Call("Server.FindNode", args, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Nodes, nil
}

func (s Server) FindNode(callArgs Args, response *Response) error {
	key := callArgs.Key
	keyInt, ok := new(big.Int).SetString(key, 16)
	if !ok {
		response.Code = 0
		response.Message = "invalid key: " + key
		return nil
	}
	response.Code = 1
	response.Message = "S"
	response.Nodes = s.routingTable.GetNearest(keyInt)
	return nil
}

func (s Server) SendStore(key string, val any, other Server) error {
	client, err := s.Contact(other)
	if err != nil {
		return err
	}

	args := Args{
		Sender: s.Node,
		Key:    key,
		Data:   val,
	}

	var resp Response
	err = client.Call("Store", args, &resp)
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}

func (s Server) FindValue(callArgs Args, response *Response) error {
	key := callArgs.Key
	keyInt, ok := new(big.Int).SetString(key, 16)
	if !ok {
		response.Code = 0
		response.Message = "invalid key: " + key
		return nil
	}
	response.Message = "S"
	response.Code = 1
	response.Nodes = s.routingTable.GetNearest(keyInt)
	return nil
}

func (s Server) Store(args Args, response *Response) error {
	s.dataStore[args.Key] = args.Data
	response.Code = 1
	response.Message = "S"
	return nil
}

func (s Server) Contact(other Server) (*bsonrpc.Client, error) {
	client, err := bsonrpc.Dial(other.Node.Host, other.Node.Port)
	if err != nil {
		return nil, fmt.Errorf("error contacting (UDP) node at %s", other.Node)
	}

	return client, nil
}

func (s Server) ContactNode(node Node) (*bsonrpc.Client, error) {
	client, err := bsonrpc.Dial(node.Host, node.Port)
	if err != nil {
		return nil, fmt.Errorf("error contacting (UDP) node at %s", node)
	}

	return client, nil
}
