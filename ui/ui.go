package ui

import (
	"fmt"

	"fullstackdev42/breaches/data"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct{}

func NewUI() *UI {
	return &UI{}
}

func (ui *UI) RunUI(people []data.Person, offset int, pageSize int, nextPage func() ([]data.Person, error), prevPage func() ([]data.Person, error)) {
	app := tview.NewApplication()
	app.EnableMouse(true) // Enable mouse support

	table := ui.RenderTable(app, people)

	// Create a new page for the data
	page := tview.NewFlex().SetDirection(tview.FlexRow)
	page.AddItem(table, 0, 1, true)

	// Add a footer with pagination
	footer := tview.NewTextView().SetText(fmt.Sprintf("Page %d", offset/pageSize+1))
	page.AddItem(footer, 1, 1, false)

	// Add key handlers for 'n' and 'p'
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'n':
			// Fetch the next page of data
			people, err := nextPage()
			if err != nil {
				fmt.Println("Error fetching data from database:", err)
				return event
			}
			// Update the table and footer
			table = ui.RenderTable(app, people)
			footer.SetText(fmt.Sprintf("Page %d", offset/pageSize+1))
		case 'p':
			// Fetch the previous page of data
			people, err := prevPage()
			if err != nil {
				fmt.Println("Error fetching data from database:", err)
				return event
			}
			// Update the table and footer
			table = ui.RenderTable(app, people)
			footer.SetText(fmt.Sprintf("Page %d", offset/pageSize+1))
		}
		return event
	})

	app.SetRoot(page, true)

	err := app.Run()
	if err != nil {
		fmt.Println("Error running application:", err)
	}
}

func (ui *UI) RenderTable(app *tview.Application, people []data.Person) *tview.Table {
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
