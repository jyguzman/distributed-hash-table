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
	servers, err := initServers(20)
	if err != nil {
		panic(err)
	}
	fmt.Println("      ", servers[2].Node.Prefix(5))
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
		//xor := servers[2].Node.Xor(servers[i].Node)
		//text := xor.Text(2)
		//fmt.Println(servers[2].Node.Port, "2 id:", servers[2].Id())
		//fmt.Println(servers[i].Node.Port, "id  :", servers[i].Node.Prefix(5))
		//fmt.Println(servers[i].Node.Port, "xor :", text)

		//fmt.Println(xor.Bit(0), xor.Bit(1), xor.Bit(2), xor.Bit(3), len(text), text[0:5])
		for j := i + 1; j < len(servers); j++ {
			err = servers[i].SendPing(servers[j])
			if err != nil {
				panic(err)
			}
		}
	}
	servers[10].DisplayRoutingTable()
	////wg := new(sync.WaitGroup)
	////wg.Add(50)
	////
	////for i := 0; i < len(servers); i++ {
	//servers[37].Lookup(servers[10].Id())
	////	wg.Done()
	////}
	////wg.Wait()

	//nodes := make([]kademlia.Node, 10)
	//bucket := kademlia.KBucket{
	//	Owner:    kademlia.Node{},
	//	Capacity: 20,
	//	Head:     nil,
	//	Tail:     nil,
	//	Size:     0,
	//	Prefix:   "",
	//}
	//for i := 0; i < len(nodes); i++ {
	//	nodes[i] = kademlia.NewNode("localhost", 8000+i, nil)
	//}
	//
	//bucket.Add(nodes[0])
	//bucket.Add(nodes[1])
	//bucket.Add(nodes[2])
	//bucket.Add(nodes[3])
	//bucket.Add(nodes[4])
	//bucket.Add(nodes[5])
	//fmt.Println(bucket.String())
}
