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

func fetchHTML(board string) (response *http.Response) {
	response, err := http.Get("https://www.ptt.cc/bbs/" + board + "/index.html")

	if err != nil {
		log.Fatal(err)
	}

	return response
}

func parseHTMLByToken(response *http.Response) {
	tokenizer := html.NewTokenizer(response.Body)

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return
		}
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "div" {
				for _, tagAttr := range token.Attr {
					if tagAttr.Key == "class" && tagAttr.Val == "title" {

					}
				}
			}
		}
	}
}

func parseHTML(response *http.Response) {
	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	findTitleDiv(doc)

}

func findTitleDiv(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "div" {
		for _, tagAttr := range node.Attr {
			if tagAttr.Key == "class" && tagAttr.Val == "title" {
				findTitleAnchor(node)
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findTitleDiv(child)
	}
}

func findTitleAnchor(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		fmt.Println(node.FirstChild.Data)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		findTitleAnchor(child)
	}
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
