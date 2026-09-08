package main

import "fmt"

type LocationProfile struct {
	Name      string
	Latitude  float64
	Longitude float64
	Timezone  string
}

func addLocationProfile(config Config, name string) error {
	latitude, longitude, timezone, err := chooseLocation()
	if err != nil {
		return err
	}

	profile := LocationProfile{
		Name:      name,
		Latitude:  latitude,
		Longitude: longitude,
		Timezone:  timezone,
	}

	if config.Profiles == nil {
		config.Profiles = make(map[string]LocationProfile)
	}

	config.Profiles[name] = profile
	config.ActiveProfile = name
	config.Latitude = latitude
	config.Longitude = longitude
	config.Timezone = timezone

	err = saveConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func listLocationProfile(config Config) {
	fmt.Println()
	fmt.Println(titleStyle.Render("Location Profiles"))
	fmt.Println()

	if len(config.Profiles) == 0 {
		fmt.Println(mutedStyle.Render("No location profiles saved."))
		return
	}

	for _, profile := range config.Profiles {
		if profile.Name == config.ActiveProfile {
			profileText := fmt.Sprintf("* %s\n", profile.Name)
			fmt.Println(activeStyle.Render(profileText))
		} else {

			fmt.Printf(" %s\n", profile.Name)
		}
	}
}

func useLocationProfile(config Config, name string) error {
	profile, ok := config.Profiles[name]
	if !ok {
		return fmt.Errorf("location profile %q not found", name)
	}
	config.ActiveProfile = name
	config.Latitude = profile.Latitude
	config.Longitude = profile.Longitude
	config.Timezone = profile.Timezone

	err := saveConfig(config)
	if err != nil {
		return err
	}
	return nil

}

func removeLocationProfile(config Config, name string) error {
	_, ok := config.Profiles[name]
	if !ok {
		return fmt.Errorf("location profile %q not found", name)
	}

	if name == config.ActiveProfile {
		return fmt.Errorf("cannot remove active location profile; switch profiles first")
	}
	delete(config.Profiles, name)

	err := saveConfig(config)
	if err != nil {
		return err
	}
	return nil
}
