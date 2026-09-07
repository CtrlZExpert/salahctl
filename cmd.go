package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "salahctl",
	Short:   "Islamic prayer times from the terminal",
	Version: "0.1.0",
}

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show today's prayer times",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		showPrayerTimes(config)
	},
}

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current prayer",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return

		}

		showCurrentPrayer(config)
	},
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show the next prayer and countdown",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		showNextPrayer(config)
	},
}

var tomorrowCmd = &cobra.Command{
	Use:   "tomorrow",
	Short: "Show tomorrow's prayer times",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		showTomorrowPrayersTimes(config)
	},
}
var remainingCmd = &cobra.Command{
	Use:   "remaining",
	Short: "Show remaining prayers for today",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		showRemainingPrayers(config)
	},
}

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "Show prayer times for the next 7 days",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		showWeeklyPrayerTimes(config)
	},
}
var notifyCheckCmd = &cobra.Command{
	Use:    "notify-check",
	Short:  "Check whether a prayer notification is due",
	Args:   cobra.NoArgs,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		err = checkPrayerNotification(config)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	},
}

var prayerCmd = &cobra.Command{
	Use:   "prayer <name>",
	Short: "Show the time for specific prayer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		showPrayer(config, args[0])
	},
}

var dateCmd = &cobra.Command{
	Use:   "date YYYY-MM-DD",
	Short: "Show prayer time for a specific date",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		showPrayerTimesByDate(config, args[0])
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Run full configuration setup",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runConfig()
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		showConfig(config)
	},
}

var locationConfigCmd = &cobra.Command{
	Use:   "location",
	Short: "Update location",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		updateLocation(config)
	},
}

var methodConfigCmd = &cobra.Command{
	Use:   "method",
	Short: "Update calculation method",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		updateCalculationMethod(config)
	},
}

var asrMethodConfigCmd = &cobra.Command{
	Use:   "asr",
	Short: "Update Asr Method",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		updateAsrMethod(config)
	},
}

var locationCmd = &cobra.Command{
	Use:   "location",
	Short: "Manage location profile",
}

var locationAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add new location profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		err = addLocationProfile(config, args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	},
}

var locationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List location profiles",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		listLocationProfile(config)
	},
}

var locationUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Use location profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		err = useLocationProfile(config, args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	},
}

var locationRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove location profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		err = removeLocationProfile(config, args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	},
}

var monthCmd = &cobra.Command{
	Use:   "month",
	Short: "Show monthly prayer times",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		showMonthlyPrayerTimes(config)
	},
}

func init() {
	rootCmd.AddCommand(
		todayCmd,
		currentCmd,
		nextCmd,
		tomorrowCmd,
		remainingCmd,
		weekCmd,
		notifyCheckCmd,
		prayerCmd,
		dateCmd,
		configCmd,
		locationCmd,
		monthCmd,
	)

	configCmd.AddCommand(
		showConfigCmd,
		locationConfigCmd,
		methodConfigCmd,
		asrMethodConfigCmd,
	)

	locationCmd.AddCommand(
		locationAddCmd,
		locationListCmd,
		locationUseCmd,
		locationRemoveCmd,
	)
}
