package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/", foo)

	log.Println("server starts...")
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log.Println(ctx)
	fmt.Fprintln(w, ctx)
}
