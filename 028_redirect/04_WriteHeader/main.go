package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {

	http.HandleFunc("/barred", barred)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/", foo)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Server starts...")
	http.ListenAndServe(":8080", nil)
}

func bar(w http.ResponseWriter, req *http.Request) {
	log.Println("[bar] Request Method : ", req.Method)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func barred(w http.ResponseWriter, req *http.Request) {
	log.Println("[barred] Request Method : ", req.Method)
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func foo(w http.ResponseWriter, req *http.Request) {
	log.Println("[foo] Request Method : ", req.Method)
}
