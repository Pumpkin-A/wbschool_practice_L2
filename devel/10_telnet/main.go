package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func getDataFromUser() (string, time.Duration) {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout in seconds")
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		log.Fatalf("usage: go-telnet --timeout=10s host port")
	}
	host := args[0]
	port := args[1]

	return host + ":" + port, *timeout
}

func main() {
	client := NewTelnetClient(getDataFromUser())
	if err := client.initConnection(); err != nil {
		log.Fatalf("error while connecting %v", err)
	}
	defer func() {
		if err := client.closeConnection(); err != nil {
			log.Fatalf("error while closing conn")
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	errors := make(chan error, 1)

	go func() {
		for {
			err := client.receieveMsg()
			if err != nil {
				errors <- err
				return
			}
		}
	}()

	go func() {
		// Чтобы ловить EOF
		inputReader := bufio.NewReader(os.Stdin)
		for {
			line, err := inputReader.ReadString('\n')
			if err != nil {
				errors <- err
				return
			}
			if err := client.sendMsg(line); err != nil {
				errors <- err
				return
			}
		}
	}()

	closed := make(chan interface{})
	go func() {
		defer close(closed)
		for {
			select {
			case <-signals:
				return
			case err := <-errors:
				fmt.Println(err)
				if err != nil {
					return
				}
			default:
				continue
			}
		}
	}()

	<-closed
}
