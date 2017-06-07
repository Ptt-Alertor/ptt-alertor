package controllers

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("public/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}
