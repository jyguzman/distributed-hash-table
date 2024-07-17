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
	//id := kademlia.HashToBigInt(kademlia.GetHash("localhost:8000"))
	//fmt.Println(id.Text(16))
	//doc := bson.M{
	//	"hello": "world",
	//	"there": bson.M{
	//		"inner": "world",
	//		"inner_two": bson.M{
	//			"inner": "world_two",
	//		},
	//	},
	//}
	docD := bson.D{
		{"hello", "world"},
		{"there", bson.D{
			{"inner", "world"},
			{"inner_two", bson.D{
				{"inner", "world_two"},
			}},
		}},
	}
	//docT := bson.D{
	//	{"hello", bson.D{
	//		{"thing", "you"},
	//	}},
	//}
	//docE := bson.D{
	//	{"hello", "world"},
	//	{"there", bson.D{
	//		{"inner", "world"},
	//		{"inner_two", bson.D{
	//			{"inner", "world_two"},
	//			{"inner_three", bson.D{
	//				{"inner", "world_two"},
	//			}},
	//		}},
	//	}},
	//}
	bytes, err := bson.Marshal(docD)
	if err != nil {
		panic(err)
	}
	fmt.Println(bytes)

	m := bson.D{}
	err = bson.Unmarshal(bytes, &m)
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
