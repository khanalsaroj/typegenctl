package util

import (
	"encoding/json"
	"fmt"

	"github.com/khanalsaroj/typegenctl/internal/result"
)

func PrintJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func PrintHuman(r *result.ValidationResult) {
	if r == nil {
		fmt.Println("No validation results available.")
		return
	}

	if len(r.Success) == 0 && len(r.Info) == 0 && len(r.Failure) == 0 {
		fmt.Println()
		fmt.Println("TypeGen System")
		fmt.Println("--------------------")
		fmt.Println("No issues detected. System is already in a clean state.")
		fmt.Println()
		return
	}

	fmt.Println("TypeGen System Check")
	fmt.Println("--------------------")

	if len(r.Success) > 0 {
		fmt.Println("\nSUCCESS")
		for _, msg := range r.Success {
			fmt.Printf("  ✓ %s\n", msg)
		}
	}

	if len(r.Info) > 0 {
		fmt.Println("\nINFO")
		for _, msg := range r.Info {
			fmt.Printf("  ℹ %s\n", msg)
		}
	}

	if len(r.Failure) > 0 {
		fmt.Println("\nFAILURE")
		for _, err := range r.Failure {
			fmt.Printf("  ✗ %s\n", err.Error())
		}
	}

	fmt.Println()
}
