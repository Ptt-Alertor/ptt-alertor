package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/robfig/cron"

	"github.com/liam-lai/ptt-alertor/hello"
	"github.com/liam-lai/ptt-alertor/jobs"
	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
	"github.com/liam-lai/ptt-alertor/myutil"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, hello.HelloWorld())
}

func boardIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	bd := new(board.Board)
	bd.Name = strings.ToUpper(params.ByName("boardName"))
	articles := bd.OnlineArticles()
	articlesJSON, err := json.Marshal(articles)
	if err != nil {
		myutil.LogJSONEncode(err, articles)
	}
	fmt.Fprintf(w, "%s", articlesJSON)
}

func userFind(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	u := new(user.User).Find(params.ByName("account"))
	uJSON, err := json.Marshal(u)
	if err != nil {
		myutil.LogJSONEncode(err, u)
	}
	fmt.Fprintf(w, "%s", uJSON)
}

func userCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	u := new(user.User)
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		myutil.LogJSONDecode(err, r.Body)
		http.Error(w, "not a json valid format", 400)
	}
	err = u.Save()
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
}

func userModify(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	u := new(user.User)
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		myutil.LogJSONDecode(err, r.Body)
		http.Error(w, "not a json valid format", 400)
	}

	if u.Profile.Account != account {
		http.Error(w, "account does not match", 400)
	}

	err = u.Update()
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

}

type myRouter struct {
	httprouter.Router
}

func (mr myRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"IP":  r.RemoteAddr,
		"URI": r.URL.Path,
	}).Info("visit")
	mr.Router.ServeHTTP(w, r)
}

func newRouter() *myRouter {
	return &myRouter{
		Router: *httprouter.New(),
	}
}

func basicAuth(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		user, password, hasAuth := r.BasicAuth()
		auth := myutil.Config("auth")
		if hasAuth && user == auth["user"] && password == auth["password"] {
			handle(w, r, params)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func main() {
	log.Info("Start Jobs")
	startJobs()
	// Web Server
	log.Info("Web Server Start on Port 9090")

	router := newRouter()
	router.GET("/", index)
	router.GET("/boards/:boardName/articles", boardIndex)

	// users apis
	router.GET("/users/:account", userFind)
	router.POST("/users", basicAuth(userCreate))
	router.PUT("/users/:account", basicAuth(userModify))

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

func init() {
	new(jobs.GenBoards).Run()
	jobs.NewFetcher().Run()
}
