package crawler

import "testing"

func Test_fetchHTML(t *testing.T) {
	type args struct {
		reqURL string
	}
	tests := []struct {
		name             string
		args             args
		wantResponseCode int
		wantErr          bool
	}{
		{"ok", args{"https://www.ptt.cc/bbs/LoL/index.html"}, 200, false},
		{"ok", args{"https://www.ptt.cc/bbs/LoL/M.1512235390.A.215.html"}, 200, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := fetchHTML(tt.args.reqURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponseCode {
				t.Errorf("fetchHTML() = %v, want %v", gotResponse.StatusCode, tt.wantResponseCode)
			}
		})
	}
}
