package user

import (
	"reflect"
	"testing"
)

func TestUser_All(t *testing.T) {
	tests := []struct {
		name   string
		u      User
		wantUs []*User
	}{
		{"ok", User{drive: new(Mock)}, []*User{
			&User{Profile: Profile{Account: "dinos80152@gmail.com"}, drive: new(Mock)},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotUs := tt.u.All(); !reflect.DeepEqual(gotUs, tt.wantUs) {
				t.Errorf("User.All() = %v, want %v", gotUs[0], tt.wantUs[0])
			}
		})
	}
}

func TestUser_Save(t *testing.T) {
	tests := []struct {
		name    string
		u       User
		wantErr bool
	}{
		{"ok", User{Profile: Profile{Account: "liam.lai@gmail.com", Email: "liam.lai@gmail.com"}, drive: new(Mock)}, false},
		{"duplicate", User{Profile: Profile{Account: "dinos80152@gmail.com", Email: "dinos80152@gmail.com"}, drive: new(Mock)}, true},
		{"not enough data", User{Profile: Profile{Account: "dinos80152@gmail.com"}, drive: new(Mock)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Save(); (err != nil) != tt.wantErr {
				t.Errorf("User.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_Update(t *testing.T) {
	tests := []struct {
		name    string
		u       User
		wantErr bool
	}{
		{"ok", User{Profile: Profile{Account: "dinos80152@gmail.com"}, drive: new(Mock)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Update(); (err != nil) != tt.wantErr {
				t.Errorf("User.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_Find(t *testing.T) {
	type args struct {
		account string
	}
	tests := []struct {
		name string
		u    User
		args args
		want User
	}{
		{"ok", User{drive: new(Mock)}, args{"dinos80152@gmail.com"},
			User{Profile: Profile{Account: "dinos80152@gmail.com"}, drive: new(Mock)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Find(tt.args.account); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
