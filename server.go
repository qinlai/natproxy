package main

import (
	"fmt"
	"log"
	"net"
	"time"

	fastway "github.com/funny/fastway/go"
	"github.com/funny/link"
	"github.com/funny/slab"
)

var gServers []*Server

func startServer(serverAddr string) {
	for serverID, proxy := range gConf.Proxys {
		endPoint, e := fastway.DialServer("tcp", serverAddr, fastway.EndPointCfg{
			ServerID:     uint32(serverID + 1),
			MemPool:      slab.NewSyncPool(64, 4*1024, 64),
			MaxPacket:    128 * 1024,
			PingInterval: 60 * time.Second,
			PingTimeout:  5 * time.Second,
			AuthKey:      gConf.AuthKey,
			SendChanSize: 100000,
			RecvChanSize: 100000,
			MsgFormat:    NewProtocolFormat(),
			TimeoutCallback: func() bool {
				log.Println("fastway disconnected")
				return false
			},
		})

		if e != nil {
			panic(e.Error())
		}

		s := NewServer(proxy.Proxy, endPoint)
		gServers = append(gServers, s)
		log.Println(fmt.Sprintf("start proxy server %d (%s)", serverID+1, proxy.Proxy))
		go s.Serv()
	}
}

func stopServer() {
	for _, s := range gServers {
		s.Stop()
	}
}

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
