package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/MSA-Software-LLC/adhan-go/pkg/calc"
	"github.com/MSA-Software-LLC/adhan-go/pkg/data"
	"github.com/MSA-Software-LLC/adhan-go/pkg/util"
	"github.com/ringsaturn/tzf"
)

type Config struct {
	Latitude  float64 `toml:"latitude"`
	Longitude float64 `toml:"longitude"`
	Timezone  string  `toml:"timezone"`
	Method    string  `toml:"method"`
	AsrMethod string  `toml:"asr_method"`
}

type LocationResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

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
func loadConfig() (Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, err
	}
	configPath := filepath.Join(configDir, "salahctl", "config.toml")
	config := Config{}

	_, err = toml.DecodeFile(configPath, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func saveConfig(config Config) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	salahctlDir := filepath.Join(configDir, "salahctl")
	configPath := filepath.Join(salahctlDir, "config.toml")

	err = os.MkdirAll(salahctlDir, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}
	return nil
}

func parseCalculationMethod(method string) (calc.CalculationMethod, error) {
	switch method {
	case "muslim_world_league":
		return calc.MUSLIM_WORLD_LEAGUE, nil
	case "north_america":
		return calc.NORTH_AMERICA, nil
	default:
		return 0, fmt.Errorf("unsupported calculation method: %q ", method)
	}
}

func parseAsrMethod(method string) (calc.AsrJuristicMethod, error) {
	switch method {
	case "standard":
		return calc.SHAFI_HANBALI_MALIKI, nil
	case "hanafi":
		return calc.HANAFI, nil
	default:
		return 0, fmt.Errorf("unsupported calculation method: %q ", method)
	}
}

func lookupLocation(location string) ([]LocationResult, error) {
	query := url.QueryEscape(location)
	requestURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=5", query)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "salahctl/1.0")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("location lookup failed: HTTP %d", response.StatusCode)
	}
	var results []LocationResult

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&results)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("location not found: %q", location)
	}

	return results, nil
}

func chooseCalculationMethod() string {
	for {
		fmt.Println()
		fmt.Println("Calculation Method")
		fmt.Println("  1. Muslim World League")
		fmt.Println("  2. North America")
		fmt.Println()
		fmt.Print("Choose method: ")

		var choiceAction string
		fmt.Scan(&choiceAction)

		validChoice, err := strconv.Atoi(choiceAction)
		if err != nil {
			fmt.Println("Invalid selection. Enter a number 1-2")
			continue
		}

		switch validChoice {
		case 1:
			return "muslim_world_league"
		case 2:
			return "north_america"
		default:
			fmt.Println("Invalid selection. Enter a number 1-2")
			continue
		}

	}
}

