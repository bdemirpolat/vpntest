package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"net"
	"vpntest/cmd"
)

var conn net.Conn

func createListener() (net.Listener, error) {
	return net.Listen("tcp", "89.252.131.88:8990")
}

func listenTCP(listener net.Listener, iface *water.Interface) {
	for {
		fmt.Println("tcp connection listening")
		var err error
		conn, err = listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go handle(conn, iface)
	}
}
func handle(conn net.Conn, iface *water.Interface) {
	message := make([]byte, 2000)
	n, err := conn.Read(message)
	if err != nil {
		log.Println("conn read error:", err)
	}
	message = message[:n]
	cmd.WritePacket(message)
	if iface != nil {
		_, err = iface.Write(message)
		if err != nil {
			log.Println("ifce write err:", err)
		}
	}
}

func listenInterface(iface *water.Interface) {
	fmt.Println("interface listening")
	packet := make([]byte, 2000)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			log.Println("ifce read error:", err)
		}
		packet = packet[:n]
		cmd.WritePacket(packet)
		log.Printf("Packet Received: % x\n", packet)
		_, err = conn.Write(packet)
		if err != nil {
			log.Println("conn write error:", err)
		}
	}
}
