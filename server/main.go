package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"net"
	"vpntest/cmd"
)

var ifce *water.Interface

func main() {
	l, err := net.Listen("tcp", "89.252.131.88:8990")
	if err != nil {
		log.Println(err)
	}
	for {
		fmt.Println("listening")
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}

		ifce, err = water.New(water.Config{
			DeviceType: water.TUN,
		})
		if err != nil {
			log.Println(err)
		}

		log.Printf("Interface Name: %s\n", ifce.Name())

		fmt.Println("incoming connection remote addr:", conn.RemoteAddr())

		out, err := cmd.RunCommand(fmt.Sprintf("sudo ip addr add %s/24 dev %s", "192.168.1.60", ifce.Name()))
		if err != nil {
			log.Println("ip addr add error:", out, err)
		}
		fmt.Println(out)

		out, err = cmd.RunCommand(fmt.Sprintf("sudo ip link set dev %s up", ifce.Name()))
		if err != nil {
			log.Println("ip link set error:", out, err)
		}
		fmt.Println(out)

		go func(tcpConn net.Conn) {
			for {
				message := make([]byte, 2000)
				_, err := conn.Read(message)
				if err != nil {
					log.Println(err)
				}
				fmt.Println("incoming message:", string(message))
				_, err = ifce.Write(message)
				if err != nil {
					log.Println(err)
				}
			}
		}(conn)

		packet := make([]byte, 2000)
		for {
			n, err := ifce.Read(packet)
			if err != nil {
				log.Println(err)
			}
			log.Printf("Packet Received: % x\n", packet[:n])
			_, err = conn.Write(packet)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
