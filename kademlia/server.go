package kademlia

import (
	"fmt"
	"go-dht/bson"
	"log"
	"math/big"
	"net"
	"net/rpc"
)

type Server struct {
	Node         Node
	rpcServer    *rpc.Server
	UdpConn      *net.UDPConn
	DataStore    map[string][]byte
	RoutingTable *RoutingTable
}

func NewServer(host string, port int) Server {
	n := NewNode(host, port)
	rpcServer := rpc.NewServer()
	s := Server{
		Node:         n,
		rpcServer:    rpcServer,
		DataStore:    make(map[string][]byte),
		RoutingTable: NewRoutingTable(n, 3),
	}
	err := rpcServer.Register(&s)
	if err != nil {
		log.Fatal(err)
	}
	s.UpdateRoutingTable(s.Node)
	return s
}

func (s Server) Listen() {
	host, port := s.Node.Host, s.Node.Port
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

func (s Server) UpdateBucketList(position int, node Node) {
	s.Node.Buckets[position].Append(node)
}

func (s Server) UpdateRoutingTable(node Node) {
	s.RoutingTable.Add(node)
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
	err = client.Call("Server.PingRpc", args, &reply)
	if err != nil {
		return err
	}
	s.UpdateRoutingTable(other.Node)
	return nil
}

func (s Server) Put(key string, value any) error {
	data, err := bson.MarshalValue(value)
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
		Key:    HashToBigInt(GetHash(key)),
		Data:   data,
	}
	fmt.Println(args)
	return nil
}

func (s Server) FindNodes(other Server, key string) ([]Node, error) {
	client, err := s.Contact(other)
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

func (s Server) GetValue(other Server, key string) ([]Node, error) {
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
	xor := new(big.Int).Xor(other.Node.ID, HashToBigInt(keyHash))
	fmt.Println(args, xor, client)
	return []Node{}, nil
}

func (s Server) Contact(other Server) (*rpc.Client, error) {
	address := fmt.Sprintf("%s:%d", other.Node.Host, other.Node.Port)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("error contacting node at %s", address)
	}
	return client, nil
}
