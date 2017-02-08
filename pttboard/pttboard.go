package pttboard

import (
	"log"
	"net/http"

	"encoding/json"

	"golang.org/x/net/html"
)

type article struct {
	Title  string
	Link   string
	Date   string
	Author string
}

var articles []article

func fetchHTML(board string) (response *http.Response) {
	response, err := http.Get("https://www.ptt.cc/bbs/" + board + "/index.html")

	if err != nil {
		log.Fatal(err)
	}

	return response
}

func parseHTML(response *http.Response) *html.Node {
	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func buildArticles(board string) []article {

	htmlNodes := parseHTML(fetchHTML(board))

	articleBlocks := traverseHTMLNode(htmlNodes, findArticleBlocks)
	targetNodes = make([]*html.Node, 0)
	articles = make([]article, len(articleBlocks))

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

func getAnchorLink(anchor *html.Node) string {
	for _, value := range anchor.Attr {
		if value.Key == "href" {
			return value.Val
		}
	}
	return ""
}

type findInHTML func(node *html.Node) *html.Node

var targetNodes []*html.Node

func traverseHTMLNode(node *html.Node, find findInHTML) []*html.Node {

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		targetNode := find(child)
		if targetNode != nil {
			targetNodes = append(targetNodes, targetNode)

		}
		traverseHTMLNode(child, find)
	}
	return targetNodes
}

func findArticleBlocks(node *html.Node) *html.Node {
	return findDivByClassName(node, "r-ent")
}

func findTitleDiv(node *html.Node) *html.Node {
	return findDivByClassName(node, "title")
}

func findAnchor(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return node
	}
	return nil
}

func findMetaDiv(node *html.Node) *html.Node {
	return findDivByClassName(node, "meta")
}

func findDateDiv(node *html.Node) *html.Node {
	return findDivByClassName(node, "date")
}

func findAuthorDiv(node *html.Node) *html.Node {
	return findDivByClassName(node, "author")
}

func findDivByClassName(node *html.Node, className string) *html.Node {
	if node.Type == html.ElementNode && node.Data == "div" {
		for _, tagAttr := range node.Attr {
			if tagAttr.Key == "class" && tagAttr.Val == className {
				return node
			}
		}
	}
	return nil
}

func FirstPage(board string) []byte {
	articles := buildArticles(board)
	articlesJSON, err := json.Marshal(articles)
	if err != nil {
		log.Fatal(err)
	}
	return articlesJSON
}
