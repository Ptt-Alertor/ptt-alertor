package controllers

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, _ := template.ParseFiles("public/index.html")
	t.Execute(w, nil)
}
