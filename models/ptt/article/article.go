package article

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/liam-lai/ptt-alertor/connections"
	"github.com/liam-lai/ptt-alertor/myutil"
	log "github.com/meifamily/logrus"
)

const prefix = "article:"

type Article struct {
	ID       int
	Code     string `json:"code,omitempty"`
	Title    string
	Link     string
	Date     string `json:"Date,omitempty"`
	Author   string `json:"Author,omitempty"`
	PushList []Push `json:"PushList,omitempty"`
}

type Push struct {
	Tag        string
	UserID     string
	Content    string
	IPDateTime string
}

type ArticleAction interface {
	MatchKeyword(keyword string) bool
}

func (a Article) ParseID(Link string) (id int) {
	reg, err := regexp.Compile("https?://www.ptt.cc/bbs/.*/M\\.(\\d+)\\..*")
	if err != nil {
		log.Fatal(err)
	}
	id, err = strconv.Atoi(reg.FindStringSubmatch(Link)[1])
	if err != nil {
		id = 0
	}
	return id
}

func (a Article) MatchKeyword(keyword string) bool {
	if strings.Contains(keyword, "&") {
		keywords := strings.Split(keyword, "&")
		for _, keyword := range keywords {
			if !matchKeyword(a.Title, keyword) {
				return false
			}
		}
		return true
	}
	if strings.HasPrefix(keyword, "regexp:") {
		return matchRegex(a.Title, keyword)
	}
	return matchKeyword(a.Title, keyword)
}

func (a Article) Create() error {
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("SADD", "articles", a.Code)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (a Article) Save() error {
	conn := connections.Redis()
	defer conn.Close()

	articleJSON, err := json.Marshal(a)
	if err != nil {
		myutil.LogJSONEncode(err, a)
		return err
	}
	_, err = conn.Do("SET", prefix+a.Code, articleJSON)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func matchRegex(title string, regex string) bool {
	pattern := strings.TrimPrefix(regex, "regexp:")
	b, err := regexp.MatchString(pattern, title)
	if err != nil {
		return false
	}
	return b
}

func matchKeyword(title string, keyword string) bool {
	if strings.HasPrefix(keyword, "!") {
		excludeKeyword := strings.Trim(keyword, "!")
		return !containKeyword(title, excludeKeyword)
	}
	return containKeyword(title, keyword)
}

func containKeyword(title string, keyword string) bool {
	return strings.Contains(strings.ToLower(title), strings.ToLower(keyword))
}

func (a Article) String() string {
	return a.Title + "\r\n" + a.Link
}
