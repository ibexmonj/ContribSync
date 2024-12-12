package commands

import "fmt"

func Help() {
	fmt.Println("csync - Contribution Sync CLI")
	fmt.Println("Available commands:")
	fmt.Println("  help          - Show this help message")
	fmt.Println("  config        - Show or modify configuration")
	fmt.Println("  reminder      - Start the daily reminder service")
}
