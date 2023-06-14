package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/set", set)
	http.HandleFunc("/read", read)
	http.HandleFunc("/expire", expire)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Server starts...")
	http.ListenAndServe(":8080", nil)
}

func expire(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("my-cookie")
	if err == http.ErrNoCookie {
		http.Redirect(w, req, "/set", http.StatusSeeOther)
		return
	}
	cookie.MaxAge = -1 // delete cookie
	http.SetCookie(w, cookie)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func read(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("my-cookie")
	if err == http.ErrNoCookie {
		http.Redirect(w, req, "/set", http.StatusSeeOther)
		return
	}
	fmt.Fprintf(w, `<h1>Your Cookie:<br>%v</h1><h1><a href="/expire">expire</a></h1>`, cookie)
}

func set(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "my-cookie",
		Value: "choco pie",
		Path:  "/",
	})
	fmt.Fprintln(w, `<h1><a href="/read">read</a></h1>`)
}

func index(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, `<h1><a href="/set">set a cookie</a></h1>`)
}
