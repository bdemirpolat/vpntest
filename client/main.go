package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"net"
	"vpntest/cmd"
)

var conn *net.TCPConn

func main() {
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	out, err := cmd.RunCommand(fmt.Sprintf("sudo ifconfig %s 10.1.0.10 10.1.0.20 up", ifce.Name()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

	out, err = cmd.RunCommand(fmt.Sprintf("route add -host 89.252.131.88 -interface %s", ifce.Name()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

	serverAddr := "89.252.131.88:8990"
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		log.Fatal("ResolveTCPAddr failed:", err.Error())
	}

	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("Dial failed:", err.Error())
	}

	fmt.Println("dial ok")

	go func() {
		for {
			reply := make([]byte, 2000)
			_, err = conn.Read(reply)
			if err != nil {
				log.Fatal("Write to server failed:", err.Error())
			}

			_, err = ifce.Write(reply)
			if err != nil {
				log.Fatal("write to interface failed:", err.Error())
			}
		}
	}()

	packet := make([]byte, 2000)
	for {
		fmt.Println("loop started")
		n, err := ifce.Read(packet)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Packet Received: % x\n", packet[:n])
	}
}
