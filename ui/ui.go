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
	Total    int // Total number of items (optional)
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

func (ui *UI) RunUI(people []data.Person, pagination *Pagination) error {
	ui.app.EnableMouse(true) // Enable mouse support

	// Add a footer with pagination and total items (if available)
	ui.footer = tview.NewTextView()
	footerText := fmt.Sprintf("Page %d", pagination.Offset/pagination.PageSize+1)
	if pagination.Total > 0 {
		footerText += fmt.Sprintf(" (Total: %d)", pagination.Total)
	}
	ui.footer.SetText(footerText)

	// Render the initial table
	ui.table = ui.RenderTable(people)
	ui.page.AddItem(ui.table, TableIndex, 1, true)
	ui.page.AddItem(ui.footer, FooterIndex, 1, false)

	// Add key handlers for 'n' (next), 'p' (previous), 's' (sort by specific column)
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'n':
			err := ui.updateTable(pagination.NextPage, pagination.Logger, pagination)
			if err != nil {
				pagination.Logger.Error("error fetching next page:", err)
				// Display error to user?
			}
		case 'p':
			err := ui.updateTable(pagination.PrevPage, pagination.Logger, pagination)
			if err != nil {
				pagination.Logger.Error("error fetching previous page:", err)
				// Display error to user?
			}
			// Add cases for other sorting keys (e.g., 's1' for sorting by first name)
		}
		return event
	})

	ui.app.SetRoot(ui.page, true)

	return ui.app.Run()
}

func (ui *UI) updateTable(fetchPage func(loggo.LoggerInterface) ([]data.Person, error), logger loggo.LoggerInterface, pagination *Pagination) error {
	// Fetch the next/previous page of data
	people, err := fetchPage(logger)
	if err != nil {
		return err // Propagate the error
	}

	// Update the table and footer
	ui.page.RemoveItem(ui.table)
	ui.table = ui.RenderTable(people)
	ui.page.AddItem(ui.table, TableIndex, 1, true)
	footerText := fmt.Sprintf("Page %d", pagination.Offset/pagination.PageSize+1)
	if pagination.Total > 0 {
		footerText += fmt.Sprintf(" (Total: %d)", pagination.Total)
	}
	ui.footer.SetText(footerText)
	return nil
}

func (ui *UI) RenderTable(people []data.Person) *tview.Table {
	t := tview.NewTable().
		SetBorders(true)

	// Add headers with alignment
	t.SetCell(0, 0, tview.NewTableCell("ID1").SetAlign(tview.AlignCenter))
	t.SetCell(0, 1, tview.NewTableCell("ID2").SetAlign(tview.AlignCenter))
	t.SetCell(0, 2, tview.NewTableCell("First Name").SetAlign(tview.AlignCenter))
	t.SetCell(0, 3, tview.NewTableCell("Last Name").SetAlign(tview.AlignCenter))
	t.SetCell(0, 4, tview.NewTableCell("Gender").SetAlign(tview.AlignCenter))
	t.SetCell(0, 5, tview.NewTableCell("Birth Place").SetAlign(tview.AlignCenter))
	t.SetCell(0, 6, tview.NewTableCell("Current Place").SetAlign(tview.AlignCenter))
	t.SetCell(0, 7, tview.NewTableCell("Job").SetAlign(tview.AlignCenter))
	t.SetCell(0, 8, tview.NewTableCell("Date").SetAlign(tview.AlignCenter))

	// Add data with potential truncation (adjust max length as needed)
	for i, person := range people {
		truncatedID1 := person.ID1
		if len(person.ID1) > 20 {
			truncatedID1 = person.ID1[:20] // Truncate to 20 characters
		}

		truncatedID2 := person.ID2
		if len(person.ID2) > 20 {
			truncatedID2 = person.ID2[:20] // Truncate to 20 characters
		}

		truncatedFirstName := person.FirstName
		if len(person.FirstName) > 20 {
			truncatedFirstName = person.FirstName[:20]
		}

		truncatedLastName := person.LastName
		if len(person.LastName) > 20 {
			truncatedLastName = person.LastName[:20]
		}

		truncatedGender := person.Gender
		if len(person.Gender) > 10 {
			truncatedGender = person.Gender[:10] // Adjust truncation length
		}

		truncatedBirthPlace := person.BirthPlace
		if len(person.BirthPlace) > 25 {
			truncatedBirthPlace = person.BirthPlace[:25]
		}

		truncatedCurrentPlace := person.CurrentPlace
		if len(person.CurrentPlace) > 25 {
			truncatedCurrentPlace = person.CurrentPlace[:25]
		}

		truncatedJob := person.Job
		if len(person.Job) > 20 {
			truncatedJob = person.Job[:20]
		}

		truncatedDate := person.Date
		if len(person.Date) > 16 {
			truncatedDate = person.Date[:16] // Adjust truncation length
		}

		t.SetCell(i+1, 0, tview.NewTableCell(truncatedID1).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 1, tview.NewTableCell(truncatedID2).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 2, tview.NewTableCell(truncatedFirstName).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 3, tview.NewTableCell(truncatedLastName).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 4, tview.NewTableCell(truncatedGender).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 5, tview.NewTableCell(truncatedBirthPlace).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 6, tview.NewTableCell(truncatedCurrentPlace).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 7, tview.NewTableCell(truncatedJob).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 8, tview.NewTableCell(truncatedDate).SetAlign(tview.AlignCenter))
	}

	return t // Return the created table
}
