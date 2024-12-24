package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func usage() {
	log.Printf("Usage: ./go-telnet [OPTION]... HOST PORT \n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	var timeout = flag.Duration("timeout", 10*time.Second, "timeout")

	var showHelp = flag.Bool("h", false, "Show help message")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	args := flag.Args()

	if len(args) != 2 {
		showUsageAndExit(1)
	}

	var host = args[0]
	var port = args[1]

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), *timeout)

	if err != nil {
		time.Sleep(*timeout)
		log.Fatalf("go-telnet: %s", err)
	}

	osSignals := make(chan os.Signal, 1)
	listenErr := make(chan error, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go req(conn, listenErr, osSignals)
	go resp(conn, listenErr, osSignals)

	select {
	case <-osSignals:
		conn.Close()
	case err = <-listenErr:
		if err != nil {
			log.Fatalf("go-telnet: %s", err)
		}
	}
}

func req(conn net.Conn, listenErr chan<- error, osSignals chan<- os.Signal) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				osSignals <- syscall.Signal(syscall.SIGQUIT)
				return
			}
			listenErr <- err
		}

		fmt.Fprint(conn, text+"\n")
	}
}

func resp(conn net.Conn, listenErr chan<- error, osSignals chan<- os.Signal) {
	for {
		reader := bufio.NewReader(conn)
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				osSignals <- syscall.Signal(syscall.SIGQUIT)
				return
			}
			listenErr <- err
		}

		fmt.Print(text)
	}
}
