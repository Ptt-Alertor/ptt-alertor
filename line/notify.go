package line

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	log "github.com/meifamily/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/liam-lai/ptt-alertor/myutil"

	"encoding/json"

	"sort"

	"strconv"

	"net/url"

	user "github.com/liam-lai/ptt-alertor/models/user/redis"
)

const notifyBotHost string = "https://notify-bot.line.me"
const notifyAPIHost string = "https://notify-api.line.me"

var params map[string]string
var clientID string
var clientSecret string
var redirectURI string

func init() {
	clientID = config["clientID"]
	clientSecret = config["clientSecret"]
	appConfig := myutil.Config("app")
	redirectURI = appConfig["host"] + "/line/notify/callback"
}

func buildQueryString(params map[string]string) (query string) {
	var keys []string
	for key, _ := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		query += fmt.Sprintf("%s=%s&", key, params[key])
	}
	return query
}

func getAuthorizeURL(lineID string) string {
	var uri = "/oauth/authorize"
	params = map[string]string{
		"response_type": "code",
		"client_id":     clientID,
		"redirect_uri":  redirectURI,
		"scope":         "notify",
		"state":         lineID,
		"response_mode": "form_post",
	}
	query := buildQueryString(params)
	return fmt.Sprintf("%s%s?%s", notifyBotHost, uri, query)
}

// CatchCallback accept line notify post request to get user code
func CatchCallback(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.FormValue("error") != "" {
		log.Error(r.FormValue("state"), r.FormValue("error"), r.FormValue("error_description"))
	}
	code := r.FormValue("code")
	lineID := r.FormValue("state")
	u := new(user.User).Find(lineID)
	accessToken, err := fetchAccessToken(code)
	if err != nil {
		log.WithError(err).Error("Fetch Access Token Failed")
	}
	u.Profile.LineAccessToken = accessToken
	err = u.Update()
	if err != nil {
		log.WithError(err).Error("User Update Failed")
	}

	if err != nil {
		Notify(accessToken, "\n連結 LINE Notify 失敗。\n請至 Ptt Alertor LINE 主頁回報區留言。")
	} else {
		Notify(accessToken, "\n請回到 Ptt Alertor 輸入「指令」查看相關功能。\nPtt Alertor: 設定看板、作者、關鍵字\nLINE Notify: 最新文章通知")
	}

	t, err := template.ParseFiles("public/notify.html")
	if err != nil {
		log.WithError(err).Error("Show notify.html Failed")
	}
	t.Execute(w, nil)
}

func fetchAccessToken(code string) (string, error) {
	type responseBody struct {
		AccessToken string `json:"access_token"`
	}
	uri := "/oauth/token"
	params = map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  redirectURI,
		"client_id":     clientID,
		"client_secret": clientSecret,
	}
	body := buildQueryString(params)
	r, err := http.Post(notifyBotHost+uri, "application/x-www-form-urlencoded", bytes.NewBufferString(body))
	if err != nil {
		log.WithError(err).Error("Post Error")
	}
	if r.StatusCode != http.StatusOK {
		err := errors.New("Get Line Access Token Error, StatusCode:" + strconv.Itoa(r.StatusCode))
		log.WithError(err).Error()
		return "", err
	}
	var rspBody responseBody
	err = json.NewDecoder(r.Body).Decode(&rspBody)
	if err != nil {
		log.WithError(err).Error("Decode Line Access Token Error")
		return "", err
	}
	return rspBody.AccessToken, nil
}

func Notify(accessToken string, message string) {
	uri := "/api/notify"
	queryStr := url.Values{}
	queryStr.Add("message", message)
	encodeQueryStr := queryStr.Encode()
	pr, err := http.NewRequest("POST", notifyAPIHost+uri, bytes.NewBufferString(encodeQueryStr))
	pr.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	pr.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	r, err := client.Do(pr)
	if err != nil {
		log.WithError(err).Error("Notify Request Failed")
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(r.Body)
		log.Error(r.Status, string(data))
	}
}
