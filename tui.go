package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/MSA-Software-LLC/adhan-go/pkg/calc"
	tea "github.com/charmbracelet/bubbletea"
)

type prayerRow struct {
	prayer calc.Prayer
	name   string
	time   time.Time
}

type model struct {
	prayerTimes    *calc.PrayerTimes
	currentPrayer  calc.Prayer
	nextPrayer     calc.Prayer
	nextPrayerTime time.Time
	remaining      time.Duration
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {

	rows := []prayerRow{
		{calc.FAJR, "Fajr:", m.prayerTimes.Fajr},
		{calc.SUNRISE, "Sunrise:", m.prayerTimes.Sunrise},
		{calc.DHUHR, "Dhuhr:", m.prayerTimes.Dhuhr},
		{calc.ASR, "Asr:", m.prayerTimes.Asr},
		{calc.MAGHRIB, "Maghrib:", m.prayerTimes.Maghrib},
		{calc.ISHA, "Isha:", m.prayerTimes.Isha},
	}
	var prayerTimes strings.Builder
	for _, row := range rows {
		prayer := formatPrayerRow(row.name, row.time)

		if row.prayer == m.currentPrayer {
			activeRow := fmt.Sprintf("%-9s %s", row.name, row.time.Format("3:04 PM"))
			prayer = activeStyle.Render(activeRow)
		}

		prayerTimes.WriteString(prayer)
		prayerTimes.WriteString("\n")
	}
	prayerList := prayerTimes.String()
	nextPrayerName := prayerName(m.nextPrayer)
	nextPrayerLine := labelStyle.Width(16).Render("Next:") + " " + nextPrayerName + " at " + m.nextPrayerTime.Format("3:04 PM")

	hours := int(m.remaining.Hours())
	minutes := int(m.remaining.Minutes()) % 60

	remainingText := fmt.Sprintf("%dh %dm", hours, minutes)

	timeRemaining := labelStyle.Width(16).Render("Time remaining:") + " " + remainingText

	title := fmt.Sprint(titleStyle.Render("salahctl"))
	quit := fmt.Sprintf(mutedStyle.Render("q quit"))
	screen := fmt.Sprintf(
		`%s


%s

%s
%s

		
%s`,
		title,
		prayerList,
		nextPrayerLine,
		timeRemaining,
		quit,
	)

	return screen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func runTUI() error {
	c, err := loadConfig()
	if err != nil {
		return err
	}

	prayerTimes, err := calculatePrayerTimes(c)
	if err != nil {
		return err
	}

	currentPrayer := prayerTimes.CurrentPrayer(time.Now())
	nextPrayer, nextPrayerTime, remaining, err := getNextPrayer(c)
	if err != nil {
		return err
	}

	m := model{
		prayerTimes:    prayerTimes,
		currentPrayer:  currentPrayer,
		nextPrayer:     nextPrayer,
		nextPrayerTime: nextPrayerTime,
		remaining:      remaining,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		return err
	}

	return nil
}

func formatPrayerRow(name string, prayerTime time.Time) string {
	return labelStyle.Width(9).Render(name) + " " + prayerTime.Format("3:04 PM")
}
