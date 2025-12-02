package expressgo

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	router     *Router
	middleware []Middleware
}

func NewServer() *Server {
	return &Server{
		router: NewRouter(),
	}
}

func (s *Server) Listen(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Print("Error opening port: ", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Print("Error reading packet: ", err)
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\nError Parsing Request"))
		return nil
	}

	parts := strings.Split(strings.TrimSpace(line), " ")
	method := parts[0]
	path := parts[1]
	httpVersion := parts[2]
	fmt.Printf("Method: %v; Path: %v; Version: %v; \n", method, path, httpVersion)

	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}

		header := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(header) == 2 {
			key := strings.TrimSpace(header[0])
			value := strings.TrimSpace(header[1])
			fmt.Printf("%s: %s \n", key, value)
			headers[key] = value
		}
	}

	var body []byte
	if contentLength, exists := headers["Content-Length"]; exists {
		length, _ := strconv.Atoi(contentLength)
		body = make([]byte, length)
		io.ReadFull(reader, body)
	}

	res := NewRes(conn)
	req := NewReq(method, path, headers, body)
	routeHandler, err := s.router.Match(method, path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	requestHandler := s.applyMiddleware(routeHandler)
	requestHandler(req, res)

	return nil
}

func (s *Server) Get(path string, handler Handler) error {
	s.router.Get[path] = handler
	return nil
}

func (s *Server) Post(path string, handler Handler) error {
	s.router.Post[path] = handler
	return nil
}

func (s *Server) Use(middleware Middleware) {
	s.middleware = append(s.middleware, middleware)
}

func (s *Server) applyMiddleware(handler Handler) Handler {
	result := handler

	for i := len(s.middleware) - 1; i >= 0; i-- {
		result = s.middleware[i](result)
	}
	return result
}
