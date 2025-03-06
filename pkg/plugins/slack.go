package plugins

import (
	"fmt"
	"github.com/ibexmonj/ContribSync/pkg/logger"
	"net/http"
	"os"
	"strings"
)

// SlackPlugin allows sending messages to Slack
type SlackPlugin struct{}

func (s *SlackPlugin) Init() error {
	logger.Logger.Info().Msg("✅ Slack plugin initialized")
	return nil
}

func (s *SlackPlugin) Info() (string, string) {
	return "slack", "Slack Plugin: Send messages to Slack channels"
}

func (s *SlackPlugin) Execute(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Usage: csync plugin exec slack send [channel] [message]")
	}

	if args[0] == "send" {
		channel := args[1]
		message := strings.Join(args[2:], " ")
		return sendSlackMessage(channel, message)
	}
	return fmt.Errorf("unknown Slack command: %s", args[0])
}

func sendSlackMessage(channel, message string) error {
	slackWebhook := os.Getenv("SLACK_WEBHOOK_URL")
	if slackWebhook == "" {
		return fmt.Errorf("SLACK_WEBHOOK_URL is not set. Please export your Slack webhook URL.")
	}

	payload := fmt.Sprintf(`{"channel": "%s", "text": "%s"}`, channel, message)
	req, err := http.NewRequest("POST", slackWebhook, strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create Slack request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Slack message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response from Slack: %s", resp.Status)
	}

	logger.Logger.Info().Str("channel", channel).Msg("✅ Message sent to Slack")
	fmt.Printf("📨 Message sent to Slack channel: %s\n", channel)
	return nil
}
