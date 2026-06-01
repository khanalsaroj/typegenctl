package cmd

import (
	"github.com/khanalsaroj/typegenctl/internal/app/usecase"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate system configuration and runtime environment",
	Long: `Perform a comprehensive validation of the TypeGen system setup.

This command verifies configuration correctness, host prerequisites, and
Docker availability without modifying any state. It is intended for
pre-flight checks prior to initialization, startup, or deployment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.Check(opts)
	},
}

func init() {
	addServiceFlags(checkCmd)
	rootCmd.AddCommand(checkCmd)
}
