package cmd

import (
	"fmt"
	"fullstackdev42/breaches/data"

	"github.com/rivo/tview"
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
	t := tview.NewTable()

	for i, person := range people {
		t.SetCell(i, 0, tview.NewTableCell(person.ID1))
		t.SetCell(i, 1, tview.NewTableCell(person.ID2))
		t.SetCell(i, 2, tview.NewTableCell(person.FirstName))
		t.SetCell(i, 3, tview.NewTableCell(person.LastName))
		t.SetCell(i, 4, tview.NewTableCell(person.Gender))
		t.SetCell(i, 5, tview.NewTableCell(person.BirthPlace))
		t.SetCell(i, 6, tview.NewTableCell(person.CurrentPlace))
		t.SetCell(i, 7, tview.NewTableCell(person.Job))
		t.SetCell(i, 8, tview.NewTableCell(person.Date))
	}

	app := tview.NewApplication()
	err := app.SetRoot(t, true).Run()
	if err != nil {
		fmt.Println("Error running application:", err)
	}
}
