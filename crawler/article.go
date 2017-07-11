package crawler

import "golang.org/x/net/html"

func findTitleDiv(node *html.Node) *html.Node {
	return findDivByClassName(node, "title")
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

func findDividerDiv(node *html.Node) *html.Node {
	return findDivByClassName(node, "r-list-sep")
}

func findOgTitleMeta(node *html.Node) *html.Node {
	return findMeta(node, "og:title")
}

func findEmailProtected(node *html.Node) *html.Node {
	n := findAnchor(node)
	if n != nil {
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "__cf_email__" {
				return n
			}
		}
	}
	return nil
}
