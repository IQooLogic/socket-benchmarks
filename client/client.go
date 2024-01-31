package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	addrName := "/tmp/echo.sock"
	conn, err := net.Dial("unix", addrName)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

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

	select {}
}