func chooseAsrMethod() string {
	for {
		fmt.Println()
		fmt.Println("Asr Method")
		fmt.Println("  1. Standard")
		fmt.Println("  2. Hanafi")
		fmt.Println()
		fmt.Print("Choose method: ")

		var choiceAction string
		fmt.Scan(&choiceAction)

		validChoice, err := strconv.Atoi(choiceAction)
		if err != nil {
			fmt.Println("Invalid selection. Enter a number 1-2")
			continue
		}

		switch validChoice {
		case 1:
			return "standard"
		case 2:
			return "hanafi"
		default:
			fmt.Println("Invalid selection. Enter a number 1-2")
			continue
		}

	}
}

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
		if len(os.Args) != 3 {
			fmt.Println("Usage: salahctl date YYYY-MM-DD")
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
		isvalidcount := validateArgCount(len(os.Args), 3, "Usage: salahctl prayer <fajr|dhuhr|asr|maghrib|isha>")
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

func calculatePrayerTimeForDate(c Config, date time.Time) (*calc.PrayerTimes, error) {

	dateComponents := data.NewDateComponents(date)
	coordinate, err := util.NewCoordinates(c.Latitude, c.Longitude)
	if err != nil {
		return nil, err
	}
	method, err := parseCalculationMethod(c.Method)
	if err != nil {
		return nil, err
	}
	params := calc.GetMethodParameters(method)
	asrMethod, err := parseAsrMethod(c.AsrMethod)
	if err != nil {
		return nil, err
	}

	params.Madhab = asrMethod

	prayerTimes, err := calc.NewPrayerTimes(coordinate, dateComponents, params)
	if err != nil {
		return nil, err
	}

	err = prayerTimes.SetTimeZone(c.Timezone)
	if err != nil {
		return nil, err
	}
	return prayerTimes, nil
}

func calculatePrayerTimes(c Config) (*calc.PrayerTimes, error) {
	now := time.Now()
	prayerTimes, err := calculatePrayerTimeForDate(c, now)
	if err != nil {
		return prayerTimes, err
	}
	return prayerTimes, nil
}

func getNextPrayer(c Config) (calc.Prayer, time.Time, time.Duration, error) {
	prayerTimes, err := calculatePrayerTimes(c)
	if err != nil {
		return calc.NO_PRAYER, time.Time{}, 0, err
	}
	now := time.Now()
	next := prayerTimes.NextPrayer(now)
	nextTime := prayerTimes.TimeForPrayer(next)
	if next == calc.NO_PRAYER {
		tomorrow := now.AddDate(0, 0, 1)
		tomorrowPrayerTimes, err := calculatePrayerTimeForDate(c, tomorrow)
		if err != nil {
			return calc.NO_PRAYER, time.Time{}, 0, err
		}
		next = calc.FAJR
		nextTime = tomorrowPrayerTimes.Fajr
	}
	remaining := nextTime.Sub(now).Truncate(time.Minute)

	return next, nextTime, remaining, nil
}

func showCurrentPrayer(c Config) {
	prayerTimes, err := calculatePrayerTimes(c)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	now := time.Now()
	current := prayerTimes.CurrentPrayer(now)
	next, nextTime, remaining, err := getNextPrayer(c)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	hours := int(remaining.Hours())
	minutes := int(remaining.Minutes()) % 60
	fmt.Println()
	fmt.Printf("Current Prayer: %s\n", prayerName(current))
	fmt.Printf("Next Prayer: %s at %s\n", prayerName(next), nextTime.Format("3:04 PM"))
	fmt.Printf("Time remaining: %dh %dm\n", hours, minutes)
}

func showNextPrayer(c Config) {
	next, nextTime, remaining, err := getNextPrayer(c)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	hours := int(remaining.Hours())
	minutes := int(remaining.Minutes()) % 60
	fmt.Println()
	fmt.Printf("Next Prayer: %s at %s\n", prayerName(next), nextTime.Format("3:04 PM"))
	fmt.Printf("Time remaining: %dh %dm\n", hours, minutes)
}

func showPrayer(config Config, prayerName string) {
	prayerTimes, err := calculatePrayerTimes(config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var selectedTime time.Time
	switch prayerName {
	case "fajr":
		selectedTime = prayerTimes.Fajr
	case "dhuhr":
		selectedTime = prayerTimes.Dhuhr
	case "asr":
		selectedTime = prayerTimes.Asr
	case "maghrib":
		selectedTime = prayerTimes.Maghrib
	case "isha":
		selectedTime = prayerTimes.Isha
	default:
		fmt.Println()
		fmt.Printf("Error: unknown prayer %q\n", prayerName)
		fmt.Println()
		fmt.Println("Usage: salahctl prayer <fajr|dhuhr|asr|maghrib|isha>")
		return
	}
	displayName := strings.ToUpper(prayerName[:1]) + prayerName[1:]
	fmt.Println(displayName, selectedTime.Format("3:04 PM"))
}

func showTomorrowPrayersTimes(c Config) {
	tomorrow := time.Now().AddDate(0, 0, 1)
	prayerTimesTomorrow, err := calculatePrayerTimeForDate(c, tomorrow)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println()
	fmt.Println("Tomorrow")
	fmt.Println()
	printPrayerTimes(prayerTimesTomorrow)
}

func showPrayerTimesByDate(c Config, dateString string) {

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		fmt.Println("Invalid date format. Use YYYY-MM-DD")
		return
	}
	prayerTimesByDate, err := calculatePrayerTimeForDate(c, date)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println()
	fmt.Println(date.Format("Monday, January 2, 2006"))
	fmt.Println()
	printPrayerTimes(prayerTimesByDate)

}

func showWeeklyPrayerTimes(c Config) {
	date := time.Now()

	for i := 0; i < 7; i++ {
		prayerTimes, err := calculatePrayerTimeForDate(c, date)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println(date.Format("Monday, January 2, 2006"))
		fmt.Println()
		printPrayerTimes(prayerTimes)
		fmt.Println()

		date = date.AddDate(0, 0, 1)
	}
}

func showRemainingPrayers(config Config) {
	prayerTimes, err := calculatePrayerTimes(config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	now := time.Now()

	fmt.Println("Remaining Prayers")
	fmt.Println()

	if now.Before(prayerTimes.Fajr) {
		fmt.Println("Fajr:", prayerTimes.Fajr.Format("3:04 PM"))
	}
	if now.Before(prayerTimes.Dhuhr) {
		fmt.Println("Dhuhr:", prayerTimes.Dhuhr.Format("3:04 PM"))
	}
	if now.Before(prayerTimes.Asr) {
		fmt.Println("Asr:", prayerTimes.Asr.Format("3:04 PM"))
	}
	if now.Before(prayerTimes.Maghrib) {
		fmt.Println("Maghrib:", prayerTimes.Maghrib.Format("3:04 PM"))
	}
	if now.Before(prayerTimes.Isha) {
		fmt.Println("Isha", prayerTimes.Isha.Format("3:04 PM"))
	}

	if now.After(prayerTimes.Isha) {
		fmt.Println("All prayers are complete for today.")
	}
	fmt.Println()

}

func printPrayerTimes(prayerTimes *calc.PrayerTimes) {
	fmt.Printf("Fajr:    %s\n", prayerTimes.Fajr.Format("3:04 PM"))
	fmt.Printf("Sunrise: %s\n", prayerTimes.Sunrise.Format("3:04 PM"))
	fmt.Printf("Dhuhr:   %s\n", prayerTimes.Dhuhr.Format("3:04 PM"))
	fmt.Printf("Asr:     %s\n", prayerTimes.Asr.Format("3:04 PM"))
	fmt.Printf("Maghrib: %s\n", prayerTimes.Maghrib.Format("3:04 PM"))
	fmt.Printf("Isha:    %s\n", prayerTimes.Isha.Format("3:04 PM"))

}

func prayerName(prayer calc.Prayer) string {
	switch prayer {
	case calc.FAJR:
		return "Fajr"
	case calc.SUNRISE:
		return "None"
	case calc.DHUHR:
		return "Dhuhr"
	case calc.ASR:
		return "Asr"
	case calc.MAGHRIB:
		return "Maghrib"
	case calc.ISHA:
		return "Isha"
	case calc.NO_PRAYER:
		return "None"
	default:
		return "Unknown"

	}
}

func showPrayerTimes(c Config) {
	prayerTimes, err := calculatePrayerTimes(c)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println()
	printPrayerTimes(prayerTimes)
}
func showConfig(config Config) {
	fmt.Printf("Latitude:   %f\n", config.Latitude)
	fmt.Printf("Longitude:   %f\n", config.Longitude)
	fmt.Printf("Timezone:    %s\n", config.Timezone)
	fmt.Printf("Method:     %s\n", config.Method)
	fmt.Printf("Asr Method: %s\n", config.AsrMethod)
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

func runConfig() {
	latitude, longitude, timezone, err := chooseLocation()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	method := chooseCalculationMethod()
	asrMethod := chooseAsrMethod()

	config := Config{
		Latitude:  latitude,
		Longitude: longitude,
		Timezone:  timezone,
		Method:    method,
		AsrMethod: asrMethod,
	}

	err = saveConfig(config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Configuration saved successfully")

}

func chooseLocation() (float64, float64, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Location (city, state/country): ")

	location, err := reader.ReadString('\n')
	if err != nil {
		return 0.0, 0.0, "", err

	}
	location = strings.TrimSpace(location)
	results, err := lookupLocation(location)
	if err != nil {
		return 0.0, 0.0, "", err

	}

	for i, result := range results {
		fmt.Printf("%d. %s\n", i+1, result.DisplayName)
	}

	fmt.Print("Choose location: ")
	choiceInput, err := reader.ReadString('\n')
	if err != nil {
		return 0.0, 0.0, "", err
	}
	choiceInput = strings.TrimSpace(choiceInput)
	choice, err := strconv.Atoi(choiceInput)
	if err != nil {
		return 0.0, 0.0, "", err
	}
	if choice < 1 || choice > len(results) {
		return 0.0, 0.0, "", errors.New("Invalid location choice")
	}
	selected := results[choice-1]

	latitude, err := strconv.ParseFloat(selected.Lat, 64)
	if err != nil {
		return 0.0, 0.0, "", err
	}
	longitude, err := strconv.ParseFloat(selected.Lon, 64)
	if err != nil {
		return 0.0, 0.0, "", err

	}
	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		return 0.0, 0.0, "", err
	}
	timezone := finder.GetTimezoneName(longitude, latitude)

	return latitude, longitude, timezone, nil
}

func updateLocation(config Config) {
	latitude, longitude, timezone, err := chooseLocation()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	config.Latitude = latitude
	config.Longitude = longitude
	config.Timezone = timezone

	err = saveConfig(config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Location has successfully been updated")

}

func updateCalculationMethod(config Config) {
	method := chooseCalculationMethod()
	config.Method = method
	err := saveConfig(config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Method has successfully been updated")
}

func updateAsrMethod(config Config) {
	asrMethod := chooseAsrMethod()
	config.AsrMethod = asrMethod
	err := saveConfig(config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Asr Method has successfully been updated")

}
