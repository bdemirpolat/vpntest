package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"net"
	"net/http"
	"vpntest/cmd"
)

var ifce *water.Interface

func main() {
	config := water.Config{
		DeviceType: water.TUN,
	}

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Println(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	out, err := cmd.RunCommand(fmt.Sprintf("sudo ip addr add %s/32 dev %s", "10.1.0.10", ifce.Name()))
	if err != nil {
		log.Println("ip addr add error:", out, err)
	}
	fmt.Println(out)

	out, err = cmd.RunCommand(fmt.Sprintf("sudo ip link set dev %s up", ifce.Name()))
	if err != nil {
		log.Println("ip link set error:", out, err)
	}
	fmt.Println(out)

	go runTestServer()

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

		go func(tcpConn net.Conn) {
			message := make([]byte, 2000)
			for {
				n, err := conn.Read(message)
				if err != nil {
					log.Println("conn read error:", err)
				}
				message = message[:n]
				cmd.WritePacket(message)
				if ifce != nil {
					_, err = ifce.Write(message)
					if err != nil {
						log.Println("ifce write err:", err)
					}
				}
			}
		}(conn)
		packet := make([]byte, 2000)
		for {
			n, err := ifce.Read(packet)
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
}

func runTestServer() {
	http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(fmt.Sprintf("hi %s", request.RemoteAddr)))
		return
	})
	http.ListenAndServe("10.1.0.10:8080", nil)
}
