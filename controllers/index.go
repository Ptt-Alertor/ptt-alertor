package controllers

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/models/top"
	"github.com/meifamily/ptt-alertor/shorturl"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("public/line.html", "public/tpls/head.tpl", "public/tpls/header.tpl", "public/tpls/command.tpl", "public/tpls/footer.tpl", "public/tpls/script.tpl")
	if err != nil {
		panic(err)
	}
	t.Execute(w, struct{ URI string }{""})
}

func LineIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("public/line.html", "public/tpls/head.tpl", "public/tpls/header.tpl", "public/tpls/command.tpl", "public/tpls/footer.tpl", "public/tpls/script.tpl")
	if err != nil {
		panic(err)
	}
	t.Execute(w, struct{ URI string }{"line"})
}

func MessengerIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("public/messenger.html", "public/tpls/head.tpl", "public/tpls/header.tpl", "public/tpls/command.tpl", "public/tpls/footer.tpl", "public/tpls/script.tpl")
	if err != nil {
		panic(err)
	}
	t.Execute(w, struct{ URI string }{"messenger"})
}

func TelegramIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("public/telegram.html", "public/tpls/head.tpl", "public/tpls/header.tpl", "public/tpls/command.tpl", "public/tpls/footer.tpl", "public/tpls/script.tpl")
	if err != nil {
		panic(err)
	}
	t.Execute(w, struct{ URI string }{"telegram"})
}

func Top(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	count := 100
	keywords := top.ListKeywordWithScore(count)
	authors := top.ListAuthorWithScore(count)
	pushsum := top.ListPushSumWithScore(count)
	t, err := template.ParseFiles("public/top.html", "public/tpls/head.tpl", "public/tpls/header.tpl", "public/tpls/footer.tpl", "public/tpls/script.tpl")
	if err != nil {
		panic(err)
	}
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
	t.Execute(w, data)
}

func Docs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("public/docs.html", "public/tpls/head.tpl", "public/tpls/script.tpl")
	if err != nil {
		panic(err)
	}
	t.Execute(w, struct{ URI string }{"docs"})
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
