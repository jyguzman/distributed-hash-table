package main

import (
	"fmt"
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
	servers, err := initServers(10)
	if err != nil {
		panic(err)
	}
	//for _, server := range servers {
	//	fmt.Println(server.Id())
	//}
	fmt.Println("      ", servers[0].Id().Text(2))
	for i := 1; i < len(servers); i++ {
		err = servers[0].SendPing(servers[i])
		if err != nil {
			panic(err)
		}
		xor := servers[0].Node.Xor(servers[i].Node)
		text := xor.Text(2)
		//bytes := xor.Bytes()
		//new_ := new(big.Int).SetBytes(bytes)
		//fmt.Println("new:", new_.Text(2))
		fmt.Println(servers[0].Node.Port, "0 id:", servers[0].Id().Text(2))
		fmt.Println(servers[i].Node.Port, "id  :", servers[i].Id().Text(2))
		fmt.Println(servers[i].Node.Port, "xor :", text)

		fmt.Println(xor.Bit(0), xor.Bit(1), xor.Bit(2), xor.Bit(3), len(text), text[0:5])
		//for j := i + 1; j < len(servers); j++ {
		//	err = servers[i].SendPing(servers[j])
		//	if err != nil {
		//		panic(err)
		//	}
		//}
	}
	servers[0].DisplayRoutingTable()
	//servers[0].Lookup(servers[9].Id())
	//if err != nil {
	//	panic(err)
	//}
}
