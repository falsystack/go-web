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
var dbUsers = map[string]user{}
var dbSession = map[string]string{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	bs, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	dbUsers["test@test.com"] = user{
		UserName: "test@test.com",
		Password: bs,
		First:    "Driven",
		Last:     "Test",
	}
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

	cookie, _ := req.Cookie("session")
	delete(dbSession, cookie.Value)

	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1, // 0未満の場合破棄される。
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, req, "/login", http.StatusSeeOther)
	return
}

func login(w http.ResponseWriter, req *http.Request) {
	if isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	u := user{}
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		pwd := req.FormValue("password")

		// ユーザをチェック
		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "存在しないユーザです。", http.StatusForbidden)
			return
		}

		// パスワードのチェック
		if err := bcrypt.CompareHashAndPassword(u.Password, []byte(pwd)); err != nil {
			http.Error(w, "パスワードが一致しませんでした。", http.StatusForbidden)
			return
		}

		// session生成
		id, _ := uuid.NewUUID()
		session := &http.Cookie{
			Name:     "session",
			Value:    id.String(),
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, session)
		dbSession[session.Value] = un

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.gohtml", u)
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

		// 存在するユーザか確認
		_, ok := dbUsers[un]
		if ok {
			http.Error(w, "既に存在するユーザです。", http.StatusForbidden)
			return
		}

		// session生成
		id, _ := uuid.NewUUID()
		session := &http.Cookie{
			Name:     "session",
			Value:    id.String(),
			Path:     "/",
			HttpOnly: true,
		}

		// bcrypt
		bs, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// 保存
		http.SetCookie(w, session)
		dbSession[session.Value] = un
		dbUsers[un] = user{
			UserName: un,
			Password: bs,
			First:    f,
			Last:     l,
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
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}

func index(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	tpl.ExecuteTemplate(w, "index.gohtml", u)
}
