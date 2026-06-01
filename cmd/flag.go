package cmd

import (
	"github.com/khanalsaroj/typegenctl/internal/domain"
	"github.com/spf13/cobra"
)

func addServiceFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(
		&opts.BackendOnly,
		domain.Backend,
		false,
		"Apply the command to the backend service only",
	)

	cmd.Flags().BoolVar(
		&opts.FrontendOnly,
		domain.Frontend,
		false,
		"Apply the command to the frontend service only",
	)
}
