package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/MSA-Software-LLC/adhan-go/pkg/calc"
	"github.com/MSA-Software-LLC/adhan-go/pkg/data"
	"github.com/MSA-Software-LLC/adhan-go/pkg/util"
)

func showMonthlyPrayerTimes(config Config) {
	now := time.Now()
	firstDay := time.Date(
		now.Year(),
		now.Month(),
		1,
		0,
		0,
		0,
		0,
		now.Location(),
	)
	fmt.Println()
	fmt.Println(titleStyle.Render("Monthly Prayer Times\n"))
	fmt.Println(mutedStyle.Render(now.Format("January 2006\n")))
	fmt.Println()

	header := fmt.Sprintf("%-8s %-9s %-9s %-9s %-9s %-9s %-9s\n", "Date", "Fajr", "Sunrise", "Dhuhr", "Asr", "Maghrib", "Isha")
	fmt.Println(headingStyle.Render(header))

	for date := firstDay; date.Month() == now.Month(); date = date.AddDate(0, 0, 1) {
		prayerTimesByDate, err := calculatePrayerTimeForDate(config, date)
		if err != nil {
			printError(err)
			return
		}
		fmt.Printf(
			"%-8s %-9s %-9s %-9s %-9s %-9s %-9s\n",
			date.Format("Jan 02"),
			prayerTimesByDate.Fajr.Format("3:04 PM"),
			prayerTimesByDate.Sunrise.Format("3:04 PM"),
			prayerTimesByDate.Dhuhr.Format("3:04 PM"),
			prayerTimesByDate.Asr.Format("3:04 PM"),
			prayerTimesByDate.Maghrib.Format("3:04 PM"),
			prayerTimesByDate.Isha.Format("3:04 PM"),
		)

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
		printError(err)
		return
	}
	now := time.Now()
	current := prayerTimes.CurrentPrayer(now)
	next, nextTime, remaining, err := getNextPrayer(c)
	if err != nil {
		printError(err)
		return
	}
	hours := int(remaining.Hours())
	minutes := int(remaining.Minutes()) % 60
	remainingStr := fmt.Sprintf("%dh %dm", hours, minutes)
	title := titleStyle.Render("Current Prayer")

	fmt.Println()
	fmt.Println(title)
	fmt.Println()
	fmt.Printf("%s %s\n",
		labelStyle.Width(17).Render("Current Prayer:"),
		activeStyle.Render(prayerName(current)),
	)
	fmt.Printf("%s %s at %s\n",
		labelStyle.Width(17).Render("Next Prayer:"),
		activeStyle.Render(prayerName(next)),
		valueStyle.Render(nextTime.Format("3:04 PM")),
	)
	fmt.Printf("%s %s\n",
		mutedStyle.Width(17).Render("Time remaining:"),
		valueStyle.Render(remainingStr),
	)
}

func showNextPrayer(c Config) {
	next, nextTime, remaining, err := getNextPrayer(c)
	if err != nil {
		printError(err)
		return
	}
	hours := int(remaining.Hours())
	minutes := int(remaining.Minutes()) % 60
	remainingStr := fmt.Sprintf("%dh %dm", hours, minutes)
	title := titleStyle.Render("Next Prayer")

	fmt.Println()
	fmt.Println(title)
	fmt.Println()
	fmt.Printf("%s %s at %s\n",
		labelStyle.Width(17).Render("Next Prayer:"),
		activeStyle.Render(prayerName(next)),
		valueStyle.Render(nextTime.Format("3:04 PM")),
	)
	fmt.Printf("%s %s\n",
		mutedStyle.Width(17).Render("Time remaining:"),
		valueStyle.Render(remainingStr),
	)
}

