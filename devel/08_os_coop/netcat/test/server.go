package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// listen to incoming udp packets
	pc, err := net.ListenPacket("udp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go serve(pc, addr, buf[:n])
	}

}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	// 0 - 1: ID
	// 2: QR(1): Opcode(4)
	// buf[2] |= 0x80 // Set QR bit
	fmt.Println(buf)

	pc.WriteTo(buf, addr)
}
