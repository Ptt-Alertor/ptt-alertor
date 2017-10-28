package myutil

import "testing"

func TestDiffMap(t *testing.T) {
	type args struct {
		mapOld []map[string]string
		mapNew []map[string]string
	}
	tests := []struct {
		name        string
		args        args
		wantDiffLen int
	}{
		{"same", args{
			[]map[string]string{{"a": "a"}, {"b": "b"}},
			[]map[string]string{{"a": "a"}, {"b": "b"}},
		}, 0},
		{"diff same len", args{
			[]map[string]string{{"a": "a"}, {"b": "b"}},
			[]map[string]string{{"a": "a"}, {"c": "c"}},
		}, 1},
		{"new len longer than old", args{
			[]map[string]string{{"a": "a"}, {"b": "b"}},
			[]map[string]string{{"a": "a"}, {"b": "b"}, {"c": "c"}},
		}, 1},
		{"old len longer than new", args{
			[]map[string]string{{"a": "a"}, {"b": "b"}, {"c": "c"}},
			[]map[string]string{{"a": "a"}, {"b": "b"}},
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDiff := DiffMap(tt.args.mapOld, tt.args.mapNew); len(gotDiff) != tt.wantDiffLen {
				t.Errorf("DiffMap() = %v, want len %v", gotDiff, tt.wantDiffLen)
			}
		})
	}
}
