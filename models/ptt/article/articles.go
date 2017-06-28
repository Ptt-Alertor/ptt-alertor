package article

import (
	"github.com/garyburd/redigo/redis"
	"github.com/liam-lai/ptt-alertor/connections"
	"github.com/liam-lai/ptt-alertor/myutil"
	log "github.com/meifamily/logrus"
)

type Articles []Article

func (as Articles) All() []*Article {
	codes := as.list()
	aps := make([]*Article, 0)
	for _, code := range codes {
		a := new(Article)
		a.Code = code
		aps = append(aps, a)
	}
	return aps
}

func (as Articles) list() []string {
	conn := connections.Redis()
	defer conn.Close()
	articleCodes, err := redis.Strings(conn.Do("SMEMBERS", "articles"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return articleCodes
}

func (as Articles) String() string {
	var content string
	for _, a := range as {
		content += "\r\n\r\n" + a.String()
	}
	return content
}
