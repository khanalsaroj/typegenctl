package cmd

import (
	"github.com/sarojkhanal/typegenctl/internal/app/usecase"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start TypeGen services",
	Long: `Start the TypeGen service stack using the resolved system configuration.
This command launches the API and Web services according to the settings
defined in typegen.yaml, ensuring that required networks, volumes, and runtime
dependencies are in place. If services are already running, the command
performs no destructive actions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.Run(opts)
	},
}

func init() {
	addServiceFlags(runCmd)
	rootCmd.AddCommand(runCmd)
}
