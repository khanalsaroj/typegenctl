package cmd

import (
	"github.com/khanalsaroj/typegenctl/internal/app/usecase"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop running TypeGen services",
	Long: `Gracefully stop running TypeGen service containers.

This command stops the API and Web containers while preserving their
configuration, networks, and volumes. Containers are not removed and can be
restarted later using the start or restart commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.Stop(opts)
	},
}

func init() {
	addServiceFlags(stopCmd)
	rootCmd.AddCommand(stopCmd)
}
