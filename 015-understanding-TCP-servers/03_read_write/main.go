package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

/*
# on terminal
telnet localhost 8080

# telnet 終了方法
^] -> quit
*/

func main() {
	li, err := net.Listen("tcp", ":8080")
	if hasErr(err) {
		log.Fatalln(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if hasErr(err) {
			log.Println(err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	// deadlineを設定することもできる
	// err := conn.SetDeadline(time.Now().Add(10 * time.Second))
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		fmt.Fprintf(conn, "貴方の声が聞こえます：%s\n", ln)
	}
	defer conn.Close()

	fmt.Println("ここまでたどり着けない")
}

func hasErr(err error) bool {
	return err != nil
}
