package article

import (
	"strings"

	log "github.com/Ptt-Alertor/logrus"
	"github.com/Ptt-Alertor/ptt-alertor/connections"
	"github.com/Ptt-Alertor/ptt-alertor/myutil"
	"github.com/garyburd/redigo/redis"
)

type Articles []Article

func (as Articles) List() []string {
	conn := connections.Redis()
	defer conn.Close()
	keys, err := redis.Strings(conn.Do("KEYS", prefix+"*"+subsSuffix))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	codes := make([]string, 0)
	for _, key := range keys {
		code := strings.TrimSuffix(strings.TrimPrefix(key, prefix), subsSuffix)
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

func (as Articles) StringWithPushSum() string {
	var content string
	for _, a := range as {
		content += "\r\n\r\n" + a.StringWithPushSum()
	}
	return content
}
