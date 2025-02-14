package ui

import (
	"fmt"

	"fullstackdev42/breaches/data"

	"github.com/gdamore/tcell/v2"
	"github.com/jonesrussell/loggo"
	"github.com/rivo/tview"
)

const (
	TableIndex   = 0
	FooterIndex  = 1
	IDLength     = 20
	NameLength   = 20
	GenderLength = 10
	PlaceLength  = 25
	JobLength    = 20
	DateLength   = 16
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

func (ui *UI) CreateDataTable() *tview.Table {
	t := tview.NewTable().SetBorders(true)

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

	return t
}

func (ui *UI) RunUI(people []data.Person, pagination *Pagination) error {
	ui.app.EnableMouse(true) // Enable mouse support

	// Add a footer with pagination and total items (if available)
	ui.footer = tview.NewTextView()
	ui.SetupFooter(pagination)

	// Create the initial table structure
	ui.table = ui.CreateDataTable()

	// Populate the table with data
	ui.table = ui.PopulateTable(ui.table, people)
	ui.page.AddItem(ui.table, TableIndex, 1, true)
	ui.page.AddItem(ui.footer, FooterIndex, 1, false)

	// Add key handlers for 'n' (next), 'p' (previous), 's' (sort by specific column)
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'n':
			err := ui.updateTable(pagination.NextPage, pagination.Logger, pagination, true)
			if err != nil {
				ui.footer.SetText(fmt.Sprintf("Error fetching next page: %v", err))
				return event
			}
		case 'p':
			err := ui.updateTable(pagination.PrevPage, pagination.Logger, pagination, false)
			if err != nil {
				ui.footer.SetText(fmt.Sprintf("Error fetching previous page: %v", err))
				return event
			}
		}
		return event
	})

	ui.app.SetRoot(ui.page, true)

	return ui.app.Run()
}

func (ui *UI) updateTable(fetchPage func(loggo.LoggerInterface) ([]data.Person, error), logger loggo.LoggerInterface, pagination *Pagination, isNext bool) error {
	people, err := ui.fetchPageData(fetchPage, logger)
	if err != nil {
		return err
	}

	ui.updatePaginationOffset(pagination, isNext)
	ui.clearTableRows()
	ui.populateTableWithData(people)
	ui.updateFooterText(pagination)

	return nil
}

func (ui *UI) fetchPageData(fetchPage func(loggo.LoggerInterface) ([]data.Person, error), logger loggo.LoggerInterface) ([]data.Person, error) {
	people, err := fetchPage(logger)
	if err != nil {
		logger.Error("error fetching page:", err)
		return nil, err
	}
	return people, nil
}

func (ui *UI) updatePaginationOffset(pagination *Pagination, isNext bool) {
	if isNext {
		pagination.Offset += pagination.PageSize
	} else {
		if pagination.Offset-pagination.PageSize >= 0 {
			pagination.Offset -= pagination.PageSize
		}
	}
}

func (ui *UI) clearTableRows() {
	for i := ui.table.GetRowCount() - 1; i > 0; i-- {
		ui.table.RemoveRow(i)
	}
}

func (ui *UI) populateTableWithData(people []data.Person) {
	ui.table = ui.PopulateTable(ui.table, people)
}

func (ui *UI) updateFooterText(pagination *Pagination) {
	ui.SetupFooter(pagination)
}

// Truncate truncates a string to the specified length.
func Truncate(s string, length int) string {
	if len(s) > length {
		return s[:length]
	}
	return s
}

func FormatPersonData(person data.Person) data.Person {
	person.ID1 = Truncate(person.ID1, IDLength)
	person.ID2 = Truncate(person.ID2, IDLength)
	person.FirstName = Truncate(person.FirstName, NameLength)
	person.LastName = Truncate(person.LastName, NameLength)
	person.Gender = Truncate(person.Gender, GenderLength)
	person.BirthPlace = Truncate(person.BirthPlace, PlaceLength)
	person.CurrentPlace = Truncate(person.CurrentPlace, PlaceLength)
	person.Job = Truncate(person.Job, JobLength)
	person.Date = Truncate(person.Date, DateLength)

	return person
}

func (ui *UI) PopulateTable(t *tview.Table, people []data.Person) *tview.Table {
	// Add data with potential truncation (adjust max length as needed)
	for i, person := range people {
		person = FormatPersonData(person)

		t.SetCell(i+1, 0, tview.NewTableCell(person.ID1).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 1, tview.NewTableCell(person.ID2).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 2, tview.NewTableCell(person.FirstName).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 3, tview.NewTableCell(person.LastName).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 4, tview.NewTableCell(person.Gender).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 5, tview.NewTableCell(person.BirthPlace).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 6, tview.NewTableCell(person.CurrentPlace).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 7, tview.NewTableCell(person.Job).SetAlign(tview.AlignCenter))
		t.SetCell(i+1, 8, tview.NewTableCell(person.Date).SetAlign(tview.AlignCenter))
	}

	return t
}

func (ui *UI) SetupFooter(pagination *Pagination) {
	totalPages := (pagination.Total + pagination.PageSize - 1) / pagination.PageSize
	currentPage := pagination.Offset/pagination.PageSize + 1
	footerText := fmt.Sprintf("Page %d/%d", currentPage, totalPages)
	if pagination.Total > 0 {
		footerText += fmt.Sprintf(" (Total: %d)", pagination.Total)
	}
	ui.footer.SetText(footerText)
}
