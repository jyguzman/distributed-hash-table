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
	ID   *big.Int
	Host string
	Port int
}

type Contact struct {
	Id   string
	Host string
	Port int
}

func NewNode(host string, port int, id *big.Int) Node {
	if id == nil {
		addressHash := GetHash(host + ":" + strconv.Itoa(port))
		id = HashToBigInt(addressHash)
	}
	return Node{Host: host, Port: port, ID: id}
}

func (n Node) ToTriple() Contact {
	return Contact{Id: n.ID.String(), Host: n.Host, Port: n.Port}
}

func FromContact(contact Contact) Node {
	id, _ := new(big.Int).SetString(contact.Id, 16)
	return Node{ID: id, Host: contact.Host, Port: contact.Port}
}

func FromTuple(tuple bson.A) Node {
	id, host, port := tuple[0].(*big.Int), tuple[1].(string), tuple[2].(int)
	return NewNode(host, port, id)
}

func (n Node) Tuple() bson.A {
	return bson.A{n.ID, n.Host, n.Port}
}

func (n Node) String() string {
	return fmt.Sprintf("(%s:%d %s)", n.Host, n.Port, n.ID.Text(16))
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
