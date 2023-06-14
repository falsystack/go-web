package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", set)
	http.HandleFunc("/read", read)
	http.HandleFunc("/readAll", readAll)
	http.HandleFunc("/abundance", abundance)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Server starts...")
	http.ListenAndServe(":8080", nil)
}

func abundance(w http.ResponseWriter, req *http.Request) {
	cookies := []*http.Cookie{
		&http.Cookie{
			Name:  "c1",
			Value: "test cookie c1",
		},
		&http.Cookie{
			Name:  "c2",
			Value: "test cookie c2",
		},
	}

	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
	log.Println("[abundance] cookies written")
}

func read(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("my-cookie")
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	log.Println("[read] your cookie : ", cookie)
}

func readAll(w http.ResponseWriter, req *http.Request) {
	cookies := req.Cookies()
	if cookies == nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	for _, cookie := range cookies {
		log.Println("[readAll] your cookie: ", cookie)
	}
}

func set(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "my-cookie",
		Value: "chocopie",
		Path:  "/",
	})
	log.Println("[set] cookie written")
}
