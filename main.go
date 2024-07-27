package main

import (
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
	servers, err := initServers(2)
	if err != nil {
		panic(err)
	}
	err = servers[0].SendPing(servers[1])
	if err != nil {
		panic(err)
	}
	////for _, server := range servers {
	////	fmt.Println(server.Id())
	////}
	////fmt.Println("      ", servers[9].Node.Prefix(1)) //Id().Text(2))
	////for i := 0; i < len(servers); i++ {
	////	fmt.Println(servers[i].Node.Port, "id  :", servers[i].Node.Prefix(5))
	////}
	//for i := 0; i < len(servers); i++ {
	//	//if i != 2 {
	//	//	err = servers[2].SendPing(servers[i])
	//	//	if err != nil {
	//	//		panic(err)
	//	//	}
	//	//	//servers[9].DisplayRoutingTable()
	//	//}
	//	//xor := servers[0].Node.Xor(servers[i].Node)
	//	//text := xor.Text(2)
	//	//fmt.Println(servers[0].Node.Port, "0 id:", servers[0].Id().Text(2))
	//	//fmt.Println(servers[i].Node.Port, "id  :", servers[i].Node.Prefix(5)) //.Id().Text(2))
	//	//fmt.Println(servers[i].Node.Port, "xor :", text)
	//	//
	//	//fmt.Println(xor.Bit(0), xor.Bit(1), xor.Bit(2), xor.Bit(3), len(text), text[0:5])
	//	for j := i + 1; j < len(servers); j++ {
	//		err = servers[i].SendPing(servers[j])
	//		if err != nil {
	//			panic(err)
	//		}
	//	}
	//}
	//servers[5].DisplayRoutingTable()
	////fmt.Println(servers[4].Prefixes())
	//servers[0].Lookup(servers[0].Id())
	////if err != nil {
	////	panic(err)
	////}
	//type Thing struct {
	//	One   string
	//	Two   float64
	//	Three []kademlia.Node
	//}
	//t := Thing{"hello", 10.5, []kademlia.Node{servers[0].Node, servers[1].Node}}
	//nodes := []kademlia.Node{servers[0].Node, servers[1].Node}
	//fmt.Println(nodes)
	//m := bson.M{
	//	"Message": "pong",
	//}
	//mBytes, err := bson.Marshal(m)
	//if err != nil {
	//	panic(err)
	//}
	//
	//type Response struct {
	//	Message string
	//	Nodes   []kademlia.Node
	//}
	//var newM Response
	////newM := bson.M{}
	//err = bson.Unmarshal(mBytes, &newM)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(newM)
}
