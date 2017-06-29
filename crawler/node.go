package crawler

import (
	"golang.org/x/net/html"
)

type findInHTML func(node *html.Node) *html.Node

var targetNodes []*html.Node

func initialTargetNodes() {
	targetNodes = make([]*html.Node, 0)
}

func traverseHTMLNode(nodes *html.Node, find findInHTML) []*html.Node {

	for child := nodes.FirstChild; child != nil; child = child.NextSibling {
		targetNode := find(child)
		if targetNode != nil {
			targetNodes = append(targetNodes, targetNode)

		}
		traverseHTMLNode(child, find)
	}
	return targetNodes
}

func findAnchor(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return node
	}
	return nil
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

func findSpanByClassName(node *html.Node, className string) *html.Node {
	if node.Type == html.ElementNode && node.Data == "span" {
		for _, tagAttr := range node.Attr {
			if tagAttr.Key == "class" && tagAttr.Val == className {
				return node
			}
		}
	}
	return nil
}

func getAnchorLink(anchor *html.Node) string {
	for _, value := range anchor.Attr {
		if value.Key == "href" {
			return value.Val
		}
	}
	return ""
}

func getTitle(title *html.Node) string {
	return title.FirstChild.Data
}

func findTitle(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "title" {
		return node
	}
	return nil
}

func findMeta(node *html.Node, property string) *html.Node {
	if node.Type == html.ElementNode && node.Data == "meta" {
		for _, tagAttr := range node.Attr {
			if tagAttr.Key == "property" && tagAttr.Val == property {
				return node
			}
		}
	}
	return nil
}

func getMetaContent(meta *html.Node) string {
	for _, value := range meta.Attr {
		if value.Key == "content" {
			return value.Val
		}
	}
	return ""
}
