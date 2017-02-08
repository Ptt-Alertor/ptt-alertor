package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

const board string = "FREE_BOX"

type article struct {
	title  string
	href   string
	date   string
	author string
}

var articles []article

func fetchHTML(board string) (response *http.Response) {
	response, err := http.Get("https://www.ptt.cc/bbs/" + board + "/index.html")

	if err != nil {
		log.Fatal(err)
	}

	return response
}

func parseHTML(response *http.Response) {
	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	articleBlocks := traverseHTMLNode(doc, findArticleBlocks)
	targetNodes = make([]*html.Node, 0)
	articles = make([]article, len(articleBlocks))

	for index, articleBlock := range articleBlocks {
		for _, titleDiv := range traverseHTMLNode(articleBlock, findTitleDiv) {
			targetNodes = make([]*html.Node, 0)

			anchors := traverseHTMLNode(titleDiv, findAnchor)

			if len(anchors) == 0 {
				articles[index].title = titleDiv.FirstChild.Data
				articles[index].href = ""
				continue
			}

			for _, anchor := range traverseHTMLNode(titleDiv, findAnchor) {
				articles[index].title = anchor.FirstChild.Data
				articles[index].href = getAnchorLink(anchor)
			}
		}
		for _, metaDiv := range traverseHTMLNode(articleBlock, findMetaDiv) {
			targetNodes = make([]*html.Node, 0)

			for _, date := range traverseHTMLNode(metaDiv, findDateDiv) {
				articles[index].date = date.FirstChild.Data
			}
			for _, author := range traverseHTMLNode(metaDiv, findAuthorDiv) {
				articles[index].author = author.FirstChild.Data
			}
		}
	}

	for _, article := range articles {
		fmt.Println(article)
	}
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

func Json() []byte {
	// fmt.Printf("%s", fetchHTML(board))
	parseHTML(fetchHTML(board))
	var json []byte
	return json
}

func main() {
	Json()
}
