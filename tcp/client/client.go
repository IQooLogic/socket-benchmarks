package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"syscall"
)

func main() {
	//cpuprofile := "cpu-client.pprof"
	//f, err := os.Create(cpuprofile)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "could not create CPU profile: %v\n", err)
	//	panic(err)
	//}
	//if err := pprof.StartCPUProfile(f); err != nil {
	//	fmt.Fprintf(os.Stderr, "could not start CPU profile: %v\n", err)
	//	panic(err)
	//}

	addrName := "localhost:8080"
	conn, err := net.Dial("tcp", addrName)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		pprof.StopCPUProfile()
		os.Exit(1)
	}()

	for {
		var message string
		if len(os.Args) > 1 {
			message = strings.Join(os.Args[1:], " ")
		} else {
			message = fmt.Sprintf("Hello, %s", addrName)
		}

		fmt.Println("Sending message:", message)
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err.Error())
			return
		}

		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading response:", err.Error())
			return
		}
		fmt.Println("Received message:", string(buf[:n]))
	}
}
