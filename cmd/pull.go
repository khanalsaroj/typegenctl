package cmd

import (
	"github.com/sarojkhanal/typegenctl/internal/app/usecase"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull required Docker images for TypeGen services",
	Long: `Fetch and verify all Docker images required to run TypeGen services.
This command ensures that the correct service images are present in the local
Docker registry. Images are pulled only when missing or outdated, making the
operation safe to run multiple times without unnecessary downloads.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.Pull(opts)
	},
}

func init() {
	addServiceFlags(pullCmd)
	rootCmd.AddCommand(pullCmd)
}
