package cmd

import (
	"bytes"
	"fmt"
	"golang.org/x/net/ipv4"
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

func WritePacket(frame []byte) {
	header, err := ipv4.ParseHeader(frame)
	if err != nil {
		fmt.Println("write packet err:", err)
	} else {
		fmt.Println("SRC:", header.Src)
		fmt.Println("DST:", header.Dst)
	}
}
