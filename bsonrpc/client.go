package bsonrpc

import (
	"fmt"
	"go-dht/bson"
	"net"
	"strconv"
)

type Client struct {
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

func Dial(host string, port int) (*Client, error) {
	addr, err := net.ResolveUDPAddr("udp", host+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn}, nil
}