func showPrayer(config Config, prayerName string) {
	prayerTimes, err := calculatePrayerTimes(config)
	if err != nil {
		printError(err)
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
		errMessage := fmt.Sprintf("Error: unknown prayer %q", prayerName)
		printErrorMessage(errMessage)
		fmt.Println()
		fmt.Println(mutedStyle.Render("Usage: salahctl prayer <fajr|dhuhr|asr|maghrib|isha>"))
		return
	}
	displayName := strings.ToUpper(prayerName[:1]) + prayerName[1:]
	title := titleStyle.Render(displayName + " Prayer")
	label := labelStyle.Width(9).Render(displayName + ":")
	value := valueStyle.Render(selectedTime.Format("3:04 PM"))

	fmt.Println()
	fmt.Println(title)
	fmt.Println()
	fmt.Printf("%s%s\n", label, value)
}

func showTomorrowPrayersTimes(c Config) {
	tomorrow := time.Now().AddDate(0, 0, 1)
	prayerTimesTomorrow, err := calculatePrayerTimeForDate(c, tomorrow)
	if err != nil {
		printError(err)
		return
	}
	fmt.Println()
	fmt.Println(titleStyle.Render("Tomorrow's Prayer Times"))
	fmt.Println()
	printPrayerTimes(prayerTimesTomorrow)
}

func showPrayerTimesByDate(c Config, dateString string) {

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		errMessage := fmt.Sprintln("Invalid date format. Use YYYY-MM-DD")
		printErrorMessage(errMessage)
		return
	}
	prayerTimesByDate, err := calculatePrayerTimeForDate(c, date)
	if err != nil {
		printError(err)
		return
	}
	fmt.Println()
	fmt.Println(titleStyle.Render(date.Format("Monday, January 2, 2006")))
	fmt.Println()
	printPrayerTimes(prayerTimesByDate)

}

func showWeeklyPrayerTimes(c Config) {
	date := time.Now()

	for i := 0; i < 7; i++ {
		prayerTimes, err := calculatePrayerTimeForDate(c, date)
		if err != nil {
			printError(err)
			return
		}

		fmt.Println(titleStyle.Render(date.Format("Monday, January 2, 2006")))
		fmt.Println()
		printPrayerTimes(prayerTimes)
		fmt.Println()

		date = date.AddDate(0, 0, 1)
	}
}

func showRemainingPrayers(config Config) {
	prayerTimes, err := calculatePrayerTimes(config)
	if err != nil {
		printError(err)
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

	fmt.Println()
	fmt.Println(titleStyle.Render("Remaining Prayers"))
	fmt.Println()

	for _, prayer := range prayers {
		if now.Before(prayer.time) {
			fmt.Println(labelStyle.Width(9).Render(prayer.name+":"),
				valueStyle.Render(prayer.time.Format("3:04 PM")),
			)
		}
	}

	if now.After(prayerTimes.Isha) {
		fmt.Println(successStyle.Render("All prayers are complete for today."))
	}

	fmt.Println()
}

func printPrayerTimes(prayerTimes *calc.PrayerTimes) {
	fmt.Printf("%s %s\n",
		labelStyle.Width(9).Render("Fajr:"),
		valueStyle.Render(prayerTimes.Fajr.Format("3:04 PM")),
	)
	fmt.Printf("%s %s\n",
		labelStyle.Width(9).Render("Sunrise:"),
		valueStyle.Render(prayerTimes.Sunrise.Format("3:04 PM")),
	)
	fmt.Printf("%s %s\n",
		labelStyle.Width(9).Render("Dhuhr:"),
		valueStyle.Render(prayerTimes.Dhuhr.Format("3:04 PM")),
	)
	fmt.Printf("%s %s\n",
		labelStyle.Width(9).Render("Asr:"),
		valueStyle.Render(prayerTimes.Asr.Format("3:04 PM")),
	)
	fmt.Printf("%s %s\n",
		labelStyle.Width(9).Render("Maghrib:"),
		valueStyle.Render(prayerTimes.Maghrib.Format("3:04 PM")),
	)
	fmt.Printf("%s %s\n",
		labelStyle.Width(9).Render("Isha:"),
		valueStyle.Render(prayerTimes.Isha.Format("3:04 PM")),
	)

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
		printError(err)
		return
	}
	fmt.Println()

	fmt.Println(titleStyle.Render("Today's Prayer Times"))
	fmt.Println()
	printPrayerTimes(prayerTimes)
}
