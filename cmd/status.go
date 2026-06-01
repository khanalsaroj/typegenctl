package cmd

import (
	"github.com/khanalsaroj/typegenctl/internal/app/usecase"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display TypeGen service status",
	Long: `Inspect and report the current runtime status of TypeGen services.

This command provides a read-only view of the API and Web service containers,
including their existence and running state. No configuration or runtime
changes are performed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.Status(opts)
	},
}

func init() {
	addServiceFlags(statusCmd)
	rootCmd.AddCommand(statusCmd)
}
