package main

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

type user struct {
	UserName string
	Password []byte
	First    string
	Last     string
}

var tpl *template.Template
var dbSessions = map[string]string{}
var dbUsers = map[string]user{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/signup", signup)

	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Server starts...")
	http.ListenAndServe(":8080", nil)
}

func signup(w http.ResponseWriter, req *http.Request) {
	if isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		pwd := req.FormValue("password")
		fn := req.FormValue("firstname")
		ln := req.FormValue("lastname")

		if _, ok := dbUsers[un]; ok {
			http.Error(w, "すでに存在するユーザーネームです。", http.StatusForbidden)
			return
		}

		id, _ := uuid.NewUUID()
		cookie := &http.Cookie{
			Name:     "session",
			Value:    id.String(),
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		dbSessions[cookie.Value] = un

		// 暗号化
		bs, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		u := user{
			UserName: un,
			Password: bs,
			First:    fn,
			Last:     ln,
		}
		dbUsers[un] = u

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func bar(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	if !isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}

func index(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	tpl.ExecuteTemplate(w, "index.gohtml", u)

}
