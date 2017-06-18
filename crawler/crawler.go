package crawler

import (
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/liam-lai/ptt-alertor/models/ptt/article"

	"golang.org/x/net/html"
)

// BuildArticles makes board's index articles to a article slice
func BuildArticles(board string) article.Articles {

	htmlNodes := parseHTML(fetchHTML(board))

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
				articles[index].Link = "https://www.ptt.cc" + getAnchorLink(anchor)
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

// CheckBoardExist use for checking board exist or not
func CheckBoardExist(board string) bool {
	response := fetchHTML(board)
	if response.StatusCode == http.StatusNotFound {
		return false
	}
	return true
}

func fetchHTML(board string) (response *http.Response) {

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Redirect")
		},
	}

	reqURL := "https://www.ptt.cc/bbs/" + board + "/index.html"
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
