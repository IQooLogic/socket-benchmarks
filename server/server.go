package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

func main() {
	cpuprofile := "cpu-server.pprof"
	f, err := os.Create(cpuprofile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create CPU profile: %v\n", err)
		panic(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Fprintf(os.Stderr, "could not start CPU profile: %v\n", err)
		panic(err)
	}

	addrName := "/tmp/echo.sock"
	os.Remove(addrName)

	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: addrName, Net: "unix"})
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}

	defer l.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(addrName)
		pprof.StopCPUProfile()
		os.Exit(1)
	}()

	fmt.Println("Listening on ", addrName)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}

		fmt.Println("New client connected", conn.RemoteAddr().String())

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Client disconnected")
			} else {
				fmt.Println("Error reading:", err.Error())
			}
			break
		}

		message := string(buf[:n])
		// message := unsafe.String(unsafe.SliceData(buf), n)
		// fmt.Println("Message received:", message)

		response := fmt.Sprintf("Hello, %s! You sent: '%s'", conn.RemoteAddr().String(), message)

		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error sending response:", err.Error())
			break
		}
	}
}
