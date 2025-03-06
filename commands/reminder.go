package commands

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/ibexmonj/ContribSync/pkg/logger"
	"github.com/spf13/cobra"
	"os/exec"
	"runtime"
	"time"

	"github.com/ibexmonj/ContribSync/config"
)

func SendDesktopNotification(title, message string) error {
	logger.Logger.Info().Str("title", title).Str("message", message).Msg("Attempting to send notification")

	err := beeep.Notify(title, message, "")
	if err != nil {
		logger.Logger.Error().Err(err).Msg("❌ Failed to send desktop notification using beeep")

		if runtime.GOOS == "darwin" {
			logger.Logger.Warn().Msg("Using macOS fallback notification")
			err = SendMacNotification(title, message)
			if err != nil {
				logger.Logger.Error().Err(err).Msg("❌ Failed to send macOS notification")
			}
		}
	}

	return err
}

func StartReminder() {
	logger.Logger.Info().Msg("⏰ Reminder service started. Press Ctrl+C to stop.")

	for {
		now := time.Now()
		currentTime := now.Format("15:04") // Format time as HH:MM

		if currentTime == config.ConfigData.Reminder.Time {
			logger.Logger.Info().
				Str("Time", config.ConfigData.Reminder.Time).
				Str("Message", config.ConfigData.Reminder.Message).
				Msg("Triggering reminder")

			fmt.Printf("\n📢 Reminder: It's %s - %s\n", config.ConfigData.Reminder.Time, config.ConfigData.Reminder.Message)

			err := SendDesktopNotification(config.ConfigData.Reminder.Title, config.ConfigData.Reminder.Message)
			if err != nil {
				logger.Logger.Error().Err(err).Msg("Failed to send notification")
				fmt.Printf("❌ Failed to send notification: %v\n", err)
			}

			time.Sleep(60 * time.Second) // Wait to avoid sending notifications every second
		} else {
			time.Sleep(10 * time.Second) // Check every 10 seconds
		}
	}
}

func ReminderCommand(args []string) {
	if len(args) > 0 && args[0] == "test" {
		TestReminder()
		return
	}

	fmt.Println("🔔 Starting the reminder service...")
	StartReminder()
}

func SendMacNotification(title, message string) error {
	notification := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	cmd := exec.Command("osascript", "-e", notification)
	err := cmd.Run()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("❌ Failed to send macOS notification via AppleScript")
	}
	return err
}

func TestReminder() {
	fmt.Println("📢 Sending test notification...")

	err := SendDesktopNotification(config.ConfigData.Reminder.Title, config.ConfigData.Reminder.Message)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to send test notification")
		fmt.Printf("❌ Failed to send notification: %v\n", err)
	} else {
		fmt.Println("✅ Test notification sent successfully!")
	}
}

func NewReminderCommand() *cobra.Command {
	var reminderCmd = &cobra.Command{
		Use:   "reminder",
		Short: "Manage reminders",
		Long:  "Start or test the reminder service.",
	}

	reminderCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start the reminder service",
		Run: func(cmd *cobra.Command, args []string) {
			StartReminder()
		},
	})

	reminderCmd.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "Send a test reminder notification",
		Run: func(cmd *cobra.Command, args []string) {
			TestReminder()
		},
	})

	return reminderCmd
}
