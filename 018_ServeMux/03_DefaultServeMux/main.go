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

	http.Handle("/cat", c)
	http.Handle("/dog", d)

	// nilを入れるとdefault serve muxが使用される。
	http.ListenAndServe(":8080", nil)
}
