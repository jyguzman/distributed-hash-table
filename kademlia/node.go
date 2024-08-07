package kademlia

import (
	"fmt"
	"go-dht/bson"
	"go-dht/pkg/util"
	"math/big"
	"strconv"
)

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
		addressHash := util.GetHash(host + ":" + strconv.Itoa(port))
		id = util.HashToBigInt(addressHash)
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

func (n Node) Equals(other Node) bool {
	return n.String() == other.String()
}

func (n Node) Xor(other Node) *big.Int {
	return new(big.Int).Xor(n.Id, other.Id)
}

func (n Node) Prefix(length int) string {
	pre := ""
	for i := 0; i < length; i++ {
		pre += strconv.Itoa(int(n.Id.Bit(159 - i)))
	}
	return pre
}
