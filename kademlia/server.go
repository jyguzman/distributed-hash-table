package kademlia

import (
	"fmt"
	"go-dht/bson"
	"go-dht/bsonrpc"
	"log"
	"net"
	"net/rpc"
)

type Server struct {
	Node          Node
	tcpRpcServer  *rpc.Server
	bsonRpcServer *bsonrpc.Server
	dataStore     map[string][]byte
	routingTable  *RoutingTable
}

func NewServer(host string, port int) (Server, error) {
	n := NewNode(host, port, nil)
	s := Server{
		Node:         n,
		dataStore:    make(map[string][]byte),
		routingTable: NewRoutingTable(n, 8),
	}
	tcpRpcServer := rpc.NewServer()
	err := tcpRpcServer.Register(&s)
	if err != nil {
		return Server{}, err
	}
	bsonRpcServer, err := bsonrpc.NewServer(host, port)
	if err != nil {
		return Server{}, err
	}
	err = bsonRpcServer.Register(&s)
	if err != nil {
		return Server{}, err
	}
	s.bsonRpcServer = bsonRpcServer
	s.tcpRpcServer = tcpRpcServer
	s.UpdateRoutingTable(s.Node)
	return s, nil
}

func (s Server) Listen() {
	host, port := s.Node.Host, s.Node.Port
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go s.tcpRpcServer.Accept(l)
	go s.bsonRpcServer.Listen()
}

func (s Server) Bootstrap(servers ...Server) {
	for _, server := range servers {
		if server.Node.ID != s.Node.ID {
			fmt.Println(s.SendPing(server))
		}
	}
}

func (s Server) UpdateRoutingTable(node Node) {
	s.routingTable.Add(node)
}

func (s Server) DisplayRoutingTable() {
	fmt.Println(s.routingTable)
}

func (s Server) Put(key string, value any) error {
	_, data, err := bson.MarshalValue(value)
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
