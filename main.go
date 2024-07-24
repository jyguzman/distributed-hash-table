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
		servers[i], err = kademlia.NewServer("localhost", int32(8000+i))
		if err != nil {
			return nil, err
		}
		servers[i].Listen()
	}
	return servers, nil
}

type TestInner struct {
	Five int64
	Six  bool
}

type Test struct {
	One   int32
	Two   float64
	Three string
	Four  TestInner
}

func main() {
	//m := bson.M{
	//	"hello": "world",
	//	"inner": bson.D{
	//		{"three", 5.0},
	//	},
	//	"array": bson.A{
	//		"hello", "world", bson.D{
	//			{"yes", "sir"},
	//		},
	//	},
	//}
	//aBytes, err := bson.Marshal(m)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(aBytes)
	//newA := bson.D{}
	//err = bson.UnmarshalValue(bson.Object, aBytes, &newA)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("newA:", newA)
	//sBytes, err := bson.Marshal(bson.BSONString("hello, world"))
	//if err != nil {
	//	panic(err)
	//}
	//var s string
	//err = bson.Unmarshal(sBytes, &s)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("s:", s)
	//m := bson.M{
	//	"hello": "world",
	//	"inner": bson.D{
	//		{"three", 5.0},
	//	},
	//	"array": bson.A{
	//		"hello", "world", bson.D{
	//			{"yes", "sir"},
	//		},
	//	},
	//}
	type Holder struct {
		Nodes []kademlia.Node
	}
	n := kademlia.NewNode("localhost", 8000, nil)
	n2 := kademlia.NewNode("localhost", 8001, nil)
	//arr := []kademlia.Node{n, n2}
	//h := Holder{arr}
	//ti := TestInner{Five: 100, Six: true}
	//t := Test{One: 5, Two: 10.0, Three: "hello, world", Four: ti}
	//tBytes, err := bson.Marshal(h)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(tBytes)
	//var nodes kademlia.Node
	//nodes := bson.A{}
	//var nodes Holder
	//err = bson.Unmarshal(tBytes, &nodes)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(nodes)
	//m := bson.M{
	//	"One":   int32(5),
	//	"Two":   5.5,
	//	"Three": "bing",
	//	"Four": bson.M{
	//		"Five": int64(100),
	//		"Six":  true,
	//	},
	//}
	m := bson.M{
		"Nodes": bson.A{
			n, n2,
		},
	}
	var t Holder
	_, err := bson.UnmarshalStruct(m, &t)
	if err != nil {
		panic(err)
	}
	fmt.Println(t)
	//fmt.Println(reflect.TypeOf(newA["array"].(bson.A)[2]))
	//for _, elem := range a {
	//	fmt.Println("elem:", elem)
	//}
	//m := bson.M{
	//	"hello":  5.10,
	//	"number": 2.5,
	//	"inner": &bson.M{
	//		"one": 1.5,
	//	},
	//	"this": 10.0,
	//}
	//mBytes, err := bson.Marshal(m)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(mBytes)
	////fmt.Println(mBytes)
	//n := &bson.M{}
	//err = n.UnmarshalBSON(mBytes)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(*n)
	//r := bson.NewReader(mBytes)
	//s, err := r.ReadDocument()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("s:", *s)
	//raw := s.Pairs["inner"]
	////p, err := bson.Marshal(raw.Data)
	////if err != nil {
	////	panic(err)
	////}
	//newR := bson.NewReader(raw.Data)
	//sRawD, err := newR.ReadDocument()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("sRawD:", *sRawD)
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
