package main

import (
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {

	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("server starts...")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if err != nil {
		id, _ := uuid.NewUUID()
		cookie = &http.Cookie{
			Name:     "session",
			Value:    id.String(),
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}

	tpl.ExecuteTemplate(w, "index.gohtml", cookie.Value)
}
