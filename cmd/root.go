package cmd

import (
	"os"

	"fullstackdev42/breaches/data"

	logger "github.com/jonesrussell/loggo"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "breaches",
	Short: "A brief description of your application",
}

// NewRootCmd now returns *cobra.Command
func NewRootCmd(
	dataHandler *data.DataHandler,
	logger *logger.LoggerInterface,
) *cobra.Command {
	// cfg := NewConfig()

	// rootCmd.PersistentFlags().BoolVarP(&cfg.Debug, "debug", "d", false, "Enable debug mode")

	importCommand := NewImportCommand(dataHandler)
	viewCommand := NewViewCommand(dataHandler, logger)

	rootCmd.AddCommand(importCommand.Command())
	rootCmd.AddCommand(viewCommand.Command())

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
