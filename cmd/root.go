package cmd

import (
	"fmt"

	"github.com/sarojkhanal/typegenctl/internal/app"
	"github.com/sarojkhanal/typegenctl/internal/version"
	"github.com/spf13/cobra"
)

var opts = &app.Options{}

var rootCmd = &cobra.Command{
	Use:   "typegenctl",
	Short: "Unified control plane for TypeGen services",
	Long: `typegenctl is a production-ready command-line interface designed to orchestrate, validate, and operate the TypeGen ecosystem.

It provides deterministic lifecycle management for frontend and backend services, configuration validation, runtime status inspection, and safe
execution of operational workflows. The CLI enforces explicit intent, predictable behavior, and clear feedback to support both local development
and enterprise deployment environments.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if opts.ShowVersion {
			fmt.Printf("TypeGen CLI (typegenctl)\n")
			fmt.Printf("  Version: %s\n", version.String())
			return nil
		}
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&opts.ConfigPath,
		"config",
		"/opt/typegen/config/typegen.yaml",
		"Path to system configuration file",
	)

	rootCmd.PersistentFlags().BoolVar(
		&opts.JSONOutput,
		"json",
		false,
		"Output results in JSON format",
	)

	rootCmd.PersistentFlags().BoolVar(
		&opts.DryRun,
		"dry-run",
		false,
		"Show what actions would be performed without executing them",
	)

	rootCmd.Flags().BoolVarP(
		&opts.ShowVersion,
		"version",
		"v",
		false,
		"Print version information",
	)
}
