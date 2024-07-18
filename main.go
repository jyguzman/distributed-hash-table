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
	//fmt.Println(nodes)
	type Inner struct {
		Three float64
	}
	type Thing struct {
		One   int32
		Two   string
		Three Inner
	}

	thing := Thing{
		One:   50,
		Two:   "jdlsk",
		Three: Inner{3.0},
	}
	bytes, err := bson.Marshal(thing)
	if err != nil {
		panic(err)
	}

	var newThing Thing
	_, err = bson.Unmarshal(bytes, &newThing)
	if err != nil {
		panic(err)
	}
	fmt.Println(newThing)
}
