package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/bar", bar)
	http.HandleFunc("/", foo)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Server Starts...")
	http.ListenAndServe(":8080", nil)
}

func bar(w http.ResponseWriter, req *http.Request) {
	log.Println("[bar] Request Method : ", req.Method)
	http.Redirect(w, req, "/", http.StatusMovedPermanently)
}

func foo(w http.ResponseWriter, req *http.Request) {
	log.Println("[foo] Request Method : ", req.Method)
}
