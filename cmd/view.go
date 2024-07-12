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
	pageSize := 20
	offset := 0

	app := tview.NewApplication()
	app.EnableMouse(true) // Enable mouse support

	pages := tview.NewPages()

	// Fetch the initial data
	people, err := v.dataHandler.FetchDataFromDB(offset, pageSize)
	if err != nil {
		fmt.Println("Error fetching data from database:", err)
		return
	}

	table := v.RenderTable(app, people)

	// Create a new page for the data
	page := tview.NewFlex().SetDirection(tview.FlexRow)
	page.AddItem(table, 0, 1, true)

	// Add a footer with pagination
	footer := tview.NewTextView().SetText(fmt.Sprintf("Page %d", offset/pageSize+1))
	page.AddItem(footer, 1, 1, false)

	pages.AddPage("main", page, true, true)

	// Handle input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'n': // Next page
			offset += pageSize
		case 'p': // Previous page
			offset -= pageSize
			if offset < 0 {
				offset = 0
			}
		}

		// Fetch the new data
		people, err := v.dataHandler.FetchDataFromDB(offset, pageSize)
		if err != nil {
			fmt.Println("Error fetching data from database:", err)
			return event
		}

		// Update the table and footer
		table.Clear()
		table = v.RenderTable(app, people)
		footer.SetText(fmt.Sprintf("Page %d", offset/pageSize+1))

		return event
	})

	app.SetRoot(pages, true)

	err = app.Run()
	if err != nil {
		fmt.Println("Error running application:", err)
	}
}

func (v *ViewCommand) RenderTable(app *tview.Application, people []data.Person) *tview.Table {
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

	return t // Return the created table
}
