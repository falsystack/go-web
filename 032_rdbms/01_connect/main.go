package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/intro?charset=utf8")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello SQL")
}
