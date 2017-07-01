package messenger

import "encoding/json"

type GetStartedObject struct {
	GetStarted `json:"get_started"`
}

type GetStarted struct {
	Payload string `json:"payload"`
}

// SetGetStartedButton setting messenger started button
// payload: GET_STARTED_PAYLOAD
func (m *Messenger) SetGetStartedButton(payload string) {
	obj := GetStartedObject{
		GetStarted{
			Payload: payload,
		},
	}
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	m.callProfileAPI(data)
}

func (m *Messenger) callProfileAPI(body interface{}) {
	url := ProfileURL + m.AccessToken
	callAPI(url, body)
}
