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
	Last     string
	First    string
}

var tpl *template.Template
var dbUsers = map[string]user{}
var dbSession = map[string]string{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))

	// mock user
	bs, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	dbUsers["test@test.com"] = user{
		UserName: "test@test.com",
		Password: bs,
		Last:     "Test",
		First:    "Driven",
	}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)

	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Server starts...")
	http.ListenAndServe(":8080", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	if isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		pwd := req.FormValue("password")

		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "一致するユーザーが存在しません。", http.StatusForbidden)
			return
		}

		// decrypt
		if err := bcrypt.CompareHashAndPassword(u.Password, []byte(pwd)); err != nil {
			http.Error(w, "パスワードが一致しません。。", http.StatusForbidden)
			return
		}

		// create session
		id, _ := uuid.NewUUID()
		cookie := &http.Cookie{
			Name:     "session",
			Value:    id.String(),
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		dbSession[cookie.Value] = un

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func signup(w http.ResponseWriter, req *http.Request) {
	if isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		pwd := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")

		if _, ok := dbUsers[un]; ok {
			http.Error(w, "すでに存在するユーザーネームです。", http.StatusForbidden)
			return
		}

		// create session
		id, _ := uuid.NewUUID()
		cookie := &http.Cookie{
			Name:     "session",
			Value:    id.String(),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)

		// bcrypt
		bs, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		// store user
		dbSession[cookie.Value] = un
		dbUsers[un] = user{
			UserName: un,
			Password: bs,
			Last:     l,
			First:    f,
		}

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
