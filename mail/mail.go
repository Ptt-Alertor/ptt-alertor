package mail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/liam-lai/ptt-alertor/models/ptt/article"
	"github.com/liam-lai/ptt-alertor/myutil"
	"gopkg.in/mailgun/mailgun-go.v1"
)

type Mail struct {
	Title
	Body
	Receiver string
}

type Title struct {
	BoardName       string
	Keyword         string
	articleQuantity int
}

type Body struct {
	Articles []article.Article
}

func (title Title) String() string {
	return "[PTTAlertor] 在 " + title.BoardName + " 版有 " + strconv.Itoa(title.articleQuantity) + " 篇關於「" + title.Keyword + "」的文章發表"
}

func (body Body) String() string {
	var content string
	for _, article := range body.Articles {
		content += article.Title + "\r\n" +
			"https://www.ptt.cc" + article.Link + "\r\n" +
			"\r\n"
	}
	return content + "Send From PTT Alertor"
}

func (mail Mail) Send() {
	mg := newMailgun()

	mail.articleQuantity = len(mail.Body.Articles)
	message := mailgun.NewMessage(
		"PttAlertor@mg.dinolai.com",
		mail.Title.String(),
		mail.Body.String(),
		mail.Receiver)
	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}

func newMailgun() mailgun.Mailgun {
	config := readConfig()

	domain := config["domain"]
	apiKey := config["apiKey"]
	publicAPIKey := config["publicAPIKey"]

	return mailgun.NewMailgun(domain, apiKey, publicAPIKey)
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
