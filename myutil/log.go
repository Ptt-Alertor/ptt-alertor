package myutil

import (
	"fmt"
	"runtime"

	"io"
	"io/ioutil"

	"strconv"

	log "github.com/meifamily/logrus"
)

func LogJSONEncode(err error, obj interface{}) {
	str := fmt.Sprintf("%#v", obj)
	log.WithError(err).WithFields(log.Fields{"object": str}).Warn("JSON Encode Error")
}

func LogJSONDecode(err error, data interface{}) {
	var bytes []byte

	if d, ok := data.(io.ReadCloser); ok {
		bytes, _ = ioutil.ReadAll(d)
	} else {
		bytes, _ = data.([]byte)
	}

	str := string(bytes)
	log.WithError(err).WithFields(log.Fields{"string": str}).Warn("JSON Decode Error")
}

func BasicRuntimeInfo() map[string]string {
	pc, fn, line, _ := runtime.Caller(1)
	info := map[string]string{
		"file":     fn,
		"function": runtime.FuncForPC(pc).Name(),
		"line":     strconv.Itoa(line),
	}
	return info
}
