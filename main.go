package main

import (
	"database/sql"
	"fmt"
	"fullstackdev42/breaches/cmd"
	"fullstackdev42/breaches/data"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/canada.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	dataHandler := data.NewDataHandler("./data/Canada.txt", db)

	rootCmd := cmd.NewRootCmd(dataHandler)

	rootCmd.Execute()
}
