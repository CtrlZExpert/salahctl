package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Error: missing a command")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  salahctl <command>")
		fmt.Println()
		fmt.Println("Run 'salahctl --help' for more information")
		return
	}
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	command := os.Args[1]
	usage := usageFor(command)
	switch command {
	case "today":
		isValidCount := validateArgCount(len(os.Args), 2, usage)
		if !isValidCount {
			return
		}
		showPrayerTimes(config)
	case "config":
		if len(os.Args) == 2 {
			runConfig()
			return

		}

		isValidCount := validateArgCount(len(os.Args), 3, usage)
		if !isValidCount {
			return
		}

		configCommand := os.Args[2]
		switch configCommand {
		case "show":
			showConfig(config)
		case "location":
			updateLocation(config)
		case "method":
			updateCalculationMethod(config)
		case "asr":
			updateAsrMethod(config)
		default:
			fmt.Printf("Error: unknown config command %q\n", configCommand)
			fmt.Println()
			fmt.Println("Usage: salahctl config [show|location|method|asr]")
		}
	case "current":
		isValidCount := validateArgCount(len(os.Args), 2, usage)
		if !isValidCount {
			return
		}
		showCurrentPrayer(config)
	case "next":
		isValidCount := validateArgCount(len(os.Args), 2, usage)
		if !isValidCount {
			return
		}
		showNextPrayer(config)
	case "tomorrow":
		isValidCount := validateArgCount(len(os.Args), 2, usage)
		if !isValidCount {
			return
		}
		showTomorrowPrayersTimes(config)
	case "date":
		isValidCount := validateArgCount(len(os.Args), 3, usage)
		if !isValidCount {
			return
		}

		showPrayerTimesByDate(config, os.Args[2])
	case "remaining":
		isValidCount := validateArgCount(len(os.Args), 2, usage)
		if !isValidCount {
			return
		}
		showRemainingPrayers(config)
	case "prayer":
		isvalidcount := validateArgCount(len(os.Args), 3, usage)
		if !isvalidcount {
			return
		}
		prayerName := os.Args[2]
		showPrayer(config, prayerName)
	case "week":
		isValidCount := validateArgCount(len(os.Args), 2, usage)
		if !isValidCount {
			return
		}
		showWeeklyPrayerTimes(config)

	case "--help", "-h":
		isValidCount := validateArgCount(len(os.Args), 2, usage)
		if !isValidCount {
			return
		}
		showHelp()
	case "--version", "-v":
		isValidCount := validateArgCount(len(os.Args), 2, "Usage: salahctl --version")
		if !isValidCount {
			return
		}
		showVersion()
	default:
		fmt.Printf("Error: unknown command %q\n", command)
		fmt.Println()
		fmt.Println("Run 'salahctl --help' for more information")
	}

}
