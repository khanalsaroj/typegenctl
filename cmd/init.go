package cmd

import (
	"github.com/khanalsaroj/typegenctl/internal/app/usecase"
	"github.com/khanalsaroj/typegenctl/internal/domain"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstrap the TypeGen runtime environment",
	Long: `Initialize and prepare the local system for running TypeGen services.

This command performs a safe, idempotent bootstrap process that:
  - Validates host prerequisites and configuration integrity
  - Creates required directories and filesystem structure
  - Pulls and verifies required Docker images
  - Creates Docker networks if they do not already exist

The operation is idempotent and can be executed multiple times without
side effects. Use the --force flag to explicitly override existing
configuration during re-initialization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.Initialization(opts)
	},
}

func init() {
	initCmd.Flags().BoolVar(
		&opts.Force,
		"force",
		false,
		"Force re-initialization by overwriting existing configuration",
	)

	initCmd.Flags().IntVar(
		&opts.BackendPort,
		domain.Backend,
		0,
		"Backend service port (applied during initial setup or when --force is specified)",
	)

	initCmd.Flags().IntVar(
		&opts.FrontendPort,
		domain.Frontend,
		0,
		"Frontend service port (applied during initial setup or when --force is specified)",
	)

	rootCmd.AddCommand(initCmd)
}
