package cmd

import (
	"github.com/sarojkhanal/typegenctl/internal/app/usecase"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart running TypeGen services",
	Long: `Restart existing TypeGen service containers in a controlled manner.

This command restarts the API and Web service containers while preserving
their configuration, networks, and volumes. Containers must already exist;
use the init or run commands to create services if they are not present.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.Restart(opts)
	},
}

func init() {
	addServiceFlags(restartCmd)
	rootCmd.AddCommand(restartCmd)
}
