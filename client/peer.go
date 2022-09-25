package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"net"
	"vpntest/cmd"
)

func connectTcpServer() (net.Conn, error) {
	conn, err := net.Dial("tcp", "89.252.131.88:8990")
	if err != nil {
		return nil, err
	}
	fmt.Println("dial ok")
	return conn, nil
}

func listenTcpConn(iface *water.Interface, conn net.Conn) {
	message := make([]byte, 2000)
	for {
		n, err := conn.Read(message)
		if err != nil {
			log.Println("read from connection failed:", err.Error())
		}
		message = message[:n]
		cmd.WritePacket(message)
		_, err = iface.Write(message)
		if err != nil {
			log.Println("write to interface failed:", err.Error())
		}
	}
}

func listenInterface(iface *water.Interface, conn net.Conn) {
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
