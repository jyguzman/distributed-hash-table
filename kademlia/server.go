package kademlia

import (
	"fmt"
	"go-dht/bsonrpc"
)

type Server struct {
	Node         Node
	rpcServer    *bsonrpc.Server
	dataStore    map[string]any
	routingTable *RoutingTable
}

func NewServer(host string, port int) (Server, error) {
	n := NewNode(host, port, nil)
	s := Server{
		Node:         n,
		dataStore:    make(map[string]any),
		routingTable: NewRoutingTable(n, Options.BucketCapacity),
	}
	s.updateRoutingTable(s.Node)

	bsonRpcServer, err := bsonrpc.NewServer(host, port)
	if err != nil {
		return Server{}, err
	}
	err = bsonRpcServer.Register(&s)
	if err != nil {
		return Server{}, err
	}
	s.rpcServer = bsonRpcServer

	return s, nil
}

func (s Server) RoutingTable() *RoutingTable {
	return s.routingTable
}

func (s Server) Listen() {
	go s.rpcServer.Listen()
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
