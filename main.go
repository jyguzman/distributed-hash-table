package main

import (
	"fmt"
	"go-dht/bson"
	"go-dht/kademlia"
	"slices"
)

func initServers(n int) ([]kademlia.Server, error) {
	servers := make([]kademlia.Server, n)
	var err error
	for i := 0; i < len(servers); i++ {
		servers[i], err = kademlia.NewServer("localhost", int32(8000+i))
		if err != nil {
			return nil, err
		}
		servers[i].Listen()
	}
	return servers, nil
}

type Test struct {
	One   int32
	Three string
}

func main() {
	t := &Test{100, "hi"}
	tBytes, errOne := bson.Marshal(t)
	if errOne != nil {
		panic(errOne)
	}
	t2 := &Test{200, "bye"}
	t2Bytes, errTwo := bson.Marshal(t2)
	if errTwo != nil {
		panic(errTwo)
	}
	vals := []*Test{t, t2}
	bytes, err := bson.Marshal(vals)
	if err != nil {
		panic(err)
	}
	fmt.Println(tBytes, len(tBytes))
	fmt.Println(t2Bytes, len(t2Bytes))
	fmt.Println(bytes)
	tBytes = append(tBytes, t2Bytes...)
	fmt.Println(slices.Equal(tBytes, bytes))
	//servers, err := initServers(10)
	//if err != nil {
	//	panic(err)
	//}
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
}
