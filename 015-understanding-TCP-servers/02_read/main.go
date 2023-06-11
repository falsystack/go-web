package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	li, err := net.Listen("tcp", ":8080")
	if hasErr(err) {
		log.Println(err)
	}
	defer li.Close()
	for {
		conn, err := li.Accept()
		if hasErr(err) {
			log.Println(err)
			continue
		}
		// go routineを利用
		// request毎 go routineを生成
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
	fmt.Println("ここまでたどりつかない")
}

func hasErr(err error) bool {
	return err != nil
}
