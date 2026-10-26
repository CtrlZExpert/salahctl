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
	viewDate       time.Time
	prayerDate     time.Time
	remaining      time.Duration
	view           viewMode
	previousView   viewMode
	width          int
	height         int
	scrollOffset   int
	selectedPrayer int
	weekOffset     int
	monthOffset    int
	selectedDay    int
}

type prayerRow struct {
	prayer calc.Prayer
	name   string
	time   time.Time
}

type viewMode int
type tickMsg time.Time

const (
	todayView = iota
	weekView
	monthView
	prayerDetailView
	helpView
)

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
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
	case helpView:
		return m.renderHelp()
	}

	return ""
}

func (m model) renderToday() string {
	prayerTimes, err := calculatePrayerTimeForDate(m.config, m.viewDate)
	if err != nil {
		return errorStyle.Render(err.Error())
	}
	rows := prayerRows(prayerTimes)
	var prayerListBuilder strings.Builder

	for i, row := range rows {
		prayer := formatPrayerRow(row.name, row.time)

		if row.prayer == m.currentPrayer && isToday(m.viewDate) {
			activeRow := fmt.Sprintf("%-9s %s", row.name, row.time.Format("3:04 PM"))
			prayer = activeRender(activeRow, m.viewDate)
		}
		if i == m.selectedPrayer {
			prayer = "> " + prayer
		} else {
			prayer = "  " + prayer
		}

		prayerListBuilder.WriteString(prayer)
		prayerListBuilder.WriteString("\n")
	}
	prayerList := prayerListBuilder.String()
	nextPrayerName := prayerName(m.nextPrayer)
	nextPrayerLine := ""

	hours := int(m.remaining.Hours())
	minutes := int(m.remaining.Minutes()) % 60

	remainingText := fmt.Sprintf("%dh %dm", hours, minutes)
	timeRemaining := ""

	title := titleStyle.Render("salahctl")
	todayHeading := headingStyle.Render("Prayer Times — " + m.viewDate.Format("January 2, 2006"))
	if isToday(m.viewDate) {
		todayHeading = headingStyle.Render("Today's Prayer Times")
		nextPrayerLine = labelStyle.Width(16).Render("Next prayer:") + " " + nextPrayerName + " at " + m.nextPrayerTime.Format("3:04 PM")
		timeRemaining = labelStyle.Width(16).Render("Time remaining:") + " " + remainingText

	}

	viewBar := m.renderViewBar()
	scrollBar := mutedStyle.Render("↑/k up • ↓/j down")
	footer := mutedStyle.Render("[ ? ] Help [ q ] Quit")
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
		footer,
	)
	return screen
}

func (m model) renderPrayerDetail() string {
	status := "Upcoming"
	now := time.Now()
	prayerTimes, err := calculatePrayerTimeForDate(m.config, m.viewDate)
	if err != nil {
		return errorStyle.Render(err.Error())

	}
	rows := prayerRows(prayerTimes)
	selected := rows[m.selectedPrayer]
	prayerHeading := headingStyle.Render(selected.name)
	prayerTime := labelStyle.Width(16).Render("Prayer time:") + " " + selected.time.Format("3:04 PM")
	title := titleStyle.Render("salahctl")
	dateLine := ""

	if !isToday(m.viewDate) {
		dateLine = m.viewDate.Format("January 2, 2006")
	}
	footer := mutedStyle.Render("[ esc ] Back [ ? ] Help [ q ] Quit")

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
	if status == "Upcoming" && isToday(m.viewDate) {
		timeUntilLine = labelStyle.Width(16).Render("Time until:") + " " + timeUntilText

	}

	hoursRemaining := int(m.remaining.Hours())
	minutesRemaining := int(m.remaining.Minutes()) % 60

	remainingText := fmt.Sprintf("%dh %dm", hoursRemaining, minutesRemaining)
	nextPrayerName := prayerName(m.nextPrayer)
	nextPrayerLine := ""

	if status == "Current" && isToday(m.viewDate) {
		nextPrayerLine = labelStyle.Width(16).Render("Next prayer:") + " " + nextPrayerName + " at " + m.nextPrayerTime.Format("3:04 PM")
		timeUntilLine = labelStyle.Width(16).Render("Time remaining:") + " " + remainingText

	}

	statusLine := ""

	if isToday(m.viewDate) {
		statusLine = labelStyle.Width(16).Render("Status:") + " " + status

	}

	screen := fmt.Sprintf(`
%s

%s
%s

%s
%s
%s
%s

%s`,

		title,
		prayerHeading,
		dateLine,
		prayerTime,
		statusLine,
		nextPrayerLine,
		timeUntilLine,
		footer)

	return screen

}

