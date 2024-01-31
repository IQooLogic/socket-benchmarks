package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	addrName := "/tmp/echo.sock"
	os.Remove(addrName)

	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: addrName, Net: "unix"})
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}

	defer l.Close()

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
		fmt.Println("Message received:", message)

		response := fmt.Sprintf("Hello, %s! You sent: '%s'", conn.RemoteAddr().String(), message)

		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error sending response:", err.Error())
			break
		}
	}
}
