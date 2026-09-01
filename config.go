package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/MSA-Software-LLC/adhan-go/pkg/calc"
)

type Config struct {
	Latitude  float64 `toml:"latitude"`
	Longitude float64 `toml:"longitude"`
	Timezone  string  `toml:"timezone"`
	Method    string  `toml:"method"`
	AsrMethod string  `toml:"asr_method"`
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
		return 0, fmt.Errorf("unsupported calculation  Asr method: %q ", method)
	}
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

func showConfig(config Config) {
	fmt.Printf("Latitude:   %f\n", config.Latitude)
	fmt.Printf("Longitude:   %f\n", config.Longitude)
	fmt.Printf("Timezone:    %s\n", config.Timezone)
	fmt.Printf("Method:     %s\n", config.Method)
	fmt.Printf("Asr Method: %s\n", config.AsrMethod)
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
