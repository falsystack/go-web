package main

import (
	"io"
	"net/http"
)

type hotdog int

func (d hotdog) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "dog dog dog")
}

type hotcat int

func (c hotcat) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "cat cat cat")
}

func main() {
	var c hotcat
	var d hotdog

	mux := http.NewServeMux()
	mux.Handle("/cat", c)
	mux.Handle("/dog/", d)

	http.ListenAndServe(":8080", mux)
}
