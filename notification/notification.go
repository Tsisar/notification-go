package notification

import (
	"github.com/Tsisar/extended-log-go/log"
	"os"
	"strings"
	"sync"
	"time"

	"fmt"
)

const messageTTL = 10 * time.Minute

var version string
var appName string
var environment string
var messageStatus sync.Map

func init() {
	version = getVersion("/app/VERSION")
	appName = getStringEnv("APP_NAME", "")
	environment = getStringEnv("ENVIRONMENT", "")
}

// Error sends an error message to Slack Channel
func Error(message string) {
	sendMsg(fmt.Sprintf("🔴 ERROR\n%s", message))
}

// Warning sends a warning message to Slack Channel
func Warning(message string) {
	sendMsg(fmt.Sprintf("🟡 WARNING\n%s", message))
}

// Info sends an info message to Slack Channel
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

// checkMessageSentStatus checks if a message has been sent and refreshes its cooldown.
func checkMessageSentStatus(message string) bool {
	now := time.Now()

	// Check if the message has been sent recently
	lastSent, exists := messageStatus.Load(message)

	if exists {
		lastSentTime := lastSent.(time.Time)
		if now.Sub(lastSentTime) < messageTTL {
			return true // Message has been sent recently
		}
	}

	// Update the message status
	messageStatus.Store(message, now)
	return false
}

func getVersion(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading version file: %v", err)
	}

	return strings.TrimSpace(string(data))
}
