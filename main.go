package main

import (
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/liam-lai/ptt-alertor/hello"
	"github.com/liam-lai/ptt-alertor/pttboard"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr + " visit: " + r.RequestURI)
	fmt.Fprintf(w, hello.HelloWorld())
}

func board(w http.ResponseWriter, r *http.Request) {
	board := strings.ToUpper(r.FormValue("board"))
	fmt.Fprintf(w, "%s", pttboard.FirstPage(board))
}

func main() {
	fmt.Println("----Web Server Start on Port 9090----")
	http.HandleFunc("/", index)
	http.HandleFunc("/articles", board)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServer ", err)
	}
}
