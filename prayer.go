package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/MSA-Software-LLC/adhan-go/pkg/calc"
	"github.com/MSA-Software-LLC/adhan-go/pkg/data"
	"github.com/MSA-Software-LLC/adhan-go/pkg/util"
)

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
	return calculatePrayerTimeForDate(c, now)

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

	prayers := []struct {
		name string
		time time.Time
	}{
		{"Fajr", prayerTimes.Fajr},
		{"Dhuhr", prayerTimes.Dhuhr},
		{"Asr", prayerTimes.Asr},
		{"Maghrib", prayerTimes.Maghrib},
		{"Isha", prayerTimes.Isha},
	}

	fmt.Println("Remaining Prayers")
	fmt.Println()

	for _, prayer := range prayers {
		if now.Before(prayer.time) {
			fmt.Println(prayer.name+":", prayer.time.Format("3:04 PM"))
		}
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
