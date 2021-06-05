package messenger

import (
	"testing"
)

func TestMessenger_SetGreetingText(t *testing.T) {
	type fields struct {
		VerifyToken string
		AccessToken string
	}
	type args struct {
		greetingStrings []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "test", fields: fields{"", ""}, args: args{[]string{"a", "b"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Messenger{
				VerifyToken: tt.fields.VerifyToken,
				AccessToken: tt.fields.AccessToken,
			}
			m.SetGreetingText(tt.args.greetingStrings)
		})
	}
}