func (m model) renderWeek() string {
	var week strings.Builder
	date := time.Now()
	date = date.AddDate(0, 0, m.weekOffset*7)
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
	scrollBar := mutedStyle.Render("←/h previous • →/l next")
	footer := mutedStyle.Render("[ ? ] Help [ q ] Quit")
	screen := fmt.Sprintf(`

%s

%s

%s

%s

%s
%s`,
		title,
		weekHeading,
		weekList,
		viewBar,
		scrollBar,
		footer)

	return screen
}

func (m model) renderMonth() string {
	now := time.Now()
	viewDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, m.monthOffset, 0)
	visibleRows := m.monthVisibleRows()
	totalDays := daysInMonth(viewDate)
	start := m.scrollOffset
	end := m.scrollOffset + visibleRows

	if end > totalDays {
		end = totalDays
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
		date := time.Date(viewDate.Year(),
			viewDate.Month(),
			day,
			0, 0, 0, 0,
			viewDate.Location(),
		)
		dayName := date.Format("Mon")

		prayerTimes, err := calculatePrayerTimeForDate(m.config, date)
		if err != nil {
			return errorStyle.Render(err.Error())
		}
		formattedTimes := formatPrayerTimeRow(prayerTimes)
		row := fmt.Sprintf("%s %s %s", dayName, date.Format("Jan 02"), formattedTimes)
		if day == m.selectedDay {
			row = "> " + row
		} else {
			row = "  " + row
		}

		renderedRow := activeRender(row, date)
		output.WriteString(renderedRow)
		output.WriteString("\n")

	}
	outputList := output.String()
	title := titleStyle.Render("salahctl")
	monthHeading := headingStyle.Render("Monthly Prayer Times  — " + viewDate.Format("January 2006"))
	viewBar := m.renderViewBar()
	scrollBar := mutedStyle.Render("←/h previous • →/l next • ↑/k up • ↓/j down")
	footer := mutedStyle.Render("[ ? ] Help [ q ] Quit")
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
		footer)

	return screen

}

