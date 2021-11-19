package raknet

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tinyfluffs/raknet/internal/raknet"
	"log"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	mtu             = 1500
	bitFlagDatagram = 0x80
)

func Listen(network, address string) (net.Listener, error) {
	if network != "raknet" {
		return net.Listen(network, address)
	}
	var err error

	l := &Listener{}
	l.id = rand.NewSource(time.Now().Unix()).Int63()
	l.conn, err = net.ListenPacket("udp", address)
	l.motd.Store([]byte{})

	if err != nil {
		return nil, err
	}
	go l.listen()
	return l, nil
}

type Listener struct {
	conn        net.PacketConn
	id          int64
	motd        atomic.Value
	incoming    chan *Conn
	connections sync.Map
}

func (l *Listener) Accept() (net.Conn, error) {
	conn, ok := <- l.incoming
	if !ok {
		return nil, errors.New("listener closed")
	}
	return conn, nil
}

func (l *Listener) Close() error {
	return l.conn.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.conn.LocalAddr()
}

func (l *Listener) ID() int64 {
	return l.id
}

func (l *Listener) MOTD(b []byte) {
	l.motd.Store(b)
}

func (l *Listener) listen() {
	b := make([]byte, mtu)
	for {
		n, addr, err := l.conn.ReadFrom(b)
		if err != nil {
			close(l.incoming)
			return
		}
		if err = l.handle(b[:n], addr); err != nil {
			log.Printf("error handling packet (addr=%v): %v\n", addr, err)
		}
	}
}

func (l *Listener) handle(b []byte, addr net.Addr) error {
	val, ok := l.connections.Load(addr)
	if !ok {
		// Offline
		buf := bytes.NewBuffer(b)
		id, err := buf.ReadByte()
		if err != nil {
			return fmt.Errorf("error reading packet ID byte: %v", err)
		}

		switch id {
		case raknet.IDUnconnectedPing, raknet.IDUnconnectedPingOpenConnections:
			return l.handleUnconnectedPing(buf, addr)
		case raknet.IDOpenConnectionRequest1:
			return l.handleOpenConnectionRequest1(buf, addr)
		case raknet.IDOpenConnectionRequest2:
			return l.handleOpenConnectionRequest2(buf, addr)
		default:
			if id&bitFlagDatagram == 0 {
				return fmt.Errorf("unknown packet received (%x): %x", id, b)
			}
		}
	}
	conn := val.(*Conn)
	select {
	case <-conn.closed:
		return nil
	default:
		return conn.receive(b)
	}
}

func (l *Listener) handleUnconnectedPing(buf *bytes.Buffer, addr net.Addr) error {
	p := &raknet.UnconnectedPing{}
	p.Unmarshal(buf)
	pong := &raknet.UnconnectedPong{
		ServerGUID:    l.id,
		SendTimestamp: p.SendTimestamp,
		MOTD:          l.motd.Load().([]byte),
	}
	_, err := l.conn.WriteTo(pong.Marshal(), addr)
	return err
}

func (l *Listener) handleOpenConnectionRequest1(buf *bytes.Buffer, addr net.Addr) error {
	return nil
}

func (l *Listener) handleOpenConnectionRequest2(buf *bytes.Buffer, addr net.Addr) error {
	return nil
}
