package crawler

import "testing"

func TestBuildArticle(t *testing.T) {
	type args struct {
		board       string
		articleCode string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test", args{"ezsoft", "M.1497363598.A.74E"}, "[推薦][自製] Ptt Alertor Ptt新文章即時通知"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := BuildArticle(tt.args.board, tt.args.articleCode)
			if got := a.Title; got != tt.want {
				t.Errorf("BuildArticle() = %v, want %v", got, tt.want)
			}
		})
	}
}
