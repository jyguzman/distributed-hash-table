package kademlia

import (
	"fmt"
	"log"
	"net"
)

type UdpClient struct {
	Conn   *net.UDPConn
	Server Server
}

func (s Server) SendUDP(message string, server Server) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", server.Node.Host, server.Node.Port))
	if err != nil {
		log.Println(err)
	}
	//_, err = s.UdpServer.WriteToUDP([]byte(message), addr)
	if err != nil {
		log.Println(err)
	}
	server.Requests <- []string{message, addr.String()}
	server.WaitGroup.Add(1)

}

func (c UdpClient) Call(method string, args any, reply any) error {
	switch method {
	case "ping":
		_, err := c.Conn.Write([]byte("pong"))
		if err != nil {
			return err
		}
	case "store":
		fmt.Println("store")
	case "find_nodes":
		fmt.Println("find_node")
	case "find_value":
		fmt.Println("find_value")
	default:
		fmt.Println("Unknown method")
	}
	return nil
}

func KrpcPing(args Server, reply Response) {

}

type Request struct {
}

type Response struct {
}

func (s Server) SendPing(server Server) error {
	conn, err := s.Connect(server)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte("ping"))
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return err
	}
	fmt.Println(n, addr)
	return nil
}

func (s Server) SendPong(server Server) {}

func (s Server) Connect(server Server) (*net.UDPConn, error) {
	sAddrString := fmt.Sprintf("%s:%d", s.Node.Host, s.Node.Port)
	otherAddrStr := fmt.Sprintf("%s:%d", server.Node.Host, server.Node.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", sAddrString)
	if err != nil {
		return nil, err
	}
	otherAddr, err := net.ResolveUDPAddr("udp", otherAddrStr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", udpAddr, otherAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
