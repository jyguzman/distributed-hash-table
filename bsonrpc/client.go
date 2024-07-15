package bsonrpc

import (
	"fmt"
	"go-dht/bson"
	"go-dht/kademlia"
	"math/big"
	"net"
	"strconv"
)

type Client struct {
	ID   *big.Int
	conn *net.UDPConn
}

func (c Client) Call(args bson.M) error {
	bytes, err := bson.Marshal(args)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(bytes)
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	n, _, err := c.conn.ReadFromUDP(buf)
	fmt.Println(string(buf[:n]))

	return nil
}

func (c Client) Ping() error {
	args := bson.M{
		"q":  "ping",
		"id": c.ID,
	}
	return c.Call(args)
}

func (c Client) Store(key string, data any) error {
	args := bson.M{
		"q":     "store",
		"id":    c.ID,
		"key":   kademlia.GetHash(key),
		"value": data,
	}
	return c.Call(args)
}

func (c Client) FindNodes(key string) ([]int, error) {
	args := bson.M{
		"q":   "find_node",
		"id":  c.ID,
		"key": kademlia.GetHash(key),
	}
	err := c.Call(args)
	if err != nil {
		return nil, err
	}
	return []int{}, nil
}

func (c Client) FindValue(key string) ([]int, error) {
	args := bson.M{
		"q":   "find_value",
		"id":  c.ID,
		"key": kademlia.GetHash(key),
	}
	err := c.Call(args)
	if err != nil {
		return nil, err
	}
	return []int{}, nil
}

func Dial(id *big.Int, host string, port int) (*Client, error) {
	addr, err := net.ResolveUDPAddr("udp", host+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	return &Client{id, conn}, nil
}
