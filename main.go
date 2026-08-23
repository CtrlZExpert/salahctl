package main

import (
	"bufio"
	"encoding/json"
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
		fmt.Println("Usage: salahctl <command>")
		return
	}
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	command := os.Args[1]
	switch command {
	case "today":
		showPrayerTimes(config)
	case "config":
		runConfig()
	default:
		fmt.Printf("Error: command %q not found\n", command)
		fmt.Println()
	}

}

func showPrayerTimes(c Config) {
	now := time.Now()
	date := data.NewDateComponents(now)
	coordinate, err := util.NewCoordinates(c.Latitude, c.Longitude)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	method, err := parseCalculationMethod(c.Method)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	params := calc.GetMethodParameters(method)
	asrMethod, err := parseAsrMethod(c.AsrMethod)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	params.Madhab = asrMethod

	prayerTimes, err := calc.NewPrayerTimes(coordinate, date, params)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = prayerTimes.SetTimeZone(c.Timezone)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Fajr:    %s\n", prayerTimes.Fajr.Format("3:04 PM"))
	fmt.Printf("Sunrise: %s\n", prayerTimes.Sunrise.Format("3:04 PM"))
	fmt.Printf("Dhuhr:   %s\n", prayerTimes.Dhuhr.Format("3:04 PM"))
	fmt.Printf("Asr:     %s\n", prayerTimes.Asr.Format("3:04 PM"))
	fmt.Printf("Maghrib: %s\n", prayerTimes.Maghrib.Format("3:04 PM"))
	fmt.Printf("Isha:    %s\n", prayerTimes.Isha.Format("3:04 PM"))

}

func runConfig() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Location (city, state/country): ")

	location, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	location = strings.TrimSpace(location)
	results, err := lookupLocation(location)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i, result := range results {
		fmt.Printf("%d. %s\n", i+1, result.DisplayName)
	}

	fmt.Print("Choose location: ")
	choiceInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	choiceInput = strings.TrimSpace(choiceInput)
	choice, err := strconv.Atoi(choiceInput)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if choice < 1 || choice > len(results) {
		fmt.Println("Error: invalid location choice")
		return
	}
	selected := results[choice-1]

	latitude, err := strconv.ParseFloat(selected.Lat, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	longitude, err := strconv.ParseFloat(selected.Lon, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	timezone := finder.GetTimezoneName(longitude, latitude)
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
