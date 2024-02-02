package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	msgCount := 1000000
	pipe := "my.pipe"
	go server(pipe)

	time.Sleep(50 * time.Millisecond)

	f, err := os.OpenFile(pipe, os.O_WRONLY, 0777)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	message := "Hello, unix!"

	start := time.Now()
	for n := 0; n < msgCount; n++ {
		f.WriteString(message)
	}
	elapsed := time.Since(start)
	totaldata := int64(msgCount * len(message))

	fmt.Printf("Sent %d msg in %d ms; throughput %d msg/sec (%d MB/sec)\n",
		msgCount, elapsed.Milliseconds(),
		(int64(msgCount)*1000000000)/elapsed.Nanoseconds(),
		(totaldata*1000)/elapsed.Nanoseconds())
}

func server(pipe string) {
	fmt.Println("Opening pipe", pipe)

	err := unix.Mkfifo(pipe, 0777)
	if err != nil {
		fmt.Println("Error opening pipe", err.Error())
	}

	f, err := os.OpenFile(pipe, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("Error opening file", err.Error())
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	for {
		_, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error reading line", err.Error())
		}
	}
}
