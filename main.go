package main

import (
	"fmt"
	"go-dht/bson"
)

type Dummy struct {
	Thing int
}

func main() {
	doc := bson.D{
		{"hello", "world"},
		{"32_number", int32(1)},
		{"64_number", int64(64)},
		{"boolean", true},
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
	//fmt.Println(doc)
	//bServer, err := bsonrpc.NewServer("localhost", 8000)
	//if err != nil {
	//	panic(err)
	//}
	//server.Listen()
	//client, err := bsonrpc.Dial("localhost", 8000)
	//if err != nil {
	//	panic(err)
	//}
	//client.Call("PING")
	//if err != nil {
	//	panic(err)
	//}
	//server := kademlia.NewServer("localhost", 8001)
	//serverTwo := kademlia.NewServer("localhost", 8002)
	//server.Listen()
	//serverTwo.Listen()
	//err := server.BsonPing(serverTwo)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(server.BsonRpcServer.ServiceMethods)
	//args := bson.M{
	//	"type": "DummyMethod",
	//}
	//err := server.BsonRpcServer.Call(args)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//val := reflect.ValueOf(args)
	//server.BsonRpcServer.ServiceMethods["BSONDummyMethod"].Method.Func.Call([]reflect.Value{server.BsonRpcServer.Service, val})
	//serverTwo := kademlia.NewServer("localhost", 8002)
	//serverThree := kademlia.NewServer("localhost", 8003)
	//serverFour := kademlia.NewServer("localhost", 8004)
	//serverFive := kademlia.NewServer("localhost", 8005)
	//serverSix := kademlia.NewServer("localhost", 8006)
	//serverSeven := kademlia.NewServer("localhost", 8007)
	//serverEight := kademlia.NewServer("localhost", 8008)
	//serverNine := kademlia.NewServer("localhost", 8009)
	//server.Listen()
	//serverTwo.Listen()
	//serverThree.Listen()
	//serverFour.Listen()
	//serverFive.Listen()
	//serverSix.Listen()
	//serverSeven.Listen()
	//serverEight.Listen()
	//serverNine.Listen()
	//err := server.Ping(serverTwo)
	//err = server.Ping(serverThree)
	//err = server.Ping(serverFour)
	//err = server.Ping(serverFive)
	//err = server.Ping(serverSix)
	//err = server.Ping(serverSeven)
	//err = server.Ping(serverEight)
	//err = server.Ping(serverNine)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(server.Node.Buckets[1].String())
	//fmt.Println("Root:")
	//fmt.Println(server.RoutingTable.String())
	//fmt.Println()
	//kNearest := server.RoutingTable.GetNearest(serverSeven.Node.ID)
	//fmt.Println("K nearest:")
	//fmt.Println(kNearest)
	//fmt.Println(server.Node.ID)
	//fmt.Println(serverTwo.Node.ID)
	//fmt.Println(serverThree.Node.ID)
	//fmt.Println(serverFour.Node.ID)
	//fmt.Println(serverFive.Node.ID)
	//fmt.Println(serverSix.Node.ID)
	//fmt.Println(serverSeven.Node.ID)
	//fmt.Println(serverEight.Node.ID)
	//fmt.Println(serverNine.Node.ID)
	//fmt.Println(server.RoutingTable.BucketPrefixes)
	//fmt.Println(serverTwo.Node.RoutingTable)
	//fmt.Println(server.Node.Buckets[1].String())
	//fmt.Println(serverTwo.Node.Buckets[1].String())
	//fmt.Println(network.RandNumber())
	//hash := network.GetHash("random key")
	//fmt.Println(network.HashToBigInt(hash), network.HashToBigInt(hash).BitLen())
	//bytes, err := bson.Serialize(nil)
	//
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(bytes)
	//rt := kademlia.NewRoutingTable(2)
	//nodeOne := kademlia.NewNode("localhost", 8001)
	//nodeTwo := kademlia.NewNode("localhost", 8002)
	//nodeThree := kademlia.NewNode("localhost", 8003)
	//nodeFour := kademlia.NewNode("localhost", 8004)
	//nodeFive := kademlia.NewNode("localhost", 8005)
	//nodeSix := kademlia.NewNode("localhost", 8006)
	//nodeSeven := kademlia.NewNode("localhost", 8007)
	//nodeEight := kademlia.NewNode("localhost", 8008)
}
