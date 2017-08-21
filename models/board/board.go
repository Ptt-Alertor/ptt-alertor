package board

import (
	"math"
	"strings"

	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/crawler"
	"github.com/meifamily/ptt-alertor/models/article"
	"github.com/meifamily/ptt-alertor/models/board/redis"
	"github.com/meifamily/ptt-alertor/myutil/maputil"
	"github.com/meifamily/ptt-alertor/rss"
)

type BoardNotExistError struct {
	Suggestion string
}

func (e BoardNotExistError) Error() string {
	return "board is not exist"
}

var driver = new(redis.Board)

type Driver interface {
	List() []string
	Exist(boardName string) bool
	GetArticles(boardName string) article.Articles
	Create(boardName string) error
	Save(boardName string, articles article.Articles) error
	Delete(boardName string) error
}

type Board struct {
	Name           string
	Articles       article.Articles
	OnlineArticles article.Articles
	NewArticles    article.Articles
	driver         Driver
}

func NewBoard() *Board {
	return &Board{
		driver: driver,
	}
}

func (bd Board) List() []string {
	return bd.driver.List()
}

func (bd Board) Exist() bool {
	return bd.driver.Exist(bd.Name)
}

func (bd Board) All() (bds []*Board) {
	boards := bd.List()
	for _, board := range boards {
		bd := NewBoard()
		bd.Name = board
		bds = append(bds, bd)
	}
	return bds
}

func (bd Board) GetArticles() (articles article.Articles) {
	return bd.driver.GetArticles(bd.Name)
}

func (bd Board) Create() error {
	return bd.driver.Create(bd.Name)
}

func (bd Board) Save() error {
	return bd.driver.Save(bd.Name, bd.Articles)
}

func (bd Board) Delete() error {
	return bd.driver.Delete(bd.Name)
}

func (bd *Board) WithArticles() {
	bd.Articles = bd.GetArticles()
}

func (bd *Board) WithNewArticles() {
	bd.NewArticles, bd.OnlineArticles = NewArticles(*bd)
}

func (bd Board) FetchArticles() (articles article.Articles) {
	articles, err := rss.BuildArticles(bd.Name)
	if err != nil {
		log.WithField("board", bd.Name).WithError(err).Error("RSS Parse Failed, Switch to HTML Crawler")
		articles, err = crawler.BuildArticles(bd.Name, -1)
		if err != nil {
			log.WithField("board", bd.Name).WithError(err).Error("HTML Parse Failed")
		}
	}
	if strings.EqualFold(bd.Name, "allpost") {
		fixLink(&articles)
	}
	return articles
}

func fixLink(articles *article.Articles) {
	for i, a := range *articles {
		preParenthesesIndex := strings.LastIndex(a.Title, "(")
		backParenthesesIndex := strings.LastIndex(a.Title, ")")
		realBoard := a.Title[preParenthesesIndex+1 : backParenthesesIndex]
		a.Link = strings.Replace(a.Link, "ALLPOST", realBoard, -1)
		(*articles)[i] = a
	}
}

func NewArticles(bd Board) (newArticles, onlineArticles article.Articles) {
	newArticles = make(article.Articles, 0)
	savedArticles := bd.driver.GetArticles(bd.Name)
	onlineArticles = bd.FetchArticles()
	if len(savedArticles) == 0 {
		return nil, onlineArticles
	}
	for _, onlineArticle := range onlineArticles {
		for index, savedArticle := range savedArticles {
			if onlineArticle.ID <= savedArticle.ID {
				break
			}
			if index == len(savedArticles)-1 {
				newArticles = append(newArticles, onlineArticle)
			}
		}
	}
	return newArticles, onlineArticles
}

func (bd Board) SuggestBoardName() string {
	names := bd.List()
	boardWeight := map[string]int{}
	chars := strings.Split(strings.ToLower(bd.Name), "")
	for _, name := range names {
		count := 0
		for _, char := range chars {
			if strings.Contains(name, char) {
				count++
			}
		}
		boardWeight[name] = count / (1 + int(math.Abs(float64(len(bd.Name)-len(name)))))
	}
	return maputil.FirstByValueInt(boardWeight)
}

func CheckBoardExist(boardName string) (bool, string) {
	bd := NewBoard()
	bd.Name = boardName
	if bd.Exist() {
		return true, ""
	}
	if crawler.CheckBoardExist(boardName) {
		bd.Create()
		return true, ""
	}

	suggestBoard := bd.SuggestBoardName()
	log.WithFields(log.Fields{
		"inputBoard":   boardName,
		"suggestBoard": suggestBoard,
	}).Warning("Board Not Exist")
	return false, suggestBoard
}
