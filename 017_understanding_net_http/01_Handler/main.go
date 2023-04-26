package main

import (
	"fmt"
	"net/http"
)

type Ramen string

func (r Ramen) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// writerを利用してresponse bodyに文字そのまま出力
	fmt.Fprintf(w, "何でも入力してください。")
}

func main() {
	var r Ramen
	http.ListenAndServe(":8080", r)
}
