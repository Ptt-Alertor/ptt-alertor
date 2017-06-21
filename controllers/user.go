package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
	"github.com/liam-lai/ptt-alertor/myutil"
)

func UserFind(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	u := new(user.User).Find(params.ByName("account"))
	uJSON, err := json.Marshal(u)
	if err != nil {
		myutil.LogJSONEncode(err, u)
	}
	fmt.Fprintf(w, "%s", uJSON)
}

func UserAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	us := new(user.User).All()
	fmt.Fprintf(w, "總人數：%d\n", len(us))
	var line, messenger, subCount, boardCount, keywordCount, authorCount, idleUser int
	for _, u := range us {
		if u.Profile.Line != "" {
			line++
		}
		if u.Profile.Messenger != "" {
			messenger++
		}
		subCount = len(u.Subscribes)
		if subCount == 0 {
			idleUser++
		}
		boardCount += subCount
		for _, s := range u.Subscribes {
			keywordCount += len(s.Keywords)
			authorCount += len(s.Authors)
		}
		fmt.Fprintf(w, "%s\n", u.Profile.Account)
	}
	fmt.Fprintf(w, "LINE: %d, Messenger: %d, count(Board): %d, count(Keyword): %d, count(Author): %d, IdleUsers: %d", line, messenger, boardCount, keywordCount, authorCount, idleUser)
}

func UserCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

func UserModify(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
