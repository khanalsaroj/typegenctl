package cmd

import (
	"github.com/sarojkhanal/typegenctl/internal/app/usecase"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Open Typegen dashboard in browser",
	Long:  "Opens the Typegen dashboard in your default web browser.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return usecase.Dashboard(opts)
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)

}
