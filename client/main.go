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

func main() {
	iface, err := createTun("10.10.10.10")
	if err != nil {
		fmt.Println("interface can not created:", err)
		return
	}

	conn, err := createListener()
	if err != nil {
		fmt.Println("udp conn create error:", err)
	}

	go listenUDP(conn, iface)
	go listenInterface(iface, conn)

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, os.Interrupt, syscall.SIGTERM)
	<-termSignal
	fmt.Println("closing")
}

func createListener() (*net.UDPConn, error) {
	return net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{89, 252, 131, 88}, Port: 8990, Zone: ""})
}

func listenUDP(listener *net.UDPConn, iface *water.Interface) {
	for {
		fmt.Println("udp connection listening")
		message := make([]byte, 1500)
		for {
			n, err := listener.Read(message)
			if err != nil {
				log.Println("conn read error:", err)
			}
			message = message[:n]
			fmt.Println("START - incoming packet from TUNNEL")
			cmd.WritePacket(message)
			fmt.Println("DONE - incoming packet from TUNNEL")
			if iface != nil {
				_, err = iface.Write(message)
				if err != nil {
					log.Println("ifce write err:", err)
				} else {
					fmt.Println("iface write done")
				}
			}
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
		fmt.Println("START - incoming packet from INTERFACE")
		cmd.WritePacket(packet)
		fmt.Println("DONE - incoming packet from INTERFACE")
		if conn != nil {
			_, err = conn.Write(packet)
			if err != nil {
				log.Println("conn write error:", err)
			}
		}
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
