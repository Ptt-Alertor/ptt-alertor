package crawler

import (
	"testing"
	"time"
)

func BenchmarkCurrentPage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CurrentPage("lol")
	}
}

func BenchmarkBuildArticles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildArticles("lol", 9697)
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
		{"same", args{time.Date(0, 01, 10, 03, 01, 0, 0, time.FixedZone("CST", 8*60*60))}, 2018},
		{"month before", args{time.Date(0, 12, 10, 03, 01, 0, 0, time.FixedZone("CST", 8*60*60))}, 2017},
		{"tomorrow", args{time.Now().AddDate(0, 0, 1)}, 2017},
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
