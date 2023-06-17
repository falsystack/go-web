package main

import (
	"github.com/google/uuid"
	"net/http"
)

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
		http.SetCookie(w, cookie)
	}

	sid := dbSessions[cookie.Value]
	return dbUsers[sid]
}

func isLoggedIn(req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err != nil {
		return false
	}

	sid := dbSessions[cookie.Value]
	_, ok := dbUsers[sid]
	return ok
}
