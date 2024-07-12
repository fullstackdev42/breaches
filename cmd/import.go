package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

type Person struct {
	ID1, ID2, FirstName, LastName, Gender, BirthPlace, CurrentPlace, Job, Date string
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import the data into an SQLite database",
	Long: `This command reads the data from the specified file and loads it into an SQLite database.
The data is stored in a table with columns corresponding to the fields of the data.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", "./canada.db")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer db.Close()

		sqlStmt := `
		CREATE TABLE IF NOT EXISTS people (
			ID1 TEXT,
			ID2 TEXT,
			FirstName TEXT,
			LastName TEXT,
			Gender TEXT,
			BirthPlace TEXT,
			CurrentPlace TEXT,
			Job TEXT,
			Date TEXT
		);
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Printf("%q: %s\n", err, sqlStmt)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			fmt.Println(err)
			return
		}
		stmt, err := tx.Prepare("INSERT INTO people(ID1, ID2, FirstName, LastName, Gender, BirthPlace, CurrentPlace, Job, Date) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		file, err := os.Open("data/Canada.txt")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, ":")
			if len(parts) >= 10 {
				_, err = stmt.Exec(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6], parts[8], parts[9])
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		tx.Commit()

		fmt.Println("Data loaded into SQLite database.")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
