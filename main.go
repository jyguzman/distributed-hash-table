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
	//servers, err := initServers(8)
	//if err != nil {
	//	panic(err)
	//}
	//for i := 1; i < len(servers); i++ {
	//	err = servers[0].SendPing(servers[i])
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//nodes, err := servers[1].SendFindNode(servers[1].Node.ID, servers[0])
	//if err != nil {
	//	panic(err)
	//}

	node := kademlia.NewNode("localhost", 8000, nil)
	nodeBytes, err := bson.Marshal(node)
	if err != nil {
		panic(err)
	}
	type Result struct {
		Nodes bson.A
	}

	var newNode kademlia.Node
	_, err = bson.Unmarshal(nodeBytes, &newNode)
	if err != nil {
		panic(err)
	}
	fmt.Println("node res:", newNode)
	//bytesTwo, err := bson.Marshal(thingMap)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(bytes)
	//fmt.Println(bytesTwo, slices.Equal(bytes, bytesTwo))
	//var newFlat First
	//_, err = bson.Unmarshal(bytes, &newFlat)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("finished:", newFlat)

}
