package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

func main() {

	http.HandleFunc("/", counter)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Server starts...")
	http.ListenAndServe(":8080", nil)
}

func counter(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("counter-cookie")
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "counter-cookie",
			Value: "0",
		}
	}

	counter, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Fatalln(err)
	}
	counter++
	cookie.Value = strconv.Itoa(counter)

	http.SetCookie(w, cookie)
	io.WriteString(w, cookie.Value)
}
