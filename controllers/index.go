package controllers

import (
	"html/template"
	"net/http"
	"strconv"

	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/models/counter"
	"github.com/meifamily/ptt-alertor/models/top"
	"github.com/meifamily/ptt-alertor/shorturl"
	"golang.org/x/net/websocket"
)

var tpls = []string{
	"public/docs.html",
	"public/top.html",
	"public/telegram.html",
	"public/messenger.html",
	"public/line.html",
	"public/tpls/head.tpl",
	"public/tpls/header.tpl",
	"public/tpls/slogan.tpl",
	"public/tpls/command.tpl",
	"public/tpls/counter.tpl",
	"public/tpls/footer.tpl",
	"public/tpls/script.tpl",
}

var templates = template.Must(template.ParseFiles(tpls...))

var host string

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	LineIndex(w, r, nil)
}

func count() (counterStrs []string) {
	count, err := counter.Alert()
	if err != nil {
		return nil
	}
	countStrs := strings.Split((strconv.Itoa(count)), "")
	for index, num := range countStrs {
		counterStrs = append(counterStrs, num)
		if backIndex := len(countStrs) - index; backIndex != 1 && backIndex%3 == 1 {
			counterStrs = append(counterStrs, ",")
		}
	}
	return counterStrs
}

func LineIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := templates.ExecuteTemplate(w, "line.html", struct {
		URI   string
		Host  string
		Count []string
	}{"line", r.Host, count()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MessengerIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := templates.ExecuteTemplate(w, "messenger.html", struct {
		URI   string
		Host  string
		Count []string
	}{"messenger", r.Host, count()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TelegramIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := templates.ExecuteTemplate(w, "telegram.html", struct {
		URI   string
		Host  string
		Count []string
	}{"telegram", r.Host, count()})
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

func WebSocket(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	websocket.Handler(CounterHandler).ServeHTTP(w, r)
}

func CounterHandler(ws *websocket.Conn) {
	conn := connections.Redis()
	defer conn.Close()
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe("alert-counter")
	defer psc.Unsubscribe("alert-counter")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			_, err := ws.Write([]byte(v.Data))
			if err != nil {
				ws.Close()
				psc.Unsubscribe("alert-counter")
				return
			}
		case redis.Subscription:
			// fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			return
		}
	}
}
