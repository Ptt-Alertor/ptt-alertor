package crawler

import (
	"testing"
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
