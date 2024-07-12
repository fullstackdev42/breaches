package cmd

import (
	"database/sql"
	"fmt"
	"fullstackdev42/breaches/data"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

type ViewCommand struct {
	dataHandler *data.DataHandler
}

func NewViewCommand(dataHandler *data.DataHandler) *ViewCommand {
	return &ViewCommand{
		dataHandler: dataHandler,
	}
}

func (v *ViewCommand) Command() *cobra.Command {
	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View the data in a sortable table",
		Long: `This command reads the data from the specified file and displays it in a sortable table.
		You can navigate through the table using the next and back buttons.`,
		Run: func(cmd *cobra.Command, args []string) {
			db, err := sql.Open("sqlite3", "./people.db")
			if err != nil {
				fmt.Println("Error opening database:", err)
				return
			}
			defer db.Close()

			people, err := v.dataHandler.FetchDataFromDB(db)
			if err != nil {
				fmt.Println("Error fetching data from database:", err)
				return
			}

			pageSize := 20
			totalPages := len(people) / pageSize
			if len(people)%pageSize != 0 {
				totalPages++
			}

			for i := 0; i < totalPages; i++ {
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"ID1", "ID2", "First Name", "Last Name", "Gender", "Birth Place", "Current Place", "Job", "Date"})

				start := i * pageSize
				end := start + pageSize
				if end > len(people) {
					end = len(people)
				}

				for _, person := range people[start:end] {
					t.AppendRow([]interface{}{person.ID1, person.ID2, person.FirstName, person.LastName, person.Gender, person.BirthPlace, person.CurrentPlace, person.Job, person.Date})
				}

				fmt.Printf("Page %d of %d\n", i+1, totalPages)
				t.Render()
				if i != totalPages-1 {
					fmt.Print("Press Enter to continue to the next page...")
					fmt.Scanln()
				}
			}
		},
	}

	return viewCmd
}
