package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/models"
	"github.com/meifamily/ptt-alertor/models/user"
	"github.com/meifamily/ptt-alertor/myutil"
)

func UserFind(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	u := models.User().Find(params.ByName("account"))
	uJSON, err := json.Marshal(u)
	if err != nil {
		myutil.LogJSONEncode(err, u)
	}
	fmt.Fprintf(w, "%s", uJSON)
}

func UserAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	us := models.User().All()

	data := struct {
		Total, Line, Messenger, Telegram, IdleUser, BlockUser         int
		SubCount, BoardCount, KeywordCount, AuthorCount, PushSumCount int
		User, Room, Group                                             int
		Users                                                         []*user.User
	}{}
	data.Users = us
	data.Total = len(us)
	for _, u := range us {
		if !u.Enable {
			data.BlockUser++
		}
		if u.Profile.Line != "" {
			data.Line++
		}
		if u.Profile.Messenger != "" {
			data.Messenger++
		}
		if u.Profile.Telegram != "" {
			data.Telegram++
		}
		switch u.Profile.Type {
		case "user", "":
			data.User++
		case "room":
			data.Room++
		case "group":
			data.Group++
		}
		data.SubCount = len(u.Subscribes)
		if data.SubCount == 0 {
			data.IdleUser++
		}
		data.BoardCount += data.SubCount
		for _, s := range u.Subscribes {
			data.KeywordCount += len(s.Keywords)
			data.AuthorCount += len(s.Authors)
			if s.PushSum.Up != 0 || s.PushSum.Down != 0 {
				data.PushSumCount++
			}
		}
	}
	t, err := template.ParseFiles("public/user.tpl")
	if err != nil {
		panic(err)
	}
	t.Execute(w, data)
}

func UserCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	u := models.User()
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		myutil.LogJSONDecode(err, r.Body)
		http.Error(w, "not a json valid format", 400)
	}
	if err := u.Save(); err != nil {
		http.Error(w, err.Error(), 400)
	}
}

func UserModify(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := params.ByName("account")
	u := models.User()
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		myutil.LogJSONDecode(err, r.Body)
		http.Error(w, "not a json valid format", 400)
	}

	if u.Profile.Account != account {
		http.Error(w, "account does not match", 400)
	}

	if err := u.Update(); err != nil {
		http.Error(w, err.Error(), 400)
	}

}
