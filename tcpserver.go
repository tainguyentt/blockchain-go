package main

import (
	"log"
	"net"
	"os"
)

func runTcpServer(connHandler func(net.Conn)) {
	tcpPort := os.Getenv("ADDR")
	server, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TCP  Server Listening on port :", tcpPort)

	defer server.Close()

	// infinite loop to accept new connections
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// new go routine
		go connHandler(conn)
	}
}
