package bsonrpc

import (
	"fmt"
	"go-dht/bson"
	"log"
	"math/big"
	"net"
)

type Server struct {
	Host string
	Port int
	ID   *big.Int
	conn *net.UDPConn
}

func NewServer(host string, port int, id *big.Int) (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{host, port, id, conn}, nil
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

func (s *Server) sendResponse(message []byte, sender *net.UDPAddr) {
	_, err := s.conn.WriteToUDP(message, sender)
	if err != nil {
		log.Printf("Error writing to UDP socket: %s", err)
	}
}

func (s *Server) Pong(addr *net.UDPAddr) {
	bytes, err := bson.Marshal(bson.M{"id": s.ID})
	if err != nil {
		log.Printf("Error marshalling response: %s", err)
		return
	}
	s.sendResponse(bytes, addr)
}
