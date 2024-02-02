package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	// err := syscall.Mkfifo("my.pipe", 0777)
	//
	// if err != nil {
	// 	fmt.Println("Error opening pipe", err.Error())
	// }
	// f, err := os.Open("my.pipe")
	// if err != nil {
	// 	fmt.Println("Error opening file", err.Error())
	// }
	// defer f.Close()
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		pprof.StopCPUProfile()
		os.Exit(1)
	}()

	t := time.NewTimer(2 * time.Minute)
	go func() {
		<-t.C
		fmt.Println("time is up!")
		t.Stop()
		pprof.StopCPUProfile()
		os.Exit(0)
	}()

	pipe := "my.pipe"

	fmt.Println("Opening pipe", pipe)

	err = unix.Mkfifo(pipe, 0777)
	if err != nil {
		fmt.Println("Error opening pipe", err.Error())
	}

	f, err = os.OpenFile(pipe, os.O_RDONLY, os.ModeNamedPipe)
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
		// else {
		// 	fmt.Print("load string:" + string(line))
		// }
	}
}
