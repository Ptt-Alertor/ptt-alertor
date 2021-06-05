package line

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/models"
)

const notifyBotHost string = "https://notify-bot.line.me"
const notifyAPIHost string = "https://notify-api.line.me"

var (
	params       map[string]string
	clientID     = os.Getenv("LINE_CLIENT_ID")
	clientSecret = os.Getenv("LINE_CLIENT_SECRET")
	redirectURI  = os.Getenv("APP_HOST") + "/line/notify/callback"
)

func buildQueryString(params map[string]string) (query string) {
	var keys []string
	for key := range params {
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
		log.WithFields(log.Fields{
			"error":       r.FormValue("error"),
			"state":       r.FormValue("state"),
			"description": r.FormValue("error_description"),
		}).Error("Get LINE Notify Callback Failed")
	}

	code, lineID := r.FormValue("code"), r.FormValue("state")
	accessToken, err := fetchAccessToken(code)
	if err != nil {
		log.WithError(err).Error("Fetch Access Token Failed")
	}

	u := models.User.Find(lineID)
	u.Profile.LineAccessToken = accessToken
	if err := u.Update(); err != nil {
		log.WithError(err).Error("User Update Failed")
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

func checkLineAccessTokenExist(lineID string) bool {
	u := models.User.Find(lineID)
	if u.Profile.LineAccessToken == "" {
		return false
	}
	return true
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
		return
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(r.Body)
		log.WithFields(log.Fields{
			"status":   r.Status,
			"response": string(data),
		}).Error("LINE Notify Failed")
	}
}
