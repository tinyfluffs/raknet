package raknet

import (
	"bytes"
	"encoding/binary"
)

type UnconnectedPong struct {
	SendTimestamp int64
	ServerGUID    int64
	MOTD          []byte
}

func (p *UnconnectedPong) ID() byte {
	return IDUnconnectedPong
}

func (p *UnconnectedPong) Marshal() []byte {
	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, IDUnconnectedPong)
	_ = binary.Write(buf, binary.BigEndian, p.SendTimestamp)
	_ = binary.Write(buf, binary.BigEndian, p.ServerGUID)
	_ = binary.Write(buf, binary.BigEndian, unconnectedMagic)
	_ = binary.Write(buf, binary.BigEndian, int16(len(p.MOTD)))
	_ = binary.Write(buf, binary.BigEndian, p.MOTD)
	return buf.Bytes()
}

func (p *UnconnectedPong) Unmarshal(buf *bytes.Buffer) {
	_ = binary.Read(buf, binary.BigEndian, &p.SendTimestamp)
	_ = binary.Read(buf, binary.BigEndian, &p.ServerGUID)

	var magic [16]byte // Don't care for it
	_ = binary.Read(buf, binary.BigEndian, &magic)

	var motdLen int16
	_ = binary.Read(buf, binary.BigEndian, &motdLen)
	p.MOTD = make([]byte, motdLen)
	_, _ = buf.Read(p.MOTD)
}
