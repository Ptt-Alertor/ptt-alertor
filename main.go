package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/liam-lai/ptt-alertor/hello"
	"github.com/liam-lai/ptt-alertor/pttboard"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println(r.RemoteAddr + " visit: " + r.URL.Path)
	fmt.Fprintf(w, hello.HelloWorld())
}

func board(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	board := strings.ToUpper(params.ByName("boardName"))
	fmt.Fprintf(w, "%s", pttboard.FirstPage(board))
}

func main() {
	fmt.Println("----Web Server Start on Port 9090----")

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/board/:boardName/articles", board)

	err := http.ListenAndServe(":9090", router)
	if err != nil {
		log.Fatal("ListenAndServer ", err)
	}
}
