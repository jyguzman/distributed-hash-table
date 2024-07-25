package kademlia

import (
	"fmt"
	"math/big"
	"net/rpc"
)

type CallArgs struct {
	Caller Node
	Key    *big.Int
	Data   []byte
	RpcId  *big.Int
}

type Reply struct {
	Recipient Node
	Message   string
	Value     []byte
	Nodes     []Node
	Code      int
}

// These are the server's functions

func (s Server) PingRpc(ca *CallArgs, reply *Reply) error {
	reply.Recipient = s.Node
	reply.Message = fmt.Sprintf("PONG %v", ca.RpcId)
	reply.Code = 1
	s.updateRoutingTable(ca.Caller)
	return nil
}

func (s Server) StoreRpc(ca *CallArgs, reply *Reply) error {
	caller, rpcId := ca.Caller, ca.RpcId
	key, value := ca.Key, ca.Data
	reply.Recipient = s.Node
	fmt.Println(caller, key, value, rpcId)
	return nil
}

func (s Server) FindNodeRpc(ca *CallArgs, reply *Reply) error {
	caller, rpcId := ca.Caller, ca.RpcId
	reply.Nodes = s.routingTable.GetNearest(ca.Key)
	reply.Recipient = s.Node
	fmt.Println(caller, rpcId)
	return nil
}

func (s Server) FindValueRpc(ca *CallArgs, reply *Reply) error {
	reply.Recipient = s.Node
	return nil
}

// Client calls these functions to request from server

func (s Server) PingTcp(other Server) error {
	client, err := s.ContactTcp(other)
	if err != nil {
		return err
	}
	args := CallArgs{
		Caller: s.Node,
		RpcId:  RandNumber(),
	}
	var reply Reply
	err = client.Call("Server.PingRpc", args, &reply)
	if err != nil {
		return err
	}
	s.updateRoutingTable(other.Node)
	return nil
}

func (s Server) StoreTcp(node Node, key string, data []byte) error {
	args := CallArgs{
		Caller: s.Node,
		RpcId:  RandNumber(),
		Key:    HashToBigInt(GetHash(key)),
		Data:   data,
	}
	fmt.Println(args)
	return nil
}

func (s Server) FindNodesTcp(other Server, key string) ([]Node, error) {
	client, err := s.ContactTcp(other)
	if err != nil {
		return nil, err
	}
	args := CallArgs{
		Caller: s.Node,
		RpcId:  RandNumber(),
		Key:    HashToBigInt(GetHash(key)),
	}
	var reply Reply
	err = client.Call("Server.FindNode", args, &reply)
	if err != nil {
		return nil, err
	}
	fmt.Println(args, client)
	return reply.Nodes, nil
}

func (s Server) FindValueTcp(other Server, key string) ([]Node, error) {
	client, err := s.Contact(other)
	if err != nil {
		return nil, err
	}
	args := CallArgs{
		Caller: s.Node,
		RpcId:  RandNumber(),
		Key:    HashToBigInt(GetHash(key)),
	}
	keyHash := GetHash(key)
	xor := new(big.Int).Xor(other.Node.Id, HashToBigInt(keyHash))
	fmt.Println(args, xor, client)
	return []Node{}, nil
}

func (s Server) ContactTcp(other Server) (*rpc.Client, error) {
	address := fmt.Sprintf("%s:%d", other.Node.Host, other.Node.Port)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("error contacting node at %s", address)
	}
	return client, nil
}
