package kademlia

import (
	"fmt"
	"math/big"
)

type CallArgs struct {
	Caller Node
	Key    string
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

func (n Node) Ping(ca *CallArgs, reply *Reply) error {
	reply.Recipient = n
	reply.Message = fmt.Sprintf("PONG %v", ca.RpcId)
	reply.Code = 1
	distance := new(big.Int).Xor(n.ID, ca.Caller.ID)
	bucket := len(distance.Bytes())*8 - distance.BitLen() + 1
	// fmt.Printf("Bucket: %d\n", bucket)
	n.UpdateBucket(bucket, ca.Caller)
	return nil
}

func (n Node) Store(ca *CallArgs, reply *Reply) error {
	caller, rpcId := ca.Caller, ca.RpcId
	key, value := ca.Key, ca.Data
	reply.Recipient = n
	fmt.Println(caller, key, value, rpcId)
	return nil
}

func (n Node) FindNode(ca *CallArgs, reply *Reply) error {
	caller, hashedKey, rpcId := ca.Caller, ca.Key, ca.RpcId
	distance := new(big.Int).Xor(caller.ID, HashToBigInt(hashedKey))
	reply.Recipient = n
	fmt.Println(caller, hashedKey, rpcId, distance)
	return nil
}

func (n Node) FindValue(ca *CallArgs, reply *Reply) error {
	reply.Recipient = n
	return nil
}
