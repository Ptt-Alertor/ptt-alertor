package board

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/liam-lai/ptt-alertor/ptt/article"
)

func TestIndex(t *testing.T) {
	tests := []struct {
		name string
		b    Board
		want []byte
	}{
		// TODO: Add test cases.
		{"FREE_BOX", Board{Name: "FREE_BOX"}, []byte("[{}]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			articleJSON := tt.b.Index()
			var articles []article.Article
			json.NewDecoder(bytes.NewReader(articleJSON)).Decode(articles)

			articleType := reflect.TypeOf(articles[0])
			if articleType.String() != "article.Article" {
				t.Errorf("FirstPage() content error")
			}

		})
	}
}
