package cmd

import (
	"fmt"
	"fullstackdev42/breaches/data"
	"fullstackdev42/breaches/ui"

	"github.com/jonesrussell/loggo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

type ViewCommand struct {
	dataHandler *data.DataHandler
}

func NewViewCommand(dataHandler *data.DataHandler, logger *loggo.LoggerInterface) *ViewCommand {
	return &ViewCommand{
		dataHandler: dataHandler,
	}
}

func (v *ViewCommand) Command() *cobra.Command {
	ui := ui.NewUI()

	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View the data in a sortable table",
		Long: `This command reads the data from the specified file and displays it in a sortable table.
		You can navigate through the table using the next and back buttons.`,
		Run: func(cmd *cobra.Command, args []string) {
			pageSize := 20
			offset := 0

			// Fetch the initial data
			people, err := v.dataHandler.FetchDataFromDB(offset, pageSize)
			if err != nil {
				fmt.Println("Error fetching data from database:", err)
				return
			}

			// Define the functions to fetch the next and previous pages
			nextPage := func() ([]data.Person, error) {
				offset += pageSize
				return v.dataHandler.FetchDataFromDB(offset, pageSize)
			}
			prevPage := func() ([]data.Person, error) {
				if offset > 0 {
					offset -= pageSize
				}
				return v.dataHandler.FetchDataFromDB(offset, pageSize)
			}

			ui.RunUI(people, offset, pageSize, nextPage, prevPage)
		},
	}

	return viewCmd
}
