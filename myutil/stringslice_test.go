package myutil

import (
	"reflect"
	"testing"
)

func TestStringSlice_Clean(t *testing.T) {
	tests := []struct {
		name   string
		ss     *StringSlice
		result *StringSlice
	}{
		{"empty", &StringSlice{"", "*"}, &StringSlice{}},
		{"mixed", &StringSlice{"", "*", "abc"}, &StringSlice{"abc"}},
		{"mixed", &StringSlice{"", "abc", "*"}, &StringSlice{"abc"}},
		{"mixed", &StringSlice{"a", "b", "c"}, &StringSlice{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ss.Clean()
			if !reflect.DeepEqual(tt.ss, tt.result) {
				t.Errorf("StringSlice.Clean() want %+v, got %+v", tt.result, tt.ss)
			}
		})
	}
}

func TestStringSlice_RemoveStringsSpace(t *testing.T) {
	tests := []struct {
		name   string
		ss     *StringSlice
		result *StringSlice
	}{
		{"empty", &StringSlice{""}, &StringSlice{""}},
		{"single", &StringSlice{"a b"}, &StringSlice{"ab"}},
		{"multiple", &StringSlice{"a b", "c d"}, &StringSlice{"ab", "cd"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ss.RemoveStringsSpace()
			if !reflect.DeepEqual(tt.ss, tt.result) {
				t.Errorf("StringSlice.RemoveStringsSpace() want %+v, got %+v", tt.result, tt.ss)
			}
		})
	}
}

func TestStringSlice_AppendNonRepeatElement(t *testing.T) {
	type args struct {
		str           string
		caseSensitive bool
	}
	tests := []struct {
		name   string
		ss     *StringSlice
		args   args
		result *StringSlice
	}{
		{"repeat, samecase", &StringSlice{"a", "b", "c"}, args{"a", true}, &StringSlice{"a", "b", "c"}},
		{"repead, diffcase", &StringSlice{"a", "b", "c"}, args{"A", false}, &StringSlice{"a", "b", "c"}},
		{"non-repeat", &StringSlice{"a", "b", "c"}, args{"d", true}, &StringSlice{"a", "b", "c", "d"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ss.AppendNonRepeatElement(tt.args.str, tt.args.caseSensitive)
			if !reflect.DeepEqual(tt.ss, tt.result) {
				t.Errorf("StringSlice.RemoveStringsSpace() want %+v, got %+v", tt.result, tt.ss)
			}
		})
	}
}

func TestStringSlice_AppendNonRepeat(t *testing.T) {
	type args struct {
		objectStrs    []string
		caseSensitive bool
	}
	tests := []struct {
		name   string
		ss     *StringSlice
		args   args
		result *StringSlice
	}{
		{"repeat, samecase", &StringSlice{"a", "b", "c"}, args{[]string{"a", "b"}, true}, &StringSlice{"a", "b", "c"}},
		{"repeat, diffcase", &StringSlice{"a", "b", "c"}, args{[]string{"A", "b"}, false}, &StringSlice{"a", "b", "c"}},
		{"non-repeat", &StringSlice{"a", "b", "c"}, args{[]string{"a", "d"}, true}, &StringSlice{"a", "b", "c", "d"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ss.AppendNonRepeat(tt.args.objectStrs, tt.args.caseSensitive)
			if !reflect.DeepEqual(tt.ss, tt.result) {
				t.Errorf("StringSlice.RemoveStringsSpace() want %+v, got %+v", tt.result, tt.ss)
			}
		})
	}
}

func TestStringSlice_DeleteElement(t *testing.T) {
	type args struct {
		str           string
		caseSensitive bool
	}
	tests := []struct {
		name   string
		ss     *StringSlice
		args   args
		result *StringSlice
	}{
		{"repeat, samecase", &StringSlice{"a", "b", "c"}, args{"a", true}, &StringSlice{"b", "c"}},
		{"repead, diffcase", &StringSlice{"a", "b", "c"}, args{"A", false}, &StringSlice{"b", "c"}},
		{"non-repeat", &StringSlice{"a", "b", "c"}, args{"d", true}, &StringSlice{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ss.DeleteElement(tt.args.str, tt.args.caseSensitive)
			if !reflect.DeepEqual(tt.ss, tt.result) {
				t.Errorf("StringSlice.RemoveStringsSpace() want %+v, got %+v", tt.result, tt.ss)
			}
		})
	}
}

func TestStringSlice_Delete(t *testing.T) {
	type args struct {
		objectStrs    []string
		caseSensitive bool
	}
	tests := []struct {
		name   string
		ss     *StringSlice
		args   args
		result *StringSlice
	}{
		{"repeat, samecase", &StringSlice{"a", "b", "c"}, args{[]string{"a", "b"}, true}, &StringSlice{"c"}},
		{"repeat, diffcase", &StringSlice{"a", "b", "c"}, args{[]string{"A", "b"}, false}, &StringSlice{"c"}},
		{"non-repeat", &StringSlice{"a", "b", "c"}, args{[]string{"a", "d"}, true}, &StringSlice{"b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ss.Delete(tt.args.objectStrs, tt.args.caseSensitive)
			if !reflect.DeepEqual(tt.ss, tt.result) {
				t.Errorf("StringSlice.RemoveStringsSpace() want %+v, got %+v", tt.result, tt.ss)
			}
		})
	}
}

func TestStringSlice_Index(t *testing.T) {
	type args struct {
		value         string
		caseSensitive bool
	}
	tests := []struct {
		name string
		ss   StringSlice
		args args
		want int
	}{
		{"found, samecase", StringSlice{"a", "b", "c"}, args{"a", true}, 0},
		{"found, diffcase", StringSlice{"a", "b", "c"}, args{"A", false}, 0},
		{"not found", StringSlice{"a", "b", "c"}, args{"d", true}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.Index(tt.args.value, tt.args.caseSensitive); got != tt.want {
				t.Errorf("StringSlice.Index() = %v, want %v", got, tt.want)
			}
		})
	}
}
