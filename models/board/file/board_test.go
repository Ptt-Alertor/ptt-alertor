package file

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBoard_GetArticles(t *testing.T) {
	tests := []struct {
		name string
		bd   Board
		want string
	}{
		// TODO: Add test cases.
		{"TestJoke", Board{}, "article.Article"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bd.GetArticles("joke")
			fmt.Println(reflect.TypeOf(got[0]))
			if reflect.TypeOf(got[0]).String() != tt.want {
				t.Errorf("Board.GetArticles() = %v, want %v", got, tt.want)
			}
		})
	}
}
