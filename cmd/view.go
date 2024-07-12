package cmd

import (
	"fullstackdev42/breaches/data"
	"fullstackdev42/breaches/ui"

	_ "github.com/mattn/go-sqlite3"
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
	ui := ui.NewUI(v.dataHandler)
	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View the data in a sortable table",
		Long: `This command reads the data from the specified file and displays it in a sortable table.
		You can navigate through the table using the next and back buttons.`,
		Run: func(cmd *cobra.Command, args []string) {
			ui.RunUI()
		},
	}

	return viewCmd
}
