package main

import (
	"testing"

	"github.com/MSA-Software-LLC/adhan-go/pkg/calc"
)

func TestParseCalculationMethod(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected calc.CalculationMethod
	}{
		{"Muslim World League", "muslim_world_league", calc.MUSLIM_WORLD_LEAGUE},
		{"North America", "north_america", calc.NORTH_AMERICA},
		{"Egyptian", "egyptian", calc.EGYPTIAN},
		{"Karachi", "karachi", calc.KARACHI},
		{"Umm Al-Qura", "umm_al_qura", calc.UMM_AL_QURA},
		{"Dubai", "dubai", calc.DUBAI},
		{"Moon Sighting Committee", "moon_sighting_committee", calc.MOON_SIGHTING_COMMITTEE},
		{"Kuwait", "kuwait", calc.KUWAIT},
		{"Qatar", "qatar", calc.QATAR},
		{"Singapore", "singapore", calc.SINGAPORE},
		{"UOIF", "uoif", calc.UOIF},
		{"Tehran", "tehran", calc.TEHRAN},
		{"Turkey", "turkey", calc.TURKEY},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method, err := parseCalculationMethod(tt.input)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if method != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, method)
			}
		})
	}
}
func TestParseAsrMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected calc.AsrJuristicMethod
	}{
		{"standard", calc.SHAFI_HANBALI_MALIKI},
		{"hanafi", calc.HANAFI},
	}

	for _, test := range tests {
		asrMethod, err := parseAsrMethod(test.input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if asrMethod != test.expected {
			t.Errorf("expected %v, got %v", test.expected, asrMethod)
		}
	}
}

func TestPrayerName(t *testing.T) {
	tests := []struct {
		input    calc.Prayer
		expected string
	}{
		{calc.FAJR, "Fajr"},
		{calc.SUNRISE, "None"},
		{calc.DHUHR, "Dhuhr"},
		{calc.ASR, "Asr"},
		{calc.MAGHRIB, "Maghrib"},
		{calc.ISHA, "Isha"},
		{calc.NO_PRAYER, "None"},
		{calc.Prayer(999), "Unknown"},
	}

	for _, test := range tests {
		actual := prayerName(test.input)

		if actual != test.expected {
			t.Errorf("expected %q, got %q", test.expected, actual)
		}
	}
}
