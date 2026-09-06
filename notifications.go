package main

import (
	"fmt"
	"os/exec"
	"time"
)

func sendNotification(title string, message string) error {
	cmd := exec.Command("notify-send", title, message)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func checkPrayerNotification(c Config) error {
	now := time.Now()
	prayerTimes, err := calculatePrayerTimes(c)
	if err != nil {
		return err
	}
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

	for _, prayer := range prayers {
		if now.Format("15:04") == prayer.time.Format("15:04") {
			message := fmt.Sprintf("%s prayer time", prayer.name)

			err := sendNotification("salahctl", message)
			if err != nil {
				return err
			}
			return nil
		}

	}
	return nil
}
