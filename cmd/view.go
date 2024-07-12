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
	logger      *loggo.LoggerInterface
}

func NewViewCommand(dataHandler *data.DataHandler, logger *loggo.LoggerInterface) *ViewCommand {
	return &ViewCommand{
		dataHandler: dataHandler,
		logger:      logger,
	}
}

func (v *ViewCommand) Command() *cobra.Command {
	userInterface := ui.NewUI()

	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View the data in a sortable table",
		Long: `This command reads the data from the specified file and displays it in a sortable table.
		You can navigate through the table using the next and back buttons.`,
		Run: func(cmd *cobra.Command, args []string) {
			pageSize := 20
			offset := 0

			v.runCommand(offset, pageSize, userInterface)
		},
	}

	return viewCmd
}

func (v *ViewCommand) runCommand(offset int, pageSize int, userInterface *ui.UI) {
	// Fetch the initial data
	people, err := v.dataHandler.FetchDataFromDB(offset, pageSize)
	if err != nil {
		fmt.Println("Error fetching data from database:", err)
		return
	}

	// Get the total number of items
	total, err := v.dataHandler.GetTotalItems()
	if err != nil {
		fmt.Println("Error fetching total number of items:", err)
		return
	}

	pagination := &ui.Pagination{
		Offset:   offset,
		PageSize: pageSize,
		NextPage: func(logger loggo.LoggerInterface) ([]data.Person, error) {
			return v.nextPage(logger, offset, pageSize)
		},
		PrevPage: func(logger loggo.LoggerInterface) ([]data.Person, error) {
			return v.prevPage(logger, offset, pageSize)
		},
		Logger: *v.logger,
		Total:  total, // Set the total number of items
	}

	userInterface.RunUI(people, pagination)
}

func (v *ViewCommand) nextPage(logger loggo.LoggerInterface, offset int, pageSize int) ([]data.Person, error) {
	logger.Debug("nextPage called")
	offset += pageSize
	return v.dataHandler.FetchDataFromDB(offset, pageSize)
}

func (v *ViewCommand) prevPage(logger loggo.LoggerInterface, offset int, pageSize int) ([]data.Person, error) {
	logger.Debug("prevPage called")

	if offset > 0 {
		offset -= pageSize
	}
	return v.dataHandler.FetchDataFromDB(offset, pageSize)
}
