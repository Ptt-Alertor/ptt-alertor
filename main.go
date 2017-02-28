package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/liam-lai/ptt-alertor/hello"
	"github.com/liam-lai/ptt-alertor/jobs"
	"github.com/liam-lai/ptt-alertor/ptt/board"
	"github.com/robfig/cron"
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
	fmt.Println("----Start Jobs----")
	startJobs()
	// Web Server
	fmt.Println("----Web Server Start on Port 9090----")

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/board/:boardName/articles", boardIndex)

	err := http.ListenAndServe(":9090", router)
	if err != nil {
		log.Fatal("ListenAndServer ", err)
	}

}

func startJobs() {
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		new(jobs.Message).Run()
		jobs.NewFetcher().Run()
	})
	c.AddJob("@every 1h", new(jobs.GenBoards))
	c.Start()
}
