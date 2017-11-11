package maputil

import (
	"testing"
)

func TestFirstByValueInt(t *testing.T) {
	type args struct {
		strs map[string]int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{
			strs: map[string]int{
				"a": 1,
			},
		}, "a"},
		{"seq", args{
			strs: map[string]int{
				"a": 1,
				"b": 3,
				"c": 5,
			},
		}, "c"},
		{"rand", args{
			strs: map[string]int{
				"a": 1,
				"b": 99,
				"c": 15,
				"d": -20,
			},
		}, "b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxIntKey(tt.args.strs); got != tt.want {
				t.Errorf("FirstByValueInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirstByValueFloat(t *testing.T) {
	type args struct {
		strs map[string]float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{
			strs: map[string]float64{
				"a": 1.64,
			},
		}, "a"},
		{"seq", args{
			strs: map[string]float64{
				"a": 1.64,
				"b": 3.17,
				"c": 5.25,
			},
		}, "c"},
		{"rand", args{
			strs: map[string]float64{
				"a": 1.00,
				"b": 99.999,
				"c": 15.1515,
				"d": -20.287,
			},
		}, "b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxFloatKey(tt.args.strs); got != tt.want {
				t.Errorf("FirstByValueFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
