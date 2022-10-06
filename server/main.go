package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vpntest/cmd"
)

func main() {
	ip := "192.168.53.1"
	iface, err := createTun(ip)
	if err != nil {
		fmt.Println("interface create err:", err)
		return
	}
	listener, err := createListener()
	if err != nil {
		fmt.Println("listener create err:", err)
		return
	}
	go runTestServer(iface.Name(), "192.168.53.2")
	go listenUDP(listener, iface)
	go listenInterface(iface)

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, os.Interrupt, syscall.SIGTERM)
	<-termSignal
	fmt.Println("closing")
}

func runTestServer(iface, ip string) {
	out, err := cmd.RunCommand(fmt.Sprintf("sudo ip addr add %s/24 dev %s", ip, iface))
	if err != nil {
		fmt.Println(out)
		return
	}
	http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(fmt.Sprintf("hi %s", request.RemoteAddr)))
		return
	})
	err = http.ListenAndServe(fmt.Sprintf("%s:8080", ip), nil)
	if err != nil {
		log.Println(err)
	}
}

func createListener() (*net.UDPConn, error) {
	return net.ListenUDP("udp", &net.UDPAddr{IP: []byte{89, 252, 131, 88}, Port: 8990, Zone: ""})
}

func listenUDP(listener *net.UDPConn, iface *water.Interface) {
	for {
		fmt.Println("udp connection listening")
		message := make([]byte, 65535)
		for {
			n, err := listener.Read(message)
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
		conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{89, 252, 131, 88}, Port: 8990, Zone: ""})
		if err == nil {
			_, err = conn.Write(packet[:n])
			if err != nil {
				log.Println("conn write error:", err)
			}
		}
		fmt.Println("START - incoming packet from INTERFACE")
		cmd.WritePacket(packet[:n])
		fmt.Println("DONE - incoming packet from INTERFACE")
		conn.Close()
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
