package main

import (
	"testing"

	"github.com/MSA-Software-LLC/adhan-go/pkg/calc"
)

func TestParseCalculationMethod(t *testing.T) {
	expected := calc.NORTH_AMERICA

	method, err := parseCalculationMethod("north_america")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if method != expected {
		t.Errorf("expected %v, got %v", expected, method)
	}
}

func TestParseCalculationMethodMuslimWorldLeague(t *testing.T) {
	expected := calc.MUSLIM_WORLD_LEAGUE

	method, err := parseCalculationMethod("muslim_world_league")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if method != expected {
		t.Errorf("expected %v, got %v", expected, method)
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

func TestUsageFor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"date", "Usage: salahctl date YYYY-MM-DD"},
		{"config", "Usage: salahctl config [show|location|method|asr]"},
		{"--help", "Usage: salahctl --help"},
		{"--version", "Usage: salahctl --version"},
		{"today", "Usage: salahctl today"},
	}

	for _, test := range tests {
		actual := usageFor(test.input)

		if actual != test.expected {
			t.Errorf("for %q: expected %q, got %q", test.input, test.expected, actual)
		}
	}
}
