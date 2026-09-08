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
	Latitude      float64                    `toml:"latitude"`
	Longitude     float64                    `toml:"longitude"`
	Timezone      string                     `toml:"timezone"`
	Method        string                     `toml:"method"`
	AsrMethod     string                     `toml:"asr_method"`
	ActiveProfile string                     `toml:"active_profile"`
	Profiles      map[string]LocationProfile `toml:"profiles"`
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
	case "egyptian":
		return calc.EGYPTIAN, nil
	case "karachi":
		return calc.KARACHI, nil
	case "umm_al_qura":
		return calc.UMM_AL_QURA, nil
	case "dubai":
		return calc.DUBAI, nil
	case "moon_sighting_committee":
		return calc.MOON_SIGHTING_COMMITTEE, nil
	case "kuwait":
		return calc.KUWAIT, nil
	case "qatar":
		return calc.QATAR, nil
	case "singapore":
		return calc.SINGAPORE, nil
	case "uoif":
		return calc.UOIF, nil
	case "tehran":
		return calc.TEHRAN, nil
	case "turkey":
		return calc.TURKEY, nil
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
		fmt.Println(headingStyle.Render("Calculation Method"))
		fmt.Println()
		fmt.Println("  1. Muslim World League")
		fmt.Println("  2. North America")
		fmt.Println("  3. Egyptian")
		fmt.Println("  4. Karachi")
		fmt.Println("  5. Umm Al-Qura")
		fmt.Println("  6. Dubai")
		fmt.Println("  7. Moon Sighting Committee")
		fmt.Println("  8. Kuwait")
		fmt.Println("  9. Qatar")
		fmt.Println(" 10. Singapore")
		fmt.Println(" 11. UOIF")
		fmt.Println(" 12. Tehran")
		fmt.Println(" 13. Turkey")
		fmt.Println()
		fmt.Print(labelStyle.Render("Choose method: "))

		var choiceAction string
		fmt.Scan(&choiceAction)

		validChoice, err := strconv.Atoi(choiceAction)
		if err != nil {
			printErrorMessage("Invalid selection. Enter a number 1-13")

			continue
		}

		switch validChoice {
		case 1:
			return "muslim_world_league"
		case 2:
			return "north_america"
		case 3:
			return "egyptian"
		case 4:
			return "karachi"
		case 5:
			return "umm_al_qura"
		case 6:
			return "dubai"
		case 7:
			return "moon_sighting_committee"
		case 8:
			return "kuwait"
		case 9:
			return "qatar"
		case 10:
			return "singapore"
		case 11:
			return "uoif"
		case 12:
			return "tehran"
		case 13:
			return "turkey"
		default:
			printErrorMessage("Invalid selection. Enter a number 1-13")
			continue
		}
	}
}
func chooseAsrMethod() string {
	for {
		fmt.Println()
		fmt.Println(headingStyle.Render("Asr Method"))
		fmt.Println()
		fmt.Println("  1. Standard")
		fmt.Println("  2. Hanafi")
		fmt.Println()
		fmt.Print(labelStyle.Render("Choose method: "))

		var choiceAction string
		fmt.Scan(&choiceAction)

		validChoice, err := strconv.Atoi(choiceAction)
		if err != nil {
			printErrorMessage("Invalid selection. Enter a number 1-2")
			continue
		}

		switch validChoice {
		case 1:
			return "standard"
		case 2:
			return "hanafi"
		default:
			printErrorMessage("Invalid selection. Enter a number 1-2")
			continue
		}

	}
}

func showConfig(config Config) {
	fmt.Println()
	fmt.Println(titleStyle.Render("Configuration"))
	fmt.Println()

	latitude := fmt.Sprintf("%f", config.Latitude)
	longitude := fmt.Sprintf("%f", config.Longitude)
	fmt.Printf("%s%s\n",
		labelStyle.Width(18).Render("Latitude:"),
		valueStyle.Render(latitude),
	)
	fmt.Printf("%s%s\n",
		labelStyle.Width(18).Render("Longitude:"),
		valueStyle.Render(longitude),
	)
	fmt.Printf("%s%s\n",
		labelStyle.Width(18).Render("Timezone:"),
		valueStyle.Render(config.Timezone),
	)
	fmt.Printf("%s%s\n",
		labelStyle.Width(18).Render("Method:"),
		valueStyle.Render(config.Method),
	)
	fmt.Printf("%s%s\n",
		labelStyle.Width(18).Render("Asr Method:"),
		valueStyle.Width(18).Render(config.AsrMethod),
	)
	fmt.Printf("%s%s\n",
		labelStyle.Width(18).Render("Active Profile:"),
		activeStyle.Width(18).Render(config.ActiveProfile),
	)
}
func runConfig() {

	fmt.Println()
	fmt.Println(titleStyle.Render("Configuration Setup"))
	fmt.Println()

	latitude, longitude, timezone, err := chooseLocation()
	if err != nil {
		printError(err)
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
		printError(err)
		return
	}

	fmt.Println(successStyle.Render("Configuration saved successfully"))

}

func updateLocation(config Config) {
	latitude, longitude, timezone, err := chooseLocation()
	if err != nil {
		printError(err)
		return
	}

	config.Latitude = latitude
	config.Longitude = longitude
	config.Timezone = timezone

	err = saveConfig(config)
	if err != nil {
		printError(err)
		return
	}
	fmt.Println(successStyle.Render("Location has successfully been updated"))

}

func updateCalculationMethod(config Config) {
	method := chooseCalculationMethod()
	config.Method = method
	err := saveConfig(config)
	if err != nil {
		printError(err)
		return
	}
	fmt.Println(successStyle.Render("Method has successfully been updated"))
}

func updateAsrMethod(config Config) {
	asrMethod := chooseAsrMethod()
	config.AsrMethod = asrMethod
	err := saveConfig(config)
	if err != nil {
		printError(err)
		return
	}
	fmt.Println(successStyle.Render("Asr Method has successfully been updated"))

}
