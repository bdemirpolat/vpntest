package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(fmt.Sprintf("hi %s", request.RemoteAddr)))
		return
	})
	http.ListenAndServe("10.1.0.20:8990", nil)
}
