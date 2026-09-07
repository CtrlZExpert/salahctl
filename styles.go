package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var primaryColor = lipgloss.Color("#10B981")   //emerald
var secondaryColor = lipgloss.Color("#6EE7B7") //light emerald
var successColor = lipgloss.Color("#22C55E")   //green
var mutedColor = lipgloss.Color("#6B7280")     //gray
var errorColor = lipgloss.Color("#EF4444")     //red

var titleStyle = lipgloss.NewStyle().
	Foreground(primaryColor).
	Bold(true)

var headingStyle = lipgloss.NewStyle().
	Foreground(secondaryColor).
	Bold(true)

var labelStyle = lipgloss.NewStyle().
	Foreground(secondaryColor)

var valueStyle = lipgloss.NewStyle()

var activeStyle = lipgloss.NewStyle().
	Foreground(successColor).
	Bold(true)

var successStyle = lipgloss.NewStyle().
	Foreground(successColor)

var errorStyle = lipgloss.NewStyle().
	Foreground(errorColor).
	Bold(true)

var mutedStyle = lipgloss.NewStyle().
	Foreground(mutedColor)

func printError(err error) {
	fmt.Println(errorStyle.Render("Error: " + err.Error()))
}

func printErrorMessage(message string) {
	fmt.Println(errorStyle.Render(message))
}
