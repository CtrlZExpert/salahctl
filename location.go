package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ringsaturn/tzf"
)

type LocationResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
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
		return 0.0, 0.0, "", errors.New("invalid location choice")
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
