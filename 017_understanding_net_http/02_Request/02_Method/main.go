package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type Natto int

var tpl *template.Template

func (n Natto) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}

	data := struct {
		Method      string
		Submissions url.Values
	}{
		req.Method,
		req.Form,
	}
	tpl.ExecuteTemplate(w, "index.gohtml", data)
}

func init() {
	tpl = template.Must(template.ParseFiles("017_understanding_net_http/02_Request/02_Method/index.gohtml"))
}

func main() {
	var n Natto
	http.ListenAndServe(":8080", n)
}
