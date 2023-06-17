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
	Role     string
}

var tpl *template.Template

var dbUsers = map[string]user{}
var dbSessions = map[string]string{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/", index)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	log.Println("server starts...")
	http.ListenAndServe(":8080", nil)
}

func logout(w http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	// get session
	cookie, _ := req.Cookie("session")

	// delete session
	delete(dbSessions, cookie.Value)

	// destroy cookie
	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, cookie)

	// redirect
	http.Redirect(w, req, "/", http.StatusSeeOther)
	return
}

func login(w http.ResponseWriter, req *http.Request) {
	if isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		pwd := req.FormValue("password")

		// check user
		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "ユーザーネームが間違いました。", http.StatusForbidden)
			return
		}

		// check password
		if err := bcrypt.CompareHashAndPassword(u.Password, []byte(pwd)); err != nil {
			http.Error(w, "パスワードが間違いました。", http.StatusForbidden)
			return
		}

		// recreate session
		id, _ := uuid.NewUUID()
		session := &http.Cookie{
			Name:     "session",
			Value:    id.String(),
			Path:     "/",
			HttpOnly: true,
		}
		dbSessions[session.Value] = un
		http.SetCookie(w, session)

		// redirect
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
		r := req.FormValue("role")

		// check user
		_, ok := dbUsers[un]
		if ok {
			http.Error(w, "既に存在するユーザーです。", http.StatusSeeOther)
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
		dbSessions[cookie.Value] = un

		// bcrpyt
		bs, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// create user
		dbUsers[un] = user{
			UserName: un,
			Password: bs,
			First:    f,
			Last:     l,
			Role:     r,
		}

		// redirect
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

	if u.Role != "007" {
		http.Error(w, "007要員だけ入れれます。", http.StatusForbidden)
		return
	}

	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}

func index(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	log.Println(u)
	tpl.ExecuteTemplate(w, "index.gohtml", u)
}
