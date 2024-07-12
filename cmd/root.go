package cmd

import (
	"fmt"
	"os"

	"fullstackdev42/breaches/data"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

type RootCmd struct {
	dataHandler *data.DataHandler
}

func NewRootCmd(dataHandler *data.DataHandler) *RootCmd {
	return &RootCmd{
		dataHandler: dataHandler,
	}
}

func (r *RootCmd) rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "breaches",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
		examples and usage of using your application. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			// You can use r.dataHandler here
		},
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.breaches.yaml)")
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cobra.OnInitialize(func() {
		initConfig(cfgFile)
	})

	viewCmd := NewViewCmd(r.dataHandler)
	importCmd := NewImportCmd(r.dataHandler)

	cmd.AddCommand(viewCmd.viewCmd())
	cmd.AddCommand(importCmd.importCmd())

	return cmd
}

func initConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".breaches")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func (r *RootCmd) Execute() {
	err := r.rootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
