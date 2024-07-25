package main

import (
	"fmt"
	"go-dht/bson"
	"go-dht/kademlia"
	"reflect"
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

type InnerInnerTest struct {
	Eight string
	Nine  int64
	Ten   []float64
}

type TestInner struct {
	Five  int64
	Six   bool
	Seven InnerInnerTest
}

type Test struct {
	One   int32
	Two   float64
	Three string
	Four  TestInner
}

func main() {
	type Holder struct {
		Nodes []kademlia.Node
	}
	n := kademlia.NewNode("localhost", 8000, nil)
	n2 := kademlia.NewNode("localhost", 8001, nil)
	h := Holder{[]kademlia.Node{n, n2}}
	hBytes, err := bson.Marshal(h)
	if err != nil {
		panic(err)
	}
	var newH Holder
	err = bson.Unmarshal(hBytes, &newH)
	if err != nil {
		panic(err)
	}
	fmt.Println(newH.Nodes[0], reflect.TypeOf(newH.Nodes[0]))
	tii := InnerInnerTest{
		Eight: "whoa",
		Nine:  109,
		Ten:   []float64{0.5, 4.6},
	}
	ti := TestInner{
		Five:  50,
		Six:   true,
		Seven: tii,
	}
	t := Test{
		One:   10,
		Two:   50.5,
		Three: "hello",
		Four:  ti,
	}
	tBytes, err := bson.Marshal(t)
	if err != nil {
		panic(err)
	}
	var newT Test
	err = bson.Unmarshal(tBytes, &newT)
	if err != nil {
		panic(err)
	}
	fmt.Println(newT)
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
