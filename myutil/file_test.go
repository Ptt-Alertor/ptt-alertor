package myutil

import "testing"

func TestFileNameAndExtension(t *testing.T) {
	type args struct {
		basefilename string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{"abc.cde.json", args{"abc.cde.json"}, "abc.cde", "json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FileNameAndExtension(tt.args.basefilename)
			if got != tt.want {
				t.Errorf("FileNameAndExtension() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FileNameAndExtension() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
