package controllers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/liam-lai/ptt-alertor/hello"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// t, err := template.ParseFiles("public/index.html")
	// if err != nil {
	// 	panic(err)
	// }
	// t.Execute(w, nil)
	fmt.Fprintf(w, hello.HelloWorld())
}
