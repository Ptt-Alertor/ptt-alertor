package controllers

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/liam-lai/ptt-alertor/myutil"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, _ := template.ParseFiles(myutil.PublicPath() + "/index.html")
	t.Execute(w, nil)
}
