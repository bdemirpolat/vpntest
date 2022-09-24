package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"net"
	"vpntest/cmd"
)

var conn net.Conn

func main() {
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	out, err := cmd.RunCommand(fmt.Sprintf("sudo ifconfig %s 10.1.0.10 10.1.0.20 up", ifce.Name()))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(out)

	conn, err = net.Dial("tcp", "89.252.131.88:8990")
	if err != nil {
		log.Println("Dial failed:", err.Error())
	}

	fmt.Println("dial ok")

	out, err = cmd.RunCommand(fmt.Sprintf("route add -host 10.1.0.10 -interface %s", ifce.Name()))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(out)

	go func() {
		message := make([]byte, 2000)
		for {
			n, err := conn.Read(message)
			if err != nil {
				log.Println("read from connection failed:", err.Error())
			}
			message = message[:n]
			cmd.WritePacket(message)
			_, err = ifce.Write(message)
			if err != nil {
				log.Println("write to interface failed:", err.Error())
			}
		}
	}()

	packet := make([]byte, 2000)
	for {
		fmt.Println("loop started")
		n, err := ifce.Read(packet)
		if err != nil {
			log.Println("ifce read error:", err)
		}
		packet = packet[:n]
		cmd.WritePacket(packet)
		log.Printf("Packet Received: % x\n", string(packet))
		_, err = conn.Write(packet)
		if err != nil {
			log.Println("conn write error:", err)
		}
	}
}
