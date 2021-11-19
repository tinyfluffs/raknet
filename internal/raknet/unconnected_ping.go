package raknet

import (
	"bytes"
	"encoding/binary"
)

type UnconnectedPing struct {
	SendTimestamp int64
	ClientGUID    int64
}

func (p *UnconnectedPing) ID() byte {
	return IDUnconnectedPing
}

func (p *UnconnectedPing) Marshal() []byte {
	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, IDUnconnectedPing)
	_ = binary.Write(buf, binary.BigEndian, p.SendTimestamp)
	_ = binary.Write(buf, binary.BigEndian, unconnectedMagic)
	_ = binary.Write(buf, binary.BigEndian, p.ClientGUID)
	return buf.Bytes()
}

func (p *UnconnectedPing) Unmarshal(buf *bytes.Buffer) {
	_ = binary.Read(buf, binary.BigEndian, &p.SendTimestamp)

	var magic [16]byte // Don't care for it
	_ = binary.Read(buf, binary.BigEndian, &magic)

	_ = binary.Read(buf, binary.BigEndian, &p.ClientGUID)
}
