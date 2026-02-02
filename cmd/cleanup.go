package cmd

import (
	"github.com/sarojkhanal/typegenctl/internal/app/usecase"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove obsolete TypeGen Docker artifacts",
	Long: `Clean up unused TypeGen Docker resources to reclaim local disk space.

This command removes outdated TypeGen Docker images and stopped containers,
retaining only the latest images required by typegen.yaml. Running containers
are never stopped or removed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.CleanUp(opts)
	},
}

func init() {
	addServiceFlags(cleanupCmd)
	rootCmd.AddCommand(cleanupCmd)
}
