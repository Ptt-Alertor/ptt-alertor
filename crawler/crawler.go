package crawler

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/models/article"

	"regexp"

	"strings"

	"strconv"

	"github.com/meifamily/ptt-alertor/models/pushsum"
	"golang.org/x/net/html"
)

const pttHostURL = "https://www.ptt.cc"

// CurrentPage find Board Last Page Number
func CurrentPage(board string) (int, error) {
	url := makeBoardURL(board, -1)
	htmlNodes, err := fetchHTML(url)
	if err != nil {
		return 0, err
	}
	paging := findNodes(htmlNodes, findPagingBlock)
	for _, page := range paging {
		anchors := findNodes(page, findAnchor)
		for _, a := range anchors {
			if strings.Contains(a.FirstChild.Data, "上頁") {
				link := getAnchorLink(a)
				re := regexp.MustCompile("\\d+")
				page, err := strconv.Atoi(re.FindString(link))
				if err != nil {
					return 0, err
				}
				return page + 1, nil
			}
		}
	}
	return 0, errors.New("Parse Currenect Page Error")
}

// BuildArticles makes board's index articles to a article slice
func BuildArticles(board string, page int) (articles article.Articles, err error) {
	reqURL := makeBoardURL(board, page)
	htmlNodes, err := fetchHTML(reqURL)
	if err != nil {
		return nil, err
	}
	articleBlocks := findNodes(htmlNodes, findArticleBlocks)
	for _, articleBlock := range articleBlocks {
		article := article.Article{}
		for _, pushCountDiv := range findNodes(articleBlock, findPushCountDiv) {
			if child := pushCountDiv.FirstChild; child != nil {
				if child := child.FirstChild; child != nil {
					if err == nil {
						article.PushSum = convertPushCount(child.Data)
					}
				}
			}
		}
		for _, titleDiv := range findNodes(articleBlock, findTitleDiv) {

			anchors := findNodes(titleDiv, findAnchor)

			if len(anchors) == 0 {
				article.Title = titleDiv.FirstChild.Data
				article.Link = ""
				continue
			}

			for _, anchor := range findNodes(titleDiv, findAnchor) {
				article.Title = anchor.FirstChild.Data
				link := pttHostURL + getAnchorLink(anchor)
				article.Link = link
				article.ID = article.ParseID(link)
			}
		}
		for _, metaDiv := range findNodes(articleBlock, findMetaDiv) {
			for _, date := range findNodes(metaDiv, findDateDiv) {
				article.Date = strings.TrimSpace(date.FirstChild.Data)
			}
			for _, author := range findNodes(metaDiv, findAuthorDiv) {
				article.Author = author.FirstChild.Data
			}
		}
		articles = append(articles, article)
		if isLastArticleBlock(articleBlock) {
			break
		}
	}
	return articles, nil
}

func isLastArticleBlock(articleBlock *html.Node) bool {
	for next := articleBlock.NextSibling; ; next = next.NextSibling {
		if next == nil {
			break
		}
		if next.Type == html.ElementNode {
			for _, attr := range next.Attr {
				if attr.Val == "r-list-sep" {
					return true
				}
				return false
			}
		}
	}
	return false
}

func convertPushCount(str string) int {
	for num, text := range pushsum.NumTextMap {
		if strings.EqualFold(str, text) {
			return num
		}
	}
	cnt, err := strconv.Atoi(str)
	if err != nil {
		cnt = 0
	}
	return cnt
}

