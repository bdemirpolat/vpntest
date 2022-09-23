package cmd

import (
	"bytes"
	"github.com/songgao/packets/ethernet"
	"log"
	"os/exec"
)

func RunCommand(command string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if stderr.String() != "" {
		return stderr.String(), err
	}
	return stdout.String(), err

}

func WritePacket(frame ethernet.Frame) {
	log.Printf("Dst: %s\n", frame.Destination())
	log.Printf("Src: %s\n", frame.Source())
	log.Printf("Ethertype: % x\n", frame.Ethertype())
	log.Printf("Payload: % x\n", frame.Payload())
}
