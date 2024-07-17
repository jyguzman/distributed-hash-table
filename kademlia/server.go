package kademlia

import (
	"fmt"
	"go-dht/bsonrpc"
	"log"
	"net"
	"net/rpc"
)

type Server struct {
	Node          Node
	tcpRpcServer  *rpc.Server
	bsonRpcServer *bsonrpc.Server
	dataStore     map[string]any
	routingTable  *RoutingTable
}

func NewServer(host string, port int) (Server, error) {
	n := NewNode(host, port, nil)
	s := Server{
		Node:         n,
		dataStore:    make(map[string]any),
		routingTable: NewRoutingTable(n, Options.BucketCapacity),
	}
	s.updateRoutingTable(s.Node)

	var tcpRpcServer *rpc.Server
	var bsonRpcServer *bsonrpc.Server

	if Options.Protocol == "tcp" {
		tcpRpcServer = rpc.NewServer()
		err := tcpRpcServer.Register(&s)
		if err != nil {
			return Server{}, err
		}
		s.tcpRpcServer = tcpRpcServer
		return s, nil
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

func (s Server) updateRoutingTable(node Node) {
	s.routingTable.Add(node)
}

func (s Server) DisplayRoutingTable() {
	fmt.Println(s.routingTable)
}

func (s Server) Put(key string, value any) error {
	return nil
}

func (s Server) Get(key string) any {
	return nil
}
