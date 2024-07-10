package main

import (
	"fmt"
	"go-dht/kademlia"
	"math/big"
	"strconv"
)

func FirstNBits(val *big.Int, n uint) string {
	res := ""
	for i := 0; i < int(n); i++ {
		res += strconv.Itoa(int(val.Bit(val.BitLen() - i - 1)))
	}
	return res
}
func main() {
	server := kademlia.NewServer("localhost", 8001)
	serverTwo := kademlia.NewServer("localhost", 8002)
	serverThree := kademlia.NewServer("localhost", 8003)
	serverFour := kademlia.NewServer("localhost", 8004)
	serverFive := kademlia.NewServer("localhost", 8005)
	serverSix := kademlia.NewServer("localhost", 8006)
	serverSeven := kademlia.NewServer("localhost", 8007)
	serverEight := kademlia.NewServer("localhost", 8008)
	serverNine := kademlia.NewServer("localhost", 8009)
	server.Listen()
	serverTwo.Listen()
	serverThree.Listen()
	serverFour.Listen()
	serverFive.Listen()
	serverSix.Listen()
	serverSeven.Listen()
	serverEight.Listen()
	serverNine.Listen()
	err := server.Ping(serverTwo)
	err = server.Ping(serverThree)
	err = server.Ping(serverFour)
	err = server.Ping(serverFive)
	err = server.Ping(serverSix)
	err = server.Ping(serverSeven)
	err = server.Ping(serverEight)
	err = server.Ping(serverNine)
	if err != nil {
		panic(err)
	}
	//fmt.Println(server.Node.Buckets[1].String())
	fmt.Println("Root:")
	fmt.Println(server.RoutingTable.String())
	//fmt.Println(serverTwo.Node.RoutingTable)
	//fmt.Println(server.Node.Buckets[1].String())
	//fmt.Println(serverTwo.Node.Buckets[1].String())
	//fmt.Println(network.RandNumber())
	//hash := network.GetHash("random key")
	//fmt.Println(network.HashToBigInt(hash), network.HashToBigInt(hash).BitLen())
	//bytes, err := protocol.Serialize(nil)
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

	//nodeOne.UpdateBucket(0, nodeOne)
	//rt.Add(0, nodeOne)
	//rt.Add(5, nodeThree)
	//rt.Add(1, nodeTwo)
	//rt.Add(5, nodeFour)
	//rt.Add(3, nodeFive)
	//rt.Add(2, nodeSix)
	//rt.Add(4, nodeSeven)
	//rt.Add(1, nodeEight)
	//fmt.Println(rt.Root)
	fmt.Println(FirstNBits(server.Node.ID, 10))
	fmt.Println(FirstNBits(serverTwo.Node.ID, 10))
	fmt.Println(FirstNBits(serverThree.Node.ID, 10))
	fmt.Println(FirstNBits(serverFour.Node.ID, 10))
	fmt.Println(FirstNBits(serverFive.Node.ID, 10))
	fmt.Println(FirstNBits(serverSix.Node.ID, 10))
	fmt.Println(FirstNBits(serverSeven.Node.ID, 10))
	fmt.Println(FirstNBits(serverEight.Node.ID, 10))
	fmt.Println(FirstNBits(serverNine.Node.ID, 10))
	//nodeOne.UpdateBucket(0, nodeThree)
	//fmt.Println(nodeOne.Buckets[0].String())
}
