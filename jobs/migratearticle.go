package jobs

import (
	"time"

	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/models"
	"github.com/meifamily/ptt-alertor/models/article"
)

var redisArticle = article.NewArticle(new(article.Redis))

type migrateArticle struct {
}

func NewMigrateArticle() *migrateArticle {
	return &migrateArticle{}
}

func (m migrateArticle) Run() {
	for _, code := range new(article.Articles).List() {
		log.WithField("code", code).Info("Article Migrating")
		m.RunSingle(code)
		time.Sleep(time.Duration(250 * time.Millisecond))
	}
	log.Info("All Article Migrated")
}

func (migrateArticle) RunSingle(code string) {
	dynamoArticle := models.Article
	a := redisArticle.Find(code)

	dynamoArticle.ID = a.ID
	dynamoArticle.Code = a.Code
	dynamoArticle.Title = a.Title
	dynamoArticle.Link = a.Link
	dynamoArticle.Date = a.Date
	dynamoArticle.Author = a.Author
	dynamoArticle.Comments = a.Comments
	dynamoArticle.LastPushDateTime = a.LastPushDateTime
	dynamoArticle.Board = a.Board
	dynamoArticle.PushSum = a.PushSum

	if err := dynamoArticle.Save(); err != nil {
		log.WithField("code", code).Error("Migrate Article Failed")
	}
}
