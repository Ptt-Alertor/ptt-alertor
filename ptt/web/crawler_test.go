package web

import (
	"reflect"
	"testing"
	"time"

	"github.com/meifamily/ptt-alertor/models/article"
	gock "gopkg.in/h2non/gock.v1"
)

func BenchmarkCurrentPage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CurrentPage("lol")
	}
}

func BenchmarkBuildArticles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FetchArticles("lol", 9697)
	}
}

func Test_getYear(t *testing.T) {
	type args struct {
		pushTime time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"same", args{time.Date(0, 01, 10, 03, 01, 0, 0, time.FixedZone("CST", 8*60*60))}, time.Now().Year()},
		{"month before", args{time.Date(0, 12, 31, 23, 59, 59, 0, time.FixedZone("CST", 8*60*60))}, time.Now().Year() - 1},
		{"tomorrow", args{time.Now().AddDate(0, 0, 1)}, time.Now().Year() - 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getYear(tt.args.pushTime); got != tt.want {
				t.Errorf("getYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkURLExist(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"found", args{"http://dinolai.com"}, true},
		{"not found", args{"http://dinolai.tw"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkURLExist(tt.args.url); got != tt.want {
				t.Errorf("checkURLExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeBoardURL(t *testing.T) {
	type args struct {
		board string
		page  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no page", args{"ezsoft", -1}, "https://www.ptt.cc/bbs/ezsoft/index.html"},
		{"page1", args{"ezsoft", 1}, "https://www.ptt.cc/bbs/ezsoft/index1.html"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeBoardURL(tt.args.board, tt.args.page); got != tt.want {
				t.Errorf("makeBoardURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeArticleURL(t *testing.T) {
	type args struct {
		board       string
		articleCode string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"M.1497363598.A.74E", args{"ezsoft", "M.1497363598.A.74E"}, "https://www.ptt.cc/bbs/ezsoft/M.1497363598.A.74E.html"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeArticleURL(tt.args.board, tt.args.articleCode); got != tt.want {
				t.Errorf("makeArticleURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fetchHTML(t *testing.T) {
	type args struct {
		reqURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ok", args{"https://www.ptt.cc/bbs/LoL/index.html"}, false},
		{"R18", args{"https://www.ptt.cc/bbs/Gossiping/index.html"}, false},
		{"not found", args{"https://www.ptt.cc/bbs/DinoLai/index.html"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fetchHTML(tt.args.reqURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBuildArticles(t *testing.T) {

	defer gock.Off()
	gock.New("https://www.ptt.cc").Get("/bbs/lol/index.html").Reply(200).BodyString(dummyBody)

	type args struct {
		board string
		page  int
	}
	tests := []struct {
		name         string
		args         args
		wantArticles article.Articles
		wantErr      bool
	}{
		{"ok", args{"lol", -1}, []article.Article{
			{
				ID:      1516285019,
				Code:    "",
				Title:   "[外絮] JTeam FB",
				Link:    "https://www.ptt.cc/bbs/LoL/M.1516285019.A.BCE.html",
				Date:    "1/18",
				Author:  "Andy7577272",
				PushSum: 2,
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArticles, err := FetchArticles(tt.args.board, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotArticles, tt.wantArticles) {
				t.Errorf("BuildArticles() = %#v, want %v", gotArticles, tt.wantArticles)
				return
			}
			if !gock.IsDone() {
				t.Errorf("BuildArticles() gock status = %v, wantErr %v", gock.IsDone(), true)
			}
		})
	}
}

func TestBuildArticle(t *testing.T) {
	defer gock.Off()
	gock.New("https://www.ptt.cc").Get("/bbs/TFSHS66th321/M.1498563199.A.35C.html").
		Reply(200).BodyString(dummyArticle)
	year := time.Now().Year()

	type args struct {
		board       string
		articleCode string
	}
	tests := []struct {
		name    string
		args    args
		want    article.Article
		wantErr bool
	}{
		{"ok", args{"TFSHS66th321", "M.1498563199.A.35C"}, article.Article{
			ID:               1498563199,
			Code:             "M.1498563199.A.35C",
			Title:            "[小葉] 公告測試",
			Link:             "https://www.ptt.cc/bbs/TFSHS66th321/M.1498563199.A.35C.html",
			LastPushDateTime: time.Date(year, 01, 02, 13, 57, 0, 0, time.FixedZone("CST", 8*60*60)),
			Board:            "TFSHS66th321",
			PushSum:          0,
			Comments: article.Comments{
				article.Comment{Tag: "→ ", UserID: "ChoDino", Content: ": 快點好嗎", DateTime: time.Date(year, 01, 01, 00, 55, 0, 0, time.FixedZone("CST", 8*60*60))},
				article.Comment{Tag: "→ ", UserID: "ChoDino", Content: ": 好了~今天先做到這~預祝空軍今天賺飽飽~睡好覺@", DateTime: time.Date(year, 01, 02, 10, 22, 0, 0, time.FixedZone("CST", 8*60*60))},
				article.Comment{Tag: "→ ", UserID: "ChoDino", Content: ": 好了~今天先做到這~預祝空軍今天賺飽飽~睡好覺@", DateTime: time.Date(year, 01, 02, 10, 26, 0, 0, time.FixedZone("CST", 8*60*60))},
				article.Comment{Tag: "→ ", UserID: "ChoDino", Content: ": timezone testing", DateTime: time.Date(year, 01, 02, 13, 57, 0, 0, time.FixedZone("CST", 8*60*60))}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchArticle(tt.args.board, tt.args.articleCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildArticle() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

var dummyBody = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">


<meta name="viewport" content="width=device-width, initial-scale=1">

<title>看板 LoL 文章列表 - 批踢踢實業坊</title>

	</head>
    <body>

<div id="topbar-container">
	<div id="topbar" class="bbs-content">
		<a id="logo" href="/">批踢踢實業坊</a>
		<span>&rsaquo;</span>
		<a class="board" href="/bbs/LoL/index.html"><span class="board-label">看板 </span>LoL</a>
		<a class="right small" href="/about.html">關於我們</a>
		<a class="right small" href="/contact.html">聯絡資訊</a>
	</div>
</div>

<div id="main-container">
	<div id="action-bar-container">
		<div class="action-bar">
			<div class="btn-group btn-group-dir">
				<a class="btn selected" href="/bbs/LoL/index.html">看板</a>
				<a class="btn" href="/man/LoL/index.html">精華區</a>
			</div>
			<div class="btn-group btn-group-paging">
				<a class="btn wide" href="/bbs/LoL/index1.html">最舊</a>
				<a class="btn wide" href="/bbs/LoL/index9851.html">&lsaquo; 上頁</a>
				<a class="btn wide disabled">下頁 &rsaquo;</a>
				<a class="btn wide" href="/bbs/LoL/index.html">最新</a>
			</div>
		</div>
	</div>

	<div class="r-list-container action-bar-margin bbs-screen">
		<div class="r-ent">
			<div class="nrec"><span class="hl f2">2</span></div>
			<div class="mark"></div>
			<div class="title">

				<a href="/bbs/LoL/M.1516285019.A.BCE.html">[外絮] JTeam FB</a>

			</div>
			<div class="meta">
				<div class="date"> 1/18</div>
				<div class="author">Andy7577272</div>
			</div>
		</div>
        <div class="r-list-sep"></div>
		<div class="r-ent">
			<div class="nrec"><span class="hl f1">爆</span></div>
			<div class="mark">M</div>
			<div class="title">

				<a href="/bbs/LoL/M.1512746508.A.54D.html">[公告] 伺服器狀況詢問/聊天/揪團/抱怨/多功能區</a>

			</div>
			<div class="meta">
				<div class="date">12/08</div>
				<div class="author">InnGee</div>
			</div>
		</div>
	</div>
</div>
    </body>
</html>
	`
var dummyArticle = `

<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">

<title>[小葉] 公告測試 - 看板 TFSHS66th321 - 批踢踢實業坊</title>
<meta name="robots" content="all">
<meta name="keywords" content="Ptt BBS 批踢踢">
<meta name="description" content="測
--
※ 發信站: 批踢踢實業坊(ptt.cc), 來自: 1.170.119.214
※ 文章網址: https://www.ptt.cc/bbs/TFSHS66th321/M.1498563199.A.35C.html
→ ChoDino: 快點好嗎                                               06/30 00:55
">
<meta property="og:site_name" content="Ptt 批踢踢實業坊">
<meta property="og:title" content="[小葉] 公告測試">
<meta property="og:description" content="測
--
※ 發信站: 批踢踢實業坊(ptt.cc), 來自: 1.170.119.214
※ 文章網址: https://www.ptt.cc/bbs/TFSHS66th321/M.1498563199.A.35C.html
→ ChoDino: 快點好嗎                                               06/30 00:55
">
<link rel="canonical" href="https://www.ptt.cc/bbs/TFSHS66th321/M.1498563199.A.35C.html">

<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.22/bbs-common.css">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.22/bbs-base.css" media="screen">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.22/bbs-custom.css">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.22/pushstream.css" media="screen">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.22/bbs-print.css" media="print">

	</head>
    <body>

<div id="fb-root"></div>
<script>(function(d, s, id) {
var js, fjs = d.getElementsByTagName(s)[0];
if (d.getElementById(id)) return;
js = d.createElement(s); js.id = id;
js.src = "//connect.facebook.net/en_US/all.js#xfbml=1";
fjs.parentNode.insertBefore(js, fjs);
}(document, 'script', 'facebook-jssdk'));</script>

<div id="topbar-container">
	<div id="topbar" class="bbs-content">
		<a id="logo" href="/">批踢踢實業坊</a>
		<span>&rsaquo;</span>
		<a class="board" href="/bbs/TFSHS66th321/index.html"><span class="board-label">看板 </span>TFSHS66th321</a>
		<a class="right small" href="/about.html">關於我們</a>
		<a class="right small" href="/contact.html">聯絡資訊</a>
	</div>
</div>
<div id="navigation-container">
	<div id="navigation" class="bbs-content">
		<a class="board" href="/bbs/TFSHS66th321/index.html">返回看板</a>
		<div class="bar"></div>
		<div class="share">
			<span>分享</span>
			<div class="fb-like" data-send="false" data-layout="button_count" data-width="90" data-show-faces="false" data-href="http://www.ptt.cc/bbs/TFSHS66th321/M.1498563199.A.35C.html"></div>

			<div class="g-plusone" data-size="medium"></div>
<script type="text/javascript">
window.___gcfg = {lang: 'zh-TW'};
(function() {
var po = document.createElement('script'); po.type = 'text/javascript'; po.async = true;
po.src = 'https://apis.google.com/js/plusone.js';
var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(po, s);
})();
</script>

		</div>
	</div>
</div>
<div id="main-container">
    <div id="main-content" class="bbs-screen bbs-content"><div class="article-metaline"><span class="article-meta-tag">作者</span><span class="article-meta-value">ChoDino ()</span></div><div class="article-metaline-right"><span class="article-meta-tag">看板</span><span class="article-meta-value">TFSHS66th321</span></div><div class="article-metaline"><span class="article-meta-tag">標題</span><span class="article-meta-value">[小葉] 公告測試</span></div><div class="article-metaline"><span class="article-meta-tag">時間</span><span class="article-meta-value">Tue Jun 27 19:33:15 2017</span></div>
測

--
<span class="f2">※ 發信站: 批踢踢實業坊(ptt.cc), 來自: 1.170.119.214
</span><span class="f2">※ 文章網址: <a href="https://www.ptt.cc/bbs/TFSHS66th321/M.1498563199.A.35C.html" target="_blank" rel="nofollow">https://www.ptt.cc/bbs/TFSHS66th321/M.1498563199.A.35C.html</a>
</span><div class="push"><span class="f1 hl push-tag">→ </span><span class="f3 hl push-userid">ChoDino</span><span class="f3 push-content">: 快點好嗎</span><span class="push-ipdatetime"> 01/01 00:55
</span></div><div class="push"><span class="f1 hl push-tag">→ </span><span class="f3 hl push-userid">ChoDino</span><span class="f3 push-content">: 好了~今天先做到這~預祝空軍今天賺飽飽~睡好覺@<a href="/cdn-cgi/l/email-protection" class="__cf_email__" data-cfemail="30467052">[email&#160;protected]</a></span><span class="push-ipdatetime"> 01/02 10:22
</span></div><div class="push"><span class="f1 hl push-tag">→ </span><span class="f3 hl push-userid">ChoDino</span><span class="f3 push-content">: 好了~今天先做到這~預祝空軍今天賺飽飽~睡好覺@<a href="/cdn-cgi/l/email-protection" class="__cf_email__" data-cfemail="c2b482a0">[email&#160;protected]</a></span><span class="push-ipdatetime"> 01/02 10:26
</span></div><div class="push"><span class="f1 hl push-tag">→ </span><span class="f3 hl push-userid">ChoDino</span><span class="f3 push-content">: timezone testing</span><span class="push-ipdatetime"> 01/02 13:57
</span></div></div>

    <div id="article-polling" data-pollurl="/poll/TFSHS66th321/M.1498563199.A.35C.html?cacheKey=2117-403609369&offset=631&offset-sig=88f7d90a2437b7ecd7731bad54988cc001b5758e" data-longpollurl="/v1/longpoll?id=253ee0098f267a4184c1b6ca84911f8e8762da0a" data-offset="631"></div>
</div>
<script data-cfasync="false" src="/cdn-cgi/scripts/af2821b0/cloudflare-static/email-decode.min.js"></script><script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-32365737-1', {
    cookieDomain: 'ptt.cc',
    legacyCookieDomain: 'ptt.cc'
  });
  ga('send', 'pageview');
</script>
<script src="//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
<script src="//images.ptt.cc/bbs/v2.22/bbs.js"></script>

    </body>
</html>
`
