package kademlia

import (
	"fmt"
	"go-dht/bsonrpc"
	"go-dht/pkg/util"
	"log"
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

func (s Server) Buckets() map[string]*KBucket {
	return s.routingTable.BucketPrefixes
}

func (s Server) Bootstrap(bootstrapper Server) {
	closestToBootstrapper, err := s.SendFindNode(s.Node.Id.Text(16), bootstrapper)
	if err != nil {
		log.Printf("could not bootstrap %s against %s %s", s.Node, bootstrapper.Node, err)
		return
	}
	s.updateRoutingTable(closestToBootstrapper...)
	s.Lookup(s.Node.Id)
	// find buckets farher away
	//
}

func (s Server) Refresh() {
	for _, bucket := range s.Buckets() {
		if bucket.shouldBeRefreshed() {
			s.Lookup(bucket.randomNum())
		}
	}
}

func (s Server) Lookup(key *big.Int) []Node {
	if key == nil {
		return nil
	}
	lu := NewLookup(s, key)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	nodes := lu.Execute()
	wg.Done()
	//wg.Wait()
	fmt.Println(nodes)
	for _, node := range nodes {
		fmt.Println(node)
	}
	return nodes
}

func (s Server) updateRoutingTable(node ...Node) {
	var m sync.Mutex
	m.Lock()
	defer m.Unlock()

	for _, n := range node {
		s.routingTable.Add(n)
	}
}

func (s Server) DisplayRoutingTable() {
	fmt.Println(s.routingTable)
}

func (s Server) Put(key string, value any) {
	nodes := s.Lookup(util.HashToBigInt(util.GetHash(key)))
	for _, n := range nodes {
		err := s.sendStore(key, value, n)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s Server) Get(key string) any {
	if s.Has(key) {
		return s.dataStore[key]
	}
	//	nodes := s.Lookup(util.HashToBigInt(util.GetHash(key)))
	//	for _, n := range nodes {
	//		nodes, err := s.SendFindValue(key, val, n)
	//		if err != nil {
	//			return err
	//		}
	//	}
	return nil
}

func (s Server) Has(key string) bool {
	val, ok := s.dataStore[key]
	return ok && val != nil
}
