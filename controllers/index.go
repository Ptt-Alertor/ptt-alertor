package controllers

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/liam-lai/ptt-alertor/shorturl"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("public/line.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func LineIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("public/line.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func MessengerIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("public/messenger.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

func Redirect(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	sha1 := params.ByName("sha1")
	url := shorturl.Original(sha1)
	if url != "" {
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	} else {
		t, err := template.ParseFiles("public/404.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, nil)
	}
}
