package kademlia

import (
	"fmt"
	"math/big"
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

func (s Server) PingRpc(ca *CallArgs, reply *Reply) error {
	reply.Recipient = s.Node
	reply.Message = fmt.Sprintf("PONG %v", ca.RpcId)
	reply.Code = 1
	s.UpdateRoutingTable(ca.Caller)
	return nil
}

func (s Server) StoreRpc(ca *CallArgs, reply *Reply) error {
	caller, rpcId := ca.Caller, ca.RpcId
	key, value := ca.Key, ca.Data
	reply.Recipient = s.Node
	fmt.Println(caller, key, value, rpcId)
	return nil
}

func (s Server) FindNode(ca *CallArgs, reply *Reply) error {
	caller, rpcId := ca.Caller, ca.RpcId
	reply.Nodes = s.RoutingTable.GetNearest(ca.Key)
	reply.Recipient = s.Node
	fmt.Println(caller, rpcId)
	return nil
}

func (s Server) FindValue(ca *CallArgs, reply *Reply) error {
	reply.Recipient = s.Node
	return nil
}
