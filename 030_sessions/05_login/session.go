package main

import (
	"github.com/google/uuid"
	"net/http"
)

func isLoggedIn(req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err != nil {
		return false
	}
	sid := dbSession[cookie.Value]
	_, ok := dbUsers[sid]
	return ok
}

func getUser(w http.ResponseWriter, req *http.Request) user {
	cookie, err := req.Cookie("session")
	if err != nil {
		id, _ := uuid.NewUUID()
		cookie = &http.Cookie{
			Name:     "session",
			Value:    id.String(),
			Path:     "/",
			HttpOnly: true,
		}
	}
	http.SetCookie(w, cookie)

	u := user{}
	if un, ok := dbSession[cookie.Value]; ok {
		u = dbUsers[un]
	}
	return u
}
