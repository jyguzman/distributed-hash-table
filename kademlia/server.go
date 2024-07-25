package kademlia

import (
	"fmt"
	"go-dht/bsonrpc"
	"log"
	"math"
	"math/big"
	"sync"
)

type Server struct {
	Node         Node
	rpcServer    *bsonrpc.Server
	dataStore    map[string]any
	routingTable *RoutingTable
}

func (s Server) Id() *big.Int {
	return s.Node.Id
}

func NewServer(host string, port int32) (Server, error) {
	n := NewNode(host, port, nil)
	s := Server{
		Node:         n,
		dataStore:    make(map[string]any),
		routingTable: NewRoutingTable(n, Options.BucketCapacity),
	}
	s.updateRoutingTable(s.Node)

	bsonRpcServer, err := bsonrpc.NewServer(host, int(port))
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

func (s Server) Prefixes() map[string]*KBucket {
	return s.routingTable.BucketPrefixes
}

func (s Server) Bootstrap(servers ...Server) {
	for _, server := range servers {
		if server.Node.Id != s.Node.Id {
			fmt.Println(s.SendPing(server))
		}
	}
}

func (s Server) Lookup(key *big.Int) {
	kClosestNodes := s.routingTable.GetNearest(key)
	fmt.Println("kClosestNodes", kClosestNodes)
	alpha, numClose := Options.Alpha, len(kClosestNodes)
	limit := int(math.Min(float64(alpha), float64(numClose)))

	wg := sync.WaitGroup{}
	wg.Add(limit)
	var nodes [][]Node
	for i := 0; i < limit; i++ {
		go func() {
			list, err := s.SendFindNode(key, kClosestNodes[i])
			if err != nil {
				log.Println(err)
			}
			//heap := &NodeHeap{Key: key}
			//for _, tuple := range resp {
			//
			//}
			nodes = append(nodes, list)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(nodes)
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
