package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

/*
net packageを通してTCPサーバーを理解
*/
func main() {
	// TCP Server
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		io.WriteString(conn, "Hello from TCP server")
		fmt.Fprintln(conn, "How is your day?")
		fmt.Fprintf(conn, "Well, I hope!")

		conn.Close()
	}
}
