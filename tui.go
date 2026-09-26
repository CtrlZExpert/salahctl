package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/MSA-Software-LLC/adhan-go/pkg/calc"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	config         Config
	prayerTimes    *calc.PrayerTimes
	currentPrayer  calc.Prayer
	nextPrayer     calc.Prayer
	nextPrayerTime time.Time
	remaining      time.Duration
	view           viewMode
	width          int
	height         int
	scrollOffset   int
	selectedPrayer int
}

type prayerRow struct {
	prayer calc.Prayer
	name   string
	time   time.Time
}

type viewMode int

const (
	todayView = iota
	weekView
	monthView
	prayerDetailView
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {

	switch m.view {
	case todayView:
		return m.renderToday()
	case weekView:
		return m.renderWeek()
	case monthView:
		return m.renderMonth()
	case prayerDetailView:
		return m.renderPrayerDetail()
	}

	return ""
}

func (m model) renderToday() string {
	rows := prayerRows(m.prayerTimes)
	var prayerTimes strings.Builder

	for i, row := range rows {
		prayer := formatPrayerRow(row.name, row.time)

		if row.prayer == m.currentPrayer {
			activeRow := fmt.Sprintf("%-9s %s", row.name, row.time.Format("3:04 PM"))
			prayer = activeStyle.Render(activeRow)
		}
		if i == m.selectedPrayer {
			prayer = "> " + prayer
		} else {
			prayer = " " + prayer
		}

		prayerTimes.WriteString(prayer)
		prayerTimes.WriteString("\n")
	}
	prayerList := prayerTimes.String()
	nextPrayerName := prayerName(m.nextPrayer)
	nextPrayerLine := labelStyle.Width(16).Render("Next prayer:") + " " + nextPrayerName + " at " + m.nextPrayerTime.Format("3:04 PM")

	hours := int(m.remaining.Hours())
	minutes := int(m.remaining.Minutes()) % 60

	remainingText := fmt.Sprintf("%dh %dm", hours, minutes)

	timeRemaining := labelStyle.Width(16).Render("Time remaining:") + " " + remainingText

	title := fmt.Sprint(titleStyle.Render("salahctl"))
	todayHeading := fmt.Sprint(headingStyle.Render("Today's Prayer Times"))

	viewBar := m.renderViewBar()
	scrollBar := mutedStyle.Render("↑/k up • ↓/j down")
	quitBar := fmt.Sprintf(mutedStyle.Render("[ q ] quit"))
	screen := fmt.Sprintf(
		`
%s

%s

%s

%s

%s

%s

%s
%s`,
		title,
		todayHeading,
		prayerList,
		nextPrayerLine,
		timeRemaining,
		viewBar,
		scrollBar,
		quitBar,
	)
	return screen
}

func (m model) renderPrayerDetail() string {
	status := "Upcoming"
	now := time.Now()
	rows := prayerRows(m.prayerTimes)
	selected := rows[m.selectedPrayer]
	prayerHeading := fmt.Sprint(headingStyle.Render(selected.name))
	prayerTime := fmt.Sprint(labelStyle.Width(16).Render("Prayer time:") + " " + selected.time.Format("3:04 PM"))
	title := fmt.Sprint(titleStyle.Render("salahctl"))
	quitBar := fmt.Sprint(mutedStyle.Render("[ esc ] Back [ q ] quit"))

	if selected.time.Before(now) {
		status = "Passed"
	}

	if selected.prayer == m.currentPrayer && selected.prayer != calc.SUNRISE {
		status = "Current"
	}

	timeUntil := selected.time.Sub(now)
	hoursUntil := int(timeUntil.Hours())
	minutesUntil := int(timeUntil.Minutes()) % 60
	timeUntilText := fmt.Sprintf("%dh %dm", hoursUntil, minutesUntil)
	timeUntilLine := ""
	if status == "Upcoming" {
		timeUntilLine = labelStyle.Width(16).Render("Time until:") + " " + timeUntilText

	}

	hoursRemaining := int(m.remaining.Hours())
	minutesRemaining := int(m.remaining.Minutes()) % 60

	remainingText := fmt.Sprintf("%dh %dm", hoursRemaining, minutesRemaining)
	nextPrayerName := prayerName(m.nextPrayer)
	nextPrayerLine := ""

	if status == "Current" {
		nextPrayerLine = labelStyle.Width(16).Render("Next prayer:") + " " + nextPrayerName + " at " + m.nextPrayerTime.Format("3:04 PM")
		timeUntilLine = labelStyle.Width(16).Render("Time remaining:") + " " + remainingText

	}

	statusLine := fmt.Sprint(labelStyle.Width(16).Render("Status:") + " " + status)

	screen := fmt.Sprintf(`
%s

%s

%s
%s
%s
%s

%s`,

		title,
		prayerHeading,
		prayerTime,
		statusLine,
		nextPrayerLine,
		timeUntilLine,
		quitBar)

	return screen

}

func (m model) renderWeek() string {
	var week strings.Builder
	date := time.Now()

	headers := []string{
		"Date",
		"Fajr",
		"Sunrise",
		"Dhuhr",
		"Asr",
		"Maghrib",
		"Isha",
	}

	for i, header := range headers {
		var renderedHeader string
		if i == 0 {
			renderedHeader = labelStyle.Width(12).Render(header)
		} else {
			renderedHeader = labelStyle.Width(10).Render(header)

		}
		week.WriteString(renderedHeader)
	}

	week.WriteString("\n")

	for i := 0; i < 7; i++ {
		prayerTimes, err := calculatePrayerTimeForDate(m.config, date)
		if err != nil {
			return errorStyle.Render(err.Error())
		}
		week.WriteString(titleStyle.Width(12).Render(date.Format("Mon Jan 2")))
		prayerTimesRow := formatPrayerTimeRow(prayerTimes)
		week.WriteString(prayerTimesRow)
		week.WriteString("\n")
		date = date.AddDate(0, 0, 1)

	}

	weekList := week.String()
	title := fmt.Sprint(titleStyle.Render("salahctl"))
	weekHeading := fmt.Sprint(headingStyle.Render("Weekly Prayer Times"))
	viewBar := m.renderViewBar()
	quitBar := fmt.Sprintf(mutedStyle.Render("[ q ] quit"))

	screen := fmt.Sprintf(`

%s

%s

%s

%s

%s

`,
		title,
		weekHeading,
		weekList,
		viewBar,
		quitBar)

	return screen
}

func (m model) renderMonth() string {
	visibleRows := m.monthVisibleRows()
	daysInMonth := daysInCurrentMonth()
	start := m.scrollOffset
	end := m.scrollOffset + visibleRows

	if end > daysInMonth {
		end = daysInMonth
	}

	var output strings.Builder
	headers := []string{
		"Date",
		"Fajr",
		"Sunrise",
		"Dhuhr",
		"Asr",
		"Maghrib",
		"Isha",
	}

	for i, header := range headers {
		var renderedHeader string
		if i == 0 {
			renderedHeader = labelStyle.Width(12).Render(header)
		} else {
			renderedHeader = labelStyle.Width(10).Render(header)

		}
		output.WriteString(renderedHeader)

	}
	output.WriteString("\n")

	for day := start + 1; day <= end; day++ {
		date := time.Date(time.Now().Year(), time.Now().Month(), day, 0, 0, 0, 0, time.Now().Location())
		dayName := date.Format("Mon")

		prayerTimes, err := calculatePrayerTimeForDate(m.config, date)
		if err != nil {
			return errorStyle.Render(err.Error())
		}
		formattedTimes := formatPrayerTimeRow(prayerTimes)
		row := fmt.Sprintf("%s %s %s", dayName, date.Format("Jan 02"), formattedTimes)
		renderedRow := activeRender(row, date)
		output.WriteString(renderedRow)
		output.WriteString("\n")

	}
	outputList := output.String()
	title := fmt.Sprint(titleStyle.Render("salahctl"))
	monthHeading := fmt.Sprint(headingStyle.Render("Monthly Prayer Times"))
	viewBar := m.renderViewBar()
	scrollBar := mutedStyle.Render("↑/k up • ↓/j down")
	quitBar := fmt.Sprintf(mutedStyle.Render("[ q ] quit"))
	screen := fmt.Sprintf(`
%s

%s

%s

%s

%s
%s`,
		title,
		monthHeading,
		outputList,
		viewBar,
		scrollBar,
		quitBar)

	return screen

}

func (m model) renderViewBar() string {
	today := "[ t ] Today"
	week := "[ w ] Week"
	month := "[ m ] Month"

	switch m.view {
	case todayView:
		today = activeStyle.Render(today)
		week = mutedStyle.Render(week)
		month = mutedStyle.Render(month)
	case weekView:
		today = mutedStyle.Render(today)
		week = activeStyle.Render(week)
		month = mutedStyle.Render(month)
	case monthView:
		today = mutedStyle.Render(today)
		week = mutedStyle.Render(week)
		month = activeStyle.Render(month)
	}

	return fmt.Sprintf("%s %s %s", today, week, month)
}

func (m model) monthVisibleRows() int {
	visibleRows := m.height - 11
	if visibleRows < 5 {
		visibleRows = 5
	}
	return visibleRows
}

func (m model) maxMonthScroll(daysInMonth int, visibleRows int) int {
	return daysInMonth - visibleRows

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "t":
			m.view = todayView
		case "w":
			m.view = weekView
		case "m":
			m.view = monthView
		case "down", "j":
			if m.view == monthView {
				visibileRows := m.monthVisibleRows()
				daysInMonth := daysInCurrentMonth()
				maxScroll := m.maxMonthScroll(daysInMonth, visibileRows)
				if m.scrollOffset < maxScroll {
					m.scrollOffset++
				}
			}
			if m.view == todayView {
				if m.selectedPrayer < 5 {
					m.selectedPrayer++
				}
			}

		case "up", "k":
			if m.view == monthView {
				if m.scrollOffset > 0 {
					m.scrollOffset--

				}

			}

			if m.view == todayView {
				if m.selectedPrayer > 0 {
					m.selectedPrayer--
				}
			}
		case "enter":
			if m.view == todayView {
				m.view = prayerDetailView
			}
		case "esc":
			if m.view == prayerDetailView {
				m.view = todayView
			}

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
		config:         c,
		prayerTimes:    prayerTimes,
		currentPrayer:  currentPrayer,
		nextPrayer:     nextPrayer,
		nextPrayerTime: nextPrayerTime,
		remaining:      remaining,
		view:           todayView,
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

func prayerRows(prayerTimes *calc.PrayerTimes) []prayerRow {
	rows := []prayerRow{
		{calc.FAJR, "Fajr:", prayerTimes.Fajr},
		{calc.SUNRISE, "Sunrise:", prayerTimes.Sunrise},
		{calc.DHUHR, "Dhuhr:", prayerTimes.Dhuhr},
		{calc.ASR, "Asr:", prayerTimes.Asr},
		{calc.MAGHRIB, "Maghrib:", prayerTimes.Maghrib},
		{calc.ISHA, "Isha:", prayerTimes.Isha},
	}

	return rows

}

func formatPrayerTimeRow(prayerTimes *calc.PrayerTimes) string {
	times := []time.Time{
		prayerTimes.Fajr,
		prayerTimes.Sunrise,
		prayerTimes.Dhuhr,
		prayerTimes.Asr,
		prayerTimes.Maghrib,
		prayerTimes.Isha,
	}

	var row strings.Builder

	for _, t := range times {
		formattedTime := valueStyle.Width(10).Render(t.Format("3:04 PM"))
		row.WriteString(formattedTime)

	}

	return row.String()
}

func daysInCurrentMonth() int {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	daysInMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, -1).Day()
	return daysInMonth
}

func activeRender(text string, date time.Time) string {
	now := time.Now()
	if date.Year() == now.Year() &&
		date.Month() == now.Month() &&
		date.Day() == now.Day() {
		return activeStyle.Render(text)
	}
	return text
}
