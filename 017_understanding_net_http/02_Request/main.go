package main

import (
	"html/template"
	"log"
	"net/http"
)

type Ramen int

var tpl *template.Template

func (r Ramen) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Print("ServeHTTP")
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	err = tpl.ExecuteTemplate(w, "index.gohtml", req.Form)
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	log.Print("init")
	tpl = template.Must(template.ParseFiles("017_understanding_net_http/02_Request/index.gohtml"))
}

func main() {
	log.Print("main")
	var r Ramen
	http.ListenAndServe(":8080", r)
}
