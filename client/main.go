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
	iface, err := createTun()
	if err != nil {
		fmt.Println("interface can not created:", err)
		return
	}

	conn, err = connectTcpServer()
	if err != nil {
		fmt.Println("can not connected to tcp server:", err)
		return
	}

	/*	err = addRoute(iface)
		if err != nil {
			fmt.Println("add route error:", err)
		}*/

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

func createTun() (*water.Interface, error) {
	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Interface Name: %s\n", iface.Name())

	out, err := cmd.RunCommand(fmt.Sprintf("sudo ifconfig %s 192.168.35.1 192.168.35.35 up", iface.Name()))
	if err != nil {
		fmt.Println(out)
		return nil, err
	}

	return iface, nil
}

func addRoute(iface *water.Interface) error {
	out, err := cmd.RunCommand(fmt.Sprintf("route add -host 192.168.35.35 -interface %s", iface.Name()))
	//out, err := cmd.RunCommand(fmt.Sprintf("sudo ifconfig %s -alias 192.168.35.35", iface.Name()))
	if err != nil {
		fmt.Println(out)
		return err
	}
	return nil
}
