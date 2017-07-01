package file

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/meifamily/ptt-alertor/models/ptt/board"
)

func TestBoard_GetArticles(t *testing.T) {
	tests := []struct {
		name string
		bd   Board
		want string
	}{
		// TODO: Add test cases.
		{"TestJoke", Board{board.Board{Name: "joke"}}, "article.Article"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bd.GetArticles()
			fmt.Println(reflect.TypeOf(got[0]))
			if reflect.TypeOf(got[0]).String() != tt.want {
				t.Errorf("Board.GetArticles() = %v, want %v", got, tt.want)
			}
		})
	}
}
