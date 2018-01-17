package user

import (
	"os"
	"reflect"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/garyburd/redigo/redis"
)

var s *miniredis.Miniredis
var err error

func TestMain(m *testing.M) {
	// setup
	s, err = miniredis.Run()
	if err != nil {
		panic(err)
	}

	connectRedis = func() redis.Conn {
		conn, err := redis.Dial("tcp", s.Addr())
		if err != nil {
			panic(err)
		}
		return conn
	}

	// run test
	v := m.Run()

	// teardown
	s.Close()
	os.Exit(v)
}

func TestRedis_List(t *testing.T) {
	s.Set("user:dinos80152", `{"account":"dinos80152"}`)

	tests := []struct {
		name         string
		r            Redis
		wantAccounts []string
	}{
		{"dinos80152", Redis{}, []string{"dinos80152"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAccounts := tt.r.List(); !reflect.DeepEqual(gotAccounts, tt.wantAccounts) {
				t.Errorf("Redis.List() = %v, want %v", gotAccounts, tt.wantAccounts)
			}
		})
	}
}

func TestRedis_Exist(t *testing.T) {
	s.Set("user:dinos80152", `{"account":"dinos80152"}`)

	type args struct {
		account string
	}
	tests := []struct {
		name string
		r    Redis
		args args
		want bool
	}{
		{"dinos80152", Redis{}, args{"dinos80152"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Exist(tt.args.account); got != tt.want {
				t.Errorf("Redis.Exist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_Save(t *testing.T) {
	type args struct {
		account string
		data    interface{}
	}
	tests := []struct {
		name    string
		r       Redis
		args    args
		wantErr bool
	}{
		{"ok", Redis{}, args{"dinos80152", `{"account":"dinos80152"}`}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Save(tt.args.account, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Redis.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedis_Update(t *testing.T) {
	type args struct {
		account string
		user    interface{}
	}
	tests := []struct {
		name    string
		r       Redis
		args    args
		wantErr bool
	}{
		{"ok", Redis{}, args{"dinos80152", User{
			Enable: true,
			Profile: Profile{
				Account: "dinos80152",
			}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Update(tt.args.account, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Redis.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedis_Find(t *testing.T) {
	type args struct {
		account string
		user    *User
	}
	tests := []struct {
		name string
		r    Redis
		args args
	}{
		{"ok", Redis{}, args{"dinos80152", &User{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Find(tt.args.account, tt.args.user)
		})
	}
}
