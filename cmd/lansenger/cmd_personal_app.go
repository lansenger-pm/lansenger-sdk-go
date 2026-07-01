package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var personalAppCmd = &cobra.Command{
	Use:   "personal-app",
	Short: "Manage personal apps/bots (4.38) - requires user_token",
}

// create
var paCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a personal app/bot",
	Args:  cobra.NoArgs,
	Run:   runPACreate,
}

var paCreateUserToken string
var paCreateName string
var paCreateAvatarID string
var paCreateDesc string

// update
var paUpdateCmd = &cobra.Command{
	Use:   "update APP_ID NAME",
	Short: "Update a personal app/bot",
	Args:  cobra.ExactArgs(2),
	Run:   runPAUpdate,
}

var paUpdateUserToken string
var paUpdateAvatarID string
var paUpdateDesc string

// info
var paFetchCmd = &cobra.Command{
	Use:   "info APP_ID",
	Short: "Fetch personal app info",
	Args:  cobra.ExactArgs(1),
	Run:   runPAFetch,
}

var paFetchUserToken string

// delete
var paDeleteCmd = &cobra.Command{
	Use:   "delete APP_ID",
	Short: "Delete a personal app/bot",
	Args:  cobra.ExactArgs(1),
	Run:   runPADelete,
}

var paDeleteUserToken string

// list
var paListCmd = &cobra.Command{
	Use:   "list",
	Short: "List personal apps/bots",
	Args:  cobra.NoArgs,
	Run:   runPAList,
}

var paListUserToken string

func init() {
	paCreateCmd.Flags().StringVar(&paCreateUserToken, "user-token", "", "User token (required)")
	paCreateCmd.Flags().StringVar(&paCreateName, "name", "", "App name")
	paCreateCmd.Flags().StringVar(&paCreateAvatarID, "avatar-id", "", "Avatar media ID")
	paCreateCmd.Flags().StringVarP(&paCreateDesc, "desc", "d", "", "App description")

	paUpdateCmd.Flags().StringVar(&paUpdateUserToken, "user-token", "", "User token (required)")
	paUpdateCmd.Flags().StringVar(&paUpdateAvatarID, "avatar-id", "", "Avatar media ID")
	paUpdateCmd.Flags().StringVarP(&paUpdateDesc, "desc", "d", "", "App description")

	paFetchCmd.Flags().StringVar(&paFetchUserToken, "user-token", "", "User token (required)")

	paDeleteCmd.Flags().StringVar(&paDeleteUserToken, "user-token", "", "User token (required)")

	paListCmd.Flags().StringVar(&paListUserToken, "user-token", "", "User token (required)")

	personalAppCmd.AddCommand(paCreateCmd)
	personalAppCmd.AddCommand(paUpdateCmd)
	personalAppCmd.AddCommand(paFetchCmd)
	personalAppCmd.AddCommand(paDeleteCmd)
	personalAppCmd.AddCommand(paListCmd)
	rootCmd.AddCommand(personalAppCmd)
}

func runPACreate(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.CreatePersonalApp(ctx, paCreateUserToken, paCreateName, paCreateAvatarID, paCreateDesc)
	checkError(err)
	outputResultFields(result, []string{"app_id", "secret", "apigw_addr", "passport_addr"})
}

func runPAUpdate(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.UpdatePersonalApp(ctx, args[0], paUpdateUserToken, args[1], paUpdateAvatarID, paUpdateDesc)
	checkError(err)
	outputResult(result)
}

func runPAFetch(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchPersonalApp(ctx, args[0], paFetchUserToken)
	checkError(err)
	outputResultFields(result, []string{"app_id", "name", "avatar_id", "description", "apigw_addr", "passport_addr"})
}

func runPADelete(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.DeletePersonalApp(ctx, args[0], paDeleteUserToken)
	checkError(err)
	outputResult(result)
}

func runPAList(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchPersonalAppList(ctx, paListUserToken)
	checkError(err)
	if result.Success && !jsonOutput {
		for _, app := range result.AppList {
			fmt.Printf("  %s  %s  %s\n",
				strFromMapPrint(app, "appId"),
				strFromMapPrint(app, "appName"),
				strFromMapPrint(app, "description"))
		}
	} else {
		outputResult(result)
	}
}

func strFromMapPrint(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
