package rss

import "testing"

func TestCheckBoardExist(t *testing.T) {
	type args struct {
		board string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"exist", args{"movie"}, true},
		{"exist", args{"movies"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckBoardExist(tt.args.board); got != tt.want {
				t.Errorf("CheckBoardExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
