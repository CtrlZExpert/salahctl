package main

import "fmt"

func validateArgCount(actual int, expected int, usage string) bool {
	if actual != expected {
		fmt.Println()
		fmt.Println(usage)
		fmt.Println()
		return false
	}
	return true

}

func usageFor(command string) string {
	if command == "date" {
		usage := fmt.Sprintf("Usage: salahctl %s YYYY-MM-DD", command)
		return usage
	}
	if command == "prayer" {
		usage := fmt.Sprintf("Usage: salahctl prayer <fajr|dhuhr|asr|maghrib|isha>")
		return usage
	}
	if command == "config" {
		return "Usage: salahctl config [show|location|method|asr]"
	}
	if command == "-h" || command == "--help" {
		return "Usage: salahctl --help"
	}
	if command == "-v" || command == "--version" {
		return "Usage: salahctl --version"
	}

	usage := fmt.Sprintf("Usage: salahctl %s", command)
	return usage
}

func showHelp() {
	fmt.Print(`salahctl - Prayer times from the command line

Usage:
  salahctl <command>

Commands:
  today    		Show today's prayer times
  tomorrow		Show tomorrow's prayer times
  current		Show the current prayer
  next			Show the next prayer and countdown
  week			Show prayer times for the next 7 days
  date YYYY-MM-DD	Show prayer times for a specific date
  config		Run full configuration setup
  config show		Show current configuration
  config location	Update location
  config method		Update calculation method
  config asr		Update Asr method
  remaining		Show remaining prayers for today
  prayer <name>		Show the time for a specific prayer

Options:
  -h, --help		Show this help
  -v, --version		Show version
`)
}

func showVersion() {
	fmt.Println("Version: 0.1.0")
	fmt.Println()
}
