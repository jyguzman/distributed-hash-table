package kademlia

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
)

type Node struct {
	ID      *big.Int
	IP      string
	Port    int
	K       int
	Buckets []KBucket
	//RoutingTable *RoutingTable
}

func NewNode(IP string, Port int) Node {
	addressHash := GetHash(IP + ":" + strconv.Itoa(Port))
	buckets := make([]KBucket, 160)
	return Node{
		IP: IP, Port: Port,
		ID:           HashToBigInt(addressHash),
		Buckets:      buckets,
		RoutingTable: NewRoutingTable(1),
	}
}

func (n Node) String() string {
	return fmt.Sprintf("(%s:%d %v)", n.IP, n.Port, n.ID)
}

func (n Node) Xor(other Node) *big.Int {
	return new(big.Int).Xor(n.ID, other.ID)
}

func (n Node) UpdateBucket(position int, node Node) {
	//bucket := n.Buckets[position]
	//bucket.Append(node)
	//n.Buckets[position].Append(node)

	n.RoutingTable.Add(node)
}

func RandNumber() *big.Int {
	limit := new(big.Int).Lsh(big.NewInt(1), 160)
	random, err := rand.Int(rand.Reader, limit)
	if err != nil {
		log.Fatal(err)
	}
	return random
}

func GetHash(item string) string {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(item))
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func HashToBigInt(hash string) *big.Int {
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		log.Fatal(err)
	}
	return new(big.Int).SetBytes(hashBytes)
}
