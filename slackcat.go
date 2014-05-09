package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type SlackMsg struct {
	//	Username string `json:"username"`
	Text string `json:"text"`
}

func (m SlackMsg) Encode() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (m SlackMsg) Post(WebhookURL string) error {
	encoded, err := m.Encode()
	if err != nil {
		return err
	}

	resp, err := http.PostForm(WebhookURL, url.Values{"payload": {encoded}})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Not OK")
	}
	return nil
}

func main() {

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Cannot read STDIN: %v", err)
	}

	msg := SlackMsg{
		Text: string(bytes),
	}

	err = msg.Post("https://ph.slack.com/services/hooks/incoming-webhook?token=VxHgickriL7nFVYNAxbCOmba")
	if err != nil {
		log.Fatalf("Post failed: %v", err)
	}
}
