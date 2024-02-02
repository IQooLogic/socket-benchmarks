package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	msgCount := 1000000
	addrName := "localhost:8080"
	go server(addrName)

	time.Sleep(50 * time.Millisecond)

	conn, err := net.Dial("tcp", addrName)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	message := "Hello, unix!"

	start := time.Now()
	for n := 0; n < msgCount; n++ {
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err.Error())
			os.Exit(1)
		}
	}
	elapsed := time.Since(start)
	totaldata := int64(msgCount * len(message))

	fmt.Printf("Sent %d msg in %d ms; throughput %d msg/sec (%d MB/sec)\n",
		msgCount, elapsed.Milliseconds(),
		(int64(msgCount)*1000000000)/elapsed.Nanoseconds(),
		(totaldata*1000)/elapsed.Nanoseconds())
}

func server(addrName string) {
	l, err := net.Listen("tcp", addrName)
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

		_ = string(buf[:n])
	}
}
