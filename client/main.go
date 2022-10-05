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

var conn net.Conn

func main() {
	iface, err := createTun("10.10.10.10")
	if err != nil {
		fmt.Println("interface can not created:", err)
		return
	}

	conn, err = connectTcpServer()
	if err != nil {
		fmt.Println("can not connected to tcp server:", err)
		return
	}

	go listenTcpConn(iface, conn)
	go listenInterface(iface, conn)

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, os.Interrupt, syscall.SIGTERM)
	<-termSignal
	fmt.Println("closing")
}

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
		fmt.Println("START - incoming packet from TUNNEL")
		cmd.WritePacket(message)
		fmt.Println("DONE - incoming packet from TUNNEL")
		_, err = iface.Write(message)
		if err != nil {
			log.Println("write to interface failed:", err.Error())
		} else {
			fmt.Println("iface write done")
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
