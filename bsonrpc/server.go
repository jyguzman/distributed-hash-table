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
	host           string
	port           int
	conn           *net.UDPConn
	serviceMethods map[string]reflect.Method
	service        reflect.Value
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
		host:           host,
		port:           port,
		conn:           conn,
		serviceMethods: make(map[string]reflect.Method),
	}, nil
}

func (s *Server) Listen() {
	fmt.Println("Listening on " + s.host + ":" + strconv.Itoa(s.port))
	for {
		reqBytes, sender, err := s.readRequest()
		if err != nil {
			log.Println("Error reading request: " + err.Error())
		} else {
			reqObj, parseErr := s.unmarshalRequest(reqBytes)
			if parseErr != nil {
				log.Println("Error parsing request: " + parseErr.Error())
			} else {
				sendErr := s.sendResponse(reqObj, sender)
				if sendErr != nil {
					log.Println("Error sending response: " + sendErr.Error())
				}
			}
		}
	}
}

func (s *Server) readRequest() ([]byte, *net.UDPAddr, error) {
	buf := make([]byte, 1024)
	n, sender, err := s.conn.ReadFromUDP(buf)
	if err != nil {
		return nil, nil, err
	}
	return buf[:n], sender, err
}

func (s *Server) unmarshalRequest(req []byte) (bson.M, error) {
	args := bson.M{}
	err, _ := bson.Unmarshal(req, args)
	if err != nil {
		return nil, err
	}

	return args, nil
}

func (s *Server) sendResponse(request bson.M, sender *net.UDPAddr) error {
	reply := bson.M{}
	err := s.callMethod(request, reply)
	if err != nil {
		return err
	}

	replyBytes, err := bson.Marshal(reply)
	if err != nil {
		return err
	}

	_, err = s.conn.WriteToUDP(replyBytes, sender)
	if err != nil {
		return err
	}

	return nil
}

func isValidMethod(serviceType reflect.Type, method reflect.Method) bool {
	return method.Type.NumIn() == 3 &&
		method.Type.In(0) == serviceType &&
		method.Type.In(1) == reflect.TypeOf(bson.M{}) &&
		method.Type.In(2) == reflect.TypeOf(bson.M{}) &&
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
	s.service = reflect.ValueOf(receiver)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		_, exists := s.serviceMethods[method.Name]
		if !exists && isValidMethod(s.service.Type(), method) {
			s.serviceMethods[method.Name] = method
		}
	}
	if len(s.serviceMethods) == 0 {
		log.Println("Warning: no methods were registered")
	}
	return nil
}

func (s *Server) callMethod(args bson.M, reply bson.M) error {
	methodName := args["q"].(string)
	method, ok := s.serviceMethods[methodName]
	if !ok {
		return fmt.Errorf("no such method: " + methodName)
	}

	fnArgs := []reflect.Value{s.service, reflect.ValueOf(args), reflect.ValueOf(reply)}
	errVal := method.Func.Call(fnArgs)[0].Interface()
	if errVal != nil {
		return errVal.(error)
	}

	return nil
}
