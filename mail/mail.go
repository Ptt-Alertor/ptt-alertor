package main

import (
	"fmt"
	"log"

	"github.com/liam-lai/ptt-alertor/myutil"

	"io/ioutil"

	"encoding/json"

	"gopkg.in/mailgun/mailgun-go.v1"
)

func main() {
	config := readConfig()

	domain := config["domain"]
	apiKey := config["apiKey"]
	publicAPIKey := config["publicAPIKey"]

	mg := mailgun.NewMailgun(domain, apiKey, publicAPIKey)

	body := "[公告] LoL 板 開始舉辦樂透!\r\n" +
		"https://www.ptt.cc/bbs/LoL/M.1486635540.A.605.html\r\n" +
		"\r\n" +
		"Send From PTT Alertor"

	message := mailgun.NewMessage(
		"PttAlertor@mg.dinolai.com",
		"[PTTAlertor] 在 LoL 版有一篇關於「樂透」的文章發表",
		body,
		"dinos80152@gmail.com")
	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}

func readConfig() map[string]string {
	projectRoot := myutil.ProjectRootPath()
	mailgunConfigJSON, err := ioutil.ReadFile(projectRoot + "/config/mailgun.json")
	if err != nil {
		log.Fatal(err)
	}

	var config map[string]string
	err = json.Unmarshal(mailgunConfigJSON, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
