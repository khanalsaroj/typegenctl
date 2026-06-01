package cmd

import (
	"github.com/khanalsaroj/typegenctl/internal/app/usecase"
	"github.com/spf13/cobra"
)

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update typegenctl to the latest version",
	Long: `Update the typegenctl CLI to the latest released version from GitHub.

This command downloads the release archive for your platform, verifies it
against the published SHA-256 checksums, and then replaces the running binary
in place. If the binary lives in a protected directory (e.g. /usr/local/bin on
Linux/macOS), re-run with elevated privileges (sudo).

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
