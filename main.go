package main

import (
	"fmt"
	"go-dht/bson"
	"go-dht/kademlia"
)

func initServers(n int) ([]kademlia.Server, error) {
	servers := make([]kademlia.Server, n)
	var err error
	for i := 0; i < len(servers); i++ {
		servers[i], err = kademlia.NewServer("localhost", 8000+i)
		if err != nil {
			return nil, err
		}
		servers[i].Listen()
	}
	return servers, nil
}

func main() {
	nodes := make([]kademlia.Node, 8)
	for i := 0; i < len(nodes); i++ {
		nodes[i] = kademlia.NewNode("localhost", 8000+i+1, nil)
	}
	closestMsg := bson.M{
		"id":   nodes[0].ID,
		"host": "localhost",
		"port": int32(nodes[0].Port),
		"nodes": bson.A{
			nodes[1].Tuple(),
			nodes[2].Tuple(),
			nodes[3].Tuple(),
		},
	}
	bytes, err := bson.Marshal(closestMsg)
	if err != nil {
		panic(err)
	}
	fmt.Println(bytes)

	m := bson.M{}
	_, err = bson.Unmarshal(bytes, &m)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
	//servers, err := initServers(2)
	//err = servers[0].SendPing(servers[1])
	//if err != nil {
	//	panic(err)
	//}
	//servers[0].DisplayRoutingTable()
}
