package cmd

import (
	"fmt"
	"fullstackdev42/breaches/data"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

type ViewCommand struct {
	dataHandler *data.DataHandler
	stopped     bool // Add stopped as a field of the ViewCommand struct
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

	app := tview.NewApplication()

	for {
		// If the application has been stopped, break the loop
		if v.stopped {
			break
		}

		people, err := v.dataHandler.FetchDataFromDB(offset, pageSize)
		if err != nil {
			fmt.Println("Error fetching data from database:", err)
			return
		}

		if len(people) == 0 {
			break
		}

		v.RenderTable(app, people)

		offset += pageSize
	}
}

func (v *ViewCommand) RenderTable(app *tview.Application, people []data.Person) {
	t := tview.NewTable()

	// Add headers
	t.SetCell(0, 0, tview.NewTableCell("ID1").SetAlign(tview.AlignCenter))
	t.SetCell(0, 1, tview.NewTableCell("ID2").SetAlign(tview.AlignCenter))
	t.SetCell(0, 2, tview.NewTableCell("First Name").SetAlign(tview.AlignCenter))
	t.SetCell(0, 3, tview.NewTableCell("Last Name").SetAlign(tview.AlignCenter))
	t.SetCell(0, 4, tview.NewTableCell("Gender").SetAlign(tview.AlignCenter))
	t.SetCell(0, 5, tview.NewTableCell("Birth Place").SetAlign(tview.AlignCenter))
	t.SetCell(0, 6, tview.NewTableCell("Current Place").SetAlign(tview.AlignCenter))
	t.SetCell(0, 7, tview.NewTableCell("Job").SetAlign(tview.AlignCenter))
	t.SetCell(0, 8, tview.NewTableCell("Date").SetAlign(tview.AlignCenter))

	// Add data
	for i, person := range people {
		t.SetCell(i+1, 0, tview.NewTableCell(person.ID1))
		t.SetCell(i+1, 1, tview.NewTableCell(person.ID2))
		t.SetCell(i+1, 2, tview.NewTableCell(person.FirstName))
		t.SetCell(i+1, 3, tview.NewTableCell(person.LastName))
		t.SetCell(i+1, 4, tview.NewTableCell(person.Gender))
		t.SetCell(i+1, 5, tview.NewTableCell(person.BirthPlace))
		t.SetCell(i+1, 6, tview.NewTableCell(person.CurrentPlace))
		t.SetCell(i+1, 7, tview.NewTableCell(person.Job))
		t.SetCell(i+1, 8, tview.NewTableCell(person.Date))
	}

	app.SetRoot(t, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			// Set the stopped flag to true when the application is stopped
			v.stopped = true
			app.Stop()
		}
		return event
	})

	err := app.Run()
	if err != nil {
		fmt.Println("Error running application:", err)
	}
}
