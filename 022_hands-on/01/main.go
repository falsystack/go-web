package main

import (
	"io"
	"net/http"
)

func root(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "root")
}

func dog(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "dog")
}

func me(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Yun Kunwoong")
}

func main() {

	http.HandleFunc("/", root)
	http.HandleFunc("/dog/", dog)
	http.HandleFunc("/me/", me)

	// use default serve mux
	http.ListenAndServe(":8080", nil)
}
