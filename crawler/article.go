package crawler

import "golang.org/x/net/html"

func findArticleBlocks(node *html.Node) *html.Node {
	return findDivByClassName(node, "r-ent")
}

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
