package myutil

import (
	"reflect"
	"strings"
	"testing"
)

func Test_splitTextByLineBreak(t *testing.T) {
	type args struct {
		text  string
		limit int
	}
	tests := []struct {
		name      string
		args      args
		wantTexts []string
	}{
		{"pass", args{
			`你好
			阿`, 5}, []string{"你好", "阿"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTexts := splitTextByLineBreak(tt.args.text, tt.args.limit); !reflect.DeepEqual(gotTexts, tt.wantTexts) {
				for i, v := range gotTexts {
					if strings.TrimSpace(v) != tt.wantTexts[i] {
						t.Errorf("splitTextByLineBreak() = %v, want %v", v, tt.wantTexts[i])
					}
				}
			}
		})
	}
}
