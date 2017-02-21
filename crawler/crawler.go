package crawler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/liam-lai/ptt-alertor/ptt/article"

	"golang.org/x/net/html"
)

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
				articles[index].Link = getAnchorLink(anchor)
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

func fetchHTML(board string) (response *http.Response) {

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Redirect")
		},
	}

	reqURL := "https://www.ptt.cc/bbs/" + board + "/index.html"
	response, err := client.Get(reqURL)

	if response.StatusCode == http.StatusNotFound {
		fmt.Println(404)
	}

	if err != nil {
		if response.StatusCode == http.StatusFound {
			req := passR18(reqURL)
			response, err = client.Do(req)
		} else {
			log.Fatal(err)
		}
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
		log.Fatal(err)
	}
	return doc
}
