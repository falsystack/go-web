package main

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
)

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/intro?charset=utf8")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", index)

	http.HandleFunc("/create", create)
	http.HandleFunc("/drop", drop)

	http.HandleFunc("/insert", insert)
	http.HandleFunc("/read", read)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", del)

	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("server starts...")
	http.ListenAndServe(":8080", nil)
}

func drop(w http.ResponseWriter, req *http.Request) {
	stmt, err := db.Prepare("drop table customer;")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		log.Println(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintln(w, "drop table", rows)
}

func del(w http.ResponseWriter, req *http.Request) {
	stmt, err := db.Prepare("delete from customer where name='Line'")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		log.Println(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintln(w, "deleted record", rows)
}

func update(w http.ResponseWriter, req *http.Request) {
	stmt, err := db.Prepare("update customer set name='Line' where name='Kakao'")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		log.Println(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintln(w, "updated record", rows)
}

func read(w http.ResponseWriter, req *http.Request) {
	// Queryは主にselect文に使われる。
	rows, err := db.Query("select * from customer;")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	names := make([]string, 0)
	for rows.Next() {
		name := ""
		rows.Scan(&name)
		names = append(names, name)
	}

	fmt.Fprintln(w, "Retrived record: ", names)
}

func insert(w http.ResponseWriter, req *http.Request) {
	stmt, err := db.Prepare("insert into customer values ('Kakao');")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		log.Println(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintln(w, "Inserted record", rows)
}

func create(w http.ResponseWriter, req *http.Request) {
	stmt, err := db.Prepare("create table customer (name varchar(20));")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		log.Println(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintln(w, "Created table customer", rows)
}

func index(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "at index")
}
