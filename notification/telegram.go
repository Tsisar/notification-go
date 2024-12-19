package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"io"
	"net/http"
)

// Telegram API URL for sending messages
const telegramAPIBaseURL = "https://api.telegram.org/bot"

var Telegram *telegram

type telegram struct {
	Enable bool   `json:"enable"`
	ChatID string `json:"chatID"`
	Token  string `json:"token"`
}

func init() {
	Telegram = &telegram{
		Enable: getBoolEnv("TELEGRAM_NOTIFICATION", false),
		ChatID: getStringEnv("TELEGRAM_CHAT_ID", ""),
		Token:  getStringEnv("TELEGRAM_BOT_TOKEN", ""),
	}
}

// SendMessage sends a message to a Telegram chat
func (t *telegram) SendMessage(message string) error {
	if !t.Enable {
		return nil
	}

	url := fmt.Sprintf("%s%s/sendMessage", telegramAPIBaseURL, t.Token)
	msg := map[string]interface{}{
		"chat_id": t.ChatID,
		"text":    message,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not marshal message: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Errorf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}
