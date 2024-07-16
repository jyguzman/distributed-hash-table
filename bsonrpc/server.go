package bsonrpc

import (
	"fmt"
	"go-dht/bson"
	"log"
	"net"
	"reflect"
	"strconv"
)

type Server struct {
	Host           string
	Port           int
	conn           *net.UDPConn
	ServiceMethods map[string]reflect.Method
	Service        reflect.Value
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
	return &Server{
		Host:           host,
		Port:           port,
		conn:           conn,
		ServiceMethods: make(map[string]reflect.Method),
	}, nil
}

func (s *Server) Listen() {
	fmt.Println("Listening on " + s.Host + ":" + strconv.Itoa(s.Port))
	for {
		buf := make([]byte, 1024)
		n, sender, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading from UDP socket: %s", err)
		}

		if string(buf[:n]) == "Ping" {
			s.sendResponse([]byte("Pong"), sender)
		}
	}
}

func (s *Server) sendResponse(message []byte, sender *net.UDPAddr) {
	_, err := s.conn.WriteToUDP(message, sender)
	if err != nil {
		log.Printf("Error writing to UDP socket: %s", err)
	}
}

func isValidMethod(serviceType reflect.Type, method reflect.Method) bool {
	return method.Type.NumIn() == 2 &&
		method.Type.In(0) == serviceType &&
		method.Type.In(1) == reflect.TypeOf(bson.M{}) &&
		method.Type.NumOut() == 1 &&
		method.Type.Out(0) == reflect.TypeOf((*error)(nil)).Elem()
}

func (s *Server) Register(receiver any) error {
	t := reflect.TypeOf(receiver)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("receiver must be a pointer to struct")
	}
	if t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("receiver must be a struct")
	}
	s.Service = reflect.ValueOf(receiver)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if isValidMethod(s.Service.Type(), method) {
			s.ServiceMethods[method.Name] = method
		}
	}
	if len(s.ServiceMethods) == 0 {
		log.Println("Warning: no methods were registered")
	}
	return nil
}

func (s *Server) Call(args bson.M) error {
	methodName := args["type"].(string)
	method, ok := s.ServiceMethods[methodName]
	if !ok {
		return fmt.Errorf("no such method: " + methodName)
	}
	val := method.Func.Call([]reflect.Value{s.Service, reflect.ValueOf(args)})
	err := val[0].Interface()
	if err != nil {
		return err.(error)
	}
	return nil
}
