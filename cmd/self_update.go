package cmd

import (
	"github.com/sarojkhanal/typegenctl/internal/app/usecase"
	"github.com/spf13/cobra"
)

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update typegenctl to the latest version",
	Long: `Update the typegenctl CLI to the latest released version from GitHub.

This command downloads and installs the newest version of typegenctl.
It requires sudo privileges to replace the existing binary in /usr/local/bin.

Examples:
  typegenctl self-update
  typegenctl self-update --check
  typegenctl self-update --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.SelfUpdate(opts)
	},
}

func init() {
	selfUpdateCmd.Flags().BoolVar(&opts.CheckUpdate, "check", false, "Check for available updates without installing")
	selfUpdateCmd.Flags().BoolVar(&opts.Force, "force", false, "Reinstall even if already on the latest version")

	rootCmd.AddCommand(selfUpdateCmd)
}
