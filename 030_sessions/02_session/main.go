package main

import (
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
)

type user struct {
	UserName string
	First    string
	Last     string
}

var tpl *template.Template
var dbUsers = map[string]user{}
var dbSessions = map[string]string{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/bar", bar)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Server starts...")
	http.ListenAndServe(":8080", nil)
}

func bar(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if hasErr(err) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	sid, ok := dbSessions[cookie.Value]
	if !ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
	
	u := dbUsers[sid]
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}

func index(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if hasErr(err) {
		id, _ := uuid.NewUUID()

		cookie = &http.Cookie{
			Name:     "session",
			Value:    id.String(),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
	}

	u := user{}
	if sid, ok := dbSessions[cookie.Value]; ok {
		u = dbUsers[sid]
	}

	// post, create user
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")
		u = user{
			UserName: un,
			First:    f,
			Last:     l,
		}

		dbUsers[un] = u
		dbSessions[cookie.Value] = un
	}

	tpl.ExecuteTemplate(w, "index.gohtml", u)

}

func hasErr(err error) bool {
	return err != nil
}
