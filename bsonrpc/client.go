package bsonrpc

import (
	"go-dht/bson"
	"net"
	"strconv"
)

type Client struct {
	conn *net.UDPConn
}

func (c Client) Call(args bson.M, reply any) error {
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
	if err != nil {
		return err
	}

	err = bson.Unmarshal(buf[:n], reply)
	if err != nil {
		return err
	}

	return nil
}

func Dial(host string, port int) (*Client, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", host+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return nil, err
	}
	return &Client{conn}, nil
}
