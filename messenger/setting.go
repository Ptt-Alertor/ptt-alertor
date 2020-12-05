package messenger

import (
	"encoding/json"
)

type greeting struct {
	GreetingObjects []greetingObject
}

type greetingObject struct {
	Locale string `json:"locale"`
	Text   string `json:"text"`
}

// SetGreetingText set messenger greeting text
// https://developers.facebook.com/docs/messenger-platform/reference/messenger-profile-api/greeting
func (m *Messenger) SetGreetingText(greetingStrings []string) {
	greeting := greeting{}
	for _, str := range greetingStrings {
		obj := greetingObject{
			Locale: "default",
			Text:   str,
		}
		greeting.GreetingObjects = append(greeting.GreetingObjects, obj)
	}
	data, err := json.Marshal(greeting)
	if err != nil {
		panic(err)
	}
	m.callProfileAPI(data)
}

func (m *Messenger) callProfileAPI(body interface{}) {
	url := profileURL + m.AccessToken
	callAPI(url, body)
}
