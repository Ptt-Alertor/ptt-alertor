package board

import (
	"math"
	"strings"

	log "github.com/Ptt-Alertor/logrus"
	"github.com/Ptt-Alertor/ptt-alertor/models/article"
	"github.com/Ptt-Alertor/ptt-alertor/myutil/maputil"
	"github.com/Ptt-Alertor/ptt-alertor/ptt/rss"
	"github.com/Ptt-Alertor/ptt-alertor/ptt/web"
)

type BoardNotExistError struct {
	Suggestion string
}

func (e BoardNotExistError) Error() string {
	return "board is not exist"
}

type Driver interface {
	GetArticles(boardName string) article.Articles
	Save(boardName string, articles article.Articles) error
	Delete(boardName string) error
}

type Cacher interface {
	List() []string
	Create(boardName string) error
	Exist(boardName string) bool
	Remove(boardName string) error
}

type Board struct {
	Name           string
	Articles       article.Articles
	OnlineArticles article.Articles
	NewArticles    article.Articles
	driver         Driver
	cacher         Cacher
}

func NewBoard(drive Driver, cache Cacher) *Board {
	return &Board{
		driver: drive,
		cacher: cache,
	}
}

func (bd Board) List() []string {
	return bd.cacher.List()
}

func (bd Board) Exist() bool {
	return bd.cacher.Exist(bd.Name)
}

func (bd Board) All() (bds []*Board) {
	boards := bd.List()
	for _, board := range boards {
		bd := NewBoard(bd.driver, bd.cacher)
		bd.Name = board
		bds = append(bds, bd)
	}
	return bds
}

func (bd Board) GetArticles() (articles article.Articles) {
	return bd.driver.GetArticles(bd.Name)
}

func (bd Board) Create() error {
	return bd.cacher.Create(bd.Name)
}

func (bd Board) Save() error {
	return bd.driver.Save(bd.Name, bd.Articles)
}

func (bd Board) Delete() error {
	if err := bd.driver.Delete(bd.Name); err != nil {
		return err
	}

	if err := bd.cacher.Remove(bd.Name); err != nil {
		return err
	}
	return nil
}

func (bd *Board) WithArticles() {
	bd.Articles = bd.GetArticles()
}

func (bd *Board) WithNewArticles() {
	bd.NewArticles, bd.OnlineArticles = newArticles(*bd)
}

func newArticles(bd Board) (newArticles, onlineArticles article.Articles) {
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

func (bd Board) FetchArticles() (articles article.Articles) {
	articles, err := rss.BuildArticles(bd.Name)
	if err != nil {
		if err == rss.ErrTooManyRequests {
			log.WithError(err).Warning("RSS Parse Failed")
			return
		}
		log.WithField("board", bd.Name).WithError(err).Error("RSS Parse Failed, Switch to HTML Crawler")
		articles, err = web.FetchArticles(bd.Name, -1)
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
	return maputil.MaxIntKey(boardWeight)
}

func CheckBoardExist(boardName string) (bool, string) {
	bd := NewBoard(new(DynamoDB), new(Redis))
	bd.Name = boardName
	if bd.Exist() {
		return true, ""
	}
	if web.CheckBoardExist(boardName) {
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
