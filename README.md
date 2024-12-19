# Notification Package

This notification package provides a simple interface to send notifications to Slack and Telegram. It also prevents
repeated notifications within a defined cooldown period (default: 10 minutes).

---

## Features

1. **Environment-specific configuration:**
    - Application name and version are loaded automatically.
    - Environment variables control notification channels.

2. **Slack and Telegram integration:**
    - Send messages to Slack and Telegram channels.

3. **Duplicate message prevention:**
    - Tracks sent messages in memory to avoid repeated notifications within a cooldown period.

4. **Notification levels:**
    - `Error`
    - `Warning`
    - `Info`

---

## Installation

1. **Add the package to your project:**

```shell
go get github.com/Tsisar/notification-go
```

2. **Import the package:**

```shell
import "github.com/Tsisar/notification-go/notification"
```

---

## Environment Variables

**Required for Telegram**

- TELEGRAM_NOTIFICATION: Enable Telegram notifications (true or false).
- TELEGRAM_CHAT_ID: Chat ID where messages will be sent.
- TELEGRAM_BOT_TOKEN: Bot token for authentication.

**Required for Slack**

- SLACK_NOTIFICATION: Enable Slack notifications (true or false).
- SLACK_CHANNEL: Slack webhook URL.

**General Configuration**

- APP_NAME: Name of the application.
- ENVIRONMENT: Environment name (e.g., production, development).

---

## Usage

### Sending Notifications

1. Error Notification

```go
notification.Error("An unexpected error occurred.")
```

2. Warning Notification

```go
notification.Warning("This is a warning message.")
```

3. Info Notification

```go
notification.Info("System is running as expected.")
```

### Example Application

package main

```go
import (
    "notification"
)

func main() {
   notification.Info("Application has started successfully.")
   notification.Warning("Low disk space detected.")
   notification.Error("Database connection failed.")
}

```

### Message Cooldown
• Default TTL: 10 minutes (messageTTL).
• Prevents duplicate notifications during the cooldown period.

Integration with Slack and Telegram

### Slack

1. Set up a Slack webhook URL.
2. Configure the following environment variables:
   • SLACK_NOTIFICATION=true
   • SLACK_CHANNEL=
3. Example message:
```go
   notification.Info("System update completed successfully.")
```


### Telegram

1. Create a Telegram bot and get its token.
2. Retrieve your chat ID using the bot.
3. Configure the following environment variables:
   • TELEGRAM_NOTIFICATION=true
   • TELEGRAM_CHAT_ID=
   • TELEGRAM_BOT_TOKEN=
4. Example message:
```go
   notification.Error("Service is down!")
   ```

## File Structure

```plaintext
notification/
├── notification.go # Core logic for notifications
├── slack.go # Slack integration
├── telegram.go # Telegram integration
```

---
## Notes

1. In-memory tracking: Message statuses are stored in memory. They will reset upon application restart.
2. Error handling: Ensure the environment variables are configured correctly to avoid runtime errors.
3. Scalability: For large-scale applications, consider external storage like Redis for message tracking.

Enjoy using the notification package!