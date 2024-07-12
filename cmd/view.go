package cmd

import (
	"fmt"
	"fullstackdev42/breaches/data"
	"os"

	"github.com/olekukonko/tablewriter"
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
			people, err := v.dataHandler.LoadDataFromFile()
			if err != nil {
				fmt.Println("Error loading data from file:", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID1", "ID2", "First Name", "Last Name", "Gender", "Birth Place", "Current Place", "Job", "Date"})

			for _, person := range people {
				table.Append([]string{person.ID1, person.ID2, person.FirstName, person.LastName, person.Gender, person.BirthPlace, person.CurrentPlace, person.Job, person.Date})
			}

			table.Render()
		},
	}

	return viewCmd
}
