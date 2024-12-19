package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"net/http"
)

const slackAPIBaseURL = "https://hooks.slack.com/services"

var Slack *slack

type slack struct {
	Enable  bool   `json:"enable"`
	Channel string `json:"channel"`
}

func init() {
	Slack = &slack{
		Enable:  getBoolEnv("SLACK_NOTIFICATION", false),
		Channel: getStringEnv("SLACK_CHANNEL", ""),
	}
}

func (s *slack) SendMessage(message string) error {
	if !s.Enable {
		return nil
	}

	url := fmt.Sprintf("%s/%s", slackAPIBaseURL, s.Channel)
	payload := map[string]string{"text": message}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal message: %v", err)
	}
	body := bytes.NewReader(payloadBytes)

	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("error closing response body %v", err)
		}
	}()

	return nil
}
