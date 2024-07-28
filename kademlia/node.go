package kademlia

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"go-dht/bson"
	"log"
	"math/big"
	"strconv"
)

type Key big.Int

func (k Key) String() string {
	return (*big.Int)(&k).Text(16)
}

func (k Key) FromString(s string, base int) error {
	i, ok := (*big.Int)(&k).SetString(s, base)
	if !ok {
		return errors.New("invalid key")
	}
	fmt.Println("I:", i, "k:", k)
	return nil
}

func (k Key) MarshalBSON() ([]byte, error) {
	return bson.Marshal(k.String())
}

func (k *Key) UnmarshalBSON(data []byte) error {
	var s string
	err := bson.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	_, ok := (*big.Int)(k).SetString(s, 16)
	if !ok {
		return errors.New("invalid key")
	}
	return nil
}

type Node struct {
	Id   *big.Int
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
	return Node{Host: host, Port: port, Id: id}
}

func (n Node) MarshalBSON() ([]byte, error) {
	m := bson.M{
		"Id":   n.Id.Text(16),
		"Host": n.Host,
		"Port": n.Port,
	}
	data, err := bson.Marshal(m)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (n *Node) UnmarshalBSON(data []byte) error {
	var contact Contact
	err := bson.Unmarshal(data, &contact)
	if err != nil {
		return err
	}
	n.Host = contact.Host
	n.Port = contact.Port
	id, ok := new(big.Int).SetString(contact.Id, 16)
	if !ok {
		return fmt.Errorf("invalid id %s", contact.Id)
	}
	n.Id = id
	return nil
}

func (n Node) String() string {
	return fmt.Sprintf("(%s:%d %s)", n.Host, n.Port, n.Id.Text(16))
}

func (n Node) Xor(other Node) *big.Int {
	return new(big.Int).Xor(n.Id, other.Id)
}

func (n Node) Prefix(length int) string {
	pre := ""
	for i := 0; i < length; i++ {
		pre += strconv.Itoa(int(n.Id.Bit(i)))
	}
	return pre
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
