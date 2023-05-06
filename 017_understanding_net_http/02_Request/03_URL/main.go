package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type Udon int

var tpl *template.Template

func (u Udon) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}

	data := struct {
		Method      string
		URL         *url.URL
		Submissions url.Values
	}{
		req.Method,
		req.URL,
		req.Form,
	}
	tpl.ExecuteTemplate(w, "index.gohtml", data)
}

func init() {
	tpl = template.Must(template.ParseFiles("017_understanding_net_http/02_Request/03_URL/index.gohtml"))
}

func main() {
	var u Udon
	http.ListenAndServe(":8080", u)
}
