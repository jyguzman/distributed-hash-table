package main

import (
	"fmt"
	"go-dht/kademlia"
	"reflect"
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
	servers, err := initServers(10)
	if err != nil {
		panic(err)
	}
	for _, server := range servers {
		fmt.Println(server.Id())
	}
	fmt.Println("      ", servers[2].Node.Prefix(1))
	for i := 0; i < len(servers); i++ {
		fmt.Println(servers[i].Node.Port, "id  :", servers[i].Node.Prefix(5))
	}
	for i := 0; i < len(servers); i++ {
		//if i != 2 {
		//	err = servers[2].SendPing(servers[i])
		//	if err != nil {
		//		panic(err)
		//	}
		//}
		xor := servers[0].Node.Xor(servers[i].Node)
		text := xor.Text(2)
		fmt.Println(servers[0].Node.Port, "0 id:", servers[0].Id().Text(2))
		fmt.Println(servers[i].Node.Port, "id  :", servers[i].Node.Prefix(5))
		fmt.Println(servers[i].Node.Port, "xor :", text)

		//fmt.Println(xor.Bit(0), xor.Bit(1), xor.Bit(2), xor.Bit(3), len(text), text[0:5])
		for j := i + 1; j < len(servers); j++ {
			err = servers[i].SendPing(servers[j])
			if err != nil {
				panic(err)
			}
		}
	}
	servers[2].DisplayRoutingTable()
	servers[0].Lookup(servers[0].Id())

	var k kademlia.Key
	err = k.FromString("790f5ea8d47ea50879916d5835ecb1e0709c18fc", 16)
	if err != nil {
		panic(err)
	}
	fmt.Println(k, reflect.TypeOf(k))
}
