package article

import (
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/myutil"
	log "github.com/meifamily/logrus"
)

type Articles []Article

func (as Articles) List() []string {
	conn := connections.Redis()
	defer conn.Close()
	keys, err := redis.Strings(conn.Do("KEYS", prefix+"*"+detailSuffix))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	codes := make([]string, 0)
	for _, key := range keys {
		code := strings.TrimSuffix(strings.TrimPrefix(key, prefix), detailSuffix)
		codes = append(codes, code)
	}
	return codes
}

func (as Articles) String() string {
	var content string
	for _, a := range as {
		content += "\r\n\r\n" + a.String()
	}
	return content
}
