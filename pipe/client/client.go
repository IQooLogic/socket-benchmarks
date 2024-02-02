package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// pprof.StopCPUProfile()
		os.Exit(1)
	}()

	t := time.NewTimer(2 * time.Minute)
	go func() {
		<-t.C
		fmt.Println("time is up!")
		t.Stop()
		// pprof.StopCPUProfile()
		os.Exit(0)
	}()

	pipe := "../server/my.pipe"
	f, err := os.OpenFile(pipe, os.O_WRONLY, 0777)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	i := 0
	for {
		fmt.Println("write string to named pipe file.")
		f.WriteString(fmt.Sprintf("test write times:%d\n", i))
		i++
	}
}
