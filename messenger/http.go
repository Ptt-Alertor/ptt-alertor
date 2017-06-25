package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MsgErr struct {
	ErrorBody `json:"error"`
}

type ErrorBody struct {
	Message      string `json:"message"`
	Type         string `json:"type"`
	Code         int    `json:"code"`
	ErrorSubCode int    `json:"error_subcode,omitempty"`
	FbtraceID    string `json:"fbtrace_id"`
}

func callAPI(url string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		msgErr := &MsgErr{}
		json.Unmarshal(body, &msgErr)
		if msgErr.Code == 551 {
			return nil
		}
		return fmt.Errorf("%s(%d): %s", res.Status, res.StatusCode, string(body))
	}
	return nil
}
