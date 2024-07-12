package ui

import (
	"fmt"

	"fullstackdev42/breaches/data"

	"github.com/gdamore/tcell/v2"
	"github.com/jonesrussell/loggo"
	"github.com/rivo/tview"
)

const (
	TableIndex  = 0
	FooterIndex = 1
)

type Pagination struct {
	Offset   int
	PageSize int
	NextPage func(loggo.LoggerInterface) ([]data.Person, error)
	PrevPage func(loggo.LoggerInterface) ([]data.Person, error)
	Logger   loggo.LoggerInterface
}

type UI struct {
	app    *tview.Application
	page   *tview.Flex
	table  *tview.Table
	footer *tview.TextView
}

func NewUI() *UI {
	return &UI{
		app:  tview.NewApplication(),
		page: tview.NewFlex().SetDirection(tview.FlexRow),
	}
}

func (ui *UI) RunUI(people []data.Person, pagination *Pagination) {
	ui.app.EnableMouse(true) // Enable mouse support

	// Add a footer with pagination
	ui.footer = tview.NewTextView()

	// Render the initial table
	ui.table = ui.RenderTable(people)
	ui.page.AddItem(ui.table, TableIndex, 1, true)
	ui.page.AddItem(ui.footer, FooterIndex, 1, false)

	// Add key handlers for 'n' and 'p'
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'n':
			ui.updateTable(pagination.NextPage, pagination.Logger, pagination.Offset, pagination.PageSize)
		case 'p':
			ui.updateTable(pagination.PrevPage, pagination.Logger, pagination.Offset, pagination.PageSize)
		}
		return event
	})

	ui.app.SetRoot(ui.page, true)

	err := ui.app.Run()
	if err != nil {
		pagination.Logger.Error("error running application:", err)
	}
}

func (ui *UI) updateTable(fetchPage func(loggo.LoggerInterface) ([]data.Person, error), logger loggo.LoggerInterface, offset int, pageSize int) {
	// Fetch the next page of data
	people, err := fetchPage(logger)
	if err != nil {
		logger.Error("error fetching data from database:", err)
		return
	}

	// Update the table and footer
	ui.page.RemoveItem(ui.table)
	ui.table = ui.RenderTable(people)
	ui.page.AddItem(ui.table, TableIndex, 1, true)
	ui.footer.SetText(fmt.Sprintf("Page %d", offset/pageSize+1))
}

func (ui *UI) RenderTable(people []data.Person) *tview.Table {
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
