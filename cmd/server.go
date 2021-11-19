package main

import (
	"fmt"
	"github.com/tinyfluffs/raknet/pkg/raknet"
	"net"
)

func main() {
	ln, err := raknet.Listen("raknet", ":19132")
	if err != nil {
		panic(err)
	}
	l := ln.(*raknet.Listener)

	motd := fmt.Sprintf(
		"MCPE;%v;%v;%v;%v;%v;%v;Minecraft Server;%v;%v;%v;%v;",
		"A Minecraft Server",
		471,
		"1.17.41",
		0,
		20,
		l.ID(),
		"Creative",
		1,
		l.Addr().(*net.UDPAddr).Port,
		l.Addr().(*net.UDPAddr).Port,
	)
	l.MOTD([]byte(motd))

	for {
		_, err = l.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("connected")
	}
}
