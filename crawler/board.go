package crawler

import "golang.org/x/net/html"

func findPushCountDiv(node *html.Node) *html.Node {
	return findDivByClassName(node, "nrec")
}

func findArticleBlocks(node *html.Node) *html.Node {
	return findDivByClassName(node, "r-ent")
}
