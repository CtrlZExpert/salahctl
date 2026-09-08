package main

import (
	"github.com/spf13/cobra"
)

func setupHelp() {
	cobra.AddTemplateFunc("title", func(text string) string {
		return titleStyle.Render(text)
	})

	cobra.AddTemplateFunc("heading", func(text string) string {
		return headingStyle.Render(text)
	})

	cobra.AddTemplateFunc("muted", func(text string) string {
		return mutedStyle.Render(text)
	})

	cobra.AddTemplateFunc("command", func(text string, width int) string {
		return labelStyle.Width(width).Render(text)
	})

	rootCmd.SetHelpTemplate(`{{title .CommandPath}}
{{with .Short}}{{.}}{{end}}

{{heading "Usage:"}}
  {{.UseLine}}
{{if .HasAvailableSubCommands}}

{{heading "Available Commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{command .Name .NamePadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{heading "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

{{heading "Global Flags:"}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

{{heading "Additional help topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{command .Name .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

{{muted (printf "Use \"%s [command] --help\" for more information about a command." .CommandPath)}}{{end}}
`)
}
