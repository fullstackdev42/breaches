package main

import (
	"database/sql"
	"fmt"
	"fullstackdev42/breaches/cmd"
	"fullstackdev42/breaches/data"

	"github.com/jonesrussell/loggo"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/canada.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Create a new logger instance
	logger, err := loggo.NewLogger("./loggo.log")
	if err != nil {
		fmt.Println("Error creating logger:", err)
		return
	}

	dataHandler := data.NewDataHandler("./data/Canada.txt", db)

	rootCmd := cmd.NewRootCmd(dataHandler, &logger)

	err = rootCmd.Execute()
	if err != nil {
		logger.Error("root command execute failed", err)
	}
}
