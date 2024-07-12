package cmd

import (
	"fmt"
	"fullstackdev42/breaches/data"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
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
			v.RunViewCommand()
		},
	}

	return viewCmd
}

func (v *ViewCommand) RunViewCommand() {
	db, err := v.dataHandler.OpenDB("./data/canada.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	pageSize := 20
	offset := 0

	for {
		people, err := v.dataHandler.FetchDataFromDB(db, offset, pageSize)
		if err != nil {
			fmt.Println("Error fetching data from database:", err)
			return
		}

		if len(people) == 0 {
			break
		}

		v.RenderTable(people)

		offset += pageSize
	}
}

func (v *ViewCommand) RenderTable(people []data.Person) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID1", "ID2", "First Name", "Last Name", "Gender", "Birth Place", "Current Place", "Job", "Date"})
	t.SetPageSize(20)

	for _, person := range people {
		t.AppendRow([]interface{}{person.ID1, person.ID2, person.FirstName, person.LastName, person.Gender, person.BirthPlace, person.CurrentPlace, person.Job, person.Date})
	}

	t.Render()
}
