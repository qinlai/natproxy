package main

import (
	"net"

	fastway "github.com/funny/fastway/go"
	"github.com/funny/link"
)

func NewClient(l net.Listener, endPoint *fastway.EndPoint, serverId uint32) *Client {
	return &Client{
		l:        l,
		endPoint: endPoint,
		serverId: serverId,
	}
}

type Client struct {
	l        net.Listener
	endPoint *fastway.EndPoint
	serverId uint32
}

func (c *Client) Serv() {
	for {
		if conn, e := c.l.Accept(); e == nil {
			inSession, e := c.endPoint.Dial(c.serverId)
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
	c.endPoint.Close()
}
