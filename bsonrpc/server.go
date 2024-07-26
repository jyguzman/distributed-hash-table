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
	serviceMethods map[string]*ServiceMethod
	service        reflect.Value
}

type ServiceMethod struct {
	Method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
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
		serviceMethods: make(map[string]*ServiceMethod),
	}, nil
}

func (s *Server) Listen() {
	fmt.Println("Listening on " + s.host + ":" + strconv.Itoa(s.port))
	for {
		reqBytes, sender, err := s.readRequest()
		if err != nil {
			log.Println("Error reading request: " + err.Error())
			return
		}
		reqObj, err := s.unmarshalRequest(reqBytes)
		if err != nil {
			log.Println("Error parsing request: " + err.Error())
			return
		}
		replyBytes, err := s.handleRequest(reqObj)
		if err != nil {
			log.Println("Error handling request: " + err.Error())
			return
		}
		sendErr := s.sendResponse(replyBytes, sender)
		if sendErr != nil {
			log.Println("Error sending response: " + sendErr.Error())
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

func (s *Server) unmarshalRequest(req []byte) (*Call, error) {
	var c Call
	err := bson.Unmarshal(req, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Server) handleRequest(request *Call) ([]byte, error) {
	serviceMethod, ok := s.serviceMethods[request.Method]
	if !ok {
		return nil, fmt.Errorf("no such method: " + request.Method)
	}

	reply := reflect.New(serviceMethod.ReplyType.Elem())
	err := s.call(*serviceMethod, request.Args, reply)
	if err != nil {
		return nil, err
	}

	replyBytes, err := bson.Marshal(reply.Elem().Interface())
	if err != nil {
		return nil, err
	}

	return replyBytes, nil
}

func (s *Server) sendResponse(response []byte, sender *net.UDPAddr) error {
	_, err := s.conn.WriteToUDP(response, sender)
	if err != nil {
		return err
	}
	return nil
}

func isValidMethod(serviceType reflect.Type, method reflect.Method) bool {
	return method.Type.NumIn() == 3 &&
		method.Type.In(0) == serviceType &&
		method.Type.In(2).Kind() == reflect.Ptr &&
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
			s.serviceMethods[method.Name] = &ServiceMethod{
				Method:    method,
				ArgType:   method.Type.In(1),
				ReplyType: method.Type.In(2),
			}
		}
	}
	if len(s.serviceMethods) == 0 {
		log.Println("Warning: no methods were registered")
	}
	return nil
}

func (s *Server) call(serviceMethod ServiceMethod, args any, reply reflect.Value) error {
	fnArgs := []reflect.Value{s.service, reflect.ValueOf(args), reply}
	errVal := serviceMethod.Method.Func.Call(fnArgs)[0].Interface()
	if errVal != nil {
		return errVal.(error)
	}
	return nil
}
