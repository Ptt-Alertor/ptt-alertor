package pttboard

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestFirstPage(t *testing.T) {
	type args struct {
		board string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{"FREE_BOX", args{"FREE_BOX"}, []byte("[{}]")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			articleJSON := FirstPage(tt.args.board)

			json.NewDecoder(bytes.NewReader(articleJSON)).Decode(articles)

			articleType := reflect.TypeOf(articles[0])
			if articleType.String() != "pttboard.article" {
				t.Errorf("FirstPage() content error")
			}

		})
	}
}
