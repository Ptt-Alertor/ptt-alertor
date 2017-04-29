package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/robfig/cron"

	ctrlr "github.com/liam-lai/ptt-alertor/controllers"
	"github.com/liam-lai/ptt-alertor/jobs"
	"github.com/liam-lai/ptt-alertor/myutil"
)

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
	router.NotFound = http.FileServer(http.Dir("public"))
	router.GET("/", ctrlr.Index)
	router.GET("/boards/:boardName/articles", ctrlr.BoardIndex)

	// users apis
	router.GET("/users/:account", ctrlr.UserFind)
	router.POST("/users", ctrlr.UserCreate)
	router.PUT("/users/:account", ctrlr.UserModify)

	router.POST("/line/callback", ctrlr.LineCallback)

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
