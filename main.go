package main

import (
	"fmt"
	"go-dht/bson"
	"go-dht/kademlia"
	"math/big"
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
	id := kademlia.HashToBigInt(kademlia.GetHash("localhost:8000"))
	fmt.Println(id.Text(16))
	doc := bson.D{
		{"id", id},
		{"null", nil},
	}
	bytes, err := bson.Marshal(doc)
	if err != nil {
		panic(err)
	}

	m := bson.M{}
	err = bson.Unmarshal(bytes, m)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
	val, ok := new(big.Int).SetString(m["id"].(string), 16)
	if !ok {
		panic("invalid id")
	}
	fmt.Println(val.Text(16))
	//servers, err := initServers(10)
	//err = servers[0].BsonPing(servers[1])
	//if err != nil {
	//	panic(err)
	//}
	//for i := 1; i < len(servers); i++ {
	//	err := servers[0].Ping(servers[i])
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//fmt.Println("Root:")
	//fmt.Println(servers[0].RoutingTable.String())
	//fmt.Println()
	//kNearest := servers[0].RoutingTable.GetNearest(servers[6].Node.ID)
	//fmt.Println("K nearest:")
	//fmt.Println(kNearest)
}
