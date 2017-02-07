package main

import "github.com/funny/link"

type Customer struct {
	outSession *link.Session
	inSession  *link.Session
}

func NewCustomer(outSession, inSession *link.Session) *Customer {
	return &Customer{
		inSession:  inSession,
		outSession: outSession,
	}
}

func (c *Customer) handleOut() {
	for {
		msg, e := c.outSession.Receive()

		if e != nil {
			break
		}

		c.inSession.Send(msg)
	}
}

func (c *Customer) handleIn() {
	for {
		msg, e := c.inSession.Receive()

		if e != nil {
			break
		}

		c.outSession.Send(msg)
	}
}
