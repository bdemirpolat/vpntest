package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	iface, err := createTun("10.1.0.10")
	if err != nil {
		fmt.Println("interface create err:", err)
		return
	}
	listener, err := createListener()
	if err != nil {
		fmt.Println("listener create err:", err)
		return
	}
	go runTestServer()
	go runTestServer2()
	go listenTCP(listener, iface)
	go listenInterface(iface)

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, os.Interrupt, syscall.SIGTERM)
	<-termSignal
	fmt.Println("closing")
}

func runTestServer() {
	http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(fmt.Sprintf("hi %s", request.RemoteAddr)))
		return
	})
	err := http.ListenAndServe("10.1.0.10:8080", nil)
	if err != nil {
		log.Println(err)
	}
}

func runTestServer2() {
	http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(fmt.Sprintf("hi %s", request.RemoteAddr)))
		return
	})
	err := http.ListenAndServe("10.1.0.20:8080", nil)
	if err != nil {
		log.Println(err)
	}
}
