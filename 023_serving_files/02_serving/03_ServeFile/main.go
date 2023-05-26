package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", dog)
	http.HandleFunc("/toby.jpg", dogPic)

	fmt.Println("サーバを動きます...")

	http.ListenAndServe(":8080", nil)
}

func dog(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Println("[dog] request dogPic")
	io.WriteString(w, `<img src="toby.jpg">`)
}

func dogPic(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[dogPic] serving picture")
	http.ServeFile(w, req, "toby.jpg")
}
