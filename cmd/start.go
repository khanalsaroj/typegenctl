package cmd

import (
	"github.com/khanalsaroj/typegenctl/internal/app/usecase"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start existing TypeGen services",
	Long: `Start previously created TypeGen service containers.

This command starts the API and Web containers using the existing runtime
configuration defined in typegen.yaml. It does not recreate containers or
modify configuration. Services must already be initialized; use the init
or run commands if the services do not yet exist.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.Start(opts)
	},
}

func init() {
	addServiceFlags(startCmd)
	rootCmd.AddCommand(startCmd)
}
