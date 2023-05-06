package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type Oden int

var tpl *template.Template

func (o Oden) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}

	data := struct {
		Method      string
		URL         *url.URL
		Submissions map[string][]string
		Header      http.Header
	}{
		req.Method,
		req.URL,
		req.Form,
		req.Header,
	}
	tpl.ExecuteTemplate(w, "index.gohtml", data)
}

func init() {
	tpl = template.Must(template.ParseFiles("017_understanding_net_http/02_Request/04_Header/index.gohtml"))
}

func main() {
	//var o Oden
	http.ListenAndServe(
		":8080",
		new(Oden),
	)
}
