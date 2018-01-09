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
		{"tomorrow", args{time.Date(0, 01, 11, 03, 01, 0, 0, time.FixedZone("CST", 8*60*60))}, 2017},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getYear(tt.args.pushTime); got != tt.want {
				t.Errorf("getYear() = %v, want %v", got, tt.want)
			}
		})
	}
}
