package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
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

	err = addRoute(iface)
	if err != nil {
		fmt.Println("add route error:", err)
	}

	go listenTcpConn(iface, conn)
	go listenInterface(iface, conn)

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, os.Interrupt, syscall.SIGTERM)
	fmt.Println("closing")
}
