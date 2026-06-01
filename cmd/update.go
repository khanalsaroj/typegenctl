package cmd

import (
	"github.com/khanalsaroj/typegenctl/internal/app/usecase"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update TypeGen services to the latest images",
	Long: `Check for and fetch updated Docker images for TypeGen services.

This command pulls the latest available service images and reports the update
status without modifying existing configuration or restarting containers.
Use the restart command to apply updates to running services.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.Update(opts)
	},
}

func init() {
	addServiceFlags(updateCmd)
	rootCmd.AddCommand(updateCmd)
}
