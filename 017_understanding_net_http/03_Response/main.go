package main

import (
	"fmt"
	"net/http"
)

type Sushi int

func (s Sushi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Yun-Kunwoong", "This is custom header")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	//w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "<h1>コードを書きます</h1>")
}

func main() {
	http.ListenAndServe(":8080", new(Sushi))
}
