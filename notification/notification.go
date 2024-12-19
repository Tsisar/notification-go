package notification

import (
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"os"
	"strings"
	"time"
)

const messageTTL = 60 * time.Minute
const messageStatusFile = "/app/message_status.json"

var version string
var appName string
var environment string
var messageStatus = make(map[string]time.Time)

func init() {
	version = getVersion("/app/VERSION")
	appName = getStringEnv("APP_NAME", "")
	environment = getStringEnv("ENVIRONMENT", "")

	loadMessageStatus()
}

func Error(message string) {
	sendMsg(fmt.Sprintf("🔴 ERROR\n%s", message))
}

func Warning(message string) {
	sendMsg(fmt.Sprintf("🟡 WARNING\n%s", message))
}

func Info(message string) {
	sendMsg(fmt.Sprintf("🟢 INFO\n%s", message))
}

func sendMsg(message string) {
	if checkMessageSentStatus(message) {
		return
	}
	app := fmt.Sprintf("%s version %s", appName, version)
	message = fmt.Sprintf("%s\n%s\nenv: %s", message, app, environment)

	err := Telegram.SendMessage(message)
	if err != nil {
		log.Errorf("Error sending message to Telegram: %v", err)
	}

	err = Slack.SendMessage(message)
	if err != nil {
		log.Errorf("Error sending message to Slack: %v", err)
	}
}

func checkMessageSentStatus(message string) bool {
	now := time.Now()

	lastSent, exists := messageStatus[message]
	if exists && now.Sub(lastSent) < messageTTL {
		return true
	}

	messageStatus[message] = now
	saveMessageStatus()

	return false
}

func loadMessageStatus() {
	data, err := os.ReadFile(messageStatusFile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Errorf("Error reading message status file: %v", err)
		}
		return
	}

	err = json.Unmarshal(data, &messageStatus)
	if err != nil {
		log.Errorf("Error unmarshalling message status: %v", err)
	}
}

func saveMessageStatus() {
	data, err := json.Marshal(messageStatus)
	if err != nil {
		log.Errorf("Error marshalling message status: %v", err)
		return
	}

	err = os.WriteFile(messageStatusFile, data, 0644)
	if err != nil {
		log.Errorf("Error writing message status file: %v", err)
	}
}

func getVersion(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Errorf("Error reading version file: %v", err)
		data = []byte("unknown")
	}

	return strings.TrimSpace(string(data))
}
