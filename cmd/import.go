package cmd

import (
	"database/sql"
	"fmt"
	"fullstackdev42/breaches/data"

	"github.com/spf13/cobra"
)

type ImportCmd struct {
	dataHandler *data.DataHandler
}

func NewImportCmd(dataHandler *data.DataHandler) *ImportCmd {
	return &ImportCmd{
		dataHandler: dataHandler,
	}
}

func (i *ImportCmd) importCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import the data into an SQLite database",
		Long: `This command reads the data from the specified file and loads it into an SQLite database.
		The data is stored in a table with columns corresponding to the fields of the data.`,
		Run: func(cmd *cobra.Command, args []string) {
			people, err := i.dataHandler.LoadDataFromFile()
			if err != nil {
				fmt.Println("Error loading data from file:", err)
				return
			}

			db, err := sql.Open("sqlite3", "./canada.db")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer db.Close()

			err = i.dataHandler.LoadDataIntoDB(db, people)
			if err != nil {
				fmt.Println("Error loading data into database:", err)
				return
			}

			fmt.Println("Data loaded into SQLite database.")
		},
	}

	return cmd
}
