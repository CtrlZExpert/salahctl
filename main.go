package main

func main() {
	setupHelp()

	err := rootCmd.Execute()
	if err != nil {
		printError(err)
	}
}
