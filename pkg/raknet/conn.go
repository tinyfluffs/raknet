package raknet

import (
	"net"
	"time"
)

type Conn struct {
	backing net.Conn
	closed chan struct{}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	return c.backing.Read(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	return c.backing.Write(b)
}

func (c *Conn) Close() error {
	return c.backing.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.backing.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.backing.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.backing.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.backing.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.backing.SetWriteDeadline(t)
}

func (c *Conn) receive(b []byte) error {
	return nil
}
