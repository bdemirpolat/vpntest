package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"vpntest/cmd"
)

var conn *net.UDPConn

func main() {
	iface, err := createTun("192.168.9.9")
	if err != nil {
		fmt.Println("interface can not created:", err)
		return
	}
	conn, err = createConn()
	if err != nil {
		fmt.Println("udp conn create error:", err)
	}

	go listenUDP(iface)
	go listenInterface(iface)

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, os.Interrupt, syscall.SIGTERM)
	<-termSignal
	fmt.Println("closing")
}

func createConn() (*net.UDPConn, error) {
	return net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{89, 252, 131, 88}, Port: 8990, Zone: ""})
}

func listenUDP(iface *water.Interface) {
	for {
		fmt.Println("udp connection listening")
		message := make([]byte, 65535)
		for {
			n, err := conn.Read(message)
			if err != nil {
				log.Println("conn read error:", err)
			}
			if iface != nil {
				_, err = iface.Write(message[:n])
				if err != nil {
					log.Println("ifce write err:", err)
				} else {
					fmt.Println("iface write done")
				}
			}
			fmt.Println("START - incoming packet from TUNNEL")
			cmd.WritePacket(message[:n])
			fmt.Println("DONE - incoming packet from TUNNEL")
		}

	}
}

func listenInterface(iface *water.Interface) {
	fmt.Println("interface listening")
	packet := make([]byte, 65535)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			log.Println("ifce read error:", err)
		}

		if err == nil {
			_, err = conn.Write(packet[:n])
			if err != nil {
				log.Println("conn write error:", err)
			}
		}
		fmt.Println("START - incoming packet from INTERFACE")
		cmd.WritePacket(packet[:n])
		fmt.Println("DONE - incoming packet from INTERFACE")

	}
}

func createTun(ip string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(config)
	if err != nil {
		return nil, err
	}
	log.Printf("Interface Name: %s\n", iface.Name())
	out, err := cmd.RunCommand(fmt.Sprintf("sudo ip addr add %s/24 dev %s", ip, iface.Name()))
	if err != nil {
		fmt.Println(out)
		return nil, err
	}

	out, err = cmd.RunCommand(fmt.Sprintf("sudo ip link set dev %s up", iface.Name()))
	if err != nil {
		fmt.Println(out)
		return nil, err
	}
	return iface, nil
}
