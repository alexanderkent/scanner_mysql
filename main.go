package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func scan(host string, port int) {

	fmt.Printf(fmt.Sprintf("%s\n", strings.Repeat("-", 70)))

	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Failed to connect %s\n", err.Error())
		return
	}

	handshakePacket := &InitialHandshakePacket{}
	err = handshakePacket.Decode(conn)
	if err != nil {
		log.Printf("Failed to decode packet: %s\n", err.Error())
		return
	}

	fmt.Printf("%s\n", target)
	fmt.Printf(handshakePacket.String())
}

func main() {

	/*
		Single target mode
		./scanner host ip
	*/
	if len(os.Args) == 3 {

		flag.Parse()
		host := flag.Arg(0)
		port, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			os.Exit(-1)
		}
		scan(host, port)
		return
	}

	/*
		Run against test chassis (containers)
	*/
	for i := 3306; i <= 3311; i++ {
		scan("localhost", i)
	}
}