func (m model) renderHelp() string {
	type helpRow struct {
		key         string
		description string
	}
	title := titleStyle.Render("salahctl")
	helpHeading := headingStyle.Render("Help")
	rows := []helpRow{
		{key: "t", description: "Today"},
		{key: "w", description: "Week"},
		{key: "m", description: "Month"},
		{key: "↑/k", description: "Move up"},
		{key: "↓/j", description: "Move down"},
		{key: "←/h", description: "Previous week/month"},
		{key: "→/l", description: "Next week/month"},
		{key: "Enter", description: "Open selected day / Prayer details"},
		{key: "Esc", description: "Back"},
		{key: "q", description: "Quit"},
		{key: "?", description: "Help"},
	}
	var help strings.Builder
	for _, item := range rows {
		renderedKey := labelStyle.Width(14).Render(item.key)
		help.WriteString(renderedKey)
		help.WriteString(item.description)
		help.WriteString("\n")
	}
	helpList := help.String()
	footer := mutedStyle.Render("[ esc ] Back [ q ] Quit ")
	return fmt.Sprintf(`
%s

%s

%s

%s`,
		title,
		helpHeading,
		helpList,
		footer,
	)

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
	case tickMsg:
		if !isToday(m.prayerDate) {

			prayerTimes, err := calculatePrayerTimes(m.config)
			if err != nil {
				return m, tickCmd()
			}

			if m.prayerDate.Year() == m.viewDate.Year() &&
				m.prayerDate.Month() == m.viewDate.Month() &&
				m.prayerDate.Day() == m.viewDate.Day() {
				m.viewDate = time.Now()
			}

			m.prayerTimes = prayerTimes
			m.prayerDate = time.Now()
			m.currentPrayer = prayerTimes.CurrentPrayer(time.Now())
		}
		m.remaining = time.Until(m.nextPrayerTime)
		if m.remaining <= 0 {
			nextPrayer, nextPrayerTime, remaining, err := getNextPrayer(m.config)
			if err != nil {
				return m, tickCmd()
			}
			m.nextPrayer = nextPrayer
			m.nextPrayerTime = nextPrayerTime
			m.remaining = remaining
			m.currentPrayer = m.prayerTimes.CurrentPrayer(time.Now())

		}
		return m, tickCmd()
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "t":
			if m.view == todayView || m.view == weekView || m.view == monthView {
				m.view = todayView
				m.viewDate = time.Now()
			}
		case "w":
			if m.view == todayView || m.view == weekView || m.view == monthView {
				m.view = weekView
			}
		case "m":
			if m.view == todayView || m.view == weekView {
				visibleRows := m.monthVisibleRows()
				m.selectedDay = time.Now().Day()
				m.scrollOffset = m.selectedDay - 1 - (visibleRows / 2)
				if m.scrollOffset < 0 {
					m.scrollOffset = 0
				}
				totalDays := daysInMonth(time.Now())
				maxScroll := m.maxMonthScroll(totalDays, visibleRows)
				if m.scrollOffset > maxScroll {
					m.scrollOffset = maxScroll
				}
				if maxScroll < 0 {
					maxScroll = 0
				}
				m.view = monthView
			}
		case "down", "j":
			if m.view == monthView {
				visibileRows := m.monthVisibleRows()
				now := time.Now()
				viewDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, m.monthOffset, 0)
				totalDays := daysInMonth(viewDate)
				maxScroll := m.maxMonthScroll(totalDays, visibileRows)
				if m.scrollOffset < maxScroll {
					m.scrollOffset++
				}
				if m.selectedDay < totalDays {
					m.selectedDay++
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
				if m.selectedDay > 1 {
					m.selectedDay--
				}

			}

			if m.view == todayView {
				if m.selectedPrayer > 0 {
					m.selectedPrayer--
				}
			}
		case "right", "l":
			if m.view == weekView {
				m.weekOffset++
			}
			if m.view == monthView {
				m.monthOffset++
				m.scrollOffset = 0
				m.selectedDay = 1
			}
		case "left", "h":
			if m.view == weekView {
				m.weekOffset--
			}
			if m.view == monthView {
				m.monthOffset--
				m.scrollOffset = 0
				m.selectedDay = 1
			}

		case "enter":
			if m.view == todayView {
				m.view = prayerDetailView
			}
			if m.view == monthView {
				selectedDate := time.Date(
					time.Now().Year(),
					time.Now().Month()+time.Month(m.monthOffset),
					m.selectedDay,
					0, 0, 0, 0,
					time.Local,
				)
				m.viewDate = selectedDate
				m.previousView = m.view
				m.view = todayView

			}
		case "esc":
			if m.view == prayerDetailView {
				m.view = todayView
			} else if m.view == helpView {
				m.view = m.previousView
			} else if m.view == todayView && m.previousView == monthView {
				m.view = m.previousView
			}

		case "?":
			if m.view != helpView {
				m.previousView = m.view
				m.view = helpView
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
		viewDate:       time.Now(),
		prayerDate:     time.Now(),
		nextPrayerTime: nextPrayerTime,
		remaining:      remaining,
		view:           todayView,
		selectedDay:    time.Now().Day(),
	}
	for i, prayer := range prayerRows(m.prayerTimes) {
		if prayer.prayer == m.currentPrayer {
			m.selectedPrayer = i
			break
		}

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

func daysInMonth(date time.Time) int {
	firstDay := time.Date(
		date.Year(),
		date.Month(),
		1,
		0, 0, 0, 0,
		date.Location(),
	)
	lastDay := firstDay.AddDate(0, 1, -1)

	return lastDay.Day()
}

func activeRender(text string, date time.Time) string {
	if isToday(date) {
		return activeStyle.Render(text)

	}
	return text
}

func isToday(date time.Time) bool {
	now := time.Now()
	if date.Year() == now.Year() &&
		date.Month() == now.Month() &&
		date.Day() == now.Day() {
		return true
	}

	return false
}
