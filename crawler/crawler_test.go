package crawler

import "testing"

func TestBuildPushList(t *testing.T) {
	type args struct {
		board       string
		articleCode string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test", args{"ezsoft", "M.1497363598.A.74E"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(BuildPushList(tt.args.board, tt.args.articleCode)); got < tt.want {
				t.Errorf("BuildPushList() = %v, want %v", got, tt.want)
			}
		})
	}
}
