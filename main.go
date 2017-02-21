package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/liam-lai/ptt-alertor/hello"
	"github.com/liam-lai/ptt-alertor/ptt/board"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println(r.RemoteAddr + " visit: " + r.URL.Path)
	fmt.Fprintf(w, hello.HelloWorld())
}

func boardIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	b := new(board.Board)
	b.Name = strings.ToUpper(params.ByName("boardName"))
	fmt.Fprintf(w, "%s", b.IndexJSON())
}

func main() {
	fmt.Println("----Web Server Start on Port 9090----")

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/board/:boardName/articles", boardIndex)

	err := http.ListenAndServe(":9090", router)
	if err != nil {
		log.Fatal("ListenAndServer ", err)
	}
}
