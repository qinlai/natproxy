package main

import "net"

type customerCodec struct {
	conn net.Conn
}

func NewCustomerCodec(conn net.Conn) *customerCodec {
	return &customerCodec{
		conn: conn,
	}
}

func (c *customerCodec) Receive() (interface{}, error) {
	buff := make([]byte, 4096)
	n, e := c.conn.Read(buff)
	if e != nil {
		return nil, e
	}

	return buff[:n], nil
}

func (c *customerCodec) Send(msg interface{}) error {
	_, e := c.conn.Write(msg.([]byte))
	return e
}

func (c *customerCodec) Close() error {
	return c.conn.Close()
}
