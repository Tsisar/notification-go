package notification

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"os"
	"path/filepath"
	"time"
)

var messageTTL time.Duration
var messageStatusFile string
var version string
var appName string
var environment string
var messageStatus = make(map[string]time.Time)

func init() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Cannot determine executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)
	messageStatusFile = filepath.Join(exeDir, "message_status.json")

	appName = getStringEnv("APP_NAME", "")
	version = getStringEnv("VERSION", "unknown")
	environment = getStringEnv("ENVIRONMENT", "")

	ttl := getIntEnv("MESSAGE_TTL", 60)
	messageTTL = time.Duration(ttl) * time.Minute

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

func hashMessage(msg string) string {
	sum := sha256.Sum256([]byte(msg))
	return fmt.Sprintf("%x", sum)
}

func sendMsg(message string) {
	hash := hashMessage(message)

	if checkMessageSentStatus(hash) {
		return
	}
	app := fmt.Sprintf("%s version %s", appName, version)
	messageToSend := fmt.Sprintf("%s\n%s\nenv: %s", message, app, environment)

	err := Telegram.SendMessage(messageToSend)
	if err != nil {
		log.Errorf("Error sending message to Telegram: %v", err)
	}

	err = Slack.SendMessage(messageToSend)
	if err != nil {
		log.Errorf("Error sending message to Slack: %v", err)
	}
}

func checkMessageSentStatus(msgHash string) bool {
	now := time.Now()

	lastSent, exists := messageStatus[msgHash]
	if exists && now.Sub(lastSent) < messageTTL {
		messageStatus[msgHash] = now
		saveMessageStatus()
		return true
	}

	messageStatus[msgHash] = now
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
