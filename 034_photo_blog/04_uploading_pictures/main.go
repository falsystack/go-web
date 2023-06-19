package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("server starts...")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
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

	if req.Method == http.MethodPost {

		mf, fh, err := req.FormFile("nf")
		if err != nil {
			log.Println(err)
		}
		defer mf.Close()

		ext := strings.Split(fh.Filename, ".")[1]

		// create sha1
		hash := sha1.New()
		io.Copy(hash, mf)
		fname := fmt.Sprintf("%x", hash.Sum(nil)) + "." + ext
		log.Println(fname, "file name")

		// create new file
		wd, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		path := filepath.Join(wd, "public", "pics", fname)
		log.Println(path, "path")
		nf, err := os.Create(path)
		if err != nil {
			log.Println(err)
		}
		defer nf.Close()

		// copy
		mf.Seek(0, 0)
		io.Copy(nf, mf)

		s := cookie.Value
		if !strings.Contains(s, fname) {
			s += "|" + fname
		}
		cookie.Value = s
		http.SetCookie(w, cookie)
	}

	xs := strings.Split(cookie.Value, "|")
	tpl.ExecuteTemplate(w, "index.gohtml", xs)

}
