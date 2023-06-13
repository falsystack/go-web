package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", foo)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func foo(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.Println(err)
	}

	// 同じキーに値が複数ある場合
	form := req.Form // map[string][]string
	values := form["q"]
	fmt.Println(values[0], values[1])

	query := req.FormValue("q")
	fmt.Fprintln(w, "Do my Search: "+query)
}
