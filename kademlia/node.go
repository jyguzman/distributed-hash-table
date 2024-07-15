package kademlia

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"go-dht/bson"
	"log"
	"math/big"
	"strconv"
)

type Node struct {
	ID      *big.Int
	Host    string
	Port    int
	K       int
	Buckets []KBucket
}

func NewNode(host string, port int) Node {
	addressHash := GetHash(host + ":" + strconv.Itoa(port))
	buckets := make([]KBucket, 160)
	return Node{
		Host: host, Port: port,
		ID:      HashToBigInt(addressHash),
		Buckets: buckets,
	}
}

func (n Node) Tuple() bson.A {
	return bson.A{n.Host, n.Port, n.ID}
}

func (n Node) String() string {
	return fmt.Sprintf("(%s:%d %v)", n.Host, n.Port, n.ID)
}

func (n Node) Xor(other Node) *big.Int {
	return new(big.Int).Xor(n.ID, other.ID)
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
