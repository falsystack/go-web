package main

import (
	"fmt"
	"io"
	"net/http"
)

// use http.Handle
func main() {
	http.Handle("/", http.HandlerFunc(root))
	http.Handle("/me/", http.HandlerFunc(me))
	http.Handle("/dog/", http.HandlerFunc(dog))

	fmt.Print("start server on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func root(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "root, 02")
}

func dog(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "dog, 02")
}

func me(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Yun Kunwoong, 02")
}
