package commands

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"os/exec"
	"time"

	"github.com/ibexmonj/ContribSync/config"
)

// SendDesktopNotification sends a cross-platform desktop notification
func SendDesktopNotification(title, message string) error {
	return beeep.Notify(title, message, "") //  icon path can be empty
}

func StartReminder(cfg *config.Config) {
	fmt.Println("Reminder is running. Press Ctrl+C to stop.")

	for {
		now := time.Now()
		currentTime := now.Format("15:04") // Format time as HH:MM

		if currentTime == cfg.Reminder.Time {
			// Print to terminal
			fmt.Printf("Reminder: It's %s %s!", cfg.Reminder.Time, cfg.Reminder.Message)

			// Send desktop notification
			err := SendDesktopNotification(cfg.Reminder.Title, cfg.Reminder.Message)
			if err != nil {
				fmt.Printf("Failed to send notification: %v\n", err)
			}

			// Avoid duplicate reminders in the same minute
			time.Sleep(60 * time.Second)
		} else {
			// Check every 10 seconds
			time.Sleep(10 * time.Second)
		}
	}
}

func ReminderCommand(cfg *config.Config, args []string) {
	if len(args) > 0 && args[0] == "test" {
		TestReminder(cfg)
		return
	}

	fmt.Println("Starting the reminder service...")
	StartReminder(cfg)
}

// SendMacNotification sends a desktop notification on macOS
func SendMacNotification(title, message string) error {
	// AppleScript command for sending a notification
	notification := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	cmd := exec.Command("osascript", "-e", notification)
	return cmd.Run()
}

// TestReminder sends a test notification using the current configuration
func TestReminder(cfg *config.Config) {
	fmt.Println("Sending test notification...")

	err := SendDesktopNotification(cfg.Reminder.Title, cfg.Reminder.Message)
	if err != nil {
		fmt.Printf("Failed to send notification: %v\n", err)
	} else {
		fmt.Println("Test notification sent successfully!")
	}
}
