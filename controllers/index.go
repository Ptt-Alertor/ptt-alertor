package controllers

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/models/top"
	"github.com/meifamily/ptt-alertor/shorturl"
)

var templates = template.Must(template.ParseFiles("public/docs.html", "public/top.html", "public/telegram.html", "public/messenger.html", "public/line.html", "public/tpls/head.tpl", "public/tpls/header.tpl", "public/tpls/command.tpl", "public/tpls/footer.tpl", "public/tpls/script.tpl"))

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := templates.ExecuteTemplate(w, "line.html", struct{ URI string }{"line"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LineIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := templates.ExecuteTemplate(w, "line.html", struct{ URI string }{"line"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MessengerIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := templates.ExecuteTemplate(w, "messenger.html", struct{ URI string }{"messenger"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TelegramIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := templates.ExecuteTemplate(w, "telegram.html", struct{ URI string }{"telegram"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Top(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	count := 100
	keywords := top.ListKeywordWithScore(count)
	authors := top.ListAuthorWithScore(count)
	pushsum := top.ListPushSumWithScore(count)
	data := struct {
		URI      string
		Keywords top.WordOrders
		Authors  top.WordOrders
		PushSum  top.WordOrders
	}{
		"top",
		keywords,
		authors,
		pushsum,
	}
	err := templates.ExecuteTemplate(w, "top.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Docs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := templates.ExecuteTemplate(w, "docs.html", struct{ URI string }{"docs"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Redirect(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	sha1 := params.ByName("sha1")
	url := shorturl.Original(sha1)
	if url != "" {
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	} else {
		t, err := template.ParseFiles("public/404.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		t.Execute(w, nil)
	}
}
