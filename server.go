package main

import (
	"log"
	"net"

	fastway "github.com/funny/fastway/go"
	"github.com/funny/link"
)

func NewServer(proxy string, endPoint *fastway.EndPoint) *Server {
	return &Server{
		proxy:    proxy,
		endPoint: endPoint,
	}
}

type Server struct {
	endPoint *fastway.EndPoint
	proxy    string
}

func (s *Server) Serv() {
	for {
		inSession, e := s.endPoint.Accept()
		if e == nil {
			conn, e := net.Dial("tcp", s.proxy)
			if e == nil {
				session := link.NewSession(NewCustomerCodec(conn), 64)
				customer := NewCustomer(session, inSession)
				go customer.handleIn()
				go customer.handleOut()
			} else {
				log.Println("dial error : ", e.Error())
			}
		}
	}
}

func (s *Server) Stop() {
	s.endPoint.Close()
}
