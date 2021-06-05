package web

import "golang.org/x/net/html"

func findPushBlocks(node *html.Node) *html.Node {
	return findDivByClassName(node, "push")
}

func findPushTag(node *html.Node) *html.Node {
	tagNode := findSpanByClassName(node, "hl push-tag")
	if tagNode != nil {
		return tagNode
	}
	return findSpanByClassName(node, "f1 hl push-tag")
}

func findPushUserID(node *html.Node) *html.Node {
	return findSpanByClassName(node, "f3 hl push-userid")
}

func findPushContent(node *html.Node) *html.Node {
	return findSpanByClassName(node, "f3 push-content")
}

func findPushIPDateTime(node *html.Node) *html.Node {
	return findSpanByClassName(node, "push-ipdatetime")
}
