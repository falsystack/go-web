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

	http.HandleFunc("/", foo)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/barred", barred)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Server Starts...")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func barred(w http.ResponseWriter, req *http.Request) {
	log.Println("[barred] Request method : ", req.Method)
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func bar(w http.ResponseWriter, req *http.Request) {
	log.Println("[bar] Request method : ", req.Method)
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

func foo(w http.ResponseWriter, req *http.Request) {
	log.Println("[foo] Request method : ", req.Method)
}
