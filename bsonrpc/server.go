package bsonrpc

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Host string
	Port int
	conn *net.UDPConn
}

func NewServer(host string, port int) (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{host, port, conn}, nil
}

func (s *Server) Listen() {
	go func() {
		{
			buf := make([]byte, 1024)
			n, sender, err := s.conn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("Error reading from UDP socket: %s", err)
			}

			_, err = s.conn.WriteToUDP([]byte("Got message from "+sender.String()+": "+string(buf[:n])), sender)
			if err != nil {
				log.Printf("Error writing to UDP socket: %s", err)
			}
		}
	}()
}
