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

type WhoaAgain struct {
	Twelve   string
	Thirteen bool
}

type InnerInnerTest struct {
	Eight  string
	Nine   int
	Ten    []float64
	Eleven []WhoaAgain
}

type TestInner struct {
	Five  uint32
	Six   bool
	Seven InnerInnerTest
}

type Test struct {
	One   int8
	Two   float64
	Three string
	Four  TestInner
}

func main() {
	//type Holder struct {
	//	Nodes []kademlia.Node
	//}
	//n := kademlia.NewNode("localhost", 8000, nil)
	//n2 := kademlia.NewNode("localhost", 8001, nil)
	//h := Holder{[]kademlia.Node{n, n2}}
	//hBytes, err := bson.Marshal(h)
	//if err != nil {
	//	panic(err)
	//}
	//var newH Holder
	//err = bson.Unmarshal(hBytes, &newH)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(newH)
	//fmt.Println(newH.Nodes[0], reflect.TypeOf(newH.Nodes[0]))
	//whoa1 := WhoaAgain{
	//	Twelve: "big bong", Thirteen: true,
	//}
	//whoa2 := WhoaAgain{
	//	Twelve: "small bong", Thirteen: false,
	//}
	//tii := InnerInnerTest{
	//	Eight:  "whoa",
	//	Nine:   109,
	//	Ten:    []float64{0.5, 4.6},
	//	Eleven: []WhoaAgain{whoa1, whoa2},
	//}
	//ti := TestInner{
	//	Five:  50,
	//	Six:   true,
	//	Seven: tii,
	//}
	//t := Test{
	//	One:   10,
	//	Two:   50.5,
	//	Three: "hello",
	//	Four:  ti,
	//}
	//tBytes, err := bson.Marshal(t)
	//if err != nil {
	//	panic(err)
	//}
	//var newT Test
	//err = bson.Unmarshal(tBytes, &newT)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(newT)
	type InnerThing struct {
		ArgTen string
	}
	type Thing struct {
		ArgEight float32
		ArgNine  InnerThing
	}
	//type Args struct {
	//	ArgOne   string
	//	ArgTwo   any
	//	ArgThree any
	//	ArgFour  any
	//	ArgFive  uint32
	//	ArgSix   any
	//	ArgSeven any
	//}
	//args := Args{
	//	ArgOne: "hello",
	//	ArgTwo: uint(8),
	//	ArgThree: WhoaAgain{
	//		Twelve:   "hello",
	//		Thirteen: false,
	//	},
	//	ArgFour: []any{1.0, 100, "hello"},
	//	ArgFive: 80,
	//	ArgSix:  1.75,
	//	ArgSeven: Thing{
	//		ArgEight: 10.9,
	//		ArgNine:  InnerThing{ArgTen: "world"},
	//	},
	//}
	args := Thing{
		ArgEight: 10.5,
		ArgNine: InnerThing{
			ArgTen: "hello, world",
		},
	}
	argsBytes, _ := bson.Marshal(args)
	fmt.Println(argsBytes)

	var t Thing
	err := bson.Unmarshal(argsBytes, &t)
	if err != nil {
		panic(err)
	}
	fmt.Println(t)
	//servers, err := initServers(2)
	//if err != nil {
	//	panic(err)
	//}
	//n := servers[0].Node
	//err = servers[0].SendPing(servers[1])
	//if err != nil {
	//	panic(err)
	//}
	//type CallArgs struct {
	//	Sender kademlia.Contact
	//	Key    string
	//}
	//call := bsonrpc.Call{
	//	Method: "methodiswhoa",
	//	Args:   CallArgs{Sender: n.ToContact(), Key: "hello"},
	//}
	//callBytes, err := bson.Marshal(call)
	//if err != nil {
	//	panic(err)
	//}
	//var c bsonrpc.Call
	//err = bson.Unmarshal(callBytes, &c)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("c:", c)
	//args := kademlia.Args{
	//	Sender: servers[0].Node,
	//}
	//argsBytes, err := bson.Marshal(args)
	//
	//fmt.Println(argsBytes)
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
