package main

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"vpntest/cmd"
)

func createTun() (*water.Interface, error) {
	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Interface Name: %s\n", iface.Name())

	out, err := cmd.RunCommand(fmt.Sprintf("sudo ifconfig %s 10.1.0.10 10.1.0.10 up", iface.Name()))
	if err != nil {
		fmt.Println(out)
		return nil, err
	}

	return iface, nil
}

func addRoute(iface *water.Interface) error {
	out, err := cmd.RunCommand(fmt.Sprintf("route add -host 10.1.0.10 -interface %s", iface.Name()))
	if err != nil {
		fmt.Println(out)
		return err
	}
	return nil
}
