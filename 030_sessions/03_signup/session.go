package main

import "net/http"

func getUser(req *http.Request) user {
	u := user{}

	cookie, err := req.Cookie("session")
	if err != nil {
		return u
	}

	if sid, ok := dbSession[cookie.Value]; ok {
		u = dbUsers[sid]
	}

	return u
}

func isLoggedIn(req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err != nil {
		return false
	}

	sid := dbSession[cookie.Value]
	_, ok := dbUsers[sid]
	return ok
}
