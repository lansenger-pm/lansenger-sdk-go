package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Health check and connection verification",
}

var healthCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check Lansenger API connection health",
	Args:  cobra.NoArgs,
	Run:   runHealthCheck,
}

func init() {
	healthCmd.AddCommand(healthCheckCmd)
	rootCmd.AddCommand(healthCmd)
}

func runHealthCheck(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	ok := client.HealthCheck(ctx)
	if ok {
		result := map[string]interface{}{
			"status":  "OK",
			"message": "Lansenger connection is healthy",
		}
		if jsonOutput {
			outputResult(result)
			return
		}
		fmt.Println("OK — Lansenger connection is healthy")
		return
	}

	result := map[string]interface{}{
		"status":  "FAIL",
		"message": "Lansenger connection is not healthy",
	}
	if jsonOutput {
		outputResult(result)
		return
	}
	fmt.Println("FAIL — Lansenger connection is not healthy")
	os.Exit(1)
}