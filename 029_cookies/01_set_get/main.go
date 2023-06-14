package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", set)
	http.HandleFunc("/read", read)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Server starts...")
	http.ListenAndServe(":8080", nil)
}

func read(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("my-cookie")
	if err != nil {
		// http.StatusText(400)もある。
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("あなたのクッキーです。", cookie.Name, cookie.Value)
}

func set(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "my-cookie",
		Value: "choco pie",
		Path:  "/",
	})
	log.Println("[set] cookie written")
}
