package myutil

import (
	"reflect"
	"strings"
	"testing"
)

func Test_SplitTextByLineBreak(t *testing.T) {
	type args struct {
		text  string
		limit int
	}
	tests := []struct {
		name      string
		args      args
		wantTexts []string
	}{
		{"pass", args{"你好阿", 5}, []string{"你好阿"}},
		{"pass", args{"你好阿\n你好喔喔", 5}, []string{"你好阿\n", "你好喔喔"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTexts := SplitTextByLineBreak(strings.Replace(tt.args.text, "\t", "", -1), tt.args.limit); !reflect.DeepEqual(gotTexts, tt.wantTexts) {
				for i, v := range gotTexts {
					if v != tt.wantTexts[i] {
						t.Errorf("splitTextByLineBreak() = %v, want %v", v, tt.wantTexts[i])
					}
				}
			}
		})
	}
}
