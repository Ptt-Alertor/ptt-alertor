package crawler

import (
	"errors"
	"net/http"
	"time"

	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/models/ptt/article"

	"regexp"

	"strings"

	"golang.org/x/net/html"
)

const pttHostURL = "https://www.ptt.cc"

// BuildArticles makes board's index articles to a article slice
func BuildArticles(board string) article.Articles {

	reqURL := makeBoardURL(board)
	htmlNodes := parseHTML(fetchHTML(reqURL))

	articleBlocks := traverseHTMLNode(htmlNodes, findArticleBlocks)
	initialTargetNodes()
	articles := make(article.Articles, len(articleBlocks))
	for index, articleBlock := range articleBlocks {
		for _, titleDiv := range traverseHTMLNode(articleBlock, findTitleDiv) {
			initialTargetNodes()

			anchors := traverseHTMLNode(titleDiv, findAnchor)

			if len(anchors) == 0 {
				articles[index].Title = titleDiv.FirstChild.Data
				articles[index].Link = ""
				continue
			}

			for _, anchor := range traverseHTMLNode(titleDiv, findAnchor) {
				articles[index].Title = anchor.FirstChild.Data
				link := pttHostURL + getAnchorLink(anchor)
				articles[index].Link = link
				articles[index].ID = articles[index].ParseID(link)
			}
		}
		for _, metaDiv := range traverseHTMLNode(articleBlock, findMetaDiv) {
			initialTargetNodes()

			for _, date := range traverseHTMLNode(metaDiv, findDateDiv) {
				articles[index].Date = date.FirstChild.Data
			}
			for _, author := range traverseHTMLNode(metaDiv, findAuthorDiv) {
				articles[index].Author = author.FirstChild.Data
			}
		}
	}
	return articles
}

// BuildArticle build article object from html
func BuildArticle(board, articleCode string) article.Article {

	reqURL := makeArticleURL(board, articleCode)
	htmlNodes := parseHTML(fetchHTML(reqURL))
	atcl := article.Article{
		Link:  reqURL,
		Code:  articleCode,
		Board: board,
	}
	nodes := traverseHTMLNode(htmlNodes, findOgTitleMeta)
	if len(nodes) > 0 {
		atcl.Title = getMetaContent(nodes[0])
	} else {
		atcl.Title = "[內文標題已被刪除]"
	}
	atcl.ID = atcl.ParseID(reqURL)
	pushBlocks := traverseHTMLNode(htmlNodes, findPushBlocks)
	initialTargetNodes()
	pushes := make([]article.Push, len(pushBlocks))
	for index, pushBlock := range pushBlocks {
		for _, pushTag := range traverseHTMLNode(pushBlock, findPushTag) {
			initialTargetNodes()
			pushes[index].Tag = pushTag.FirstChild.Data
		}
		for _, pushUserID := range traverseHTMLNode(pushBlock, findPushUserID) {
			initialTargetNodes()
			pushes[index].UserID = pushUserID.FirstChild.Data
		}
		for _, pushContent := range traverseHTMLNode(pushBlock, findPushContent) {
			initialTargetNodes()
			content := pushContent.FirstChild.Data
			for n := pushContent.FirstChild.NextSibling; n != nil; n = n.NextSibling {
				if n.FirstChild != nil {
					content += n.FirstChild.Data
				}
				if n.NextSibling != nil {
					content += n.NextSibling.Data
				}
			}
			pushes[index].Content = content
		}
		for _, pushIPDateTime := range traverseHTMLNode(pushBlock, findPushIPDateTime) {
			initialTargetNodes()
			dateTime, err := parseDateTime(pushIPDateTime.FirstChild.Data)
			if err != nil {
				log.WithFields(log.Fields{
					"datetime": dateTime,
					"board":    board,
					"code":     articleCode,
				}).WithError(err).Error("Parse DateTime Error")
			}
			pushes[index].DateTime = dateTime
			if index == len(pushBlocks)-1 {
				atcl.LastPushDateTime = pushes[index].DateTime
			}
		}
	}
	atcl.PushList = pushes
	return atcl
}

func parseDateTime(ipdatetime string) (time.Time, error) {
	re, _ := regexp.Compile("(\\d+\\.\\d+\\.\\d+\\.\\d+)?\\s*(.*)")
	ipdatetime = strings.TrimSpace(ipdatetime)
	subMatches := re.FindStringSubmatch(ipdatetime)
	dateTime := strings.TrimSpace(subMatches[len(subMatches)-1])
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation("01/02 15:04", dateTime, loc)
	if err != nil {
		return t, err
	}
	t = t.AddDate(getYear(t), 0, 0)
	return t, nil
}

func getYear(pushTime time.Time) int {
	t := time.Now()
	if t.Month() == 1 && pushTime.Month() == 12 {
		return t.Year() - 1
	}
	return t.Year()
}

// CheckBoardExist use for checking board exist or not
func CheckBoardExist(board string) bool {
	reqURL := makeBoardURL(board)
	response := fetchHTML(reqURL)
	if response.StatusCode == http.StatusNotFound {
		return false
	}
	return true
}

// CheckArticleExist user for checking article exist or not
func CheckArticleExist(board, articleCode string) bool {
	reqURL := makeArticleURL(board, articleCode)
	response := fetchHTML(reqURL)
	if response.StatusCode == http.StatusNotFound {
		return false
	}
	return true
}

func makeBoardURL(board string) string {
	return pttHostURL + "/bbs/" + board + "/index.html"
}

func makeArticleURL(board, articleCode string) string {
	return pttHostURL + "/bbs/" + board + "/" + articleCode + ".html"
}

func fetchHTML(reqURL string) (response *http.Response) {

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Redirect")
		},
	}

	response, err := client.Get(reqURL)

	if response.StatusCode == http.StatusNotFound {
		log.WithField("url", reqURL).Warn("Fetched URL Not Found")
	}

	if err != nil && response.StatusCode == http.StatusFound {
		req := passR18(reqURL)
		response, err = client.Do(req)
	}

	if err != nil {
		log.WithField("url", reqURL).Error("Fetch URL Failed")
	}

	return response
}

func passR18(reqURL string) (req *http.Request) {

	req, _ = http.NewRequest("GET", reqURL, nil)

	over18Cookie := http.Cookie{
		Name:       "over18",
		Value:      "1",
		Domain:     "www.ptt.cc",
		Path:       "/",
		RawExpires: "Session",
		MaxAge:     0,
		HttpOnly:   false,
	}

	req.AddCookie(&over18Cookie)

	return req
}

func parseHTML(response *http.Response) *html.Node {
	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Error(err)
	}
	return doc
}