// BuildArticle build article object from html
func BuildArticle(board, articleCode string) (article.Article, error) {
	reqURL := makeArticleURL(board, articleCode)
	htmlNodes, err := fetchHTML(reqURL)
	if err != nil {
		return article.Article{}, err
	}
	atcl := article.Article{
		Link:  reqURL,
		Code:  articleCode,
		Board: board,
	}
	nodes := findNodes(htmlNodes, findOgTitleMeta)
	if len(nodes) > 0 {
		atcl.Title = getMetaContent(nodes[0])
	} else {
		atcl.Title = "[內文標題已被刪除]"
	}
	atcl.ID = atcl.ParseID(reqURL)
	pushBlocks := findNodes(htmlNodes, findPushBlocks)
	pushes := make([]article.Comment, len(pushBlocks))
	for index, pushBlock := range pushBlocks {
		for _, pushTag := range findNodes(pushBlock, findPushTag) {
			pushes[index].Tag = pushTag.FirstChild.Data
		}
		for _, pushUserID := range findNodes(pushBlock, findPushUserID) {
			pushes[index].UserID = pushUserID.FirstChild.Data
		}
		for _, pushContent := range findNodes(pushBlock, findPushContent) {
			content := pushContent.FirstChild.Data
			for n := pushContent.FirstChild.NextSibling; n != nil; n = n.NextSibling {
				if findEmailProtected(n) != nil {
					break
				}
				if n.FirstChild != nil {
					content += n.FirstChild.Data
				}
				if n.NextSibling != nil {
					content += n.NextSibling.Data
				}
			}
			pushes[index].Content = content
		}
		for _, pushIPDateTime := range findNodes(pushBlock, findPushIPDateTime) {
			ipdatetime := strings.TrimSpace(pushIPDateTime.FirstChild.Data)
			if ipdatetime == "" {
				break
			}
			dateTime, err := parseDateTime(ipdatetime)
			if err != nil {
				log.WithFields(log.Fields{
					"ipdatetime": ipdatetime,
					"board":      board,
					"code":       articleCode,
				}).WithError(err).Error("Parse DateTime Error")
			}
			pushes[index].DateTime = dateTime
			if index == len(pushBlocks)-1 {
				atcl.LastPushDateTime = pushes[index].DateTime
			}
		}
	}
	atcl.Comments = pushes
	return atcl, nil
}

func parseDateTime(ipdatetime string) (time.Time, error) {
	re, _ := regexp.Compile("(\\d+\\.\\d+\\.\\d+\\.\\d+)?\\s*(.*)")
	subMatches := re.FindStringSubmatch(ipdatetime)
	dateTime := strings.TrimSpace(subMatches[len(subMatches)-1])
	loc := time.FixedZone("CST", 8*60*60)
	t, err := time.ParseInLocation("01/02 15:04", dateTime, loc)
	if err != nil {
		return time.Time{}, err
	}
	t = t.AddDate(getYear(t), 0, 0)
	return t, nil
}

func getYear(pushTime time.Time) int {
	t := time.Now()
	if t.Month() == 1 && pushTime.Month() == 12 {
		return t.Year() - 1
	}
	return t.Year()
}

// CheckBoardExist use for checking board exist or not
func CheckBoardExist(board string) bool {
	return checkURLExist(makeBoardURL(board, -1))
}

// CheckArticleExist user for checking article exist or not
func CheckArticleExist(board, articleCode string) bool {
	return checkURLExist(makeArticleURL(board, articleCode))
}

func checkURLExist(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}

func makeBoardURL(board string, page int) string {
	var pageStr string
	if page < 0 {
		pageStr = ""
	} else {
		pageStr = strconv.Itoa(page)
	}
	return pttHostURL + "/bbs/" + board + "/index" + pageStr + ".html"
}

func makeArticleURL(board, articleCode string) string {
	return pttHostURL + "/bbs/" + board + "/" + articleCode + ".html"
}

// URLNotFoundError is an error type present 404 Not Found
type URLNotFoundError struct {
	URL string
}

func (u URLNotFoundError) Error() string {
	return "Fetched URL Not Found"
}

var errRedirect = errors.New("Redirect")

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return errRedirect
	},
	Timeout: 30 * time.Second,
}

func fetchHTML(reqURL string) (doc *html.Node, err error) {
	resp, err := client.Get(reqURL)
	if err != nil && resp == nil {
		log.WithField("url", reqURL).WithError(err).Error("Fetch URL Failed")
		return nil, err
	}
	defer resp.Body.Close()

	if err == nil && resp.StatusCode == http.StatusNotFound {
		err = URLNotFoundError{reqURL}
		return nil, err
	}

	if uerr, ok := err.(*url.Error); ok && uerr.Err == errRedirect {
		req := passR18(reqURL)
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
	}

	doc, err = html.Parse(resp.Body)
	if err != nil {
		log.WithError(err).Error("Crawler Fetch HTML Failed")
		return nil, err
	}
	return doc, nil
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
