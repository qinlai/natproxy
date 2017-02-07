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

var gClients []*Client
var gRestartClientChan chan bool
var gEndPoint *fastway.EndPoint

func startClient(addr string) {
	e := doStartClient(addr)
	if e != nil {
		panic(fmt.Sprintf("can't connect to gateway, please check gateway have started,error:%s", e.Error()))
	}

	go func() {
		gRestartClientChan = make(chan bool)
	L:
		for {
			select {
			case s := <-gRestartClientChan:
				if !s {
					break L
				}
				log.Println("restart client")
				stopClient()
				for doStartClient(addr) != nil {
					time.Sleep(time.Second * 2)
				}
			}
		}
	}()
}

func reStartClient() {
	select {
	case gRestartClientChan <- true:
	default:
	}
}

func doStartClient(addr string) error {
	var e error
	gEndPoint, e = fastway.DialClient("tcp", addr, fastway.EndPointCfg{
		MemPool:      slab.NewSyncPool(64, 4*1024, 64),
		MaxPacket:    128 * 1024,
		PingInterval: 5 * time.Second,
		PingTimeout:  5 * time.Second,
		SendChanSize: 512,
		RecvChanSize: 512,
		MsgFormat:    NewProtocolFormat(),
		TimeoutCallback: func() bool {
			reStartClient()
			return false
		},
	})

	if e != nil {
		return e
	}

	gClients = nil
	for serverID, proxy := range gConf.Proxys {
		l, e := net.Listen("tcp", fmt.Sprintf(":%d", proxy.Port))
		if e != nil {
			return e
		}
		c := NewClient(l, uint32(serverID+1))
		gClients = append(gClients, c)
		log.Println(fmt.Sprintf("start proxy port %d to server %d", proxy.Port, serverID+1))
		go c.Serv()
	}

	return nil
}

func stopClient() {
	for _, c := range gClients {
		c.Stop()
	}
	select {
	case gRestartClientChan <- false:
	default:
	}
}

func NewClient(l net.Listener, serverId uint32) *Client {
	return &Client{
		l:        l,
		serverId: serverId,
	}
}

type Client struct {
	l        net.Listener
	serverId uint32
}

func (c *Client) Serv() {
	for {
		if conn, e := c.l.Accept(); e == nil {
			inSession, e := gEndPoint.Dial(c.serverId)
			if e == nil {
				session := link.NewSession(NewCustomerCodec(conn), 64)

				customer := NewCustomer(session, inSession)
				go customer.handleIn()
				go customer.handleOut()
			}
		}
	}
}

func (c *Client) Stop() {
	c.l.Close()
}
