package crawler

import (
	"errors"
	"net/http"

	log "github.com/meifamily/logrus"

	"github.com/liam-lai/ptt-alertor/models/ptt/article"

	"golang.org/x/net/html"
)

const pttHostURL = "https://www.ptt.cc"

// BuildArticles makes board's index articles to a article slice
func BuildArticles(board string) article.Articles {

	reqURL := makeBoardURL(board)
	htmlNodes := parseHTML(fetchHTML(reqURL))

	articleBlocks := traverseHTMLNode(htmlNodes, findArticleBlocks)
	targetNodes = make([]*html.Node, 0)
	articles := make(article.Articles, len(articleBlocks))
	for index, articleBlock := range articleBlocks {
		for _, titleDiv := range traverseHTMLNode(articleBlock, findTitleDiv) {
			targetNodes = make([]*html.Node, 0)

			anchors := traverseHTMLNode(titleDiv, findAnchor)

			if len(anchors) == 0 {
				articles[index].Title = titleDiv.FirstChild.Data
				articles[index].Link = ""
				continue
			}

			for _, anchor := range traverseHTMLNode(titleDiv, findAnchor) {
				articles[index].Title = anchor.FirstChild.Data
				link := "https://www.ptt.cc" + getAnchorLink(anchor)
				articles[index].Link = link
				articles[index].ID = articles[index].ParseID(link)
			}
		}
		for _, metaDiv := range traverseHTMLNode(articleBlock, findMetaDiv) {
			targetNodes = make([]*html.Node, 0)

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
		Title: getTitle(traverseHTMLNode(htmlNodes, findTitle)[0]),
		Link:  reqURL,
	}
	atcl.ID = atcl.ParseID(reqURL)
	pushBlocks := traverseHTMLNode(htmlNodes, findPushBlocks)
	pushes := make([]article.Push, len(pushBlocks))
	for index, pushBlock := range pushBlocks {
		for _, pushTag := range traverseHTMLNode(pushBlock, findPushTag) {
			pushes[index].Tag = pushTag.FirstChild.Data
		}
		for _, pushUserID := range traverseHTMLNode(pushBlock, findPushUserID) {
			pushes[index].UserID = pushUserID.FirstChild.Data
		}
		for _, pushContent := range traverseHTMLNode(pushBlock, findPushContent) {
			pushes[index].Content = pushContent.FirstChild.Data
		}
		for _, pushIPDateTime := range traverseHTMLNode(pushBlock, findPushIPDateTime) {
			pushes[index].IPDateTime = pushIPDateTime.FirstChild.Data
		}
	}
	atcl.PushList = pushes
	return atcl
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
