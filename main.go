package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/liam-lai/ptt-alertor/hello"
	"github.com/liam-lai/ptt-alertor/pttboard"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr + " visit: " + r.RequestURI)
	fmt.Fprintf(w, hello.HelloWorld())
}

func board(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", pttboard.FirstPage("FREE_BOX"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/board", board)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServer ", err)
	}
}
