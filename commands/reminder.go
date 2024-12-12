package commands

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/ibexmonj/ContribSync/config"
)

func StartReminder(reminderTime string) {
	fmt.Println("Reminder is running. Press Ctrl+C to stop.")

	for {
		now := time.Now()
		currentTime := now.Format("15:04") // Format time as HH:MM

		if currentTime == reminderTime {
			// Print to terminal
			fmt.Printf("Reminder: It's %s! Don't forget to log your contributions.\n", reminderTime)

			// Send desktop notification
			err := SendMacNotification("Contribution Reminder", "Don't forget to log your contributions!")
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

func ReminderCommand(config *config.Config) {
	fmt.Println("Starting the reminder service...")
	StartReminder(config.Reminder.Time)
}

// SendMacNotification sends a desktop notification on macOS
func SendMacNotification(title, message string) error {
	// AppleScript command for sending a notification
	notification := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	cmd := exec.Command("osascript", "-e", notification)
	return cmd.Run()
}
