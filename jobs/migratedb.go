package jobs

import (
	"time"

	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/models"
	"github.com/meifamily/ptt-alertor/models/article"
	"github.com/meifamily/ptt-alertor/models/board"
)

var redisArticle = article.NewArticle(new(article.Redis))
var redisBoard = board.NewBoard(new(board.Redis), new(board.Redis))

type migrateDB struct {
}

func NewMigrateDB() *migrateDB {
	return &migrateDB{}
}

func (m migrateDB) Run() {
	m.migrateBoards()
}

func (m migrateDB) migrateBoards() {
	for _, boardName := range models.Board.List() {
		log.WithField("board", boardName).Info("Board Migrating")
		m.migrateBoard(boardName)
		time.Sleep(time.Duration(50 * time.Millisecond))
	}
	log.Info("All Board Migrated")
}

func (migrateDB) migrateBoard(boardName string) {
	redisBoard.Name = boardName

	dynamoBoard := models.Board
	dynamoBoard.Name = boardName
	dynamoBoard.Articles = redisBoard.GetArticles()

	if err := dynamoBoard.Save(); err != nil {
		log.WithField("board", boardName).Error("Migrate Board Failed")
	}
}

func (m migrateDB) migrateArticles() {
	for _, code := range new(article.Articles).List() {
		log.WithField("code", code).Info("Article Migrating")
		m.migrateArticle(code)
		time.Sleep(time.Duration(250 * time.Millisecond))
	}
	log.Info("All Article Migrated")
}

func (migrateDB) migrateArticle(code string) {
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
