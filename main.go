package main

import (
	"fullstackdev42/breaches/cmd"
	"fullstackdev42/breaches/data"
)

func main() {
	// Create a new data handler
	dataHandler := data.NewDataHandler("data/Canada.txt")

	rootCmd := cmd.NewRootCmd(dataHandler)

	rootCmd.Execute()
}
