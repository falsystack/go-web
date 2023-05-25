package main

import (
	"io"
	"net/http"
)

type Natto int

func (n Natto) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/dog":
		io.WriteString(res, "doggy doggy doggy")
	case "/cat":
		io.WriteString(res, "kitty kitty kitty")
	}
}

func main() {
	var n Natto
	http.ListenAndServe(":8080", n)
}
