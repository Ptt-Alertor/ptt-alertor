package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/liam-lai/ptt-alertor/hello"
)

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.RemoteAddr + " visit: " + r.RequestURI)
    fmt.Fprintf(w, hello.HelloWorld())
}

func main() {
    http.HandleFunc("/", index)
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal("ListenAndServer ", err)
    }
}
