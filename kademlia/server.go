package kademlia

import (
	"fmt"
	"go-dht/protocol"
	"log"
	"math/big"
	"net"
	"net/rpc"
)

type Server struct {
	Node         Node
	rpcServer    *rpc.Server
	DataStore    map[string][]byte
	RoutingTable *RoutingTable
}

func NewServer(IP string, Port int) Server {
	n := NewNode(IP, Port)
	rpcServer := rpc.NewServer()
	err := rpcServer.Register(&n)
	if err != nil {
		log.Fatal(err)
	}
	s := Server{Node: n, rpcServer: rpcServer}
	//s.Node.RoutingTable.Add(s.Node)
	return s
}

func (s Server) Listen() {
	host, port := s.Node.IP, s.Node.Port
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go s.rpcServer.Accept(l)
}

func (s Server) Bootstrap(servers ...Server) {
	for _, server := range servers {
		if server.Node.ID != s.Node.ID {
			fmt.Println(s.Ping(server))
		}
	}
}

func (s Server) UpdateBucket(position int, server Server) {
	//bucket := n.Buckets[position]
	//bucket.Append(node)
	//n.Buckets[position].Append(node)

	s.RoutingTable.Add(server)
}

func (s Server) Ping(other Server) error {
	client, err := s.Contact(other)
	if err != nil {
		return err
	}
	args := CallArgs{
		Caller: s.Node,
		RpcId:  RandNumber(),
	}
	var reply Reply
	err = client.Call("Node.Ping", args, &reply)
	if err != nil {
		return err
	}
	distance := new(big.Int).Xor(s.Node.ID, other.Node.ID)
	bucket := len(distance.Bytes())*8 - distance.BitLen() + 1
	s.Node.UpdateBucket(bucket, other.Node)
	return nil
}

func (s Server) Put(key string, value any) error {
	data, err := protocol.Serialize(value)
	if err != nil {
		return err
	}
	keyHash := GetHash(key)
	fmt.Println(keyHash, data)
	return nil
}

func (s Server) Get(key string) any {
	return nil
}

func (s Server) Store(node Node, key string, data []byte) error {

	args := CallArgs{
		Caller: s.Node,
		RpcId:  RandNumber(),
		Key:    key,
		Data:   data,
	}
	fmt.Println(args)
	return nil
}

func (s Server) FindNode(other Server, key string) ([]Node, error) {
	client, err := s.Contact(other)
	if err != nil {
		return nil, err
	}
	args := CallArgs{
		Caller: s.Node,
		RpcId:  RandNumber(),
		Key:    GetHash(key),
	}
	//xor := new(big.Int).Xor(node.ID, HashToBigInt(keyHash))
	var reply Reply
	err = client.Call("Node.FindNode", args, &reply)
	if err != nil {
		return nil, err
	}
	fmt.Println(args, client)
	return []Node{}, nil
}

func (s Server) FindValue(other Server, key string) ([]Node, error) {
	client, err := s.Contact(other)
	if err != nil {
		return nil, err
	}
	args := CallArgs{
		Caller: s.Node,
		RpcId:  RandNumber(),
		Key:    key,
	}
	keyHash := GetHash(key)
	xor := new(big.Int).Xor(other.Node.ID, HashToBigInt(keyHash))
	fmt.Println(args, xor, client)
	return []Node{}, nil
}

func (s Server) Contact(other Server) (*rpc.Client, error) {
	address := fmt.Sprintf("%s:%d", other.Node.IP, other.Node.Port)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("error contacting node at %s", address)
	}
	return client, nil
}
