package bsonrpc

import (
	"go-dht/bson"
	"net"
	"strconv"
)

type Client struct {
	conn *net.UDPConn
}

type Call struct {
	Method string
	Args   any
}

func (c Client) Call(methodName string, args any, reply any) error {
	call := Call{Method: methodName, Args: args}

	bytes, err := bson.Marshal(call)
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
