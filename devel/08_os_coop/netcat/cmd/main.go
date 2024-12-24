package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// nc host port

func usage() {
	log.Printf("Usage: ./nc [OPTION] HOST PORT \n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	var protocol string = "tcp"
	var udp = flag.Bool("u", false, "Подключаться по UDP (вместо TCP)")

	var showHelp = flag.Bool("h", false, "Show help message")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	if *udp {
		protocol = "udp"
	}

	args := flag.Args()

	if len(args) != 2 {
		showUsageAndExit(1)
	}

	var host = args[0]
	var port = args[1]

	conn, err := net.Dial(protocol, net.JoinHostPort(host, port))
	if err != nil {
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
			listenErr <- err
		}

		fmt.Print(text)
	}
}
