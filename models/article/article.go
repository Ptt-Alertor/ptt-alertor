package article

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"time"

	"fmt"

	"github.com/garyburd/redigo/redis"
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/models/pushsum"
	"github.com/meifamily/ptt-alertor/myutil"
)

const prefix = "article:"
const detailSuffix = ":detail"
const subsSuffix = ":subs"

type Article struct {
	ID               int    `json:"ID,omitempty"`
	Code             string `json:"code,omitempty"`
	Title            string
	Link             string
	Date             string    `json:"Date,omitempty"`
	Author           string    `json:"Author,omitempty"`
	Comments         Comments  `json:"pushList,omitempty"` // rename json key to comments
	LastPushDateTime time.Time `json:"lastPushDateTime,omitempty"`
	Board            string    `json:"board,omitempty"`
	PushSum          int       `json:"pushSum,omitempty"`
}

type Comment struct {
	Tag      string
	UserID   string
	Content  string
	DateTime time.Time
}

func (c Comment) String() string {
	// 推 ChoDino: 推文推文
	return fmt.Sprintf("%s %s%s", c.Tag, c.UserID, c.Content)
}

type Comments []Comment

func (cs Comments) String() string {
	var content string
	for _, p := range cs {
		content += "\n" + p.String()
	}
	return content
}

func (a Article) ParseID(Link string) (id int) {
	reg, err := regexp.Compile("https?://www.ptt.cc/bbs/.*/[GM]\\.(\\d+)\\..*")
	if err != nil {
		log.Fatal(err)
	}
	strs := reg.FindStringSubmatch(Link)
	if len(strs) < 2 {
		return 0
	}
	id, err = strconv.Atoi(strs[1])
	if err != nil {
		return 0
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

// Exist check article exist or not
func (a Article) Exist() (bool, error) {
	conn := connections.Redis()
	defer conn.Close()

	bl, err := redis.Bool(conn.Do("HEXISTS", prefix+a.Code+detailSuffix, "board"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return bl, err
}

func (a Article) Find(code string) Article {
	conn := connections.Redis()
	defer conn.Close()

	aMap, err := redis.StringMap(conn.Do("HGETALL", prefix+code+detailSuffix))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	a.Board = aMap["board"]
	err = json.Unmarshal([]byte(aMap["content"]), &a)
	if err != nil {
		log.WithField("code", code).Error("Article Content Unmarshal Failed")
		myutil.LogJSONDecode(err, aMap["content"])
	}
	return a
}

func (a Article) Save() error {
	conn := connections.Redis()
	defer conn.Close()

	articleJSON, err := json.Marshal(a)
	if err != nil {
		myutil.LogJSONEncode(err, a)
		return err
	}
	_, err = conn.Do("HMSET", prefix+a.Code+detailSuffix, "board", a.Board, "content", articleJSON)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (a Article) Destroy() error {
	conn := connections.Redis()
	defer conn.Close()

	_, err := conn.Do("DEL", prefix+a.Code+detailSuffix, prefix+a.Code+subsSuffix)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (a Article) AddSubscriber(account string) error {
	conn := connections.Redis()
	defer conn.Close()

	_, err := conn.Do("SADD", prefix+a.Code+subsSuffix, account)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (a Article) Subscribers() ([]string, error) {
	conn := connections.Redis()
	defer conn.Close()

	accounts, err := redis.Strings(conn.Do("SMEMBERS", prefix+a.Code+subsSuffix))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return accounts, err
}

func (a Article) RemoveSubscriber(sub string) error {
	conn := connections.Redis()
	defer conn.Close()

	_, err := conn.Do("SREM", prefix+a.Code+subsSuffix, sub)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (a Article) String() string {
	return a.Title + "\r\n" + a.Link
}

func (a Article) StringWithPushSum() string {
	sumStr := strconv.Itoa(a.PushSum)
	if text, ok := pushsum.NumTextMap[a.PushSum]; ok {
		sumStr = text
	}
	return fmt.Sprintf("%s %s\r\n%s", sumStr, a.Title, a.Link)
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
