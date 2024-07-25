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
	Port int32
}

type Contact struct {
	Id   string
	Host string
	Port int32
}

func NewNode(host string, port int32, id *big.Int) Node {
	if id == nil {
		addressHash := GetHash(host + ":" + strconv.Itoa(int(port)))
		id = HashToBigInt(addressHash)
	}
	return Node{Host: host, Port: port, ID: id}
}

func (n Node) ToTriple() Contact {
	return Contact{Host: n.Host, Port: n.Port, Id: n.ID.String()}
}

func FromContact(contact Contact) Node {
	id, _ := new(big.Int).SetString(contact.Id, 16)
	return Node{ID: id, Host: contact.Host, Port: contact.Port}
}

func NodeFromTuple(tuple bson.A) Node {
	fmt.Println("tuple", tuple)
	host, port, id := tuple[0].(string), tuple[1].(int32), tuple[2].(*big.Int)
	return NewNode(host, port, id)
}

func NodeFromMap(arrMap bson.M) (Node, error) {
	result := bson.A(make([]any, 3))
	for i, _ := range arrMap {
		idx, err := strconv.Atoi(i)
		if err != nil {
			return Node{}, err
		}
		result[idx] = arrMap[i]
	}
	idInt, ok := new(big.Int).SetString(result[2].(string), 16)
	if !ok {
		return Node{}, fmt.Errorf("invalid id %s", result[2])
	}
	result[2] = idInt
	return NodeFromTuple(result), nil
}

//func (x *big.Int) MarshalBSONValue() (Type, []byte, error) {
//	str := x.Text(16)
//	return bson.BSONString(str).MarshalBSONValue()
//}

func (n Node) MarshalBSON() ([]byte, error) {
	m := bson.M{
		"Id":   n.ID.Text(16),
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
	fmt.Println("contact", contact)
	n.Host = contact.Host
	n.Port = contact.Port
	id, ok := new(big.Int).SetString(contact.Id, 16)
	if !ok {
		return fmt.Errorf("invalid id %s", contact.Id)
	}
	n.ID = id
	return nil
}

func (n Node) Tuple() bson.A {
	return bson.A{n.Host, n.Port, n.ID}
}

func (n Node) String() string {
	return fmt.Sprintf("(%s:%d %s)", n.Host, n.Port, n.ID.Text(16))
}

func (n Node) Xor(other Node) *big.Int {
	return new(big.Int).Xor(n.ID, other.ID)
}

func (n Node) Prefix(length int) string {
	pre := ""
	for i := 0; i < length; i++ {
		pre += strconv.Itoa(int(n.ID.Bit(i)))
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
