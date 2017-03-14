package mail

import (
	"strconv"

	log "github.com/Sirupsen/logrus"

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
		log.WithError(err).Error("Sent Email Failed")
	} else {
		log.WithFields(log.Fields{
			"ID":   id,
			"Resp": resp,
		}).Info("Sent Email")
	}
}

func newMailgun() mailgun.Mailgun {
	config := myutil.Config("mailgun")

	domain := config["domain"]
	apiKey := config["apiKey"]
	publicAPIKey := config["publicAPIKey"]

	return mailgun.NewMailgun(domain, apiKey, publicAPIKey)
}
