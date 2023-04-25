package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	li, err := net.Listen("tcp", ":8080")
	if hasErr(&err) {
		log.Fatalln(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if hasErr(&err) {
			log.Fatalln(err)
			continue
		}
		// go routineを利用
		// request事 go routineを生成
		go handle(conn)
	}

}

func handle(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
	}
	defer conn.Close()

	fmt.Println("ここまでたどり着けない")
}

func hasErr(err *error) bool {
	return err != nil
}
