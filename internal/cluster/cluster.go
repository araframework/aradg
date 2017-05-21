package cluster

import (
	"log"
	"net"
)

var listener net.Listener

func Start() {
	go listen()
	go join()
}

func listen() {
	// Listen on TCP port all interfaces.
	listener, err := net.Listen("tcp", ":2800")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func join() {
	conn, err := net.Dial("tcp", "127.0.0.1:2801")
	if err != nil {
		log.Fatal(err)
	}

	conn.
}

func Stop() {
	if listener != nil {
		listener.Close()
	}
}

func handleConnection(conn net.Conn) {

}
